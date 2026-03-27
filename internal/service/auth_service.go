package service

import (
	"backend/internal/domain/user"
	"backend/internal/dto"
	"backend/internal/infrastructure/security/jwt"
	"backend/internal/repository"
	"backend/pkg/log"
	"backend/pkg/utils"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthServiceInterface interface {
	VerifyFirebaseIDToken(idToken string) (*auth.Token, error)
	HandleOAuth(idToken, ipAddress string) (*dto.AuthResponse, error)
	LoginManual(req dto.LoginRequest, ipAddress string) (*dto.AuthResponse, error)
	VerifyOTP(userID uuid.UUID, otp, ipAddress string) (*dto.AuthResponse, error)
	ResendOTP(userID uuid.UUID, ipAddress string) (*dto.ResendOTPResponse, error)
	Logout(tokenStr, ipAddress string) error
}

type AuthService struct {
	PermissionForUser   RolePermissionServiceInterface
	CompanyInfoForUser  CompanyServiceInterface
	EmployeeService     EmployeeServiceInterface
	LogAction           LogServiceInterface
	FirebaseClient      *auth.Client
	JwtService          jwt.JWTServiceInterface
	RedisClient         *redis.Client
	UserRepo            repository.UserRepository
	DefaultUserPassword string
}

func NewAuthService(permissionForUser RolePermissionServiceInterface, companyInfoForUser CompanyServiceInterface, employeeService EmployeeServiceInterface, logAction LogServiceInterface, jwtService jwt.JWTServiceInterface, rc *redis.Client, fc *auth.Client, userRepo repository.UserRepository, defaultPass string,
) AuthServiceInterface {
	return &AuthService{
		PermissionForUser:   permissionForUser,
		CompanyInfoForUser:  companyInfoForUser,
		EmployeeService:     employeeService,
		LogAction:           logAction,
		JwtService:          jwtService,
		RedisClient:         rc,
		FirebaseClient:      fc,
		UserRepo:            userRepo,
		DefaultUserPassword: defaultPass,
	}
}

func (s *AuthService) getOrCreateUser(ctx context.Context, email string) (*user.UserModel, bool, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	u, err := s.UserRepo.FindByEmail(ctx, email)
	if err == nil {
		return u, false, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, err
	}

	newUser := &user.UserModel{
		Email:           email,
		Password:        s.DefaultUserPassword,
		PasswordDefault: s.DefaultUserPassword,
		IsVerified:      false,
		NeedsProfile:    true,
	}

	err = s.UserRepo.Create(ctx, newUser)
	if err != nil {
		u2, err2 := s.UserRepo.FindByEmail(ctx, email)
		if err2 == nil {
			return u2, false, nil
		}
		return nil, false, err
	}

	return newUser, true, nil
}

func (s *AuthService) VerifyFirebaseIDToken(idToken string) (*auth.Token, error) {
	ctx := context.Background()

	token, err := s.FirebaseClient.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *AuthService) HandleOAuth(idToken, ipAddress string) (*dto.AuthResponse, error) {
	ctx := context.Background()

	fbToken, err := s.VerifyFirebaseIDToken(idToken)
	if err != nil {
		return nil, fmt.Errorf("invalid firebase token: %w", err)
	}

	email, ok := fbToken.Claims["email"].(string)
	if !ok || email == "" {
		return nil, fmt.Errorf("email not found in firebase token")
	}
	email = strings.ToLower(strings.TrimSpace(email))

	name, _ := fbToken.Claims["name"].(string)

	fullName := name
	if fullName == "" {
		parts := strings.Split(email, "@")
		if len(parts) > 0 {
			fullName = strings.Title(parts[0])
		} else {
			fullName = "User"
		}
	}

	currentUser, isNew, err := s.getOrCreateUser(ctx, email)
	if err != nil {
		return nil, err
	}

	if isNew {
		currentUser.NeedsProfile = true
		if err := s.UserRepo.Update(ctx, currentUser); err != nil {
			return nil, err
		}
	}

	otp, err := utils.GenerateOTP(6)
	if err != nil {
		return nil, fmt.Errorf("failed to generate OTP: %w", err)
	}

	if err := utils.SaveOTPToRedis(s.RedisClient, currentUser.ID, otp); err != nil {
		return nil, fmt.Errorf("failed to save OTP redis: %w", err)
	}

	go func(email, otpCode string) {
		if err := utils.SendOTPEmail(email, otpCode, utils.OTPTTL); err != nil {
			logger.Log.Error("Failed to send OTP email", zap.Error(err))
		}
	}(currentUser.Email, otp)

	s.LogAction.LogAction(ctx, currentUser.ID, uuid.Nil, "OAuth Login Requested", ipAddress)

	var (
		employeeID string
		roleName   string
		userName   = fullName
	)

	if len(currentUser.Employees) > 0 {
		e := currentUser.Employees[0]

		employeeID = e.ID.String()

		if e.FullName != "" {
			userName = e.FullName
		}

		if e.Role != nil {
			roleName = e.Role.Name
		}
	}

	perms, err := s.PermissionForUser.GetPermissionsForUser(ctx, currentUser)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions: %w", err)
	}

	// 🔹 10. Get Company Info (SAFE)
	var companyInfo *dto.CompanyInfo

	if currentUser.CompanyID != nil {
		companyInfo, err = s.CompanyInfoForUser.GetCompanyForUser(ctx, currentUser)
		if err != nil {
			logger.Log.Warn("failed to get company info", zap.Error(err))
			companyInfo = nil
		}
	}

	return &dto.AuthResponse{
		AccessToken:  "",
		RefreshToken: "",
		OtpRequest:   true,
		NeedsProfile: currentUser.NeedsProfile,
		UserInfo: &dto.UserInfo{
			ID:         currentUser.ID.String(),
			EmployeeID: employeeID,
			Name:       userName,
			Email:      currentUser.Email,
			Role:       roleName,
			Permission: perms,
			Verified:   currentUser.IsVerified,
		},
		CompanyInfo: companyInfo,
	}, nil
}

func (s *AuthService) LoginManual(req dto.LoginRequest, ipAddress string) (*dto.AuthResponse, error) {
	ctx := context.Background()

	email := strings.ToLower(req.Email)

	user, err := s.UserRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	isMatch := false

	if req.Password == user.PasswordDefault {
		isMatch = true
	} else {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err == nil {
			isMatch = true
		}
	}

	if !isMatch {
		logger.Log.Fatal("Invalid email or password",
			zap.String("user_id", user.ID.String()),
			zap.String("password", req.Password),
		)
	}

	otp, err := utils.GenerateOTP(6)
	if err != nil {
		logger.Log.Fatal("Failed to generate OTP",
			zap.Error(err),
		)
	}

	if err := utils.SaveOTPToRedis(s.RedisClient, user.ID, otp); err != nil {
		logger.Log.Fatal("Failed to save OTP redis",
			zap.Error(err),
		)
	}

	go func(email, otpCode, userID string) {
		if err := utils.SendOTPEmail(email, otpCode, utils.OTPTTL); err != nil {
			logger.Log.Fatal("Failed to send OTP email",
				zap.Error(err),
			)
		}
	}(user.Email, otp, user.ID.String())

	s.LogAction.LogAction(ctx, user.ID, uuid.Nil, "Manual Login Requested", ipAddress)

	var (
		employeeID string
		roleName   string
		fullName   string
	)

	if len(user.Employees) > 0 {
		e := user.Employees[0]

		employeeID = e.ID.String()
		fullName = e.FullName

		if e.Role != nil {
			roleName = e.Role.Name
		}
	}

	perms, err := s.PermissionForUser.GetPermissionsForUser(ctx, user)
	if err != nil {
		logger.Log.Fatal("Failed to get permissions",
			zap.Error(err),
		)
	}

	companyInfo, err := s.CompanyInfoForUser.GetCompanyForUser(ctx, user)
	if err != nil {
		logger.Log.Fatal("Failed to get company info",
			zap.Error(err),
		)
	}

	return &dto.AuthResponse{
		AccessToken:  "",
		RefreshToken: "",
		OtpRequest:   true,
		NeedsProfile: user.NeedsProfile,
		UserInfo: &dto.UserInfo{
			ID:         user.ID.String(),
			EmployeeID: employeeID,
			Name:       fullName,
			Email:      user.Email,
			Role:       roleName,
			Permission: perms,
			Verified:   user.IsVerified,
		},
		CompanyInfo: companyInfo,
	}, nil
}

func (s *AuthService) VerifyOTP(userID uuid.UUID, otp, ipAddress string) (*dto.AuthResponse, error) {
	const maxAttempts = 5
	ctx := context.Background()

	// 🔹 1. Get OTP dari Redis
	otpData, err := utils.GetOTP(s.RedisClient, userID.String())
	if err != nil || otpData == nil {
		return nil, fmt.Errorf("otp expired or invalid")
	}

	// 🔹 2. Check attempts
	if otpData.Attempts >= maxAttempts {
		return nil, fmt.Errorf("too many attempts")
	}

	// 🔹 3. Validate OTP
	if otpData.Code != otp {
		_ = utils.IncrementOTPAttempts(s.RedisClient, userID.String())
		return nil, fmt.Errorf("invalid otp")
	}

	// 🔹 4. Get user
	user, err := s.UserRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// 🔥 5. MARK VERIFIED
	if !user.IsVerified {
		if err := s.UserRepo.UpdateVerified(ctx, user.ID); err != nil {
			return nil, err
		}
		user.IsVerified = true
	}

	// 🔹 6. DELETE OTP
	if s.RedisClient != nil {
		_ = s.RedisClient.Del(ctx, "otp:"+user.ID.String()).Err()
	}

	// 🔥 7. CHECK USER STATUS (KEY LOGIC)

	// ==============================
	// 🟢 USER BARU → TEMP TOKEN
	// ==============================
	if user.NeedsProfile {

		tempToken, err := s.JwtService.GenerateTempToken(user.ID.String())
		if err != nil {
			return nil, err
		}

		s.LogAction.LogAction(ctx, user.ID, uuid.Nil, "OTP Verified (Pre-Register)", ipAddress)

		return &dto.AuthResponse{
			AccessToken:  tempToken,
			RefreshToken: "",
			OtpRequest:   false,
			NeedsProfile: true,
			UserInfo: &dto.UserInfo{
				ID:       user.ID.String(),
				Email:    user.Email,
				Verified: user.IsVerified,
			},
			CompanyInfo: nil,
		}, nil
	}

	// ==============================
	// 🟢 USER LAMA → FULL TOKEN
	// ==============================

	var (
		fullName   string
		roleName   string
		employeeID string
		roleID     uuid.UUID
		empID      uuid.UUID
	)

	if len(user.Employees) > 0 {
		e := user.Employees[0]

		fullName = e.FullName
		employeeID = e.ID.String()
		empID = e.ID

		if e.Role != nil {
			roleName = e.Role.Name
			roleID = e.RoleID
		}
	} else {
		fullName = user.Email
	}

	perms, err := s.PermissionForUser.GetPermissionsForUser(ctx, user)
	if err != nil {
		return nil, err
	}

	companyInfo, err := s.CompanyInfoForUser.GetCompanyForUser(ctx, user)
	if err != nil {
		companyInfo = nil
	}

	permissionNames := make([]string, len(perms))
	for i, p := range perms {
		permissionNames[i] = p.Name
	}

	companyIDStr := ""
	if user.CompanyID != nil {
		companyIDStr = user.CompanyID.String()
	}

	// 🔹 Access Token
	accessToken, err := s.JwtService.GenerateToken(jwt.CustomClaims{
		UserID:      user.ID.String(),
		Role:        roleName,
		CompanyID:   companyIDStr,
		RoleKey:     roleID.String(),
		EmployeeID:  empID.String(),
		Permissions: permissionNames,
	}, 24*time.Hour)
	if err != nil {
		return nil, err
	}

	// 🔹 Refresh Token
	refreshToken, err := s.JwtService.GenerateToken(jwt.CustomClaims{
		UserID: user.ID.String(),
		Role:   roleName,
	}, 7*24*time.Hour)
	if err != nil {
		return nil, err
	}

	s.LogAction.LogAction(ctx, user.ID, uuid.Nil, "Login Success", ipAddress)

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		OtpRequest:   false,
		NeedsProfile: false,
		UserInfo: &dto.UserInfo{
			ID:         user.ID.String(),
			EmployeeID: employeeID,
			Name:       fullName,
			Email:      user.Email,
			Role:       roleName,
			Permission: perms,
			Verified:   user.IsVerified,
		},
		CompanyInfo: companyInfo,
	}, nil
}

func (s *AuthService) ResendOTP(userID uuid.UUID, ipAddress string) (*dto.ResendOTPResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := s.UserRepo.FindByID(ctx, userID)
	if err != nil {
		logger.Log.Fatal("Failed to find user",
			zap.Error(err),
		)
	}

	// ===== Rate Limit (Redis) =====
	limitKey := "otp:resend:" + userID.String()

	count, err := s.RedisClient.Get(ctx, limitKey).Int()
	if err != nil && err != redis.Nil {
		logger.Log.Fatal("Failed to get OTP limit",
			zap.Error(err),
		)
	}

	if count >= 3 {
		logger.Log.Fatal("OTP limit exceeded",
			zap.Int("count", count),
		)
	}

	if err := s.RedisClient.Incr(ctx, limitKey).Err(); err != nil {
		logger.Log.Fatal("Failed to increment OTP limit",
			zap.Error(err),
		)
	}

	if err := s.RedisClient.Expire(ctx, limitKey, time.Minute).Err(); err != nil {
		logger.Log.Fatal("Failed to expire OTP limit",
			zap.Error(err),
		)
	}

	otp, err := utils.GenerateOTP(6)
	if err != nil {
		logger.Log.Fatal("Failed to generate OTP",
			zap.Error(err),
		)
	}

	if err := utils.SaveOTPToRedis(s.RedisClient, userID, otp); err != nil {
		logger.Log.Fatal("Failed to save OTP",
			zap.Error(err),
		)
	}

	s.LogAction.LogAction(ctx, userID, uuid.Nil, "Resend OTP", ipAddress)

	go func(email, otpCode, uid string) {
		defer func() {
			if r := recover(); r != nil {
				logger.Log.Fatal("Failed to send OTP email",
					zap.Error(fmt.Errorf("%v", r)),
				)
			}
		}()

		if err := utils.SendOTPEmail(email, otpCode, utils.OTPTTL); err != nil {
			logger.Log.Fatal("Failed to send OTP email",
				zap.Error(err),
			)
		}
	}(user.Email, otp, user.ID.String())

	return &dto.ResendOTPResponse{
		Message: "OTP berhasil dikirim",
	}, nil
}

func (s *AuthService) Logout(tokenStr, ipAddress string) error {
	ctx := context.Background()

	claims, err := s.JwtService.ValidateToken(tokenStr)
	if err != nil {
		logger.Log.Fatal("Failed to validate token",
			zap.Error(err),
		)
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		logger.Log.Fatal("Failed to parse user ID",
			zap.Error(err),
		)
	}

	key := "blacklist:" + claims.ID

	exists, err := s.RedisClient.Exists(ctx, key).Result()
	if err != nil {
		logger.Log.Fatal("Failed to check blacklist token",
			zap.Error(err),
		)
	}

	if exists > 0 {
		logger.Log.Fatal("Token already blacklisted",
			zap.String("token_id", claims.ID),
		)
	}

	ttl := time.Until(claims.ExpiresAt.Time)
	if ttl <= 0 {
		logger.Log.Fatal("Token already expired",
			zap.String("token_id", claims.ID),
		)
	}

	if err := s.RedisClient.Set(ctx, key, "true", ttl).Err(); err != nil {
		logger.Log.Fatal("Failed to blacklist token",
			zap.Error(err),
		)
	}

	user, err := s.UserRepo.FindByID(ctx, userID)
	if err == nil {
		var companyID uuid.UUID
		if user.CompanyID != nil {
			companyID = *user.CompanyID
		}

		s.LogAction.LogAction(ctx, userID, companyID, "Logout", ipAddress)
	}

	return nil
}

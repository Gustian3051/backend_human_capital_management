package service

import (
	"backend/internal/domain/employee"
	"backend/internal/dto"
	"backend/internal/infrastructure/database"
	"backend/internal/infrastructure/security/jwt"
	"backend/internal/repository"
	"backend/pkg/utils"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/google/uuid"
	redis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type RegisterServiceInterface interface {
	Register(ctx context.Context, req dto.RegisterRequest, ipAddress string) (*dto.AuthResponse, error)
}

type RegisterService struct {
	DB                   *gorm.DB
	JwtService           jwt.JWTServiceInterface
	RedisClient          *redis.Client
	UserRepo             repository.UserRepository
	CompanyRepo          repository.CompanyRepository
	EmployeeRepo         repository.EmployeeRepository
	LogAction            LogServiceInterface
	Enforcer             *casbin.Enforcer
	RolePermissionSeeder database.RolePermissionSeederInterface
}

func NewRegisterService(db *gorm.DB, jwtService jwt.JWTServiceInterface, redisClient *redis.Client, userRepo repository.UserRepository, companyRepo repository.CompanyRepository, employeeRepo repository.EmployeeRepository, logAction LogServiceInterface, enforcer *casbin.Enforcer, rolePermissionSeeder database.RolePermissionSeederInterface) *RegisterService {
	return &RegisterService{
		DB:                   db,
		JwtService:           jwtService,
		RedisClient:          redisClient,
		UserRepo:             userRepo,
		CompanyRepo:          companyRepo,
		EmployeeRepo:         employeeRepo,
		LogAction:            logAction,
		Enforcer:             enforcer,
		RolePermissionSeeder: rolePermissionSeeder,
	}
}


func (s *RegisterService) Register(
	ctx context.Context,
	req dto.RegisterRequest,
	ipAddress string,
) (*dto.AuthResponse, error) {

	// 🔹 1. Claims validation
	claims, ok := ctx.Value("claims").(*jwt.CustomClaims)
	if !ok || claims == nil {
		return nil, fmt.Errorf("unauthorized")
	}

	if claims.Role != "pre-register" {
		return nil, fmt.Errorf("invalid token for register")
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id")
	}

	enforcer := s.Enforcer
	var resp *dto.AuthResponse

	err = s.CompanyRepo.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		user, err := s.UserRepo.FindByID(ctx, userID)
		if err != nil {
			return fmt.Errorf("user not found: %w", err)
		}

		if !user.NeedsProfile {
			return fmt.Errorf("user already completed profile")
		}

		company, err := s.CompanyRepo.GetOrCreate(tx, req.CompanyName)
		if err != nil {
			return err
		}

		adminRole, err := s.RolePermissionSeeder.SeedDefaultRolesTx(tx, company.ID)
		if err != nil {
			return err
		}

		if adminRole.ID == uuid.Nil {
			return fmt.Errorf("invalid admin role")
		}

		if err := s.RolePermissionSeeder.AssignPermissionsToAdminTx(tx, adminRole.ID); err != nil {
			return err
		}
		user.CompanyID = &company.ID
		user.NeedsProfile = false

		if err := s.UserRepo.Update(ctx, user); err != nil {
			return err
		}

		existingEmp, err := s.EmployeeRepo.FindByUserAndCompanyTx(tx, user.ID, company.ID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		var employeeModel *employee.EmployeeModel

		fullName := req.Name
		if fullName == "" {
			if len(user.Employees) > 0 && user.Employees[0].FullName != "" {
				fullName = user.Employees[0].FullName
			}
		}
		if fullName == "" {
			fullName = strings.Split(user.Email, "@")[0]
		}

		if existingEmp != nil {
			employeeModel = existingEmp
		} else {
			newEmp := &employee.EmployeeModel{
				ID:        uuid.New(),
				UserID:    user.ID,
				CompanyID: company.ID,
				FullName:  fullName,
				RoleID:    adminRole.ID,
			}

			if err := s.EmployeeRepo.CreateTx(tx, newEmp); err != nil {
				return err
			}

			employeeModel = newEmp
		}

		// 🔹 8. Update employee count
		totalEmp, err := s.CompanyRepo.CountEmployees(tx, company.ID)
		if err != nil {
			return err
		}

		if err := tx.Model(company).
			Update("employee_count", totalEmp).Error; err != nil {
			return err
		}

		// 🔹 9. Casbin
		roleKey := "role:admin_" + company.ID.String()
		userKey := "user:" + user.ID.String()
		domain := "company:" + company.ID.String()

		if !enforcer.HasPolicy(roleKey, domain, "*", "*") {
			if _, err := enforcer.AddPolicy(roleKey, domain, "*", "*"); err != nil {
				return err
			}
		}

		if !enforcer.HasGroupingPolicy(userKey, roleKey, domain) {
			if _, err := enforcer.AddGroupingPolicy(userKey, roleKey, domain); err != nil {
				return err
			}
		}

		if err := enforcer.SavePolicy(); err != nil {
			return err
		}

		// 🔹 10. Generate token
		token, err := s.JwtService.GenerateToken(jwt.CustomClaims{
			UserID:      user.ID.String(),
			CompanyID:   company.ID.String(),
			Role:        "admin",
			RoleKey:     roleKey,
			EmployeeID:  employeeModel.ID.String(),
			Permissions: []string{"*"},
		}, 24*time.Hour)
		if err != nil {
			return err
		}

		// 🔹 11. Cleanup
		if s.RedisClient != nil {
			_ = s.RedisClient.Del(ctx, "user_register:"+user.ID.String()).Err()
		}

		s.LogAction.LogAction(ctx, user.ID, company.ID, "Complete Profile", ipAddress)

		// 🔹 12. Response
		resp = &dto.AuthResponse{
			AccessToken:  token,
			NeedsProfile: false,
			UserInfo: &dto.UserInfo{
				ID:         user.ID.String(),
				EmployeeID: employeeModel.ID.String(),
				Name:       employeeModel.FullName,
				Email:      user.Email,
				Role:       "admin",
				Verified:   user.IsVerified,
			},
			CompanyInfo: &dto.CompanyInfo{
				ID:            company.ID.String(),
				CompanyName:   company.Name,
				IsActive:      true,
				EmployeeCount: int(totalEmp),
				StartDate:     time.Now().Format("2006-01-02"),
				EndDate:       time.Now().AddDate(0, 0, 30).Format("2006-01-02"),
				PackageName:   "Basic",
			},
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// async email
	go func() {
		_ = utils.SendCompanyNotificationEmail(*resp.CompanyInfo, resp.UserInfo.Email)
	}()

	return resp, nil
}
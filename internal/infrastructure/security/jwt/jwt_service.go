package jwt

import (
	logger "backend/pkg/log"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service struct {
	secret []byte
}

func NewJWTService(secret string) *Service {
	return &Service{
		secret: []byte(secret),
	}
}

func (s *Service) GenerateToken(claims CustomClaims, duration time.Duration) (string, error) {
	claims.ID = uuid.NewString()
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(duration))
	claims.IssuedAt = jwt.NewNumericDate(time.Now())

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *Service) ValidateToken(tokenStr string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logger.Log.Fatal("Failed to validate JWT",
				zap.Error(errors.New("invalid signing method")),
			)
		}
		return s.secret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		logger.Log.Fatal("Failed to validate JWT",
			zap.Error(errors.New("invalid token")),
		)
	}

	return claims, nil
}

func (s *Service) GenerateResetToken(userID, email string) (string, error) {
	claims := ResetPasswordClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "reset-password",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

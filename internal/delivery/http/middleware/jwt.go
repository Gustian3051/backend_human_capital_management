package middleware

import (
	"context"
	"net/http"
	"strings"

	"backend/pkg/log"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	jwtinfra "backend/internal/infrastructure/security/jwt"
)

func JWTMiddleware(jwtService *jwtinfra.Service, redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		// ===== Get Authorization Header =====
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Log.Warn("Missing authorization header")

			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			logger.Log.Warn("Invalid authorization header")

			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header",
			})
			c.Abort()
			return
		}

		// ===== Validate JWT =====
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			logger.Log.Warn("Invalid token",
				zap.Error(err),
			)

			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			c.Abort()
			return
		}

		// ===== Check Token Blacklist (Redis) =====
		if redisClient != nil {
			ctx := context.Background()

			key := "blacklist:" + claims.ID
			exists, err := redisClient.Exists(ctx, key).Result()
			if err != nil {
				logger.Log.Error("Failed to check token blacklist",
					zap.Error(err),
				)

				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "internal server error",
				})
				c.Abort()
				return
			}

			if exists > 0 {
				logger.Log.Warn("Token has been revoked",
					zap.String("jti", claims.ID),
				)

				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "token revoked",
				})
				c.Abort()
				return
			}
		}

		// ===== Inject Context =====
		c.Set("user_id", claims.UserID)
		c.Set("employee_id", claims.EmployeeID)
		c.Set("company_id", claims.CompanyID)
		c.Set("role_id", claims.RoleID)
		c.Set("role", claims.Role)
		c.Set("permissions", claims.Permissions)
		c.Set("jti", claims.ID)
		c.Set("access_token", tokenString)

		c.Next()
	}
}
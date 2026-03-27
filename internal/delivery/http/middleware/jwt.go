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

		// ===== 1. Get Authorization Header =====
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header format",
			})
			return
		}

		tokenString := parts[1]

		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			logger.Log.Warn("Invalid token", zap.Error(err))

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}

		if redisClient != nil {
			ctx := context.Background()

			key := "blacklist:" + claims.ID
			exists, err := redisClient.Exists(ctx, key).Result()
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "internal server error",
				})
				return
			}

			if exists > 0 {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "token revoked",
				})
				return
			}
		}

		path := c.FullPath()

		if claims.Role == "pre-register" {
			if path != "/api/v1/register" {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"error": "access denied: complete profile required",
				})
				return
			}
		}

		if claims.Role != "pre-register" && path == "/api/v1/register" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "register not allowed",
			})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("employee_id", claims.EmployeeID)
		c.Set("company_id", claims.CompanyID)
		c.Set("role_key", claims.RoleKey)
		c.Set("role", claims.Role)
		c.Set("permissions", claims.Permissions)
		c.Set("jti", claims.ID)
		c.Set("access_token", tokenString)

		ctx := context.WithValue(c.Request.Context(), "claims", claims)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
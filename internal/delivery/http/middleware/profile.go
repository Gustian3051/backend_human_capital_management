package middleware

import (
	"context"
	"net/http"
	"time"

	"backend/pkg/log"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func CheckProfileMiddleware(db *gorm.DB, redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		// ===== Get user from context (JWT Middleware) =====
		userID := c.GetString("user_id")
		if userID == "" {
			logger.Log.Warn("Unauthorized: missing user_id")

			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			c.Abort()
			return
		}

		ctx := context.Background()
		cacheKey := "user_profile_complete:" + userID

		// ===== Check Redis Cache =====
		val, err := redisClient.Get(ctx, cacheKey).Result()

		if err == redis.Nil {
			// ===== Cache miss → check DB =====

			var needsProfile bool

			err := db.Table("users").
				Select("needs_profile").
				Where("id = ?", userID).
				Scan(&needsProfile).Error

			if err != nil {
				logger.Log.Error("Failed to check profile from DB",
					zap.Error(err),
					zap.String("user_id", userID),
				)

				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "failed to check profile",
				})
				c.Abort()
				return
			}

			if needsProfile {
				// cache result
				_ = redisClient.Set(ctx, cacheKey, "false", time.Hour).Err()

				logger.Log.Warn("Profile not complete",
					zap.String("user_id", userID),
				)

				c.JSON(http.StatusForbidden, gin.H{
					"error": "complete your profile first",
				})
				c.Abort()
				return
			}

			// cache success
			_ = redisClient.Set(ctx, cacheKey, "true", time.Hour).Err()

		} else if err != nil {
			// ===== Redis error =====
			logger.Log.Error("Redis error",
				zap.Error(err),
			)

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
			c.Abort()
			return

		} else if val == "false" {
			// ===== Cached: profile incomplete =====
			logger.Log.Warn("Profile not complete (cached)",
				zap.String("user_id", userID),
			)

			c.JSON(http.StatusForbidden, gin.H{
				"error": "complete your profile first",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
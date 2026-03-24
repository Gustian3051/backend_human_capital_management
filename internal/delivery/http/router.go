package http

import (
	"backend/config"
	"backend/internal/delivery/http/handler"
	"backend/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
	"github.com/casbin/casbin/v2"

	jwtinfra "backend/internal/infrastructure/security/jwt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(cfg *config.Config, enforcer *casbin.Enforcer, db *gorm.DB, redisClient *redis.Client, jwtService *jwtinfra.Service) *gin.Engine {
	r := gin.Default()

	// Middleware
	r.Use(middleware.CORSMiddleware())

	// Health Check (public)
	r.GET("/health", handler.HealthCheck)

	// Swagger (hanya dev)
	if cfg.App.Env != "production" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// API versioning
	api := r.Group("/api/v1")
	api.Use(
		middleware.JWTMiddleware(jwtService, redisClient),
		middleware.CheckProfileMiddleware(db, redisClient),
		middleware.RBACMiddleware(enforcer),
	)

	// endpoint
	// {
	// api.GET("/users/:id", userHandler.GetUser)
	// }

	return r
}
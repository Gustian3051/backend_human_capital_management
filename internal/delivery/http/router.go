package http

import (
	"backend/config"
	"backend/internal/delivery/http/handler"
	"backend/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(cfg *config.Config) *gin.Engine {
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
	// api := r.Group("/api/v1")
	// {
	
	// endpoint
	// api.GET("/users/:id", userHandler.GetUser)
	// }

	return r
}
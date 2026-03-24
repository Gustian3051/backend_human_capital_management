package main

import (
	"context"
	"time"

	"backend/config"
	"backend/internal/delivery/http"
	"backend/internal/infrastructure/database"
	firebaseinfra "backend/internal/infrastructure/firebase"
	redisinfra "backend/internal/infrastructure/redis"
	storageinfra "backend/internal/infrastructure/storage"
	"backend/pkg/log"

	"go.uber.org/zap"
)

// @title           HCM Backend API
// @version         1.0
// @description     Human Capital Management API Documentation
// @host            localhost:8080
// @BasePath        /api/v1
func main() {
	// ===== Context =====
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// ===== Init Logger =====
	logger.Init(true)
	defer logger.Sync()

	// ===== Load Config =====
	cfg := config.LoadConfig()

	logger.Init(cfg.App.Debug)

	logger.Log.Info("Starting application",
		zap.String("env", cfg.App.Env),
		zap.String("app", cfg.App.Name),
	)

	// ===== Infrastructure =====
	db := database.NewPostgresDB(cfg.DB)
	rdb := redisinfra.NewRedisClient(ctx, cfg.Redis)
	minioClient := storageinfra.NewMinioClient(ctx, cfg.Storage)
	firebaseApp := firebaseinfra.NewFirebaseApp(ctx, cfg.Firebase)

	// ===== Log Infrastructure Status =====
	logger.Log.Info("Infrastructure initialized",
		zap.Bool("database", db != nil),
		zap.Bool("redis", rdb != nil),
		zap.Bool("storage", minioClient != nil),
		zap.Bool("firebase", firebaseApp != nil),
	)

	// ===== Router =====
	router := http.NewRouter(cfg)

	// ===== Run Server =====
	addr := ":" + cfg.App.Port

	logger.Log.Info("Server running",
		zap.String("address", addr),
	)

	if err := router.Run(addr); err != nil {
		logger.Log.Fatal("Failed to start server",
			zap.Error(err),
		)
	}

	// ===== Cleanup =====
	defer func() {
		if rdb != nil {
			_ = rdb.Close()
		}
	}()
}

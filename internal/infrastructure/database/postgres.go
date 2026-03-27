package database

import (
	"fmt"
	"time"

	"backend/config"
	"backend/pkg/log"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB(cfg config.DBConfig) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.Port,
		cfg.SSLMode,
	)

	var db *gorm.DB
	var err error

	for i := 0; i < 10; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}

		logger.Log.Warn("Database not ready, retrying...",
			zap.Int("attempt", i+1),
			zap.Error(err),
		)

		time.Sleep(2 * time.Second)
	}

	if err != nil {
		logger.Log.Fatal("Failed to connect to database after retries", zap.Error(err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Log.Fatal("Failed to get sqlDB:", zap.Error(err))
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(2 * time.Minute)

	for i := 0; i < 5; i++ {
		if err := sqlDB.Ping(); err == nil {
			break
		}
		logger.Log.Warn("Ping failed, retrying...", zap.Int("attempt", i+1))
		time.Sleep(1 * time.Second)
	}

	logger.Log.Info("Database connected successfully")

	return db
}
package redis

import (
	"context"
	"fmt"
	"time"

	"backend/config"
	"backend/pkg/log"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func NewRedisClient(ctx context.Context, cfg config.RedisConfig) *redis.Client {
	if !cfg.Enabled {
		logger.Log.Warn("Redis disabled")
		return nil
	}

	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		logger.Log.Fatal("Failed to connect to Redis",
			zap.Error(err),
			zap.String("addr", addr),
		)
	}

	logger.Log.Info("Redis connected",
		zap.String("addr", addr),
	)

	return client
}
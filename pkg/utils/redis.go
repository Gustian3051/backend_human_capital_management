package utils

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"backend/pkg/log"
)


type OTPData struct {
    Code     string
    Attempts int
}

func SaveOTPToRedis(rdb *redis.Client, userID uuid.UUID, otp string) error {
    ctx := context.Background()
    key := fmt.Sprintf("otp:%s", userID.String())

    err := rdb.HSet(ctx, key, map[string]interface{}{
        "code":     otp,
        "attempts": 0,
    }).Err()
    if err != nil {
		logger.Log.Error("Failed to save OTP", zap.Error(err))
    }

    if err := rdb.Expire(ctx, key, 10*time.Minute).Err(); err != nil {
        logger.Log.Error("Failed to set OTP TTL", zap.Error(err))
    }

    return nil
}

func GetOTP(client *redis.Client, userID string) (*OTPData, error) {
    key := "otp:" + userID
    data, err := client.HGetAll(context.Background(), key).Result()
    if err != nil || len(data) == 0 {
		logger.Log.Error("Failed to get OTP", zap.Error(err))
    }
    attempts, _ := strconv.Atoi(data["attempts"])
    return &OTPData{
        Code:     data["code"],
        Attempts: attempts,
    }, nil
}


func SaveOTP(client *redis.Client, userID string, code string) error {
    return client.HSet(context.Background(), "otp:"+userID, map[string]interface{}{
        "code":     code,
        "attempts": 0,
    }).Err()
}

func IncrementOTPAttempts(client *redis.Client, userID string) error {
    return client.HIncrBy(context.Background(), "otp:"+userID, "attempts", 1).Err()
}

func DeleteOTP(client *redis.Client, userID string) error {
    return client.Del(context.Background(), "otp:"+userID).Err()
}

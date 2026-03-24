package storage

import (
	"context"
	"time"

	"backend/config"
	"backend/pkg/log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

func NewMinioClient(ctx context.Context, cfg config.StorageConfig) *minio.Client {
	if cfg.Endpoint == "" {
		logger.Log.Fatal("MinIO endpoint is required")
	}

	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		logger.Log.Fatal("Failed to connect to MinIO",
			zap.Error(err),
		)
	}

	logger.Log.Info("MinIO connected",
		zap.String("endpoint", cfg.Endpoint),
	)

	ensureBucket(ctx, client, cfg.BucketName)

	return client
}

func ensureBucket(ctx context.Context, client *minio.Client, bucket string) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		logger.Log.Fatal("Failed to check bucket",
			zap.Error(err),
			zap.String("bucket", bucket),
		)
	}

	if !exists {
		err = client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			logger.Log.Fatal("Failed to create bucket",
				zap.Error(err),
				zap.String("bucket", bucket),
			)
		}

		logger.Log.Info("Bucket created",
			zap.String("bucket", bucket),
		)
	} else {
		logger.Log.Info("Bucket already exists",
			zap.String("bucket", bucket),
		)
	}
}
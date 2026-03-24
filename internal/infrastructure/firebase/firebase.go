package firebase

import (
	"context"
	"os"
	"time"

	"backend/config"
	"backend/pkg/log"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	"go.uber.org/zap"
)

func NewFirebaseApp(ctx context.Context, cfg config.FirebaseConfig) *firebase.App {
	if cfg.CredentialsFile == "" {
		logger.Log.Fatal("Firebase credentials file is required")
	}

	if _, err := os.Stat(cfg.CredentialsFile); err != nil {
		logger.Log.Fatal("Firebase credentials file not found",
			zap.String("path", cfg.CredentialsFile),
			zap.Error(err),
		)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	opt := option.WithCredentialsFile(cfg.CredentialsFile)

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		logger.Log.Fatal("Failed to initialize Firebase",
			zap.Error(err),
		)
	}

	logger.Log.Info("Firebase initialized")

	return app
}
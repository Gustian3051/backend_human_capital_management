package config

import (
	"strings"

	"backend/pkg/log"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	App      AppConfig
	DB       DBConfig
	Redis    RedisConfig
	Storage  StorageConfig
	SMTP     SMTPConfig
	Firebase FirebaseConfig
	JWT      JWTConfig
	Auth     AuthConfig
}

var cfg *Config

func LoadConfig() *Config {
	if cfg != nil {
		return cfg
	}

	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		logger.Log.Warn("No .env file found, using environment variables",
			zap.Error(err),
		)
	} else {
		logger.Log.Info("Config loaded from .env")
	}

	cfg = &Config{
		App:      LoadAppConfig(),
		DB:       LoadDBConfig(),
		Redis:    LoadRedisConfig(),
		Storage:  LoadStorageConfig(),
		SMTP:     LoadSMTPConfig(),
		Firebase: LoadFirebaseConfig(),
		JWT:      LoadJWTConfig(),
		Auth:     LoadAuthConfig(),
	}

	logger.Log.Info("Configuration initialized",
		zap.String("app_name", cfg.App.Name),
		zap.String("env", cfg.App.Env),
	)

	return cfg
}
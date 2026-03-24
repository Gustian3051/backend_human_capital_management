package config

import "github.com/spf13/viper"

type JWTConfig struct {
	Secret    string
	ExpiresIn string
}

func LoadJWTConfig() JWTConfig {
	return JWTConfig{
		Secret:    viper.GetString("JWT_SECRET"),
		ExpiresIn: viper.GetString("JWT_EXPIRES_IN"),
	}
}
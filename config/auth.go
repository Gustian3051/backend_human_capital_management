package config

import "github.com/spf13/viper"

type AuthConfig struct {
	DefaultPassword string
}

func LoadAuthConfig() AuthConfig {
	return AuthConfig{
		DefaultPassword: viper.GetString("DEFAULT_USER_PASSWORD"),
	}
}

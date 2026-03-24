package config

import "github.com/spf13/viper"

type RedisConfig struct {
	Enabled  bool
	Host     string
	Port     string
	Password string
	DB       int
}

func LoadRedisConfig() RedisConfig {
	return RedisConfig{
		Enabled:  viper.GetBool("REDIS_ENABLED"),
		Host:     viper.GetString("REDIS_HOST"),
		Port:     viper.GetString("REDIS_PORT"),
		Password: viper.GetString("REDIS_PASSWORD"),
		DB:       viper.GetInt("REDIS_DB"),
	}
}
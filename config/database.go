package config

import (
	"time"

	"github.com/spf13/viper"
)

type DBConfig struct {
	Driver          string
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

func LoadDBConfig() DBConfig {
	return DBConfig{
		Driver:          viper.GetString("DB_DRIVER"),
		Host:            viper.GetString("DB_HOST"),
		Port:            viper.GetString("DB_PORT"),
		User:            viper.GetString("DB_USER"),
		Password:        viper.GetString("DB_PASSWORD"),
		Name:            viper.GetString("DB_NAME"),
		SSLMode:         viper.GetString("DB_SSLMODE"),
		MaxOpenConns:    viper.GetInt("DB_MAX_OPEN_CONNS"),
		MaxIdleConns:    viper.GetInt("DB_MAX_IDLE_CONNS"),
		ConnMaxLifetime: viper.GetDuration("DB_CONN_MAX_LIFETIME"),
	}
}
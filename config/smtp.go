package config

import "github.com/spf13/viper"

type SMTPConfig struct {
	Host     string
	Port     string
	User     string
	Pass     string
	FromName string
}

func LoadSMTPConfig() SMTPConfig {
	return SMTPConfig{
		Host:     viper.GetString("SMTP_HOST"),
		Port:     viper.GetString("SMTP_PORT"),
		User:     viper.GetString("SMTP_USER"),
		Pass:     viper.GetString("SMTP_PASS"),
		FromName: viper.GetString("SMTP_FROM_NAME"),
	}
}
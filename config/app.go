
package config

import "github.com/spf13/viper"

type AppConfig struct {
	Name     string
	Env      string
	Port     string
	Debug    bool
	Timezone string
}

func LoadAppConfig() AppConfig {
	return AppConfig{
		Name:     viper.GetString("APP_NAME"),
		Env:      viper.GetString("APP_ENV"),
		Port:     viper.GetString("APP_PORT"),
		Debug:    viper.GetBool("APP_DEBUG"),
		Timezone: viper.GetString("APP_TIMEZONE"),
	}
}
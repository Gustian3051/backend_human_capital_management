package config

import "github.com/spf13/viper"

type StorageConfig struct {
	Provider   string
	Endpoint   string
	AccessKey  string
	SecretKey  string
	BucketName string
	PublicURL  string
	UseSSL     bool
}

func LoadStorageConfig() StorageConfig {
	return StorageConfig{
		Provider:   viper.GetString("STORAGE_PROVIDER"),
		Endpoint:   viper.GetString("MINIO_ENDPOINT"),
		AccessKey:  viper.GetString("MINIO_ACCESS_KEY"),
		SecretKey:  viper.GetString("MINIO_SECRET_KEY"),
		BucketName: viper.GetString("MINIO_BUCKET_NAME"),
		PublicURL:  viper.GetString("MINIO_PUBLIC_URL"),
		UseSSL:     viper.GetBool("MINIO_USE_SSL"),
	}
}
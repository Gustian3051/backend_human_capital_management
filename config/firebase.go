package config

import "github.com/spf13/viper"

type FirebaseConfig struct {
	CredentialsFile string
}

func LoadFirebaseConfig() FirebaseConfig {
	return FirebaseConfig{
		CredentialsFile: viper.GetString("FIREBASE_CREDENTIALS_FILE"),
	}
}
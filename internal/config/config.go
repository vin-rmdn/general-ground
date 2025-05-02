package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func SetupEnvironment() error {
	viper.SetConfigType("dotenv")
	viper.SetConfigName(".env")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read environment variables: %w", err)
	}

	return nil
}

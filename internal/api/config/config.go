package config

import (
	"errors"
	"time"

	"github.com/spf13/viper"
)

// Config configuration struct for loading application configuration.
type Config struct {
	Port            string
	DatabaseURL     string
	MaxConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
	Environment     string
	CertFile        string
	KeyFile         string
}

// Load Loads the configuration from a .env file in the root directory.
func Load() (*Config, error) {
	viper.SetConfigFile(".env")

	// actual reading takes place here. it's saved in viper.
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return nil, err
		}
	}

	viper.AutomaticEnv()

	return &Config{
		Port:            viper.GetString("PORT"),
		DatabaseURL:     viper.GetString("DATABASE_URL"),
		MaxConns:        viper.GetInt32("DB_MAX_CONNS"),
		MaxConnLifetime: viper.GetDuration("DB_CONN_LIFETIME"),
		MaxConnIdleTime: viper.GetDuration("DB_CONN_IDLETIME"),
		Environment:     viper.GetString("ENVIRONMENT"),
		CertFile:        viper.GetString("CERT_FILE"),
		KeyFile:         viper.GetString("KEY_FILE"),
	}, nil
}

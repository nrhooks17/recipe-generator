// Package config provides configuration management for the recipe generator application.
package config

import (
	"errors"
	"time"

	"github.com/spf13/viper"
)

// Config configuration struct for loading application configuration.
// It holds all the necessary configuration parameters for the application,
// including server settings and database connection parameters.
type Config struct {
	Port            string        // HTTP server port
	DatabaseURL     string        // Connection string for the database
	MaxConns        int32         // Maximum number of database connections
	MaxConnLifetime time.Duration // Maximum lifetime of a database connection
	MaxConnIdleTime time.Duration // Maximum idle time for a database connection
	Environment     string        // Application environment (development, production, etc.)
	CertFile        string        // Path to SSL certificate file
	KeyFile         string        // Path to SSL key file
}

// Load loads the configuration from a .env file in the root directory.
// It uses viper to read environment variables and configuration files,
// returning a populated Config struct or an error if loading fails.
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

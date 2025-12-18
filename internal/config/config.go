package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

// DatabaseConfig holds Postgres connection details
type DatabaseConfig struct {
	Host     string `mapstructure:"DB_HOST" validate:"required"`
	Port     string `mapstructure:"DB_PORT" validate:"required"`
	User     string `mapstructure:"DB_USER" validate:"required"`
	Password string `mapstructure:"DB_PASSWORD" validate:"required"`
	Name     string `mapstructure:"DB_NAME" validate:"required"`
	SSLMode  string `mapstructure:"DB_SSLMODE" validate:"required,oneof=disable require verify-full"`
}

// JWTConfig holds paths to the RSA keys and token settings
type JWTConfig struct {
	PrivateKeyPath string `mapstructure:"JWT_PRIVATE_KEY_PATH" validate:"required"`
	PublicKeyPath  string `mapstructure:"JWT_PUBLIC_KEY_PATH" validate:"required"`
	Issuer         string `mapstructure:"JWT_ISSUER" validate:"required"`
	ExpirationTime int    `mapstructure:"JWT_EXPIRATION_HOURS" validate:"required,min=1"`
}

// Config holds all configuration for the application
type Config struct {
	Environment string         `mapstructure:"ENVIRONMENT" validate:"required"`
	ServerPort  string         `mapstructure:"SERVER_PORT" validate:"required"`
	Database    DatabaseConfig `mapstructure:",squash"`
	JWT         JWTConfig      `mapstructure:",squash"`
}

// LoadConfig loads the configurations from the .env file
func LoadConfig() (config *Config, err error) {
	viper.AddConfigPath("./")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	// Defaults
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("JWT_ISSUER", "auth-service")
	viper.SetDefault("JWT_EXPIRATION_HOURS", 24)

	// Read Config
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Unmarshal
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	// Validate
	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return config, nil
}

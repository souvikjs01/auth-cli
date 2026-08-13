package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Database DatabaseConfig
	Auth     AuthConfig
	TOTP     TOTPConfig
}

type DatabaseConfig struct {
	User     string
	Password string
	Name     string
	Host     string
	Port     string
}

type AuthConfig struct {
	SessionTimeout   time.Duration
	MaxLoginAttempts int
	LockoutDuration  time.Duration
}

type TOTPConfig struct {
	Issuer string
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	// Read environment variables.
	viper.AutomaticEnv()

	_ = viper.ReadInConfig()

	sessionTimeout, err := time.ParseDuration(
		viper.GetString("SESSION_TIMEOUT"),
	)
	if err != nil {
		return nil, err
	}

	lockoutDuration, err := time.ParseDuration(
		viper.GetString("LOCKOUT_DURATION"),
	)
	if err != nil {
		return nil, err
	}

	return &Config{
		Database: DatabaseConfig{
			User:     viper.GetString("POSTGRES_USER"),
			Password: viper.GetString("POSTGRES_PASSWORD"),
			Name:     viper.GetString("POSTGRES_DB"),
			Host:     viper.GetString("POSTGRES_HOST"),
			Port:     viper.GetString("POSTGRES_PORT"),
		},

		Auth: AuthConfig{
			SessionTimeout:   sessionTimeout,
			MaxLoginAttempts: viper.GetInt("MAX_LOGIN_ATTEMPTS"),
			LockoutDuration:  lockoutDuration,
		},

		TOTP: TOTPConfig{
			Issuer: viper.GetString("TOTP_ISSUER"),
		},
	}, nil
}

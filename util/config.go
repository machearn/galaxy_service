package util

import (
	"time"

	"github.com/spf13/viper"
)

// Config is a struct that holds the configuration for the application.
type Config struct {
	DBDriver             string        `mapstructure:"DB_DRIVER"`
	DBSource             string        `mapstructure:"DB_SOURCE"`
	ServerAddress        string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	TestGmailAddress     string        `mapstructure:"TEST_GMAIL_ADDRESS"`
	TestGmailPassword    string        `mapstructure:"TEST_GMAIL_PASSWORD"`
	MigrateURL           string        `mapstructure:"MIGRATE_URL"`
	RedisAddress         string        `mapstructure:"REDIS_ADDRESS"`
}

func LoadConfig(configPath string) (Config, error) {
	var config Config
	v := viper.New()
	v.AddConfigPath(configPath)
	v.SetConfigName("app")
	v.SetConfigType("env")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return config, err
	}

	if err := v.Unmarshal(&config); err != nil {
		return config, err
	}

	return config, nil
}

package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	App   AppConfig
	DB    DBConfig
	Token TokenConfig
}

type AppConfig struct {
	Port string
	Env  string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type TokenConfig struct {
	AccessTokenSecret  string
	RefreshTokenSecret string
	TempTokenSecret    string
}

func (d DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Seoul",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode,
	)
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		App: AppConfig{
			Port: getEnv("APP_PORT", "8080"),
			Env:  getEnv("APP_ENV", "development"),
		},
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "sns_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Token: TokenConfig{
			AccessTokenSecret:  getEnv("ACCESS_TOKEN_SECRET", "accessTokenSecret"),
			RefreshTokenSecret: getEnv("REFRESH_TOKEN_SECRET", "refreshTokenSecret"),
			TempTokenSecret:    getEnv("TEMP_TOKEN_SECRET", "accessTokenSecret"),
		},
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

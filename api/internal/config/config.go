package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

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

// 로그 출력 시 패스워드 노출 방지
func (d DBConfig) String() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Name, d.SSLMode,
	)
}

func Load() (*Config, error) {
	loadDotEnv()

	var errs []error

	accessSecret, err := requiredEnv("ACCESS_TOKEN_SECRET")
	if err != nil {
		errs = append(errs, err)
	}
	refreshSecret, err := requiredEnv("REFRESH_TOKEN_SECRET")
	if err != nil {
		errs = append(errs, err)
	}
	tempSecret, err := requiredEnv("TEMP_TOKEN_SECRET")
	if err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	cfg := &Config{
		App: AppConfig{
			Port: getEnv("PORT", "8080"),
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
			AccessTokenSecret:  accessSecret,
			RefreshTokenSecret: refreshSecret,
			TempTokenSecret:    tempSecret,
		},
	}

	return cfg, nil
}

// 환경별 .env 파일 로드 (.env.development, .env.production 등)
// 파일이 없으면 무시하고, production은 실제 환경변수를 사용
func loadDotEnv() {
	env := os.Getenv("APP_ENV")

	log.Printf("env : %+v", env)

	if env == "production" {
		loadSecretsFromVolume("/var/secrets/planet")
		return
	}

	if env == "" {
		env = "develop"
	}

	envFile := fmt.Sprintf(".env.%s", env)
	if err := godotenv.Load(envFile); err != nil {
		_ = godotenv.Load(".env")
	}
}

func loadSecretsFromVolume(secretDir string) {
	entries, err := os.ReadDir(secretDir)
	if err != nil {
		log.Printf("failed to read secret directory: %v", err)
		return
	}
	files, _ := os.ReadDir("/")
	for _, f := range files {
		log.Println("/", f.Name())
	}

	files, _ = os.ReadDir("/var")
	for _, f := range files {
		log.Println("/var/", f.Name())
	}

	files, _ = os.ReadDir("/etc")
	for _, f := range files {
		log.Println("/etc/", f.Name())
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		path := filepath.Join(secretDir, entry.Name())
		log.Printf("path %v", path)
		data, err := os.ReadFile(path)
		if err != nil {
			log.Printf("failed to read secret %s: %v", entry.Name(), err)
			continue
		}
		log.Printf("data %v", data)
		value := strings.TrimSpace(string(data))
		log.Printf("value %v", value)

		if err := os.Setenv(entry.Name(), value); err != nil {
			log.Printf("failed to set env %s: %v", entry.Name(), err)
		}
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func requiredEnv(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return "", fmt.Errorf("required environment variable %q is not set", key)
	}
	return v, nil
}

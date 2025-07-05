package config

import (
	"os"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	DBTimeZone string

	AccessTokenSecret  string
	RefreshTokenSecret string
}

func LoadConfig() *Config {
	return &Config{
		DBHost:             getEnv("DB_HOST", "db"),
		DBPort:             getEnv("DB_PORT", "5432"),
		DBUser:             getEnv("DB_USER", "postgres"),
		DBPassword:         getEnv("DB_PASSWORD", "p4ssw0rd"),
		DBName:             getEnv("DB_NAME", "users"),
		DBSSLMode:          getEnv("DB_SSLMODE", "disable"),
		DBTimeZone:         getEnv("DB_TIMEZONE", "Europe/Moscow"),
		AccessTokenSecret:  getEnv("ACCESS_TOKEN_SECRET", "s3cr3t"),
		RefreshTokenSecret: getEnv("REFRESH_TOKEN_SECRET", "s3cr3t123"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

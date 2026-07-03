package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port          string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPass        string
	DBName        string
	DBSSL         string
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int
	SessionTTL    time.Duration
}

func Load() Config {
	return Config{
		Port:          getEnv("PORT", "8080"),
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBUser:        getEnv("DB_USER", "postgres"),
		DBPass:        getEnv("DB_PASSWORD", "postgres"),
		DBName:        getEnv("DB_NAME", "reservation_system"),
		DBSSL:         getEnv("DB_SSLMODE", "disable"),
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDB:       getEnvAsInt("REDIS_DB", 0),
		SessionTTL:    getEnvAsDuration("SESSION_TTL", 24*time.Hour),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		parsed, err := strconv.Atoi(value)
		if err == nil {
			return parsed
		}
	}

	return fallback
}

func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		parsed, err := time.ParseDuration(value)
		if err == nil {
			return parsed
		}
	}

	return fallback
}

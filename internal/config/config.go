package config

import "os"

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	RedisAddr  string
	RedisPass  string
	JWTSecret  string
	Port       string
}

func Load() *Config {
	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "bobobox"),
		DBPassword: getEnv("DB_PASSWORD", "bobobox"),
		DBName:     getEnv("DB_NAME", "bobobox"),
		RedisAddr:  getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPass:  getEnv("REDIS_PASSWORD", "bobobox"),
		JWTSecret:  getEnv("JWT_SECRET", "secret"),
		Port:       getEnv("PORT", ":8080"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return fallback
}

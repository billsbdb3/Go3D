package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	RedisAddr  string
	ServerPort string
	ServerHost string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "library3d"),
		DBPassword: getEnv("DB_PASSWORD", "dev123"),
		DBName:     getEnv("DB_NAME", "library3d"),
		RedisAddr:  getEnv("REDIS_ADDR", "localhost:6379"),
		ServerPort: getEnv("SERVER_PORT", "3000"),
		ServerHost: getEnv("SERVER_HOST", "192.168.3.26"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

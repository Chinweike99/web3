package config

import (
	"log"
	"os"
	"strconv"
	"github.com/joho/godotenv"
)

type Config struct {
	DBHost   string
	DBPort   string
	DBUser   string
	DBPassword string
	DBName   string
	DBSSLMode string

	JWTSecret string
	JwtAccessExpiration int
	JwtRefreshExpiration int

	ServerPort string
}


var AppConfig *Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	accessExp, _ := strconv.Atoi(getEnv("JWT_ACCESS_EXPIRATION", "15"))
	refreshExp, _ := strconv.Atoi(getEnv("JWT_REFRESH_EXPIRATION", "1440"))

	AppConfig = &Config{
		DBHost:   getEnv("DB_HOST", "localhost"),
		DBPort:   getEnv("DB_PORT", "5432"),
		DBUser:   getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:   getEnv("DB_NAME", "auth_db"),
		DBSSLMode: getEnv("DB_SSLMODE", "disable"),
		JWTSecret: getEnv("JWT_SECRET", "default-secret-key"),
		JwtAccessExpiration: accessExp,
		JwtRefreshExpiration: refreshExp,
		ServerPort: getEnv("SERVER_PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
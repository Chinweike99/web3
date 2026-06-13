package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	JWTSecret            string
	JwtAccessExpiration  int
	JwtRefreshExpiration int

	ServerPort string

	// Email settings
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	SMTPFrom    string
	AppURL	   string
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
		DBHost:               getEnv("DB_HOST", "localhost"),
		DBPort:               getEnv("DB_PORT", "5432"),
		DBUser:               getEnv("DB_USER", "postgres"),
		DBPassword:           getEnv("DB_PASSWORD", ""),
		DBName:               getEnv("DB_NAME", "auth_db"),
		DBSSLMode:            getEnv("DB_SSLMODE", "disable"),
		JWTSecret:            getEnv("JWT_SECRET", "default-secret-key"),
		JwtAccessExpiration:  accessExp,
		JwtRefreshExpiration: refreshExp,
		ServerPort:           getEnv("SERVER_PORT", "8080"),
		SMTPHost:             getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:             getEnvAsInt("SMTP_PORT", 587),
		SMTPUser:             getEnv("SMTP_USER", ""),
		SMTPPassword:         getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:             getEnv("SMTP_FROM", ""),
		AppURL:               getEnv("APP_URL", "http://localhost:8080"),
	}
}

func getEnvAsInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intVal, err := strconv.Atoi(value); err == nil {
            return intVal
        }
    }
    return defaultValue
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
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
	SMTPFrom     string
	AppURL       string

	// Rate Limiting Settings
	// Rate Limiting Settings
	RateLimitEnabled        bool
	RateLimitPerMinute      int
	RateLimitPerHour        int
	RateLimitPerDay         int
	RateLimitLoginPerMin    int
	RateLimitRegisterPerMin int
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
		DBHost:                  getEnv("DB_HOST", "localhost"),
		DBPort:                  getEnv("DB_PORT", "5432"),
		DBUser:                  getEnv("DB_USER", "postgres"),
		DBPassword:              getEnv("DB_PASSWORD", ""),
		DBName:                  getEnv("DB_NAME", "auth_db"),
		DBSSLMode:               getEnv("DB_SSLMODE", "disable"),
		JWTSecret:               getEnv("JWT_SECRET", "default-secret-key"),
		JwtAccessExpiration:     accessExp,
		JwtRefreshExpiration:    refreshExp,
		ServerPort:              getEnv("SERVER_PORT", "8080"),
		SMTPHost:                getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:                getEnvAsInt("SMTP_PORT", 587),
		SMTPUser:                getEnv("SMTP_USER", ""),
		SMTPPassword:            getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:                getEnv("SMTP_FROM", ""),
		AppURL:                  getEnv("APP_URL", "http://localhost:8080"),
		RateLimitEnabled:        getEnvAsBool("RATE_LIMIT_ENABLED", true),
		RateLimitPerMinute:      getEnvAsInt("RATE_LIMIT_PER_MINUTE", 60),
		RateLimitPerHour:        getEnvAsInt("RATE_LIMIT_PER_HOUR", 1000),
		RateLimitPerDay:         getEnvAsInt("RATE_LIMIT_PER_DAY", 5000),
		RateLimitLoginPerMin:    getEnvAsInt("RATE_LIMIT_LOGIN_PER_MIN", 5),
		RateLimitRegisterPerMin: getEnvAsInt("RATE_LIMIT_REGISTER_PER_MIN", 3),
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

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

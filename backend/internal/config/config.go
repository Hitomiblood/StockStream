package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// API Externa
	ExternalAPIURL   string
	ExternalAPIToken string

	// Database
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSchema   string
	DBSSLMode  string

	// Server
	APIPort string
	APIHost string

	// Configuración
	FetchInterval int
	LogLevel      string
}

func Load() *Config {
	// Cargar .env desde la raíz del proyecto
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using environment variables")
	} else {
		log.Println("✅ Loaded .env from current directory")
	}

	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "26257"))
	fetchInterval, _ := strconv.Atoi(getEnv("FETCH_INTERVAL", "3600"))

	return &Config{
		ExternalAPIURL:   getEnv("EXTERNAL_API_URL", ""),
		ExternalAPIToken: getEnv("EXTERNAL_API_TOKEN", ""),
		DBHost:           getEnv("DB_HOST", "localhost"),
		DBPort:           dbPort,
		DBUser:           getEnv("DB_USER", "root"),
		DBPassword:       getEnv("DB_PASSWORD", ""),
		DBName:           getEnv("DB_NAME", "stockdb"),
		DBSchema:         getEnv("DB_SCHEMA", "public"),
		DBSSLMode:        getEnv("DB_SSLMODE", "disable"),
		APIPort:          getEnv("API_PORT", "8080"),
		APIHost:          getEnv("API_HOST", "localhost"),
		FetchInterval:    fetchInterval,
		LogLevel:         getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

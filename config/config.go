package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	AppPort    string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("failed to load .env file: %v", err)
	}

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnvAsInt("DB_PORT", 5432),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "app"),
		AppPort:    getEnv("APP_PORT", "8080"),
	}
}

func LoadTest() *Config {
	if err := godotenv.Overload("../../.env.test"); err != nil {
		log.Printf("failed to load test .env file: %v", err)
	}

	return &Config{
		DBHost:     getEnv("TEST_DB_HOST", "test-postgres"),
		DBPort:     getEnvAsInt("TEST_DB_PORT", 5432),
		DBUser:     getEnv("TEST_DB_USER", "test_user"),
		DBPassword: getEnv("TEST_DB_PASSWORD", "test_password"),
		DBName:     getEnv("TEST_DB_NAME", "pr_service_test_db"),
		AppPort:    getEnv("TEST_APP_PORT", "8081"),
	}
}

func (c *Config) GetDSN() string {
	return "host=" + c.DBHost +
		" port=" + strconv.Itoa(c.DBPort) +
		" user=" + c.DBUser +
		" password=" + c.DBPassword +
		" dbname=" + c.DBName +
		" sslmode=disable"
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

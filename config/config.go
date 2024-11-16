package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server     ServerConfig
	JWT        JWTConfig
	API        APIConfig
	Env        string
	TestAccess TestUserLoginConfig
}

type ServerConfig struct {
	Host string
	Port string
}

type JWTConfig struct {
	Secret        string
	ExpirationDur time.Duration
}

type APIConfig struct {
	BaseURL string
}

type TestUserLoginConfig struct {
	TestUsername string
	TestPassword string
}

var AppConfig Config

func Init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	AppConfig = Config{
		Server: ServerConfig{
			Host: getEnv("HOST", "localhost"),
			Port: getEnv("PORT", "8080"),
		},

		JWT: JWTConfig{
			Secret:        getEnvRequired("JWT_SECRET"),
			ExpirationDur: time.Duration(getEnvAsInt("JWT_EXPIRATION_HOURS", 120)) * time.Hour,
		},
		API: APIConfig{
			BaseURL: getEnv("API_BASE_URL", "http://0.0.0.0"),
		},
		Env: getEnv("ENVIRONMENT", "development"),

		TestAccess: TestUserLoginConfig{
			TestUsername: getEnv("TEST_USERNAME", ""),
			TestPassword: getEnv("TEST_PASSWORD", ""),
		},
	}
	validateConfig()
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvRequired(key string) string {
	value, exists := os.LookupEnv(key)

	if !exists {
		log.Fatalf("Required environment variable %s is not set", key)
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	strValue, exists := os.LookupEnv(key)

	if !exists {
		return defaultValue
	}

	value, err := strconv.Atoi(strValue)

	if err != nil {
		log.Printf("Warning: Invalid integer value for %s, using default: %d", key, defaultValue)
		return defaultValue
	}

	return value
}

func validateConfig() {
	if AppConfig.JWT.Secret == "" {
		log.Fatalf("JWT secret cant be empty")
	}
}

func (c *ServerConfig) GetServerAddress() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

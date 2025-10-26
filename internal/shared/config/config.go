package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Server   ServerConfig
	JWT      JWTConfig
	CORS     CORSConfig
}

// AppConfig holds application-level configuration
type AppConfig struct {
	Environment string // development, staging, production
	Port        int
	LogLevel    string // debug, info, warn, error
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Host               string
	Port               int
	Name               string
	User               string
	Password           string
	SSLMode            string
	MaxConnections     int
	MaxIdleConnections int
	ConnectionLifetime time.Duration
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

// JWTConfig holds JWT authentication configuration
type JWTConfig struct {
	Secret string
	Expiry time.Duration
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists (for local development)
	// In production, environment variables should be set by the system
	_ = godotenv.Load() // Ignore error if .env doesn't exist

	cfg := &Config{
		App: AppConfig{
			Environment: getEnv("APP_ENV", "development"),
			Port:        getEnvAsInt("APP_PORT", 8080),
			LogLevel:    getEnv("APP_LOG_LEVEL", "info"),
		},
		Database: DatabaseConfig{
			Host:               getEnv("DB_HOST", "localhost"),
			Port:               getEnvAsInt("DB_PORT", 5432),
			Name:               getEnv("DB_NAME", "logistics_db"),
			User:               getEnv("DB_USER", "postgres"),
			Password:           getEnv("DB_PASSWORD", ""),
			SSLMode:            getEnv("DB_SSL_MODE", "disable"),
			MaxConnections:     getEnvAsInt("DB_MAX_CONNECTIONS", 25),
			MaxIdleConnections: getEnvAsInt("DB_MAX_IDLE_CONNECTIONS", 5),
			ConnectionLifetime: getEnvAsDuration("DB_CONNECTION_LIFETIME", 5*time.Minute),
		},
		Server: ServerConfig{
			ReadTimeout:     getEnvAsDuration("SERVER_READ_TIMEOUT", 10*time.Second),
			WriteTimeout:    getEnvAsDuration("SERVER_WRITE_TIMEOUT", 10*time.Second),
			IdleTimeout:     getEnvAsDuration("SERVER_IDLE_TIMEOUT", 120*time.Second),
			ShutdownTimeout: getEnvAsDuration("SERVER_SHUTDOWN_TIMEOUT", 30*time.Second),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", ""),
			Expiry: getEnvAsDuration("JWT_EXPIRY", 24*time.Hour),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnvAsSlice("CORS_ALLOWED_ORIGINS", []string{"*"}),
			AllowedMethods: getEnvAsSlice("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			AllowedHeaders: getEnvAsSlice("CORS_ALLOWED_HEADERS", []string{"Content-Type", "Authorization"}),
		},
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	validator := NewValidator()

	// Validate App config
	validator.Required("APP_ENV", c.App.Environment)
	validator.OneOf("APP_ENV", c.App.Environment, []string{"development", "staging", "production"})
	validator.Range("APP_PORT", c.App.Port, 1, 65535)
	validator.OneOf("APP_LOG_LEVEL", c.App.LogLevel, []string{"debug", "info", "warn", "error"})

	// Validate Database config
	validator.Required("DB_HOST", c.Database.Host)
	validator.Range("DB_PORT", c.Database.Port, 1, 65535)
	validator.Required("DB_NAME", c.Database.Name)
	validator.Required("DB_USER", c.Database.User)
	validator.Required("DB_PASSWORD", c.Database.Password)
	validator.OneOf("DB_SSL_MODE", c.Database.SSLMode, []string{"disable", "require", "verify-full"})
	validator.Min("DB_MAX_CONNECTIONS", c.Database.MaxConnections, 1)
	validator.Min("DB_MAX_IDLE_CONNECTIONS", c.Database.MaxIdleConnections, 1)

	// Validate JWT config (if production)
	if c.App.Environment == "production" {
		validator.Required("JWT_SECRET", c.JWT.Secret)
		validator.MinLength("JWT_SECRET", c.JWT.Secret, 32)
	}

	return validator.Error()
}

// GetDatabaseURL returns the PostgreSQL connection string
func (c *DatabaseConfig) GetDatabaseURL() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.Name,
		c.SSLMode,
	)
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}

// Helper functions to read environment variables with defaults

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	// Split by comma
	var result []string
	for _, v := range strings.Split(valueStr, ",") {
		result = append(result, strings.TrimSpace(v))
	}
	return result
}

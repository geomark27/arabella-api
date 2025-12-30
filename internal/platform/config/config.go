package config

import (
	"fmt"
	"os"
	"strings"
)

// Config contiene toda la configuración de la aplicación
type Config struct {
	Port        string
	Environment string
	LogLevel    string

	// CORS
	CorsAllowedOrigins []string

	// Base de datos
	DatabaseURL string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBSSLMode   string

	// JWT (para futuras implementaciones)
	JWTSecret        string
	JWTRefreshSecret string
}

// Load carga la configuración desde variables de entorno
func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),

		// CORS
		CorsAllowedOrigins: parseCorsOrigins(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:8080")),

		// Base de datos
		DatabaseURL: getEnv("DATABASE_URL", ""),
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5433"),
		DBUser:      getEnv("DB_USER", "postgres"),
		DBPassword:  getEnv("DB_PASSWORD", ""),
		DBName:      getEnv("DB_NAME", ""),
		DBSSLMode:   getEnv("DB_SSLMODE", "disable"),

		// JWT
		JWTSecret:        getEnv("JWT_SECRET", "your-default-secret-change-in-production"),
		JWTRefreshSecret: getEnv("JWT_REFRESH_SECRET", "your-refresh-secret-change-in-production"),
	}
}

// getEnv obtiene una variable de entorno o retorna un valor por defecto
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// parseCorsOrigins parsea la lista de orígenes CORS desde una cadena separada por comas
func parseCorsOrigins(origins string) []string {
	if origins == "" {
		return []string{}
	}

	result := make([]string, 0)
	for _, origin := range strings.Split(origins, ",") {
		trimmed := strings.TrimSpace(origin)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}

// IsDevelopment retorna true si el entorno es desarrollo
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction retorna true si el entorno es producción
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// GetDBConnectionString construye y retorna la cadena de conexión para PostgreSQL
func (c *Config) GetDBConnectionString() string {
	if c.DatabaseURL != "" {
		return c.DatabaseURL
	}

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode)
}

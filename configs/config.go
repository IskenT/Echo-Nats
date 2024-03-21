package configs

import (
	"os"
	"sync"
)

type Config struct {
	once       sync.Once
	HttpServer struct {
		Port        string
		MetricsPort string
	}

	Postgres struct {
		User     string
		Password string
		Host     string
		Port     string
		DBName   string
		SSLMode  string
		DSN      string
	}

	Redis struct {
		Host string
		Port string
	}
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	cfg.once.Do(func() {
		// Initialize HTTP server configuration
		cfg.HttpServer.Port = getEnv("HTTP_ADDR", "")
		cfg.HttpServer.MetricsPort = getEnv("METRICS_PORT", "")

		// Initialize Postgres configuration
		cfg.Postgres.User = getEnv("POSTGRES_USER", "")
		cfg.Postgres.Password = getEnv("POSTGRES_PASSWORD", "")
		cfg.Postgres.Host = getEnv("POSTGRES_HOST", "")
		cfg.Postgres.Port = getEnv("POSTGRES_PORT", "")
		cfg.Postgres.DBName = getEnv("POSTGRES_DB", "")
		cfg.Postgres.SSLMode = getEnv("POSTGRES_SSLMODE", "")
		cfg.Postgres.DSN = getEnv("POSTGRES_DSN", "")

		// Initialize Redis configuration
		cfg.Redis.Host = getEnv("REDIS_HOST", "")
		cfg.Redis.Port = getEnv("REDIS_PORT", "")
	})
	return cfg, nil
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

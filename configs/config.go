package configs

import (
	"fmt"
	"os"
	"strconv"
	"sync"
)

type Config struct {
	once       sync.Once
	HttpServer struct {
		Port        string
		MetricsPort string
	}

	Postgres struct {
		User         string
		Password     string
		Host         string
		Port         string
		DBName       string
		SSLMode      string
		MaxOpenConns int
		MaxIdleConns int
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
		cfg.Postgres.MaxOpenConns = getEnvAsInt("POSTGRES_MAX_OPEN_CONNS", 10)
		cfg.Postgres.MaxIdleConns = getEnvAsInt("POSTGRES_MAX_IDLE_CONNS", 10)

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

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

func (c *Config) GetPgDsn() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Postgres.Host,
		c.Postgres.Port,
		c.Postgres.User,
		c.Postgres.Password,
		c.Postgres.DBName,
		c.Postgres.SSLMode)
}

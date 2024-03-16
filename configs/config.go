package configs

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	HttpServer struct {
		Port        string `envconfig:"HTTP_ADDR"`
		MetricsPort string `envconfig:"METRICS_PORT"`
	}

	Postgres struct {
		User         string `envconfig:"POSTGRES_USER"`
		Password     string `envconfig:"POSTGRES_PASSWORD"`
		Host         string `envconfig:"POSTGRES_HOST"`
		Port         string `envconfig:"POSTGRES_PORT"`
		DBName       string `envconfig:"POSTGRES_DB"`
		SSLMode      string `envconfig:"POSTGRES_SSLMODE"`
		MaxOpenConns int    `envconfig:"POSTGRES_MAX_OPEN_CONNS"`
		MaxIdleConns int    `envconfig:"POSTGRES_MAX_IDLE_CONNS"`
	}

	Redis struct {
		Host string `envconfig:"REDIS_HOST"`
		Port string `envconfig:"REDIS_PORT"`
	}
}

var (
	once sync.Once
)

func LoadConfig() (*Config, error) {
	var config Config

	once.Do(func() {
		err := envconfig.Process("", &config)
		if err != nil {
			log.Fatal(err)
		}

		configBytes, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Configuration:", string(configBytes))
	})

	return &config, nil
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

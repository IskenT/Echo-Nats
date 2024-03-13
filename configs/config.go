package configs

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	HttpServer struct {
		Port        string `json:"port"`
		MetricsPort string `json:"metrics_port"`
	}

	Postgres struct {
		User         string `json:"user"`
		Password     string `json:"password"`
		Host         string `json:"host"`
		Port         string `json:"port"`
		DBName       string `json:"db_name"`
		SSLMode      string `json:"ssl_mode"`
		MaxOpenConns int    `json:"maxOpenConns"`
		MaxIdleConns int    `json:"maxIdleConns"`
	}

	Logger struct {
		Lvl string `json:"lvl"`
	}

	Redis struct {
		Host string `json:"host"`
		Port string `json:"port"`
	}
}

func LoadConfig(path string) (*Config, error) {
	var config *Config

	configFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		return nil, err
	}

	return config, nil
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

package main

import (
	"log"
	"rest_clickhouse/cmd/providers"
	"rest_clickhouse/configs"
)

const configPath = "configs/config.json"

func main() {
	cnf, err := configs.LoadConfig(configPath)
	if err != nil {
		log.Panic(err)
	}

	logger, err := providers.ProvideConsoleLogger(cnf)
	if err != nil {
		log.Panic(err)
	}

	db, closeDb, err := providers.ProvidePostgres(cnf, logger)
	if err != nil {
		log.Panic(err)
	}

	redisClient, err := providers.ProvideRedis(cnf)
	if err != nil {
		log.Panic(err)
	}

}

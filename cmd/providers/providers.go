package providers

import (
	"fmt"
	"os"
	"rest_clickhouse/configs"
	postgres "rest_clickhouse/pkg/db"
	"rest_clickhouse/pkg/logger"
	"rest_clickhouse/pkg/logger/zerolog"

	"github.com/go-redis/redis"
)

func ProvidePostgres(cnf *configs.Config, logger logger.Logger) (*postgres.DB, func(), error) {
	repo, err := postgres.NewDBConnection(cnf, logger)
	if err != nil {
		return nil, nil, err
	}

	closer := func() {
		repo.Close()
	}

	return repo, closer, nil
}

func ProvideConsoleLogger(cnf *configs.Config) (logger.Logger, error) {
	return zerolog.NewZeroLog(os.Stderr, cnf.Logger.Lvl)
}

func ProvideRedis(cnf *configs.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cnf.Redis.Host, cnf.Redis.Port),
		Password: "",
		DB:       0,
	})
	_, err := client.Ping().Result()
	return client, err
}

func ProvideNats(cnf *configs.Config) (queue.PubSub, error) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, err
	}

	natsClient := natsClient.NewNatsClient(nc)

	return natsClient, nil
}

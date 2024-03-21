package providers

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"rest_clickhouse/configs"
	"rest_clickhouse/internal/infrastructure/http"
	goods_service "rest_clickhouse/internal/infrastructure/http"
	"rest_clickhouse/internal/infrastructure/queue"
	nats_client "rest_clickhouse/internal/infrastructure/queue/nats"
	postgres "rest_clickhouse/pkg/db"
	"rest_clickhouse/pkg/logger"
	"rest_clickhouse/pkg/logger/zerolog"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"

	"github.com/go-redis/redis"
	"github.com/nats-io/nats.go"
)

func ProvideHTTPServer(config *configs.Config, goodsService goods_service.GoodsService, logger logger.Logger) http.HTTPServer {
	return http.NewEchoHTTPServer(config.HttpServer.Port, goodsService, logger)
}

func ProvidePostgres(ctx context.Context, cnf *configs.Config, logger logger.Logger) (*postgres.DB, func(), error) {
	repo, err := postgres.NewDBConnection(ctx, cnf, logger)
	if err != nil {
		return nil, nil, err
	}

	closer := func() {
		repo.Close()
	}

	return repo, closer, nil
}

func ProvideConsoleLogger(cnf *configs.Config) (logger.Logger, error) {
	return zerolog.NewZeroLog(os.Stderr)
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

func ProvideQueue(cnf *configs.Config) (queue.PubSub, error) {
	nc, err := nats.Connect(fmt.Sprintf("nats://%s:4222", os.Getenv("NATS_HOST")))
	if err != nil {
		return nil, err
	}

	nats_client := nats_client.NewNatsClient(nc)

	return nats_client, nil
}

func ProvideClickhouse(cnf *configs.Config) *sql.DB {
	conn := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{"127.0.0.1:9000"},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "",
			Password: "",
		},

		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		DialTimeout: time.Second * 30,
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		Debug:                false,
		BlockBufferSize:      10,
		MaxCompressionBuffer: 10240,
	})
	conn.SetMaxIdleConns(5)
	conn.SetMaxOpenConns(10)
	conn.SetConnMaxLifetime(time.Hour)

	return conn
}

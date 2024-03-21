package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"rest_clickhouse/cmd/providers"
	"rest_clickhouse/configs"
	goods_service "rest_clickhouse/internal/infrastructure/http"
	eventQueue "rest_clickhouse/internal/infrastructure/queue/nats"
	repository "rest_clickhouse/internal/infrastructure/repository"
	"rest_clickhouse/internal/infrastructure/usecase/interactors"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("failed to load .env file: %w", err)
	}

	cnf, err := configs.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger, err := providers.ProvideConsoleLogger(cnf)
	if err != nil {
		return fmt.Errorf("failed to provide console logger: %w", err)
	}

	db, closeDB, err := providers.ProvidePostgres(ctx, cnf, logger)
	if err != nil {
		return fmt.Errorf("failed to provide postgres: %w", err)
	}

	redisClient, err := providers.ProvideRedis(cnf)
	if err != nil {
		return fmt.Errorf("failed to provide redis client: %w", err)
	}

	queue, err := providers.ProvideQueue(cnf)
	if err != nil {
		return fmt.Errorf("failed to provide queue: %w", err)
	}

	goodsRepository := repository.NewGoodsRepository(ctx, db, redisClient, logger)
	goodsInteractor := interactors.NewGoodsInteractor(goodsRepository, redisClient, queue, logger)
	goodService := goods_service.NewGoodsService(goodsInteractor, logger)

	clickHouseConn := providers.ProvideClickhouse(cnf)
	logRepo := repository.NewLogsRepository(clickHouseConn, logger)
	eventListener := eventQueue.NewEventListener(ctx, queue, logRepo, logger)
	go eventListener.ListenTopic()

	server := providers.ProvideHTTPServer(cnf, goodService, logger)

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
		fmt.Println("Terminating the app")
		fmt.Println("Shutdown workers")
		cancel()

		fmt.Println("Close DB")
		closeDB()

		fmt.Println("Stop Server")
		server.Stop(ctx)
	}()

	server.Start()

	return nil
}

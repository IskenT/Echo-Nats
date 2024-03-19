package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"rest_clickhouse/cmd/providers"
	"rest_clickhouse/configs"
	httpControllers "rest_clickhouse/internal/infrastructure/interfaces"
	eventQueue "rest_clickhouse/internal/infrastructure/queue/nats"
	repository "rest_clickhouse/internal/infrastructure/repository"
	"rest_clickhouse/internal/infrastructure/usecase/interactors"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(fmt.Errorf("failed to load .env file: %w", err))
	}

	cnf, err := configs.LoadConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to load config: %w", err))
	}

	logger, err := providers.ProvideConsoleLogger(cnf)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to provide console logger: %w", err))
	}

	db, closeDB, err := providers.ProvidePostgres(cnf, logger)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to provide postgres: %w", err))
	}

	redisClient, err := providers.ProvideRedis(cnf)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to provide redis client: %w", err))
	}

	queue, err := providers.ProvideQueue(cnf)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to provide queue: %w", err))
	}

	ctx, cancel := context.WithCancel(context.Background())
	goodsRepository := repository.NewGoodsRepository(db, redisClient, logger)
	goodsInteractor := interactors.NewGoodsInteractor(goodsRepository, redisClient, queue, logger)
	goodController := httpControllers.NewGoodsController(goodsInteractor, logger)

	clickHouseConn := providers.ProvideClickhouse(cnf)
	logRepo := repository.NewLogsRepository(clickHouseConn, logger)
	eventListener := eventQueue.NewEventListener(ctx, queue, logRepo, logger)
	go eventListener.ListenTopic()

	server := providers.ProvideHTTPServer(cnf, goodController, logger)

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
}

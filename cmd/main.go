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
	queue2 "rest_clickhouse/internal/infrastructure/queue/nats"
	repository2 "rest_clickhouse/internal/infrastructure/repository"
	interactors "rest_clickhouse/internal/infrastructure/usecase/interractors"
	"syscall"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Panic("No .env file found")
	}
}

func main() {
	cnf, err := configs.LoadConfig()
	if err != nil {
		log.Panic(err)
	}

	logger, err := providers.ProvideConsoleLogger(cnf)
	if err != nil {
		log.Panic(err)
	}

	db, closeDB, err := providers.ProvidePostgres(cnf, logger)
	if err != nil {
		log.Panic(err)
	}

	redisClient, err := providers.ProvideRedis(cnf)
	if err != nil {
		log.Panic(err)
	}

	queue, err := providers.ProvideQueue(cnf)
	if err != nil {
		log.Panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	goodsRepository := repository2.NewGoodsRepository(db, redisClient, logger)
	goodsInteractor := interactors.NewGoodsInteractor(goodsRepository, redisClient, queue, logger)
	goodController := httpControllers.NewGoodsController(goodsInteractor, logger)

	clickHouseConn := providers.ProvideClickhouse(cnf)
	logRepo := repository2.NewLogsRepository(clickHouseConn)
	eventListener := queue2.NewEventListener(ctx, queue, logRepo, logger)
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

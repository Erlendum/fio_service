package main

import (
	"context"
	"errors"
	"fio_service/internal/config"
	"fio_service/internal/server"
	cache "fio_service/pkg/cache/redis"
	"fio_service/pkg/database"
	"fio_service/pkg/kafka"
	"fio_service/pkg/logger"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	config   *config.Config
	logger   *logger.Logger
	server   *server.Server
	consumer *kafka.Consumer
	producer *kafka.Producer
}

func (a *App) Init() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}
	a.config = cfg

	lg := logger.New(a.config.Logger.Path, a.config.Logger.Level)
	if lg == nil {
		log.Fatal("can't create logger")
	}
	a.logger = lg

	db, err := database.OpenDB(cfg.Database.Driver, cfg.Database.DSN)
	if err != nil {
		a.logger.Fatal(err)
	}

	memCache, err := cache.NewRedisCache(context.Background(), cfg.Redis)
	if err != nil {
		a.logger.Fatalf("error mem cache init: %v", err)
	}

	a.producer, err = kafka.NewProducer(cfg.Kafka.Brokers, *a.logger)
	if err != nil {
		a.logger.Fatalf("error creating Kafka producer: %v", err)
	}

	a.consumer, err = kafka.NewConsumer(a.config.Kafka.Brokers, *a.logger)
	if err != nil {
		a.logger.Fatalf("error creating Kafka consumer: %v", err)
	}

	fmt.Println(db)
	fmt.Println(memCache)

	a.server = server.NewServer(cfg, nil)

}

func main() {
	var a App
	a.Init()

	go func() {
		if err := a.server.Run(); !errors.Is(err, http.ErrServerClosed) {
			a.logger.Fatalf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	a.logger.Println("server started ", a.config.Server.Port)

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := a.server.Stop(ctx); err != nil {
		a.logger.Fatalf("failed to stop server %v", err)
	}
	a.logger.Println("server stopped ")
	if err := a.logger.Close(); err != nil {
		a.logger.Fatalf("failed to close logger %v", err)
	}
}

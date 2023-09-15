package main

import (
	"context"
	"database/sql"
	"errors"
	"fio_service/internal/config"
	myHttp "fio_service/internal/delivery/http"
	"fio_service/internal/repository"
	"fio_service/internal/repository/postgres_repository"
	"fio_service/internal/server"
	"fio_service/internal/service"
	"fio_service/internal/service/serviceImpl"
	"fio_service/pkg/cache"
	redisCache "fio_service/pkg/cache/redis"
	"fio_service/pkg/database"
	"fio_service/pkg/kafka"
	"fio_service/pkg/logger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	config       *config.Config
	logger       *logger.Logger
	server       *server.Server
	repositories *appRepositoryFields
	services     *service.Services
	hander       *myHttp.Handler
}

type appRepositoryFields struct {
	personRepository repository.PersonRepository
}

func (a *App) initServices(r *appRepositoryFields, c *cache.Cache, producer *kafka.Producer, consumer *kafka.Consumer) *service.Services {
	f := &service.Services{
		Person: serviceImpl.NewPersonServiceImplementation(r.personRepository, a.logger, *c, a.config.Redis.Ttl),
		Kafka:  serviceImpl.NewKafkaSerivce(producer, consumer, r.personRepository),
	}

	return f
}

func (a *App) initPostgresRepositories(db *sql.DB) *appRepositoryFields {
	if db == nil {
		return nil
	}
	f := &appRepositoryFields{
		personRepository: postgres_repository.CreatePersonPostgresRepository(db),
	}

	return f
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

	memCache, err := redisCache.NewRedisCache(context.Background(), cfg.Redis)
	if err != nil {
		a.logger.Fatalf("error mem cache init: %v", err)
	}

	producer, err := kafka.NewProducer(cfg.Kafka.Brokers, *a.logger)
	if err != nil {
		a.logger.Fatalf("error creating Kafka producer: %v", err)
	}

	consumer, err := kafka.NewConsumer(a.config.Kafka.Brokers, *a.logger)
	if err != nil {
		a.logger.Fatalf("error creating Kafka consumer: %v", err)
	}

	a.repositories = a.initPostgresRepositories(db)
	a.services = a.initServices(a.repositories, &memCache, producer, consumer)

	a.hander = myHttp.NewHandler(a.services, a.logger)

	a.server = server.NewServer(cfg, a.hander.Init())

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

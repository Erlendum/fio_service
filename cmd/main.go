package main

import (
	"context"
	"errors"
	"fio_service/internal/config"
	"fio_service/internal/server"
	cache "fio_service/pkg/cache/redis"
	"fio_service/pkg/database"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}
	db, err := database.OpenDB(cfg.Database.Driver, cfg.Database.DSN)
	if err != nil {
		log.Fatal(err)
	}

	memCache, err := cache.NewRedisCache(context.Background(), cfg.Redis)
	if err != nil {
		log.Fatalf("error mem cache init: %v", err)
	}
	fmt.Println(db)
	fmt.Println(memCache)

	srv := server.NewServer(cfg, nil)
	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	log.Println("server started ", cfg.Server.Port)

	quit := make(chan os.Signal, 1)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		log.Fatalf("failed to stop server %v", err)
	}
}

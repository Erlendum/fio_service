package main

import (
	"context"
	"fio_service/internal/config"
	cache "fio_service/pkg/cache/redis"
	"fio_service/pkg/database"
	"fmt"
	"log"
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
}

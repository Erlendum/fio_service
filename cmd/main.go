package main

import (
	"fmt"
	"log"
	"test/internal/config"
	"test/pkg/database"
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
	fmt.Println(db)
}

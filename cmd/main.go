package main

import (
	"Learning/internal/config"
	"Learning/internal/db"
	"Learning/internal/di"
	"context"
	"fmt"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	ctx := context.Background()
	pool, err := db.New(ctx, cfg.PgDSN)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer pool.Close()

	r := di.InitRouter(pool, cfg)

	addr := ":" + cfg.Port
	fmt.Println("Listening on", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}

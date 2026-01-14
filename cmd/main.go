package main

import (
	"Learning/internal/config"
	"Learning/internal/db"
	"Learning/internal/di"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	cfg, err := config.Load()
	if err != nil {
		slog.Error("error in config load", "error", err)
	}

	ctx := context.Background()
	pool, err := db.New(ctx, cfg.PgDSN)
	if err != nil {
		slog.Error("error in db load", "error", err)
	}
	defer pool.Close()

	r := di.InitRouter(pool, cfg)

	addr := ":" + cfg.Port
	fmt.Println("Listening on", addr)
	err = http.ListenAndServe(addr, r)
	if err != nil {
		slog.Error("error on app startup", "error", err)
	}
}

package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"skin-pro-app/internal/config"
	"skin-pro-app/internal/db"
	"skin-pro-app/internal/di"
	"syscall"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start HTTP server",
	Run:   runServe,
}

func runServe(cmd *cobra.Command, args []string) {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("config load failed", "error", err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	pool, err := db.New(ctx, cfg.PgDSN)
	if err != nil {
		slog.Error("db connection failed", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	r := di.InitRouter(pool, cfg)

	addr := ":" + cfg.Port
	fmt.Println("Listening on", addr)

	go func() {
		if err := http.ListenAndServe(addr, r); err != nil {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	slog.Info("shutting down")
}

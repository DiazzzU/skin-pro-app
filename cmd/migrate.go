package main

import (
	"log/slog"
	"os"
	"skin-pro-app/internal/config"
	"skin-pro-app/internal/db"

	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Run:   runMigrate,
}

func runMigrate(cmd *cobra.Command, args []string) {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("config load failed", "error", err)
		os.Exit(1)
	}

	if err := db.RunMigrations(cfg.PgDSN); err != nil {
		slog.Error("migrations failed", "error", err)
		os.Exit(1)
	}

	slog.Info("migrations completed successfully")
}

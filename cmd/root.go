package main

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "Skin Pro Backend",
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(serveCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		slog.Error("failed to execute command", "error", err)
		os.Exit(1)
	}
}

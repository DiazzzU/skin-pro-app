package db

import (
	"errors"
	"fmt"
	"log/slog"
	"skin-pro-app/migrations"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

func RunMigrations(dsn string) error {
	source, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return fmt.Errorf("failed to create migration source: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", source, dsn)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer m.Close()

	version, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return fmt.Errorf("failed to get version: %w", err)
	}

	slog.Info("migration state", "version", version, "dirty", dirty)

	if dirty {
		return fmt.Errorf("dirty state at version %d, manual fix required", version)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			slog.Info("migrations: no changes")
			return nil
		}
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	newVersion, _, _ := m.Version()
	slog.Info("migrations applied", "version", newVersion)

	return nil
}

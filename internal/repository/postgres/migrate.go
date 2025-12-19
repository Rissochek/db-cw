package postgres

import (
	"fmt"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
)

func (pg *Postgres) RunMigrations(migrationsPath string) error {
	driver, err := postgres.WithInstance(pg.conn.DB, &postgres.Config{})
	if err != nil {
		zap.S().Panicf("failed to create postgres driver: %v", err)
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	// Получаем абсолютный путь к миграциям
	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		zap.S().Panicf("failed to get absolute path for migrations: %v", err)
		return fmt.Errorf("failed to get absolute path for migrations: %w", err)
	}

	migrationsURL := fmt.Sprintf("file://%s", filepath.ToSlash(absPath))

	m, err := migrate.NewWithDatabaseInstance(
		migrationsURL,
		"postgres",
		driver,
	)
	if err != nil {
		zap.S().Panicf("failed to create migrate instance: %v", err)
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		zap.S().Panicf("failed to get migration version: %v", err)
		return fmt.Errorf("failed to get migration version: %w", err)
	}

	if dirty {
		zap.S().Warnf("Database is in dirty state (version: %d). Attempting to force version...", version)
		if err := m.Force(int(version)); err != nil {
			zap.S().Panicf("failed to force migration version: %v", err)
			return fmt.Errorf("failed to force migration version: %w", err)
		}
		zap.S().Info("Successfully forced migration version, continuing...")
	}

	zap.S().Info("Starting database migrations...")
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			zap.S().Info("No migrations to apply, database is up to date")
			return nil
		}
		zap.S().Panicf("failed to run migrations: %v", err)
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	finalVersion, finalDirty, err := m.Version()
	if err != nil {
		zap.S().Warnf("failed to get final migration version: %v", err)
	} else {
		if finalDirty {
			zap.S().Warnf("Migrations completed but database is in dirty state. Version: %d", finalVersion)
		} else {
			zap.S().Infof("Migrations completed successfully. Current version: %d", finalVersion)
		}
	}

	return nil
}

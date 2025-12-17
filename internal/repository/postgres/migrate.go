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
		zap.S().Errorf("failed to create postgres driver: %v", err)
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	// Получаем абсолютный путь к миграциям
	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		zap.S().Errorf("failed to get absolute path for migrations: %v", err)
		return fmt.Errorf("failed to get absolute path for migrations: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", absPath),
		"postgres",
		driver,
	)
	if err != nil {
		zap.S().Errorf("failed to create migrate instance: %v", err)
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	zap.S().Info("Starting database migrations...")
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			zap.S().Info("No migrations to apply, database is up to date")
			return nil
		}
		zap.S().Errorf("failed to run migrations: %v", err)
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	version, dirty, err := m.Version()
	if err != nil {
		zap.S().Warnf("failed to get migration version: %v", err)
	} else {
		zap.S().Infof("Migrations completed successfully. Current version: %d, dirty: %v", version, dirty)
	}

	return nil
}
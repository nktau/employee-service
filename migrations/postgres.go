package migrations

import (
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

func RunPostgresMigration(logger *zap.Logger, migrationPATH, DatabaseURL string) {

	migration, err := migrate.New(migrationPATH, DatabaseURL)
	if err != nil {
		logger.Fatal("cannot create new migrate instance", zap.Error(err))
	}

	if err = migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.Fatal("failed to run migrate up", zap.Error(err))
	}

	logger.Info("db migrated successfully")
}

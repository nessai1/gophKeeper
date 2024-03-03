package metastorage

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/nessai1/gophkeeper/internal/service/config"

	_ "github.com/jackc/pgx/v5/stdlib"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type PSQLMetaStorage struct {
	config config.PSQLMetaStorageConfig
	db     *sql.DB
}

func NewPSQLMetaStorage(cfg config.PSQLMetaStorageConfig) (*PSQLMetaStorage, error) {
	db, err := sql.Open(
		"pgx",
		fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.Host,
			cfg.User,
			cfg.Password,
			cfg.DBName,
		),
	)

	if err != nil {
		return nil, fmt.Errorf("cannot open psql connection: %w", err)
	}

	err = initPSQLMigrations(db)
	if err != nil && !errors.As(migrate.ErrNoChange, &err) {
		return nil, fmt.Errorf("cannot init psql migrations: %w", err)
	}

	return &PSQLMetaStorage{
		config: cfg,
		db:     db,
	}, nil
}

func initPSQLMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	migrations, err := migrate.NewWithDatabaseInstance("file:migrations/psql", "postgres", driver)
	if err != nil {
		return fmt.Errorf("error while create migrate DB instance: %s", err.Error())
	}

	if err = migrations.Up(); err != nil {
		return fmt.Errorf("cannot up migration: %w", err)
	}

	return nil
}

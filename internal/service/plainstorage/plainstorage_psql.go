package plainstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/nessai1/gophkeeper/internal/service/config"

	_ "github.com/jackc/pgx/v5/stdlib"

	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/google/uuid"
)

type PSQLPlainStorage struct {
	config config.PSQLPlainStorageConfig
	db     *sql.DB
}

func NewPSQLPlainStorage(cfg config.PSQLPlainStorageConfig) (*PSQLPlainStorage, error) {
	db, err := sql.Open(
		"pgx",
		fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.Host,
			cfg.Port,
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

	return &PSQLPlainStorage{
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

func (s *PSQLPlainStorage) GetUserByLogin(ctx context.Context, login string) (*User, error) {
	var user User
	err := s.db.QueryRowContext(
		ctx,
		"SELECT uuid, login, password FROM users WHERE login = $1",
		login,
	).Scan(&user.UUID, &user.Login, &user.PasswordHash)

	if err != nil && errors.Is(sql.ErrNoRows, err) {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, fmt.Errorf("cannot get user by login: %w", err)
	}

	return &user, nil
}

func (s *PSQLPlainStorage) GetUserByUUID(ctx context.Context, uuid string) (*User, error) {
	var user User
	err := s.db.QueryRowContext(
		ctx,
		"SELECT uuid, login, password FROM users WHERE uuid = $1",
		uuid,
	).Scan(&user.UUID, &user.Login, &user.PasswordHash)

	if err != nil && errors.Is(sql.ErrNoRows, err) {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, fmt.Errorf("cannot get user by uuid: %w", err)
	}

	return &user, nil
}

func (s *PSQLPlainStorage) CreateUser(ctx context.Context, login string, password string) (*User, error) {
	userUUID := uuid.New().String()
	_, err := s.db.ExecContext(ctx, "INSERT INTO users (uuid, login, password) VALUES ($1, $2, $3)", userUUID, login, password)

	if err != nil {
		return nil, fmt.Errorf("cannot create user: %w", err)
	}

	return &User{
		UUID:         userUUID,
		Login:        login,
		PasswordHash: password,
	}, nil
}

func (s *PSQLPlainStorage) GetUserSecretsByType(ctx context.Context, userUUID string, secretType SecretType) ([]SecretMetadata, error) {
	//TODO implement me
	panic("implement me")
}

func (s *PSQLPlainStorage) GetPlainSecretByUUID(ctx context.Context, secretUUID string) (*PlainSecret, error) {
	//TODO implement me
	panic("implement me")
}

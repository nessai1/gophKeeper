package plainstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/nessai1/gophkeeper/internal/service/config"
	"go.uber.org/zap"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/google/uuid"
)

type PSQLPlainStorage struct {
	config config.PSQLPlainStorageConfig
	db     *sql.DB
	logger *zap.Logger
}

func NewPSQLPlainStorage(cfg config.PSQLPlainStorageConfig, l *zap.Logger) (*PSQLPlainStorage, error) {
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

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("cannot ping psql connection: %w", err)
	}

	err = initPSQLMigrations(db)
	if err != nil && !errors.As(migrate.ErrNoChange, &err) {
		return nil, fmt.Errorf("cannot init psql migrations: %w", err)
	}

	return &PSQLPlainStorage{
		config: cfg,
		db:     db,
		logger: l,
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

func (s *PSQLPlainStorage) GetUserSecretsMetadataByType(ctx context.Context, userUUID string, secretType SecretType) ([]SecretMetadata, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT uuid, name, created, updated FROM secret_metadata WHERE owner_uuid = $1 AND type = $2", userUUID, secretType)
	if err != nil {
		return nil, fmt.Errorf("cannot get user secrets by type: %w", err)
	}

	secrets := make([]SecretMetadata, 0)
	defer func() {
		err := rows.Close()
		if err != nil {
			s.logger.Error("Cannot close rows in secrets list query", zap.Error(err))
		}
	}()

	for rows.Next() {
		secret := SecretMetadata{
			UserUUID: userUUID,
			Type:     secretType,
		}

		err := rows.Scan(&secret.UUID, &secret.Name, &secret.Created, &secret.Updated)
		if err != nil {
			s.logger.Error("Error while fetch row from secrets list query", zap.Error(err))

			continue
		}

		secrets = append(secrets, secret)
	}

	return secrets, nil
}

func (s *PSQLPlainStorage) GetPlainSecretByUUID(ctx context.Context, secretUUID string) (*PlainSecret, error) {
	//TODO implement me
	panic("implement me")
}

func (s *PSQLPlainStorage) AddSecretMetadata(ctx context.Context, userUUID string, name string, dataType SecretType) (*SecretMetadata, error) {
	dataUUID := uuid.New().String()
	_, err := s.db.ExecContext(ctx, "INSERT INTO secret_metadata (uuid, owner_uuid, name, type) VALUES ($1, $2, $3, $4)", dataUUID, userUUID, name, dataType)

	if err != nil {
		return nil, fmt.Errorf("cannot create secret metadata: %w", err)
	}

	return &SecretMetadata{
		UUID:     dataUUID,
		UserUUID: userUUID,
		Name:     name,
		Type:     dataType,
	}, nil
}

func (s *PSQLPlainStorage) RemoveSecretByUUID(ctx context.Context, secretUUID string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM secret_metadata WHERE uuid = $1", secretUUID)

	if err != nil {
		return fmt.Errorf("cannot remove secret metadata: %w", err)
	}

	return nil
}

func (s *PSQLPlainStorage) GetUserSecretByName(ctx context.Context, userUUID string, secretName string, secretType SecretType) (*PlainSecret, error) {
	var (
		secretUUID       string
		created, updated time.Time
	)

	err := s.db.
		QueryRowContext(ctx, "SELECT uuid, created, updated FROM secret_metadata WHERE owner_uuid = $1 AND name = $2 AND type = $3", userUUID, secretName, secretType).
		Scan(&secretUUID, &created, &updated)

	if errors.Is(sql.ErrNoRows, err) {
		return nil, ErrSecretNotFound
	} else if err != nil {
		return nil, fmt.Errorf("error while get secret metadata: %w", err)
	}

	var (
		dbSecretContent string
		content         []byte
	)
	err = s.db.QueryRowContext(ctx, "SELECT data FROM plain_secret WHERE uuid = $1", secretUUID).Scan(&dbSecretContent)
	if err != nil && !errors.Is(sql.ErrNoRows, err) {
		return nil, fmt.Errorf("error while get secret content: %w", err)
	}

	if dbSecretContent != "" {
		content = []byte(dbSecretContent)
	}

	return &PlainSecret{
		Metadata: SecretMetadata{
			UUID:     secretUUID,
			UserUUID: userUUID,
			Name:     secretName,
			Type:     secretType,
			Created:  created,
			Updated:  updated,
		},
		Data: content,
	}, nil
}

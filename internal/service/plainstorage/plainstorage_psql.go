package plainstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/nessai1/gophkeeper/internal/service/config"
	"github.com/nessai1/gophkeeper/pkg/postgrescodes"
	"go.uber.org/zap"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/google/uuid"
)

type PSQLPlainStorage struct {
	config config.PSQLPlainStorageConfig
	db     *sqlx.DB
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
		db:     sqlx.NewDb(db, "pgx"),
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

func (s *PSQLPlainStorage) InTransaction(ctx context.Context, transaction func() error) error {
	tx, txErr := s.db.BeginTx(ctx, nil)
	if txErr != nil {
		return fmt.Errorf("cannot start intransaction: %w", txErr)
	}

	err := transaction()
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			s.logger.Error("Cannot rollback intransaction", zap.Error(txErr))

			return errors.Join(err, fmt.Errorf("cannot rollback intransaction: %w", txErr))
		}

		return err
	}

	txErr = tx.Commit()
	if txErr != nil {
		s.logger.Error("Cannot commit intransaction", zap.Error(txErr))

		return fmt.Errorf("cannot commit intransaction: %w", txErr)
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
		return nil, ErrEntityNotFound
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
		return nil, ErrEntityNotFound
	} else if err != nil {
		return nil, fmt.Errorf("cannot get user by uuid: %w", err)
	}

	return &user, nil
}

func (s *PSQLPlainStorage) CreateUser(ctx context.Context, login string, password string) (*User, error) {
	userUUID := uuid.New().String()
	_, err := s.db.ExecContext(ctx, "INSERT INTO users (uuid, login, password) VALUES ($1, $2, $3)", userUUID, login, password)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == postgrescodes.PostgresErrCodeUniqueViolation {
				return nil, ErrEntityAlreadyExists
			}
		}

		return nil, fmt.Errorf("cannot create user: %w", err)
	}

	return &User{
		UUID:         userUUID,
		Login:        login,
		PasswordHash: password,
	}, nil
}

func (s *PSQLPlainStorage) GetUserSecretsMetadataByType(ctx context.Context, userUUID string, secretType SecretType) ([]SecretMetadata, error) {
	query := `SELECT uuid, owner_uuid, name, type, created, updated FROM secret_metadata WHERE owner_uuid = ? AND type = ?`
	query, args, err := sqlx.In(query, userUUID, secretType)
	if err != nil {
		s.logger.Error("Error preparing list secrets query", zap.Error(err))

		return nil, fmt.Errorf("error preparing list secrets query: %w", err)
	}

	query = s.db.Rebind(query)

	var secrets []SecretMetadata
	err = s.db.SelectContext(ctx, &secrets, query, args...)
	if err != nil {
		s.logger.Error("Error while get list of secrets", zap.Error(err))

		return nil, fmt.Errorf("error while get list of secrets")
	}

	return secrets, nil
}

func (s *PSQLPlainStorage) AddSecretMetadata(ctx context.Context, userUUID string, name string, dataType SecretType) (*SecretMetadata, error) {
	dataUUID := uuid.New().String()
	_, err := s.db.ExecContext(ctx, "INSERT INTO secret_metadata (uuid, owner_uuid, name, type) VALUES ($1, $2, $3, $4)", dataUUID, userUUID, name, dataType)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == postgrescodes.PostgresErrCodeUniqueViolation {
				return nil, ErrEntityAlreadyExists
			}
		}

		return nil, fmt.Errorf("cannot create secret metadata: %w", err)
	}

	return &SecretMetadata{
		UUID:     dataUUID,
		UserUUID: userUUID,
		Name:     name,
		Type:     dataType,
	}, nil
}

func (s *PSQLPlainStorage) UpdateSecretMetadataUUID(ctx context.Context, userUUID string, oldUUID string, newUUID string, dataType SecretType) error {
	res, err := s.db.ExecContext(ctx, "UPDATE secret_metadata SET uuid = $1, updated = now() WHERE owner_uuid = $2 AND type = $3 AND uuid = $4", newUUID, userUUID, dataType, oldUUID)
	if err != nil {
		return fmt.Errorf("cannot make update query: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if rowsAffected == 0 {
		return ErrEntityNotFound
	}

	return nil
}

func (s *PSQLPlainStorage) RemoveSecretByUUID(ctx context.Context, secretUUID string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM secret_metadata WHERE uuid = $1", secretUUID)

	if err != nil {
		return fmt.Errorf("cannot remove secret metadata: %w", err)
	}

	return nil
}

func (s *PSQLPlainStorage) AddPlainSecret(ctx context.Context, userUUID string, name string, dataType SecretType, data []byte) (*PlainSecret, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot start create secret transaction: %w", err)
	}

	md, err := s.AddSecretMetadata(ctx, userUUID, name, dataType)
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			s.logger.Error("Cannot rollback create secret transaction", zap.Error(err))
		}

		if errors.Is(ErrEntityAlreadyExists, err) {
			return nil, ErrEntityAlreadyExists
		}

		return nil, fmt.Errorf("cannot create metadata of plain secret: %w", err)
	}

	_, err = s.db.ExecContext(ctx, "INSERT INTO plain_secret (uuid, data) VALUES ($1, $2)", md.UUID, data)
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			s.logger.Error("Cannot rollback create secret transaction", zap.Error(err))
		}

		return nil, fmt.Errorf("cannot insert plain secret data to table: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("cannot commit create secret transaction: %w", err)
	}

	return &PlainSecret{
		Metadata: *md,
		Data:     data,
	}, nil
}

func (s *PSQLPlainStorage) UpdatePlainSecretByName(ctx context.Context, ownerUUID string, name string, data []byte) error {
	var secretUUID string
	err := s.db.QueryRowContext(ctx, "SELECT uuid FROM secret_metadata WHERE owner_uuid = $1 AND name = $2", ownerUUID, name).Scan(&secretUUID)
	if err != nil && errors.Is(sql.ErrNoRows, err) {
		return ErrEntityNotFound
	} else if err != nil {
		return fmt.Errorf("error while get secret metadata: %w", err)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("cannot start update transaction: %w", err)
	}

	_, err = s.db.ExecContext(ctx, "UPDATE secret_metadata SET updated = now() WHERE owner_uuid = $1 AND name = $2", ownerUUID, name)
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			s.logger.Error("Cannot rollback update transaction", zap.Error(txErr))
		}

		return fmt.Errorf("error while update metadata: %w", err)
	}

	if data == nil {
		_, err = s.db.ExecContext(ctx, "DELETE FROM plain_secret WHERE uuid = $1", secretUUID)
	} else {
		_, err = s.db.ExecContext(ctx, "INSERT INTO plain_secret (uuid, data) VALUES ($1, $2) ON CONFLICT (uuid) DO UPDATE SET data = $2", secretUUID, data)
	}

	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			s.logger.Error("Cannot rollback update transaction", zap.Error(txErr))
		}

		return fmt.Errorf("cannot update secret data: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("cannot commit secret update: %w", err)
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
		return nil, ErrEntityNotFound
	} else if err != nil {
		return nil, fmt.Errorf("error while get secret metadata: %w", err)
	}

	var content []byte
	err = s.db.QueryRowContext(ctx, "SELECT data FROM plain_secret WHERE uuid = $1", secretUUID).Scan(&content)
	if err != nil && !errors.Is(sql.ErrNoRows, err) {
		return nil, fmt.Errorf("error while get secret content: %w", err)
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

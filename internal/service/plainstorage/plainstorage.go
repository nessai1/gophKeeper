package plainstorage

import (
	"context"
	"errors"
	"time"
)

const (
	SecretTypeCredentials = iota
	SecretTypeCard
	SecretTypeText
	SecretTypeMedia
)

type SecretType uint8

type PlainStorage interface {
	InTransaction(ctx context.Context, transaction func() error) error

	GetUserByLogin(ctx context.Context, login string) (*User, error)
	GetUserByUUID(ctx context.Context, uuid string) (*User, error)
	CreateUser(ctx context.Context, login string, password string) (*User, error)
	GetUserSecretsMetadataByType(ctx context.Context, userUUID string, secretType SecretType) ([]SecretMetadata, error)

	AddSecretMetadata(ctx context.Context, userUUID string, secretUUID, name string, dataType SecretType) (*SecretMetadata, error)
	AddPlainSecret(ctx context.Context, userUUID string, name string, dataType SecretType, data []byte) (*PlainSecret, error)

	UpdateSecretMetadataUUID(ctx context.Context, userUUID string, oldUUID string, newUUID string, dataType SecretType) error

	UpdatePlainSecretDataByName(ctx context.Context, ownerUUID string, name string, dataType SecretType, data []byte) error
	RemoveSecretByUUID(ctx context.Context, secretUUID string) error

	GetUserSecretByName(ctx context.Context, userUUID string, secretName string, secretType SecretType) (*PlainSecret, error)
}

var ErrEntityNotFound = errors.New("entity not found")
var ErrEntityAlreadyExists = errors.New("entity already exists")

type PlainSecret struct {
	Metadata SecretMetadata
	Data     []byte
}

type User struct {
	UUID         string
	Login        string
	PasswordHash string
}

type SecretMetadata struct {
	UUID     string `db:"uuid"`
	UserUUID string `db:"owner_uuid"`

	// Name for plains - title, for media - filename
	Name string `db:"name"`

	Type SecretType `db:"type"`

	Created time.Time `db:"created"`
	Updated time.Time `db:"updated"`
}

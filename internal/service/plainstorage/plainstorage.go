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
	GetUserByLogin(ctx context.Context, login string) (*User, error)
	GetUserByUUID(ctx context.Context, login string) (*User, error)
	CreateUser(ctx context.Context, login string, password string) (*User, error)
	GetUserSecretsMetadataByType(ctx context.Context, userUUID string, secretType SecretType) ([]SecretMetadata, error)

	AddSecretMetadata(ctx context.Context, userUUID string, name string, dataType SecretType) (*SecretMetadata, error)
	AddPlainSecret(ctx context.Context, userUUID string, name string, dataType SecretType, data []byte) (*PlainSecret, error)

	UpdatePlainSecretByName(ctx context.Context, ownerUUID string, name string, data []byte) error
	RemoveSecretByUUID(ctx context.Context, secretUUID string) error

	GetUserSecretByName(ctx context.Context, userUUID string, secretName string, secretType SecretType) (*PlainSecret, error)
}

var ErrSecretNotFound = errors.New("secret not found")
var ErrUserNotFound = errors.New("user not found")

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
	UUID     string
	UserUUID string

	// Name for plains - title, for media - filename
	Name string

	Type SecretType

	Created time.Time
	Updated time.Time
}

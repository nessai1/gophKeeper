package plainstorage

import (
	"context"
	"github.com/google/uuid"
	"slices"
	"time"
)

// In-memory storage, using for tests only

type MemoryStorage struct {
	Users      []User
	SecretList []PlainSecret
}

func (m *MemoryStorage) InTransaction(_ context.Context, transaction func() error) error {
	return transaction()
}

func (m *MemoryStorage) GetUserByLogin(_ context.Context, login string) (*User, error) {
	for _, v := range m.Users {
		if v.Login == login {
			return &v, nil
		}
	}

	return nil, ErrEntityNotFound
}

func (m *MemoryStorage) GetUserByUUID(_ context.Context, userUUID string) (*User, error) {
	for _, v := range m.Users {
		if v.UUID == userUUID {
			return &v, nil
		}
	}

	return nil, ErrEntityNotFound
}

func (m *MemoryStorage) CreateUser(_ context.Context, login string, password string) (*User, error) {
	for _, v := range m.Users {
		if v.Login == login {
			return nil, ErrEntityAlreadyExists
		}
	}

	newUser := User{
		UUID:         uuid.New().String(),
		Login:        login,
		PasswordHash: password,
	}

	m.Users = append(m.Users, newUser)

	return &newUser, nil
}

func (m *MemoryStorage) GetUserSecretsMetadataByType(_ context.Context, userUUID string, secretType SecretType) ([]SecretMetadata, error) {
	rs := make([]SecretMetadata, 0)
	for _, val := range m.SecretList {
		if val.Metadata.UserUUID == userUUID && val.Metadata.Type == secretType {
			rs = append(rs, val.Metadata)
		}
	}

	return rs, nil
}

func (m *MemoryStorage) AddSecretMetadata(_ context.Context, userUUID string, secretUUID, name string, dataType SecretType) (*SecretMetadata, error) {
	for _, v := range m.SecretList {
		if v.Metadata.UserUUID == userUUID && v.Metadata.Name == name && v.Metadata.Type == dataType {
			return nil, ErrEntityAlreadyExists
		}
	}

	secret := PlainSecret{
		Metadata: SecretMetadata{
			UUID:     secretUUID,
			UserUUID: userUUID,
			Name:     name,
			Type:     dataType,
			Created:  time.Now(),
			Updated:  time.Now(),
		},
		Data: nil,
	}

	m.SecretList = append(m.SecretList, secret)

	return &secret.Metadata, nil
}

func (m *MemoryStorage) AddPlainSecret(_ context.Context, userUUID string, name string, dataType SecretType, data []byte) (*PlainSecret, error) {
	for _, v := range m.SecretList {
		if v.Metadata.UserUUID == userUUID && v.Metadata.Name == name && v.Metadata.Type == dataType {
			return nil, ErrEntityAlreadyExists
		}
	}

	secretUUID := uuid.New().String()

	secret := PlainSecret{
		Metadata: SecretMetadata{
			UUID:     secretUUID,
			UserUUID: userUUID,
			Name:     name,
			Type:     dataType,
			Created:  time.Now(),
			Updated:  time.Now(),
		},
		Data: data,
	}

	m.SecretList = append(m.SecretList, secret)

	return &secret, nil
}

func (m *MemoryStorage) UpdateSecretMetadataUUID(_ context.Context, userUUID string, oldUUID string, newUUID string, dataType SecretType) error {
	var secret *SecretMetadata
	for i, v := range m.SecretList {
		if v.Metadata.UserUUID == userUUID && v.Metadata.UUID == oldUUID && v.Metadata.Type == dataType {
			secret = &m.SecretList[i].Metadata
		}

		if v.Metadata.UserUUID == userUUID && v.Metadata.UUID == newUUID && v.Metadata.Type == dataType {
			return ErrEntityAlreadyExists
		}
	}

	if secret == nil {
		return ErrEntityNotFound
	}

	secret.UUID = newUUID
	secret.Updated = time.Now()
	return nil
}

func (m *MemoryStorage) UpdatePlainSecretDataByName(_ context.Context, ownerUUID string, name string, secretType SecretType, data []byte) error {
	var secret *PlainSecret
	for i, v := range m.SecretList {
		if v.Metadata.UserUUID == ownerUUID && v.Metadata.Name == name && v.Metadata.Type == secretType {
			secret = &m.SecretList[i]
		}
	}

	if secret == nil {
		return ErrEntityNotFound
	}

	secret.Metadata.Updated = time.Now()
	secret.Data = data

	return nil
}

func (m *MemoryStorage) RemoveSecretByUUID(_ context.Context, secretUUID string) error {
	pos := -1
	for i, v := range m.SecretList {
		if v.Metadata.UUID == secretUUID {
			pos = i
		}
	}

	if pos == -1 {
		return ErrEntityNotFound
	}

	m.SecretList = slices.Delete(m.SecretList, pos, pos+1)

	return nil
}

func (m *MemoryStorage) GetUserSecretByName(_ context.Context, userUUID string, secretName string, secretType SecretType) (*PlainSecret, error) {
	for _, v := range m.SecretList {
		if v.Metadata.UserUUID == userUUID && v.Metadata.Name == secretName && v.Metadata.Type == secretType {
			return &v, nil
		}
	}

	return nil, ErrEntityNotFound
}

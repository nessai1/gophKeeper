package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	pb "github.com/nessai1/gophkeeper/api/proto"
	"github.com/nessai1/gophkeeper/internal/service/plainstorage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestServer_SecretList(t *testing.T) {
	server, _, plain, err := NewTestServer()
	require.NoError(t, err)

	testUser := plainstorage.User{
		UUID:         uuid.New().String(),
		Login:        "testUser",
		PasswordHash: "somesecrethash",
	}

	secondUser := plainstorage.User{
		UUID:         uuid.New().String(),
		Login:        "secondUser",
		PasswordHash: "somesecrethash",
	}
	plain.Users = append(plain.Users, testUser, secondUser)

	testUserTextSecret := plainstorage.PlainSecret{
		Metadata: plainstorage.SecretMetadata{
			UUID:     uuid.New().String(),
			UserUUID: testUser.UUID,
			Name:     "testUserTextSecret1",
			Type:     plainstorage.SecretTypeText,
			Created:  time.Now(),
			Updated:  time.Now(),
		},
		Data: nil,
	}

	testUserCardSecret := plainstorage.PlainSecret{
		Metadata: plainstorage.SecretMetadata{
			UUID:     uuid.New().String(),
			UserUUID: testUser.UUID,
			Name:     "testUserCardSecret1",
			Type:     plainstorage.SecretTypeCard,
			Created:  time.Now(),
			Updated:  time.Now(),
		},
		Data: nil,
	}

	secondUserTextSecret := plainstorage.PlainSecret{
		Metadata: plainstorage.SecretMetadata{
			UUID:     uuid.New().String(),
			UserUUID: secondUser.UUID,
			Name:     "secondUserTextSecret1",
			Type:     plainstorage.SecretTypeText,
			Created:  time.Now(),
			Updated:  time.Now(),
		},
		Data: nil,
	}

	plain.SecretList = append(plain.SecretList, testUserTextSecret, testUserCardSecret, secondUserTextSecret)

	ctx := context.WithValue(context.Background(), UserContextKey, &testUser)

	resp, err := server.SecretList(ctx, &pb.SecretListRequest{
		SecretType: plainstorage.SecretTypeText,
	})
	require.NoError(t, err)
	require.Equal(t, 1, len(resp.Secrets))
	secret, err := translateSecret(testUserTextSecret)
	require.NoError(t, err)
	assert.Equal(t, secret, resp.Secrets[0])

	resp, err = server.SecretList(ctx, &pb.SecretListRequest{
		SecretType: plainstorage.SecretTypeCredentials,
	})
	require.NoError(t, err)
	assert.Empty(t, resp.Secrets)

	ctx = context.WithValue(context.Background(), UserContextKey, &secondUser)
	resp, err = server.SecretList(ctx, &pb.SecretListRequest{
		SecretType: plainstorage.SecretTypeText,
	})
	require.NoError(t, err)
	require.Equal(t, 1, len(resp.Secrets))
	secret, err = translateSecret(secondUserTextSecret)
	require.NoError(t, err)
	assert.Equal(t, secret, resp.Secrets[0])
}

func TestServer_SecretSetGet(t *testing.T) {
	server, _, plain, err := NewTestServer()
	require.NoError(t, err)

	testUser := plainstorage.User{
		UUID:         uuid.New().String(),
		Login:        "testUser",
		PasswordHash: "somesecrethash",
	}

	secondUser := plainstorage.User{
		UUID:         uuid.New().String(),
		Login:        "secondUser",
		PasswordHash: "somesecrethash",
	}
	plain.Users = append(plain.Users, testUser, secondUser)

	tests := []struct {
		name       string
		setRequest pb.SecretSetRequest
		user       plainstorage.User
		alreadyHas bool
	}{
		{
			name: "Set text secret",
			setRequest: pb.SecretSetRequest{
				SecretType: pb.SecretType_TEXT,
				Name:       "mytextsecret",
				Content:    []byte("my text content"),
			},
			user: testUser,
		},
		{
			name: "Set text secret with same name",
			setRequest: pb.SecretSetRequest{
				SecretType: pb.SecretType_TEXT,
				Name:       "mytextsecret",
				Content:    []byte("another content"),
			},
			user:       testUser,
			alreadyHas: true,
		},
		{
			name: "Set text secret with same name but another user",
			setRequest: pb.SecretSetRequest{
				SecretType: pb.SecretType_TEXT,
				Name:       "mytextsecret",
			},
			user: secondUser,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), UserContextKey, &tt.user)
			getRes, err := server.SecretGet(ctx, &pb.SecretGetRequest{
				SecretType: tt.setRequest.SecretType,
				Name:       tt.setRequest.Name,
			})

			if tt.alreadyHas {
				require.NoError(t, err)
				assert.Equal(t, tt.setRequest.Name, getRes.Secret.Name)
				assert.Equal(t, tt.setRequest.SecretType, getRes.Secret.SecretType)
			} else {
				require.Error(t, err)
			}

			_, err = server.SecretSet(ctx, &tt.setRequest)
			if tt.alreadyHas {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				getRes, err = server.SecretGet(ctx, &pb.SecretGetRequest{
					SecretType: tt.setRequest.SecretType,
					Name:       tt.setRequest.Name,
				})

				require.NoError(t, err)
				assert.Equal(t, tt.setRequest.SecretType, getRes.Secret.SecretType)
				assert.Equal(t, tt.setRequest.Name, getRes.Secret.Name)
				assert.Equal(t, tt.setRequest.Content, getRes.Secret.Content)
			}
		})
	}
}

func TestServer_SecretUpdate(t *testing.T) {
	server, _, plain, err := NewTestServer()
	require.NoError(t, err)

	testUser := plainstorage.User{
		UUID:         uuid.New().String(),
		Login:        "testUser",
		PasswordHash: "somesecrethash",
	}

	secondUser := plainstorage.User{
		UUID:         uuid.New().String(),
		Login:        "secondUser",
		PasswordHash: "somesecrethash",
	}
	plain.Users = append(plain.Users, testUser, secondUser)

	testUserTextSecret := plainstorage.PlainSecret{
		Metadata: plainstorage.SecretMetadata{
			UUID:     uuid.New().String(),
			UserUUID: testUser.UUID,
			Name:     "testUserTextSecret1",
			Type:     plainstorage.SecretTypeText,
			Created:  time.Now(),
			Updated:  time.Now(),
		},
		Data: []byte("somedata"),
	}

	secondUserTextSecret := plainstorage.PlainSecret{
		Metadata: plainstorage.SecretMetadata{
			UUID:     uuid.New().String(),
			UserUUID: secondUser.UUID,
			Name:     "secondUserTextSecret1",
			Type:     plainstorage.SecretTypeText,
			Created:  time.Now(),
			Updated:  time.Now(),
		},
		Data: []byte("someeseconduserdata"),
	}
	plain.SecretList = append(plain.SecretList, testUserTextSecret, secondUserTextSecret)

	tests := []struct {
		name          string
		updateRequest pb.SecretUpdateRequest
		user          plainstorage.User
		secretExists  bool
	}{
		{
			name: "Update test secret",
			updateRequest: pb.SecretUpdateRequest{
				SecretType: pb.SecretType_TEXT,
				Name:       "testUserTextSecret1",
				Content:    []byte("my new text"),
			},
			user:         testUser,
			secretExists: true,
		},
		{
			name: "Update another user text",
			updateRequest: pb.SecretUpdateRequest{
				SecretType: pb.SecretType_TEXT,
				Name:       "secondUserTextSecret1",
				Content:    []byte("my new text"),
			},
			user:         testUser,
			secretExists: false,
		},
		{
			name: "Update another user text",
			updateRequest: pb.SecretUpdateRequest{
				SecretType: pb.SecretType_TEXT,
				Name:       "secondUserTextSecret1",
				Content:    []byte("my new text"),
			},
			user:         testUser,
			secretExists: false,
		},
		{
			name: "Update not existing text",
			updateRequest: pb.SecretUpdateRequest{
				SecretType: pb.SecretType_TEXT,
				Name:       "someText",
				Content:    []byte("my new text"),
			},
			user:         testUser,
			secretExists: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), UserContextKey, &tt.user)

			_, err := server.SecretUpdate(ctx, &tt.updateRequest)
			if tt.secretExists {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				return
			}

			res, err := server.SecretGet(ctx, &pb.SecretGetRequest{
				SecretType: tt.updateRequest.SecretType,
				Name:       tt.updateRequest.Name,
			})

			require.NoError(t, err)
			assert.Equal(t, tt.updateRequest.Content, res.Secret.Content)
			assert.Equal(t, tt.updateRequest.Name, res.Secret.Name)
			assert.Equal(t, tt.updateRequest.SecretType, res.Secret.SecretType)
		})
	}
}

func TestServer_SecretDelete(t *testing.T) {
	server, _, plain, err := NewTestServer()
	require.NoError(t, err)

	testUser := plainstorage.User{
		UUID:         uuid.New().String(),
		Login:        "testUser",
		PasswordHash: "somesecrethash",
	}

	secondUser := plainstorage.User{
		UUID:         uuid.New().String(),
		Login:        "secondUser",
		PasswordHash: "somesecrethash",
	}
	plain.Users = append(plain.Users, testUser, secondUser)

	ctx := context.WithValue(context.Background(), UserContextKey, &testUser)
	_, err = server.SecretSet(ctx, &pb.SecretSetRequest{
		SecretType: pb.SecretType_TEXT,
		Name:       "some_secret",
		Content:    []byte("my first new text"),
	})
	require.NoError(t, err)

	ctx = context.WithValue(context.Background(), UserContextKey, &secondUser)
	_, err = server.SecretSet(ctx, &pb.SecretSetRequest{
		SecretType: pb.SecretType_TEXT,
		Name:       "some_secret",
		Content:    []byte("my another text"),
	})
	require.NoError(t, err)

	tests := []struct {
		name          string
		deleteRequest pb.SecretDeleteRequest
		user          plainstorage.User
		secretExists  bool
	}{
		{
			name: "Delete exising secret",
			deleteRequest: pb.SecretDeleteRequest{
				SecretType: pb.SecretType_TEXT,
				SecretName: "some_secret",
			},
			user:         testUser,
			secretExists: true,
		},
		{
			name: "Delete secret again",
			deleteRequest: pb.SecretDeleteRequest{
				SecretType: pb.SecretType_TEXT,
				SecretName: "some_secret",
			},
			user:         testUser,
			secretExists: false,
		},
		{
			name: "Delete secret with same name for another user",
			deleteRequest: pb.SecretDeleteRequest{
				SecretType: pb.SecretType_TEXT,
				SecretName: "some_secret",
			},
			user:         secondUser,
			secretExists: true,
		},
		{
			name: "Delete secret with same name for another user again",
			deleteRequest: pb.SecretDeleteRequest{
				SecretType: pb.SecretType_TEXT,
				SecretName: "some_secret",
			},
			user:         secondUser,
			secretExists: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), UserContextKey, &tt.user)
			if tt.secretExists {
				getRes, err := server.SecretGet(ctx, &pb.SecretGetRequest{
					SecretType: tt.deleteRequest.SecretType,
					Name:       tt.deleteRequest.SecretName,
				})
				require.NoError(t, err)
				assert.Equal(t, tt.deleteRequest.SecretType, getRes.Secret.SecretType)
				assert.Equal(t, tt.deleteRequest.SecretName, getRes.Secret.Name)
			}

			_, err := server.SecretDelete(ctx, &tt.deleteRequest)

			if !tt.secretExists {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			_, err = server.SecretGet(ctx, &pb.SecretGetRequest{
				SecretType: tt.deleteRequest.SecretType,
				Name:       tt.deleteRequest.SecretName,
			})
			assert.Error(t, err)
		})
	}
}

func translateSecret(secret plainstorage.PlainSecret) (*pb.Secret, error) {
	st, err := translatePlainStorageSecretTypeToGRPC(secret.Metadata.Type)
	if err != nil {
		return nil, err
	}

	return &pb.Secret{
		SecretType:      st,
		Name:            secret.Metadata.Name,
		CreateTimestamp: secret.Metadata.Created.Unix(),
		UpdateTimestamp: secret.Metadata.Updated.Unix(),
		Content:         nil,
	}, nil
}

func translatePlainStorageSecretTypeToGRPC(secretType plainstorage.SecretType) (pb.SecretType, error) {
	switch secretType {
	case plainstorage.SecretTypeCredentials:
		return pb.SecretType_CREDENTIALS, nil
	case plainstorage.SecretTypeCard:
		return pb.SecretType_CREDIT_CARD, nil
	case plainstorage.SecretTypeText:
		return pb.SecretType_TEXT, nil
	case plainstorage.SecretTypeMedia:
		return pb.SecretType_MEDIA, nil
	}

	return 0, fmt.Errorf("undefined secret type %d", secretType)
}

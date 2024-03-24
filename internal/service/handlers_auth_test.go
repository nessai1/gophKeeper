package service

import (
	"context"
	pb "github.com/nessai1/gophkeeper/api/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestServer_RegisterAndLogin(t *testing.T) {
	s, _, plain, err := NewTestServer()
	require.NoError(t, err)

	userLogin := "login"
	userPassword := "password"

	_, err = s.Login(context.TODO(), &pb.UserCredentialsRequest{
		Login:    userLogin,
		Password: userPassword,
	})

	require.Error(t, err)

	resp, err := s.Register(context.TODO(), &pb.UserCredentialsRequest{
		Login:    userLogin,
		Password: userPassword,
	})

	require.NoError(t, err)

	userUUID, err := fetchUUID(resp.Token, s.config.SecretToken)
	require.NoError(t, err)
	user, err := plain.GetUserByUUID(context.TODO(), userUUID)
	require.NoError(t, err)
	assert.Equal(t, userLogin, user.Login)
	assert.Equal(t, hashPassword(userPassword, s.config.Salt), user.PasswordHash)

	_, err = s.Login(context.TODO(), &pb.UserCredentialsRequest{
		Login:    userLogin,
		Password: "some_another_pass",
	})
	require.Error(t, err)

	resp, err = s.Login(context.TODO(), &pb.UserCredentialsRequest{
		Login:    userLogin,
		Password: userPassword,
	})
	require.NoError(t, err)

	anotherUUID, err := fetchUUID(resp.Token, s.config.SecretToken)
	require.NoError(t, err)
	assert.Equal(t, anotherUUID, userUUID)
}

package service

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGenerateSignAndFetch(t *testing.T) {
	userUUID := uuid.New().String()
	secret := "somesecret"
	sign, err := generateSign(userUUID, secret)
	require.NoError(t, err)

	fetchedUUID, err := fetchUUID(sign, secret)
	require.NoError(t, err)
	assert.Equal(t, userUUID, fetchedUUID)

	fetchedUUID, err = fetchUUID("someShitString", secret)
	assert.Error(t, err)
}

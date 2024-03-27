package encrypt

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type userCredentials struct {
	login    string
	password string
}

func TestEncryptDecryptAES256(t *testing.T) {
	tests := []struct {
		name            string
		incomingMessage string

		sender   userCredentials
		receiver userCredentials
	}{
		{
			name:            "Same user send and receive short data",
			incomingMessage: "Some little message :)",
			sender: userCredentials{
				login:    "user1",
				password: "upass1",
			},
			receiver: userCredentials{
				login:    "user1",
				password: "upass1",
			},
		},
		{
			name:            "Same user send and receive big data",
			incomingMessage: "Some big message :)Some big message :)Some big message :)Some big message :)Some big message :)Some big message :)Some big message :)Some big message :)Some big message :)Some big message :)Some big message :)Some big message :)Some big message :)Some big message :)Some big message :)Some big message :)",
			sender: userCredentials{
				login:    "user1",
				password: "upass1",
			},
			receiver: userCredentials{
				login:    "user1",
				password: "upass1",
			},
		},
		{
			name:            "Stranger try to decrypt data",
			incomingMessage: "some message lol",
			sender: userCredentials{
				login:    "user1",
				password: "upass1",
			},
			receiver: userCredentials{
				login:    "user2",
				password: "upass2",
			},
		},
		{
			name:            "Stranger try to decrypt big data",
			incomingMessage: "some message lolsome message lolsome message lolsome message lolsome message lolsome message lolsome message lolsome message lolsome message lolsome message lolsome message lolsome message lolsome message lolsome message lolsome message lolsome message lol",
			sender: userCredentials{
				login:    "user1",
				password: "upass1",
			},
			receiver: userCredentials{
				login:    "user2",
				password: "upass2",
			},
		},
		{
			name: "User with wrong password try to decrypt data",
			sender: userCredentials{
				login:    "user1",
				password: "upass1",
			},
			receiver: userCredentials{
				login:    "user1",
				password: "upass2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inKey := BuildAESKey(tt.sender.login, tt.sender.password)

			encryptedMsg, err := EncryptAES256([]byte(tt.incomingMessage), inKey)
			require.NoError(t, err)

			outKey := BuildAESKey(tt.receiver.login, tt.receiver.password)
			decryptedMsg, err := DecryptAES256(encryptedMsg, outKey)

			if tt.sender.login == tt.receiver.login && tt.sender.password == tt.receiver.password {
				require.NoError(t, err)
				assert.Equal(t, []byte(tt.incomingMessage), decryptedMsg)
			} else {
				require.Error(t, err)
				assert.NotEqual(t, []byte(tt.incomingMessage), decryptedMsg)
			}
		})
	}
}

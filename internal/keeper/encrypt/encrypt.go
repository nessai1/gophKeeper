package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"fmt"
	"math/rand"
)

func BuildAESKey(login string, password string) [32]byte {
	return sha256.Sum256([]byte(login + password))
}

// EncryptAES256 encrypt data by AES256 algorithm
func EncryptAES256(data []byte, passphrase [32]byte) ([]byte, error) {
	c, err := aes.NewCipher(passphrase[:])
	if err != nil {
		return nil, fmt.Errorf("cannot create cipher for encrypt by AES256: %w", err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, fmt.Errorf("cannot create GCM for encrypt by AES256: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, fmt.Errorf("cannot read nonce for encrypt by AES256: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	return ciphertext, nil
}

// DecryptAES256 decrypt data by AES256 algorithm
func DecryptAES256(ciphertext []byte, passphrase [32]byte) ([]byte, error) {
	key := passphrase[:]
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("cannot create new cipher while decrypt AES256: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("cannot create new gcm while decrypt AES256: %w", err)
	}

	// Since we know the ciphertext is actually nonce+ciphertext
	// And len(nonce) == NonceSize(). We can separate the two.
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot open gcm while decrypt AES256: %w", err)
	}

	return plaintext, nil
}

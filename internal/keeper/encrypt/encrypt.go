package encrypt

import "crypto/sha256"

func BuildAESKey(login string, password string) [32]byte {
	return sha256.Sum256([]byte(login + password))
}

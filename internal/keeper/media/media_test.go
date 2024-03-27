package media

import (
	"bytes"
	"context"
	"fmt"
	"github.com/nessai1/gophkeeper/internal/keeper/encrypt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestEncryptDecryptFile(t *testing.T) {
	someFileContent := []byte("some big content for write to filesome big content for write to filesome big content for write to filesome big content for write to filesome big content for write to filesome big content for write to filesome big content for write to filesome big content for write to filesome big content for write to filesome big content for write to filesome big content for write to filesome big content for write to filesome big content for write to filesome big content for write to filesome big content for write to file")

	source, err := os.CreateTemp(os.TempDir(), "encrypt_test_source")
	require.NoError(t, err)

	_, err = source.Write(someFileContent)
	require.NoError(t, err)

	key := encrypt.BuildAESKey("somelogin", "somepassword")

	encryptDest, err := EncryptFile(context.TODO(), source, filepath.Join(os.TempDir(), "encrypt_text_dest"+time.Now().String()), key)
	require.NoError(t, err)

	decryptDest, err := DecryptFile(context.TODO(), encryptDest, filepath.Join(os.TempDir(), "decrypt_text_dest"+time.Now().String()), key)
	require.NoError(t, err)

	bf := bytes.Buffer{}
	_, err = bf.ReadFrom(decryptDest)
	require.NoError(t, err)

	fmt.Println(string(bf.Bytes()))
	fmt.Println(string(someFileContent))
	assert.Equal(t, someFileContent, bf.Bytes())
}

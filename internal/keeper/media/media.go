package media

import (
	"context"
	"errors"
	"fmt"
	"github.com/nessai1/gophkeeper/internal/keeper/encrypt"
	"io"
	"os"
)

// ContentSectionSize size of content block
const ContentSectionSize = 256

// BlockSize size of encrypted content block with trailer
const BlockSize = 284

// EncryptFile encrypt file by AES algorithm and returns encrypted file descriptor
func EncryptFile(_ context.Context, file *os.File, destination string, key [32]byte) (*os.File, error) {

	output, err := os.OpenFile(destination, os.O_CREATE|os.O_RDWR|os.O_EXCL, 0666)
	if err != nil {
		return nil, fmt.Errorf("cannot create file %s for encrypt: %w", destination, err)
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		output.Close()
		return nil, fmt.Errorf("cannot seek file to the start while encrypt: %w", err)
	}

	var bf []byte
	for {
		bf = make([]byte, ContentSectionSize)
		n, err := file.Read(bf)
		if n != ContentSectionSize {
			bf = bf[:n]
		}

		if errors.Is(err, io.EOF) {
			break
		}

		content, err := encrypt.EncryptAES256(bf, key)
		if err != nil {
			output.Close()
			return nil, fmt.Errorf("error while encrypt file content: %w", err)
		}

		_, err = output.Write(content)
		if err != nil {
			output.Close()
			return nil, fmt.Errorf("error while write encrypted content to the file: %w", err)
		}
	}

	_, err = output.Seek(0, 0)
	if err != nil {
		output.Close()
		return nil, fmt.Errorf("cannot seek output file: %w", err)
	}

	return output, nil
}

// DecryptFile decrypt file by AES algorithm and returns decrypted file descriptor
func DecryptFile(_ context.Context, file *os.File, destination string, key [32]byte) (*os.File, error) {
	output, err := os.OpenFile(destination, os.O_CREATE|os.O_RDWR|os.O_EXCL, 0666)
	if err != nil {
		return nil, fmt.Errorf("cannot create file %s for decrypt: %w", destination, err)
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		output.Close()
		return nil, fmt.Errorf("cannot seek file to the start while decrypt: %w", err)
	}

	var bf []byte
	for {
		bf = make([]byte, BlockSize)
		n, err := file.Read(bf)
		if n != BlockSize {
			bf = bf[:n]
		}

		if errors.Is(err, io.EOF) {
			break
		}

		content, err := encrypt.DecryptAES256(bf, key)
		if err != nil {
			output.Close()
			return nil, fmt.Errorf("error while decrypt file content: %w", err)
		}

		_, err = output.Write(content)
		if err != nil {
			output.Close()
			return nil, fmt.Errorf("error while write decrypted content to the file: %w", err)
		}
	}

	_, err = output.Seek(0, 0)
	if err != nil {
		output.Close()
		return nil, fmt.Errorf("cannot seek output file: %w", err)
	}

	return output, nil
}

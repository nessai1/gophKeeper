package mediastorage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// File storage contains files in folder, using for tests only

type MediaStorageLocal struct {
	StorageDir string
}

func (m *MediaStorageLocal) StartUpload(_ context.Context, key string) (MultipartUpload, error) {
	_, err := os.Stat(filepath.Join(m.StorageDir, key))
	if err != nil && errors.Is(err, os.ErrNotExist) {
		f, err := os.Open(filepath.Join(m.StorageDir, key))
		if err != nil {
			return nil, err
		}

		return &MultipartLocal{File: f}, nil
	}

	if err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("file exists")
}

func (m *MediaStorageLocal) StartDownload(_ context.Context, key string) (io.ReadCloser, error) {
	if _, err := os.Stat(filepath.Join(m.StorageDir, key)); errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("file %s doesnt exists", key)
	}

	f, err := os.Open(filepath.Join(m.StorageDir, key))

	return f, err
}

func (m *MediaStorageLocal) Delete(_ context.Context, key string) error {
	if _, err := os.Stat(filepath.Join(m.StorageDir, key)); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("file %s doesnt exists", key)
	}

	return os.Remove(filepath.Join(m.StorageDir, key))
}

type MultipartLocal struct {
	File *os.File
}

func (m *MultipartLocal) Upload(_ context.Context, content []byte) error {
	_, err := m.File.Write(content)
	return err
}

func (m *MultipartLocal) Complete(_ context.Context) error {
	return m.File.Close()
}

func (m *MultipartLocal) Abort(_ context.Context) error {
	n := m.File.Name()
	err := m.File.Close()
	if err != nil {
		return err
	}

	return os.Remove(n)
}

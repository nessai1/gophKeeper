package service

import (
	"github.com/nessai1/gophkeeper/internal/service/config"
	"github.com/nessai1/gophkeeper/internal/service/mediastorage"
	"github.com/nessai1/gophkeeper/internal/service/plainstorage"
	"go.uber.org/zap"
	"os"
)

func NewTestServer() (*Server, *mediastorage.MediaStorageLocal, *plainstorage.MemoryStorage, error) {
	path, err := os.MkdirTemp(os.TempDir(), "test_media_storage")
	if err != nil {
		return nil, nil, nil, err
	}
	media := mediastorage.MediaStorageLocal{StorageDir: path}
	plain := plainstorage.MemoryStorage{
		Users:      make([]plainstorage.User, 0),
		SecretList: make([]plainstorage.PlainSecret, 0),
	}

	s := Server{
		plainStorage: &plain,
		mediaStorage: &media,
		logger:       zap.NewNop(),
		config: config.Config{
			Address:     "localhost",
			SecretToken: "somesecret",
			Salt:        "somesalt",
		},
	}

	return &s, &media, &plain, nil
}

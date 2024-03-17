package connector

import (
	"context"
	"github.com/nessai1/gophkeeper/internal/keeper/secret"
	"io"
	"os"
)

type ServiceConnector interface {
	Ping(ctx context.Context) (answer string, error error)

	Register(ctx context.Context, login string, password string) (token string, err error)
	Login(ctx context.Context, login string, password string) (token string, err error)

	SetAuthToken(token string)

	UploadMedia(ctx context.Context, name string, reader io.Reader) (string, error)
	DownloadMedia(ctx context.Context, name string, dest string) (*os.File, error)

	ListSecret(ctx context.Context, secretType secret.SecretType) ([]secret.Secret, error)
}

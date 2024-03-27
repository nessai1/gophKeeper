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

	UploadMedia(ctx context.Context, name string, reader io.Reader, replace bool) (string, error)
	DownloadMedia(ctx context.Context, name string, dest string) (*os.File, error)

	ListSecret(ctx context.Context, secretType secret.SecretType) ([]secret.Secret, error)
	SetSecret(ctx context.Context, name string, secretType secret.SecretType, data []byte) error
	UpdateSecret(ctx context.Context, name string, secretType secret.SecretType, data []byte) error
	RemoveSecret(ctx context.Context, name string, secretType secret.SecretType) error
	GetSecret(ctx context.Context, name string, secretType secret.SecretType) ([]byte, error)
}

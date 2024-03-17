package mediastorage

import (
	"context"
	"io"
)

type MediaStorage interface {
	StartUpload(ctx context.Context, key string) (MultipartUpload, error)
	StartDownload(ctx context.Context, key string) (io.ReadCloser, error)
}

type MultipartUpload interface {
	Upload(ctx context.Context, content []byte) error
	Complete(ctx context.Context) error
	Abort(ctx context.Context) error
}

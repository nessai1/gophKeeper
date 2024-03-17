package performer

import (
	"context"
	"fmt"
	"github.com/nessai1/gophkeeper/internal/keeper/connector"
	"github.com/nessai1/gophkeeper/internal/keeper/media"
	"github.com/nessai1/gophkeeper/internal/keeper/session"
	"go.uber.org/zap"
	"io"
	"os"
	"path/filepath"
	"time"
)

type secretMediaPerformer struct {
	conn    connector.ServiceConnector
	session session.Session
	logger  *zap.Logger
	workDir string
}

// Set encrypt by AES file with path == name and save it on external service
func (p *secretMediaPerformer) Set(ctx context.Context, name string) error {
	file, err := os.Open(name)
	if err != nil {
		return fmt.Errorf("cannot find file '%s': %w", name, err)
	}

	encryptedFile, err := media.EncryptFile(ctx, file, filepath.Join(os.TempDir(), filepath.Base(name)+"_"+time.Now().String()+".encrypted"), p.session.SecretKey)
	if err != nil {
		p.logger.Error("cannot encrypt media file", zap.Error(err), zap.String("filepath", name))

		return fmt.Errorf("cannot encrypt file, see logs")
	}

	encryptedFileInfo, err := encryptedFile.Stat()
	if err != nil {
		p.logger.Error("cannot get encrypted file info", zap.Error(err), zap.String("filepath", encryptedFileInfo.Name()))

		return fmt.Errorf("cannot get encrypt file, see logs")
	}

	encryptedFile.Seek(0, 0)
	id, err := p.conn.UploadMedia(ctx, filepath.Base(name), encryptedFile)
	if err != nil {
		p.logger.Error("Errror while upload new media", zap.String("filename", filepath.Base(name)), zap.Error(err))

		return fmt.Errorf("cannot send media file to server: %w", err)
	}

	p.logger.Info("client set media secret", zap.String("login", p.session.Login), zap.String("filepath", name), zap.Int64("size", encryptedFileInfo.Size()))
	fmt.Printf("Media file %s successful saved! Encrypted size: %d bytes\tUUID: %s\n", filepath.Base(name), encryptedFileInfo.Size(), id)

	copyFile, err := os.Create(filepath.Join(p.workDir, "media", filepath.Base(name)))
	if err != nil {
		p.logger.Error("cannot copy file to media dir", zap.Error(err))

		return nil
	}

	defer copyFile.Close()
	file.Seek(0, 0)
	_, err = io.Copy(copyFile, file)
	if err != nil {
		p.logger.Error("cannot copy file to media dir", zap.Error(err))
	}

	return nil
}

func (p *secretMediaPerformer) Get(ctx context.Context, name string) error {
	//TODO implement me
	panic("implement me")
}

func (p *secretMediaPerformer) Update(ctx context.Context, name string) error {
	//TODO implement me
	panic("implement me")
}

func (p *secretMediaPerformer) Delete(ctx context.Context, name string) error {
	//TODO implement me
	panic("implement me")
}

func (p *secretMediaPerformer) List(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

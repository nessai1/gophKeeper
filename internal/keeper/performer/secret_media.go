package performer

import (
	"context"
	"fmt"
	"github.com/nessai1/gophkeeper/internal/keeper/connector"
	"github.com/nessai1/gophkeeper/internal/keeper/media"
	"github.com/nessai1/gophkeeper/internal/keeper/secret"
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
	return p.uploadFile(ctx, name, false)
}

func (p *secretMediaPerformer) uploadFile(ctx context.Context, name string, replace bool) error {
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
	id, err := p.conn.UploadMedia(ctx, filepath.Base(name), encryptedFile, replace)
	if err != nil {
		p.logger.Error("Errror while upload new media", zap.String("filename", filepath.Base(name)), zap.Error(err))

		return fmt.Errorf("cannot send media file to server: %w", err)
	}

	p.logger.Info("client set media secret", zap.String("login", p.session.Login), zap.String("filepath", name), zap.Int64("size", encryptedFileInfo.Size()))
	fmt.Printf("Media file %s successful saved! Encrypted size: %d bytes\tUUID: %s\n", filepath.Base(name), encryptedFileInfo.Size(), id)

	copyFile, err := os.Create(filepath.Join(p.workDir, "media", filepath.Base(name)))
	if err != nil {
		p.logger.Error("cannot copy file to media dir", zap.Error(err))

		return fmt.Errorf("cannot copy file to media dir: %w", err)
	}

	defer func() {
		err := copyFile.Close()
		if err != nil {
			p.logger.Error("cannot close copy file for media set", zap.Error(err))
		}
	}()

	file.Seek(0, 0)
	_, err = io.Copy(copyFile, file)
	if err != nil {
		p.logger.Error("cannot copy file to media dir", zap.Error(err))
	}

	return nil
}

func (p *secretMediaPerformer) Get(ctx context.Context, name string) error {
	f, err := p.conn.DownloadMedia(ctx, name, filepath.Join(os.TempDir(), filepath.Base(name)+"_"+time.Now().String()+".encrypted"))
	if err != nil {
		p.logger.Error("Cannot download media", zap.String("login", p.session.Login), zap.String("filename", name), zap.Error(err))

		return fmt.Errorf("cannot download media %s: %w", name, err)
	}

	defer func() {
		err := f.Close()
		if err != nil {
			p.logger.Error("cannot close encrypted file for media download", zap.Error(err))
		}
	}()

	decryptedFile, err := media.DecryptFile(ctx, f, filepath.Join(p.workDir, "media", name), p.session.SecretKey)
	if err != nil {
		p.logger.Error("Cannot decrypt downloaded media", zap.String("login", p.session.Login), zap.String("filename", name), zap.Error(err))

		return fmt.Errorf("cannot decrypt downloaded file: %w", err)
	}

	err = decryptedFile.Close()
	if err != nil {
		p.logger.Error("Cannot close downloaded media", zap.String("login", p.session.Login), zap.String("filename", name), zap.Error(err))

		return fmt.Errorf("cannot decrypt downloaded file: %w", err)
	}

	fmt.Printf("\033[32mFile %s successfuly download and stored in %s/media directory\033[0m\n", name, p.workDir)

	return nil
}

func (p *secretMediaPerformer) Update(ctx context.Context, name string) error {
	if p.isFileLoaded(name) {
		err := p.Delete(ctx, name)
		if err != nil {
			p.logger.Error("Cannot delete old file", zap.Error(err))
		}
	}

	return p.uploadFile(ctx, name, true)
}

func (p *secretMediaPerformer) Delete(ctx context.Context, name string) error {
	err := p.conn.RemoveSecret(ctx, name, secret.SecretTypeMedia)
	if err != nil {
		return fmt.Errorf("cannot remove media secret: %w", err)
	}

	if p.isFileLoaded(name) {
		err = os.Remove(filepath.Join(p.workDir, "media", name))
		p.logger.Error("Cannot remove file from local storage", zap.Error(err), zap.String("filename", name))
	}

	fmt.Printf("\033[32mFile %s successfuly removed!\033[0m\n", name)

	return nil
}

func (p *secretMediaPerformer) List(ctx context.Context) error {
	secrets, err := p.conn.ListSecret(ctx, secret.SecretTypeMedia)
	if err != nil {
		p.logger.Error("Cannot list media secrets", zap.Error(err))

		return fmt.Errorf("cannot list media secrets: %w", err)
	}

	if len(secrets) == 0 {
		fmt.Printf("No media found\n")

		return nil
	}

	printable := make([]printableSecret, len(secrets))
	for i, v := range secrets {

		name := v.Name
		if p.isFileLoaded(name) {
			name = fmt.Sprintf("\033[32m%s\033[0m", name)
		} else {
			name = fmt.Sprintf("\033[33m%s\033[0m", name)
		}

		printable[i] = printableSecret{
			Name:        name,
			Create_time: v.Created.String(),
			Update_time: v.Updated.String(),
		}
	}

	fmt.Printf("\033[32mexists in media \t \033[33mstored only in service\033[0m\n")
	printSecrets(printable)

	return nil
}

func (p *secretMediaPerformer) isFileLoaded(name string) bool {
	_, err := os.Stat(filepath.Join(p.workDir, "media", name))
	if err != nil {
		return false
	}

	return true
}

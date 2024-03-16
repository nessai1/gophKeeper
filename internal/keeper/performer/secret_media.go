package performer

import (
	"context"
	"fmt"
	"github.com/nessai1/gophkeeper/internal/keeper/connector"
	"go.uber.org/zap"
	"os"
)

type secretMediaPerformer struct {
	conn      connector.ServiceConnector
	sessional Sessional
	logger    *zap.Logger
	workDir   string
}

// Set encrypt by AES file with path == name and save it on external service
func (p *secretMediaPerformer) Set(ctx context.Context, name string) error {
	file, err := os.Open(name)
	if err != nil {
		return fmt.Errorf("cannot find file '%s': %w", name, err)
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

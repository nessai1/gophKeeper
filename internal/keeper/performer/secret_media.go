package performer

import (
	"context"
	"fmt"
	"github.com/nessai1/gophkeeper/internal/keeper/connector"
	"go.uber.org/zap"
)

type secretMediaPerformer struct {
	conn      connector.ServiceConnector
	sessional Sessional
	logger    *zap.Logger
	workDir   string
}

func (p *secretMediaPerformer) Set(ctx context.Context, name string) error {
	fmt.Printf("action set for media.....\n")
	// TODO: find file in media dir; AES encrypt by session secret key; send to server;
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

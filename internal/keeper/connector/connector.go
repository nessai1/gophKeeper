package connector

import "context"

type ServiceConnector interface {
	Ping(ctx context.Context) (answer string, error error)

	Register(ctx context.Context, login string, password string) (token string, err error)
	Login(ctx context.Context, login string, password string) (token string, err error)
}

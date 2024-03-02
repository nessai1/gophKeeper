package connector

import "context"

type ServiceConnector interface {
	Ping(ctx context.Context) (answer string, error error)
}

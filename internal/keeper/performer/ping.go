package performer

import (
	"context"
	"fmt"
	"github.com/nessai1/gophkeeper/internal/keeper/connector"
	"go.uber.org/zap"
)

type Ping struct {
}

func (p Ping) GetName() string {
	return "ping"
}

func (p Ping) GetStruct() string {
	return "ping"
}

func (p Ping) GetDescription() string {
	return "ping to the external service"
}

func (p Ping) GetDetailDescription() string {
	return "Ping to the external service. If service exists and works - it returns 'pong!' answer"
}

func (p Ping) Execute(conn connector.ServiceConnector, _ Sessional, logger *zap.Logger, _ []string, _ string) (requireExit bool, err error) {
	answer, err := conn.Ping(context.TODO())
	if err != nil {
		logger.Error("Error while ping service", zap.Error(err))
		return false, fmt.Errorf("ping error: %w", err)
	} else {
		fmt.Printf("Answer: %s\n", answer)
	}

	return false, nil
}

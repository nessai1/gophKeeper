package performer

import (
	"context"
	"fmt"
	"github.com/nessai1/gophkeeper/internal/keeper/connector"
	"github.com/nessai1/gophkeeper/pkg/command"
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

func (p Ping) Execute(input command.Readable, output command.Writable, conn connector.ServiceConnector, sessional Sessional, logger *zap.Logger, args []string) (requireExit bool, err error) {
	answer, err := conn.Ping(context.TODO())
	if err != nil {
		logger.Error("Error while ping service", zap.Error(err))
		return false, fmt.Errorf("ping error: %w", err)
	} else {
		output.Write([]byte(fmt.Sprintf("Answer: %s\n", answer)))
	}
	return false, nil
}

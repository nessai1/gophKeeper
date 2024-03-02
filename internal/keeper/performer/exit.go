package performer

import (
	"github.com/nessai1/gophkeeper/internal/keeper/connector"
	"github.com/nessai1/gophkeeper/pkg/command"
	"go.uber.org/zap"
)

type Exit struct {
}

func (p Exit) GetName() string {
	return "exit"
}

func (p Exit) GetStruct() string {
	return "exit"
}

func (p Exit) GetDescription() string {
	return "exit from application"
}

func (p Exit) GetDetailDescription() string {
	return "Exits the application, that's all :)"
}

func (p Exit) Execute(input command.Readable, output command.Writable, conn connector.ServiceConnector, sessional Sessional, logger *zap.Logger, args []string) (requireExit bool, err error) {
	return true, nil
}

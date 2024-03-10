package performer

import (
	"github.com/nessai1/gophkeeper/internal/keeper/connector"
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

func (p Exit) Execute(_ connector.ServiceConnector, _ Sessional, _ *zap.Logger, _ []string, _ string) (requireExit bool, err error) {
	return true, nil
}

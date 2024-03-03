package performer

import (
	"fmt"
	"github.com/nessai1/gophkeeper/internal/keeper/connector"
	"go.uber.org/zap"
)

type Register struct {
}

func (p Register) GetName() string {
	return "register"
}

func (p Register) GetStruct() string {
	return "register"
}

func (p Register) GetDescription() string {
	return "register in external service"
}

func (p Register) GetDetailDescription() string {
	return "Register in external service.\nFor registration client need has unauthorized session before this command, use 'logout' command for this"
}

func (p Register) Execute(_ connector.ServiceConnector, sessional Sessional, _ *zap.Logger, _ []string) (requireExit bool, err error) {
	if sessional.GetSession() != nil {
		return false, fmt.Errorf("you need to unauthorize by 'logout' before")
	}

	return true, nil
}

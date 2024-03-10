package performer

import (
	"fmt"
	"github.com/nessai1/gophkeeper/internal/keeper/connector"
	"go.uber.org/zap"
)

type Secret struct {
}

func (p Secret) GetName() string {
	return "register"
}

func (p Secret) GetStruct() string {
	return "secret [type] [action] [name]"
}

func (p Secret) GetDescription() string {
	return "Manipulate with user secrets (login-password, credit card, text, media)"
}

func (p Secret) GetDetailDescription() string {
	return `
	Hui
	jopa
`
}

func (p Secret) Execute(conn connector.ServiceConnector, sessional Sessional, logger *zap.Logger, _ []string) (requireExit bool, err error) {
	if sessional.GetSession() == nil {
		return false, fmt.Errorf("for working with secrets you need to be authorized")
	}

	s := sessional.GetSession()
	fmt.Printf("working with secrets of user %s", s.Login)
	return false, nil
}

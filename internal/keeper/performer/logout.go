package performer

import (
	"fmt"
	"github.com/nessai1/gophkeeper/internal/keeper/connector"
	"go.uber.org/zap"
)

type Logout struct {
}

func (p Logout) GetName() string {
	return "logout"
}

func (p Logout) GetStruct() string {
	return "logout"
}

func (p Logout) GetDescription() string {
	return "logout from current session"
}

func (p Logout) GetDetailDescription() string {
	return "Logout from current session."
}

func (p Logout) Execute(conn connector.ServiceConnector, sessional Sessional, logger *zap.Logger, _ []string, _ string) (requireExit bool, err error) {
	if sessional.GetSession() == nil {
		return false, fmt.Errorf("you already logouted")
	}

	sessionLogin := sessional.GetSession().Login
	sessional.SetSession(nil)
	fmt.Printf("\033[32mSuccessful logout from %s!\033[0m\n", sessionLogin)

	return false, nil
}

package performer

import (
	"context"
	"fmt"
	"github.com/nessai1/gophkeeper/internal/keeper/connector"
	"github.com/nessai1/gophkeeper/internal/keeper/session"
	"github.com/nessai1/gophkeeper/pkg/command"
	"go.uber.org/zap"
)

type Login struct {
}

func (p Login) GetName() string {
	return "login"
}

func (p Login) GetStruct() string {
	return "login"
}

func (p Login) GetDescription() string {
	return "login in external service"
}

func (p Login) GetDetailDescription() string {
	return "Login in external service.\nFor login client need has unauthorized session before this command, use 'logout' command for this"
}

func (p Login) Execute(conn connector.ServiceConnector, sessional Sessional, logger *zap.Logger, _ []string, _ string) (requireExit bool, err error) {
	if sessional.GetSession() != nil {
		return false, fmt.Errorf("you need to unauthorize by 'logout' before")
	}

	login, err := command.AskText("Enter login")
	if err != nil {
		return false, fmt.Errorf("cannot read login for register: %w", err)
	}

	password, err := command.AskSecret("Enter password")
	if err != nil {
		return false, fmt.Errorf("cannot read password for register: %w", err)
	}

	t, err := conn.Login(context.TODO(), login, password)
	if err != nil {
		logger.Error("Got service error while login", zap.Error(err))
		return false, fmt.Errorf("error while execute login command: %w", err)
	}

	s := session.NewSession(login, password, t)
	sessional.SetSession(&s)
	fmt.Printf("\033[32mSuccessful login as %s!\033[0m\n", s.Login)

	return false, nil
}

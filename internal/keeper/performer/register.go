package performer

import (
	"context"
	"fmt"
	"github.com/nessai1/gophkeeper/internal/keeper/connector"
	"github.com/nessai1/gophkeeper/internal/keeper/session"
	"github.com/nessai1/gophkeeper/pkg/command"
	"go.uber.org/zap"
	"strings"
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

func (p Register) Execute(conn connector.ServiceConnector, sessional Sessional, logger *zap.Logger, _ []string) (requireExit bool, err error) {
	if sessional.GetSession() != nil {
		return false, fmt.Errorf("you need to unauthorize by 'logout' before")
	}

	login, err := command.AskText("Enter login")
	if err != nil {
		return false, fmt.Errorf("cannot read login for register: %w", err)
	}

	if strings.TrimSpace(login) == "" {
		return false, fmt.Errorf("login can't be empty")
	}

	password, err := command.AskSecret("Enter password")
	if err != nil {
		return false, fmt.Errorf("cannot read password for register: %w", err)
	}

	if strings.TrimSpace(password) == "" {
		return false, fmt.Errorf("password can't be empty")
	}

	passwordRep, err := command.AskSecret("Repeat password")
	if err != nil {
		return false, fmt.Errorf("cannot read password for register: %w", err)
	}

	if password != passwordRep {
		return false, fmt.Errorf("passwords are not equal")
	}

	t, err := conn.Register(context.TODO(), login, password)
	if err != nil {
		logger.Error("Got service error while register", zap.Error(err))
		return false, fmt.Errorf("error while execute register command: %w", err)
	}

	s := session.NewSession(login, password, t)
	sessional.SetSession(&s)
	fmt.Printf("\033[32mSuccessful register and login as %s!\033[0m\n", s.Login)

	return false, nil
}

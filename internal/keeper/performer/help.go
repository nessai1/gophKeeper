package performer

import (
	"fmt"
	"github.com/nessai1/gophkeeper/internal/keeper/connector"
	"github.com/nessai1/gophkeeper/pkg/command"
	"go.uber.org/zap"
)

type Help struct {
}

func (p Help) GetName() string {
	return "help"
}

func (p Help) GetStruct() string {
	return "help [command]"
}

func (p Help) GetDescription() string {
	return "get help information about command"
}

func (p Help) GetDetailDescription() string {
	return "Get common help information and list of available commands\nIf has argument [command] - returns detail description about concrete command"
}

func (p Help) Execute(input command.Readable, output command.Writable, conn connector.ServiceConnector, sessional Sessional, logger *zap.Logger, args []string) (requireExit bool, err error) {
	if len(args) > 2 {
		return false, fmt.Errorf("command has too many arguments (%d, requires 1)", len(args)-1)
	}

	if len(args) == 1 {
		p.printCommandsList(output)
	} else {
		err := p.printCommandDetails(output, args[1])
		if err != nil {
			return false, fmt.Errorf("command details error: %w", err)
		}
	}

	return false, nil
}

func (p Help) printCommandsList(output command.Writable) {
	output.Write([]byte("List of commands\n--------\n"))
	for _, val := range AvailablePerformers {
		output.Write([]byte(fmt.Sprintf("%s\t%s\n", val.GetStruct(), val.GetDescription())))
	}
}

func (p Help) printCommandDetails(output command.Writable, commandName string) error {
	for _, val := range AvailablePerformers {
		if val.GetName() == commandName {
			output.Write([]byte(fmt.Sprintf("Command: %s\nPattern: %s\n\n%s\n", val.GetName(), val.GetStruct(), val.GetDetailDescription())))
			return nil
		}
	}

	return fmt.Errorf("undefined command '%s'", commandName)
}

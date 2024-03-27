package performer

import (
	"fmt"
	"github.com/nessai1/gophkeeper/internal/keeper/connector"
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

func (p Help) Execute(_ connector.ServiceConnector, _ Sessional, _ *zap.Logger, args []string, _ string) (requireExit bool, err error) {
	if len(args) > 2 {
		return false, fmt.Errorf("command has too many arguments (%d, requires 1)", len(args)-1)
	}

	if len(args) == 1 {
		p.printCommandsList()
	} else {
		err := p.printCommandDetails(args[1])
		if err != nil {
			return false, fmt.Errorf("command details error: %w", err)
		}
	}

	return false, nil
}

func (p Help) printCommandsList() {
	fmt.Print("List of commands\n--------\n")
	for _, val := range AvailablePerformers {
		fmt.Printf("%s\t%s\n", val.GetName(), val.GetDescription())
	}
}

func (p Help) printCommandDetails(commandName string) error {
	val, ok := AvailablePerformers[commandName]
	if !ok {
		return fmt.Errorf("undefined command '%s'", commandName)
	}

	fmt.Printf("Command: %s\nPattern: %s\n\n%s\n", val.GetName(), val.GetStruct(), val.GetDetailDescription())
	return nil
}

package command

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Command info about command
type Command struct {
	// Name first word of given command line
	Name string
	// Args arguments of written command, separated by spaces
	Args []string
}

// ReadCommand prompts user to enter a command to input, writes command anchor to output
func ReadCommand() (*Command, error) {
	fmt.Print("> ")

	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("cannot read command: %w", err)
	}

	strs := strings.Split(text, " ")
	strs[len(strs)-1] = strings.Trim(strs[len(strs)-1], "\n")

	return &Command{
		Name: strings.TrimSpace(strs[0]),
		Args: strs,
	}, nil
}

//
//func AskSecret(output Writable, welcomeText string) {
//	term.ReadPassword()
//	term.NewTerminal(output)
//}

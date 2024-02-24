package command

import (
	"fmt"
	"io"
	"strings"
)

// Readable interface for read data from source, usually its os.Stdout wrapper, like bufio.Reader
type Readable interface {
	io.Reader
	ReadString(delim byte) (string, error)
}

// Writable interface for write commands help info, like anchors or welcome messages (ex. 'Write password: ')
type Writable interface {
	io.Writer
}

// Command info about command
type Command struct {
	// Name first word of given command line
	Name string
	// Args arguments of written command, separated by spaces
	Args []string
}

// ReadCommand prompts user to enter a command to input, writes command anchor to output
func ReadCommand(input Readable, output Writable) (*Command, error) {
	_, err := output.Write([]byte("> "))
	if err != nil {
		return nil, fmt.Errorf("cannot write command anchor: %w", err)
	}

	text, err := input.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("cannot read command: %w", err)
	}

	strs := strings.Split(text, " ")

	return &Command{
		Name: strings.TrimSpace(strs[0]),
		Args: strs,
	}, nil
}

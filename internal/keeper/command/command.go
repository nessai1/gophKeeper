package command

import (
	"fmt"
	"io"
	"strings"
)

type Readable interface {
	io.Reader
	ReadString(delim byte) (string, error)
}

type Writable interface {
	io.Writer
}

type Command struct {
	Name string
	Args []string
}

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

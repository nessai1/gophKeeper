package command

import (
	"bufio"
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReadCommand(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		expected Command
	}{
		{
			name:    "Empty command",
			command: "",
			expected: Command{
				Name: "",
				Args: []string{"\n"},
			},
		},
		{
			name:    "Single word",
			command: "hello",
			expected: Command{
				Name: "hello",
				Args: []string{"hello\n"},
			},
		},
		{
			name:    "Two words",
			command: "hello word",
			expected: Command{
				Name: "hello",
				Args: []string{"hello", "word\n"},
			},
		},
		{
			name:    "Three words",
			command: "hello word --iLove=go",
			expected: Command{
				Name: "hello",
				Args: []string{"hello", "word", "--iLove=go\n"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputCmd := []byte(tt.command + "\n")
			reader := bytes.NewReader(inputCmd)
			writer := bytes.Buffer{}

			cmd, err := ReadCommand(bufio.NewReader(reader), &writer)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, *cmd)
		})
	}
}

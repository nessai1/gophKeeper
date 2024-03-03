package command

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
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
				Args: []string{""},
			},
		},
		{
			name:    "Single word",
			command: "hello",
			expected: Command{
				Name: "hello",
				Args: []string{"hello"},
			},
		},
		{
			name:    "Two words",
			command: "hello word",
			expected: Command{
				Name: "hello",
				Args: []string{"hello", "word"},
			},
		},
		{
			name:    "Three words",
			command: "hello word --iLove=go",
			expected: Command{
				Name: "hello",
				Args: []string{"hello", "word", "--iLove=go"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rollback, err := mockStd()
			require.NoError(t, err)
			defer rollback()

			_, err = os.Stdin.Write([]byte(tt.command + "\n"))
			require.NoError(t, err)

			_, err = os.Stdin.Seek(0, 0)
			require.NoError(t, err)

			cmd, err := ReadCommand()
			require.NoError(t, err)
			assert.Equal(t, tt.expected, *cmd)
		})
	}
}

func mockStd() (rollback func(), err error) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "mockstd")
	if err != nil {
		return nil, fmt.Errorf("cannot create temp dir for mocking std: %w", err)
	}

	mockIn, err := os.CreateTemp(tempDir, "mockIn")
	if err != nil {
		return nil, fmt.Errorf("cannot create mocked stdin: %w", err)
	}

	mockOut, err := os.CreateTemp(tempDir, "mockOut")
	if err != nil {
		return nil, fmt.Errorf("cannot create mocked out: %w", err)
	}

	defaultIn, defaultOut := os.Stdin, os.Stdout
	rollback = func() {
		os.Stdin = defaultIn
		os.Stdout = defaultOut
	}

	os.Stdin = mockIn
	os.Stdout = mockOut
	return rollback, nil
}

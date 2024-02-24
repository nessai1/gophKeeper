package keeper

import (
	"bufio"
	"fmt"
	"github.com/nessai1/gophkeeper/internal/keeper/command"
	"os"
	"time"
)

func Listen() error {
	printGreetMessage(applicationInfo{Version: "0.0.1", BuildDate: time.Now()})

	reader := bufio.NewReader(os.Stdin)
	writer := os.Stdout
	for {
		fmt.Println()
		cmd, err := command.ReadCommand(reader, writer)
		if err != nil {
			return fmt.Errorf("error while listen command by keeper: %w", err)
		}

		if cmd.Name == "exit" {
			return nil
		}

		fmt.Printf("Command: %s\t", cmd.Name)
	}
}

type applicationInfo struct {
	Version   string
	BuildDate time.Time
}

const greetMsg = `
  _  __                                    
 | |/ /   ___    ___   _ __     ___   _ __ 
 | ' /   / _ \  / _ \ | '_ \   / _ \ | '__|
 | . \  |  __/ |  __/ | |_) | |  __/ | |   
 |_|\_\  \___|  \___| | .__/   \___| |_|   
                      |_|
`

func printGreetMessage(info applicationInfo) {
	fmt.Println("\033[34m", greetMsg, "\033[0m")
	fmt.Printf("Welcome to the Keeper!\n\n")
	fmt.Printf("Version: v%s\n", info.Version)
	fmt.Printf("Build date: %s\n", info.BuildDate.String())
}

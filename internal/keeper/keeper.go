package keeper

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/nessai1/gophkeeper/internal/keeper/connector"
	"github.com/nessai1/gophkeeper/internal/keeper/session"
	"github.com/nessai1/gophkeeper/internal/logger"
	"go.uber.org/zap"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/nessai1/gophkeeper/pkg/command"
)

func Listen() error {
	writer := os.Stdout

	printGreetMessage(writer, applicationInfo{Version: "0.0.1", BuildDate: time.Now()})

	cfg, err := fetchConfig()
	if err != nil {
		return fmt.Errorf("cannot start listen application: %w", err)
	}

	app, err := NewApplication(writer, os.Stdin, cfg)
	if err != nil {
		return fmt.Errorf("cannot start listen application: %w", err)
	}

	if err = app.Run(); err != nil {
		return fmt.Errorf("error while run application: %w", err)
	}

	return nil
}

type ServiceConnector interface {
	Ping(ctx context.Context) (answer string, error error)
}

type Application struct {
	config Config

	connector ServiceConnector

	input  io.Reader
	output io.Writer

	session *session.Session

	logger *zap.Logger
}

type WorkDir struct {
}

func NewApplication(input io.Reader, output io.Writer, config Config) (*Application, error) {

	err := createKeeperDataDir(config.WorkDir)
	if err != nil {
		return nil, fmt.Errorf("cannot create keeper data dir: %w", err)
	}

	logDir := filepath.Join(config.WorkDir, "logs")
	logFile, err := logger.OpenLogFile(logDir)
	if err != nil {
		return nil, fmt.Errorf("error while create log file to dir (%s) for new application: %w", logDir, err)
	}

	loggerInstance, err := logger.BuildLogger(config.Mode, logFile)
	if err != nil {
		stat, _ := logFile.Stat()
		return nil, fmt.Errorf("error while build logger file %s with mode %s: %w", stat.Name(), config.Mode, err)
	}

	gRPCConnector, err := connector.CreateGRPCConnector(config.ServerAddr)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to the external service: %w", err)
	}

	return &Application{
		input:     input,
		output:    output,
		config:    config,
		logger:    loggerInstance,
		connector: gRPCConnector,
	}, nil
}

func createKeeperDataDir(dir string) error {
	err := os.Mkdir(dir, 0777)
	if err != nil && !errors.As(err, &os.ErrExist) {
		return fmt.Errorf("cannot credate data dir: %w", err)
	}

	logsPath := filepath.Join(dir, "logs")

	err = os.Mkdir(logsPath, 0777)
	if err != nil && !errors.As(err, &os.ErrExist) {
		return fmt.Errorf("cannot create logs dir: %w", err)
	}

	return nil
}

func (a *Application) Run() error {
	reader := bufio.NewReader(os.Stdin)

	currentSession := loadCurrentSession()
	if currentSession != nil {
		fmt.Fprintf(a.output, "\033[32mYou are authorized as '%s'!\033[0m", currentSession.Login)
	} else {
		if a.config.ServerAddr == "" {
			return fmt.Errorf("cannot start application without active session then server addr is empty")
		}
		fmt.Fprintf(a.output, "\033[33mYou are not authorized. Use 'login' / 'register' to authorize in service\033[0m")
	}

	a.logger.Info("Application was started", zap.Bool("with_session", currentSession != nil))

	for {
		_, err := a.output.Write([]byte("\n"))
		if err != nil {
			return err
		}

		cmd, err := command.ReadCommand(reader, a.output)
		if err != nil {
			return fmt.Errorf("error while listen command by keeper: %w", err)
		}

		if cmd.Name == "exit" {
			return nil
		}

		if cmd.Name == "ping" {
			answer, err := a.connector.Ping(context.TODO())
			if err != nil {
				a.output.Write([]byte("error while ping; see logs"))
				a.logger.Error("error while ping service", zap.Error(err))
			} else {
				a.output.Write([]byte(fmt.Sprintf("Answer: %s", answer)))
			}
		}
	}
}

func loadCurrentSession() *session.Session {
	return nil
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

func printGreetMessage(target io.Writer, info applicationInfo) {
	fmt.Fprintf(target, "\033[34m"+greetMsg+"\033[0m")
	fmt.Fprintf(target, "Welcome to the Keeper!\n\n")
	fmt.Fprintf(target, "Version: v%s\n", info.Version)
	fmt.Fprintf(target, "Build date: %s\n", info.BuildDate.String())
}

package keeper

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/nessai1/gophkeeper/internal/keeper/connector"
	"github.com/nessai1/gophkeeper/internal/keeper/performer"
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

type Application struct {
	config Config

	connector connector.ServiceConnector

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

	var (
		found       bool
		requireExit bool
	)
	for {
		_, err := a.output.Write([]byte("\n"))
		if err != nil {
			return err
		}

		cmd, err := command.ReadCommand(reader, a.output)
		if err != nil {
			return fmt.Errorf("error while listen command by keeper: %w", err)
		}
		if cmd.Name == "" {
			continue
		}

		found = false
		err = nil
		for _, p := range performer.AvailablePerformers {
			if p.GetName() == cmd.Name {
				found = true

				requireExit, err = p.Execute(
					reader,
					a.output,
					a.connector,
					a,
					a.logger,
					cmd.Args,
				)
			}
		}
		if !found {
			a.output.Write([]byte(fmt.Sprintf("\033[31mCommand '%s' not found!\033[0m\n", cmd.Name)))
		}

		if err != nil {
			a.output.Write([]byte(fmt.Sprintf("\033[31mError: %s\033[0m\n", err.Error())))
		}

		if requireExit {
			if err == nil {
				a.output.Write([]byte("Bye!\n"))

				return nil
			} else {
				a.output.Write([]byte("\033[31mThe application was interrupted by an error"))

				return err
			}
		}
	}
}

func (a *Application) getPerformers() {
	if a.session != nil {

	} else {

	}
}

func (a *Application) SetSession(s *session.Session) {
	a.session = s
}

func (a *Application) GetSession() *session.Session {
	return a.session
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

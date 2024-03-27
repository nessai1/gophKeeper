package keeper

import (
	"errors"
	"fmt"
	"github.com/nessai1/gophkeeper/internal/keeper/connector"
	"github.com/nessai1/gophkeeper/internal/keeper/performer"
	"github.com/nessai1/gophkeeper/internal/keeper/session"
	"github.com/nessai1/gophkeeper/internal/logger"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"time"

	"github.com/nessai1/gophkeeper/pkg/command"
)

func Listen() error {
	printGreetMessage(applicationInfo{Version: "0.0.1", BuildDate: time.Now()})

	cfg, err := fetchConfig()
	if err != nil {
		return fmt.Errorf("cannot start listen application: %w", err)
	}

	app, err := NewApplication(cfg)
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

	session *session.Session

	logger *zap.Logger
}

type WorkDir struct {
}

func NewApplication(config Config) (*Application, error) {

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

	gRPCConnector, err := connector.CreateGRPCConnector(config.ServerAddr, config.Certificate)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to the external service: %w", err)
	}

	return &Application{
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

	mediaPath := filepath.Join(dir, "media")

	err = os.Mkdir(mediaPath, 0777)
	if err != nil && !errors.As(err, &os.ErrExist) {
		return fmt.Errorf("cannot create media dir: %w", err)
	}

	return nil
}

func (a *Application) Run() error {
	currentSession, err := session.LoadLocalSession(a.config.WorkDir)
	if err != nil {
		a.logger.Error("cannot load local session for user", zap.Error(err))
	}

	if currentSession != nil {
		a.SetSession(currentSession)
		fmt.Printf("\033[32mYou are authorized as '%s'!\033[0m", currentSession.Login)
	} else {
		if a.config.ServerAddr == "" {
			return fmt.Errorf("cannot start application without active session then server addr is empty")
		}
		fmt.Printf("\033[33mYou are not authorized. Use 'login' / 'register' to authorize in service\033[0m")
	}
	fmt.Println()

	a.logger.Info("Application was started", zap.Bool("with_session", currentSession != nil))

	var requireExit bool
	for {
		cmd, err := command.ReadCommand()
		if err != nil {
			return fmt.Errorf("error while listen command by keeper: %w", err)
		}
		if cmd.Name == "" {
			continue
		}
		err = nil

		p, ok := performer.AvailablePerformers[cmd.Name]
		if !ok {
			fmt.Printf("\033[31mCommand '%s' not found!\033[0m\n", cmd.Name)
			continue
		}

		requireExit, err = p.Execute(
			a.connector,
			a,
			a.logger,
			cmd.Args,
			a.config.WorkDir,
		)

		if err != nil {
			fmt.Printf("\033[31mError: %s\033[0m\n", err.Error())
		}

		if requireExit {
			if a.GetSession() != nil {
				saveErr := session.SaveLocalSession(a.config.WorkDir, *a.GetSession())
				if saveErr != nil {
					a.logger.Error("Error while save local session", zap.Error(saveErr))
				}
			}
			if err == nil {
				fmt.Printf("Bye!\n")

				return nil
			} else {
				fmt.Printf("\033[31mThe application was interrupted by an error")

				return err
			}
		}
	}
}

func (a *Application) SetSession(s *session.Session) {
	a.session = s

	if s != nil {
		a.connector.SetAuthToken(s.AuthToken)
	} else {
		a.connector.SetAuthToken("")
	}
}

func (a *Application) GetSession() *session.Session {
	return a.session
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
	fmt.Printf("\033[34m" + greetMsg + "\033[0m")
	fmt.Printf("Welcome to the Keeper!\n\n")
	fmt.Printf("Version: v%s\n", info.Version)
	fmt.Printf("Build date: %s\n", info.BuildDate.String())
}

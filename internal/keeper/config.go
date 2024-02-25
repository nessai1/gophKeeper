package keeper

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/nessai1/gophkeeper/internal/logger"
	"os"
)

var ErrConfigFileNotFound = errors.New("config file not found")

type Config struct {
	WorkDir    string                  `json:"work_dir"`
	Mode       logger.ApplicationLevel `json:"mode"`
	ServerAddr string                  `json:"server"`
}

func fetchConfig() (Config, error) {
	configPath := flag.String("cfg", "./keeper_config.json", "Path to JSON config of application")

	workDir := flag.String("d", "./keeperData", "Path to keeper data, like log files, user content etc.")
	isDevMode := flag.Bool("dev", false, "Enable dev mode of application for complicated logs")
	serverAddr := flag.String("s", "", "Address of keeper server")

	flag.Parse()

	mode := logger.LevelProduction
	if *isDevMode {
		mode = logger.LevelDev
	}

	cfg := Config{
		WorkDir:    *workDir,
		Mode:       mode,
		ServerAddr: *serverAddr,
	}

	fileCfg, err := readFileConfig(*configPath)
	if err != nil && !errors.Is(ErrConfigFileNotFound, err) {
		return Config{}, fmt.Errorf("cannot fetch file config: %w", err)
	}

	if cfg.WorkDir == "" {
		cfg.WorkDir = fileCfg.WorkDir
	}

	if cfg.ServerAddr == "" {
		cfg.ServerAddr = fileCfg.ServerAddr
	}

	if fileCfg.Mode != "" {
		cfg.Mode = fileCfg.Mode
	}

	return cfg, nil
}

func readFileConfig(configPath string) (Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return Config{}, ErrConfigFileNotFound
	}

	buffer := bytes.Buffer{}
	n, err := buffer.ReadFrom(file)
	if n == 0 {
		return Config{}, fmt.Errorf("config file %s was empty", configPath)
	}

	if err != nil {
		return Config{}, fmt.Errorf("error while read config file %s: %w", configPath, err)
	}

	var cfg Config
	err = json.Unmarshal(buffer.Bytes(), &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("cannot unmarshal config file %s: %w", configPath, err)
	}

	return cfg, nil
}

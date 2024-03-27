package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"sync"
)

type Config struct {
	Address     string `json:"service_address"`
	SecretToken string `json:"secret_token"`
	Salt        string `json:"salt"`

	TLSCredentials *TLSCredentials `json:"tls_credentials"`

	PlainStorageConfig *PlainStorageConfig `json:"plain_storage"`

	S3Config *S3Config `json:"s3"`

	FileConfigPath string
}

type TLSCredentials struct {
	// Path to server crt file
	Crt string `json:"crt"`

	// Path to server key file
	Key string `json:"key"`
}

type S3Config struct {
	URL           string `json:"url"`
	PartitionID   string `json:"partition_id"`
	SigningRegion string `json:"signing_region"`

	Credentials *S3Credentials `json:"credentials"`
}

type PlainStorageConfig struct {
	PSQLStorage *PSQLPlainStorageConfig `json:"postgres"`
}

type PSQLPlainStorageConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

type S3Credentials struct {
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
}

var (
	once   sync.Once
	cfg    Config
	cfgErr error
)

func FetchConfig() (Config, error) {
	once.Do(func() {
		cfg, cfgErr = loadConfig()
	})

	return cfg, cfgErr
}

func loadConfig() (Config, error) {
	flagConfig := fetchFlagConfig()

	fileConfig, err := fetchFileConfig(flagConfig.FileConfigPath)
	if err != nil {
		return Config{}, fmt.Errorf("cannot fetch file config: %w", err)
	}

	config, err := mergeConfigs(flagConfig, fileConfig)
	if err != nil {
		return Config{}, fmt.Errorf("cannot merge fetched configs: %w", err)
	}

	return config, nil
}

func fetchFileConfig(filePath string) (Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return Config{}, fmt.Errorf("cannot open config file: %w", err)
	}

	b := bytes.Buffer{}
	n, err := b.ReadFrom(file)
	if err != nil {
		return Config{}, fmt.Errorf("cannot read config file: %w", err)
	}

	if n == 0 {
		return Config{}, errors.New("config file can't be empty")
	}

	config := Config{}

	err = json.Unmarshal(b.Bytes(), &config)
	if err != nil {
		return Config{}, fmt.Errorf("cannot unmarshal file config: %w", err)
	}

	return config, err
}

func fetchFlagConfig() Config {
	address := flag.String("s", "", "Address of keeper service")
	fileConfig := flag.String("c", "service_config.json", "Path to service configuration file")
	flag.Parse()

	return Config{
		Address:        *address,
		FileConfigPath: *fileConfig,
	}
}

func mergeConfigs(flagConfig, fileConfig Config) (Config, error) {
	fileConfig.FileConfigPath = flagConfig.FileConfigPath

	if flagConfig.Address != "" {
		fileConfig.Address = flagConfig.Address
	}

	if fileConfig.Address == "" {
		return Config{}, errors.New("any one config must have service address")
	}

	if fileConfig.SecretToken == "" {
		return Config{}, errors.New("any one config must have secret token")
	}

	if fileConfig.S3Config == nil {
		return Config{}, errors.New("service must have storage config (S3)")
	}

	return fileConfig, nil
}

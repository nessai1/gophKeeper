package logger

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ApplicationLevel string

var ErrUndefinedAppLevel = errors.New("unexpected applicationLevel")

const (
	LevelDev        ApplicationLevel = "dev"
	LevelProduction ApplicationLevel = "prod"
)

type Target interface {
	io.Writer
	Sync() error
}

func BuildLogger(level ApplicationLevel, target Target) (*zap.Logger, error) {
	atom := zap.NewAtomicLevel()
	logLevel, err := translateApplicationLevelToLogLevel(level)
	if err != nil {
		return nil, fmt.Errorf("cannot build logger: %w", err)
	}

	atom.SetLevel(logLevel)
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.RFC3339TimeEncoder

	return zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.Lock(target),
		atom,
	)), nil
}

func OpenLogFile(logDir string) (*os.File, error) {
	filePath := filepath.Join(logDir, time.Now().Format("2006-01-02")+".log")
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)

	if err != nil {
		return nil, fmt.Errorf("cannot open log file %s: %w", filePath, err)
	}

	return file, nil
}

func translateApplicationLevelToLogLevel(level ApplicationLevel) (zapcore.Level, error) {
	switch level {
	case LevelProduction:
		return zapcore.ErrorLevel, nil
	case LevelDev:
		return zapcore.DebugLevel, nil
	default:
		return 0, fmt.Errorf("cannot translate appLevel (%d) to logLevel: %w", level, ErrUndefinedAppLevel)
	}
}

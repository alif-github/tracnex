package util

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func ConfigZap(level zapcore.Level, output []string) (logger *zap.Logger) {
	cfg := zap.Config{
		Encoding:    "json",
		OutputPaths: output,
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:     "timestamp",
			EncodeTime:  zapcore.RFC3339NanoTimeEncoder,
			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,
		},
	}
	cfg.Level = zap.NewAtomicLevelAt(level)

	logger, err := cfg.Build()
	if err != nil {
		os.Exit(9)
	}
	return
}

type logger struct {
	LoggerDebugLevel *zap.Logger
}

var Logging logger

func SetLogger(logFileName []string) {
	Logging.LoggerDebugLevel = ConfigZap(zap.DebugLevel, logFileName)
}

func LogInfo(data []zap.Field) {
	Logging.LoggerDebugLevel.Info("", data...)
}

func LogError(data []zap.Field) {
	Logging.LoggerDebugLevel.Error("", data...)
}

package logger

import (
	"fmt"

	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger = zap.NewNop()

func Init(cl config.Logging) error {
	lvl, err := zap.ParseAtomicLevel(cl.Level)
	if err != nil {
		return fmt.Errorf("could not parse logging level: %w", err)
	}

	cfg := zap.NewProductionConfig()
	if cl.OutputToFile {
		cfg.OutputPaths = []string{"stderr", cl.OutputFilePath}
	}
	if cl.ErrorOutputToFile {
		cfg.ErrorOutputPaths = []string{"stderr", cl.ErrorOutputFilePath}
	}

	cfg.Level = lvl
	cfg.EncoderConfig = zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "ts",
		FunctionKey:    zapcore.OmitKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000Z"), // Формат ISO 8601
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
	}

	//cfg.Encoding = "json"

	zl, err := cfg.Build()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}
	// устанавливаем синглтон
	Log = zl
	return nil
}

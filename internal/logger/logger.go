package logger

import (
	"log/slog"
	"os"
)

func GetLogger() *slog.Logger {
	handlerOpts := &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stderr, handlerOpts))
	slog.SetDefault(logger)

	return logger
}

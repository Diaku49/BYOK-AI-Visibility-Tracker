package logger

import (
	"log/slog"
	"os"
)

func NewLogger() *slog.Logger {
	var handler slog.Handler
	opts := &slog.HandlerOptions{
		AddSource: false, // Set to true in dev if you want file/line numbers
		Level:     slog.LevelInfo,
	}

	if os.Getenv("ENV") == "production" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger) // Optional

	return logger
}

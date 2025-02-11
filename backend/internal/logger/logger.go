// Package logger provides constructor for log/slog package that defines config based on app environment.
package logger

import (
	"log/slog"
	envpkg "markup/internal/domain/enums/env"
	"os"
)

// New creates new slog instance based on app environment.
func New(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envpkg.Local:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case envpkg.Dev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case envpkg.Prod:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}

	return log
}

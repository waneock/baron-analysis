package logger

import (
	"log/slog"
	"os"
)

type Env string

const (
	Local Env = "local"
	Dev   Env = "dev"
	Prod  Env = "Prod"
)

type Config struct {
	Env       Env
	Level     slog.Level
	AddSource bool
}

func New(cfg Config) *slog.Logger {
	options := &slog.HandlerOptions{
		Level:     cfg.Level,
		AddSource: cfg.AddSource,
	}

	var handler slog.Handler

	switch cfg.Env {
	case Prod:
		handler = slog.NewJSONHandler(os.Stdout, options)
	default:
		handler = slog.NewTextHandler(os.Stdout, options)
	}

	return slog.New(handler)
}

func MustLoad(env string) *slog.Logger {
	switch Env(env) {
	case Prod:
		return New(Config{
			Env:       Prod,
			Level:     slog.LevelInfo,
			AddSource: false,
		})

	case Dev:
		return New(Config{
			Env:       Dev,
			Level:     slog.LevelDebug,
			AddSource: true,
		})

	default:
		return New(Config{
			Env:       Local,
			Level:     slog.LevelDebug,
			AddSource: true,
		})
	}
}

package logger

import (
	"errors"
	"log/slog"
	"os"
)

var InvalidEnvErr = errors.New("invalid environment variable")

func SetupLogger(env string) (*slog.Logger, error) {
	switch env {
	case "dev":
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})), nil
	case "prod":
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})), nil
	default:
		return nil, InvalidEnvErr
	}
}

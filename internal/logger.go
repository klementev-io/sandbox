package internal

import (
	"errors"
	"log/slog"
	"os"
)

func SetupLogger(level string, format string) error {
	var lvl slog.Level
	if err := lvl.UnmarshalText([]byte(level)); err != nil {
		return err
	}

	opts := &slog.HandlerOptions{
		Level: lvl,
	}

	var h slog.Handler
	switch format {
	case "json":
		h = slog.NewJSONHandler(os.Stdout, opts)
	case "text":
		h = slog.NewTextHandler(os.Stdout, opts)
	default:
		return errors.New("invalid logger format: " + format)
	}

	slog.SetDefault(slog.New(h))

	return nil
}

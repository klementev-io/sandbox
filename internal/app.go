package internal

import (
	"context"
	"fmt"
	"log/slog"
)

func Run(cfg Config) error {
	if err := SetupLogger(cfg.Logger.Level, cfg.Logger.Format); err != nil {
		return fmt.Errorf("failed to setup logger: %w", err)
	}

	slog.Default().InfoContext(context.Background(), "starting application")

	return nil
}

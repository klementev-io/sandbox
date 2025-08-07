package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/klementev-io/sandbox/internal"
	"github.com/klementev-io/sandbox/internal/config"
)

func run() int {
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Default().ErrorContext(context.Background(), "could not load config", "error", err)
		return 1
	}
	err = internal.Run(cfg)
	if err != nil {
		slog.Default().ErrorContext(context.Background(), "failed to run application", "error", err)
		return 1
	}
	return 0
}

func main() {
	os.Exit(run())
}

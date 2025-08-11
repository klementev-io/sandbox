package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/klementev-io/sandbox/internal"
	"github.com/klementev-io/sandbox/internal/config"
)

func run() int {
	cfg, err := internal.LoadConfig[config.Cfg]()
	if err != nil {
		slog.Default().ErrorContext(context.Background(), "could not load config", "error", err)
		return 1
	}

	if err = internal.Run(cfg); err != nil {
		slog.Default().ErrorContext(context.Background(), "failed to run application", "error", err)
		return 1
	}
	return 0
}

func main() {
	os.Exit(run())
}

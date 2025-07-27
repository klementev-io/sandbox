package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/klementev-io/sandbox/internal"
)

func main() {
	cfg, err := internal.LoadConfig()
	if err != nil {
		slog.Default().ErrorContext(context.Background(), "could not load config", "error", err)
		os.Exit(1)
	}

	err = internal.Run(cfg)
	if err != nil {
		slog.Default().ErrorContext(context.Background(), "failed to run application", "error", err)
		os.Exit(1)
	}
}

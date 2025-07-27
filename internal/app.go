package internal

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"golang.org/x/sync/errgroup"
)

const (
	readHeaderTimeout = 5 * time.Second
	interruptTimeout  = 2 * time.Second
)

func Run(cfg Config) error {
	if err := SetupLogger(cfg.Logger.Level, cfg.Logger.Format); err != nil {
		return fmt.Errorf("failed to setup logger: %w", err)
	}

	router := http.NewServeMux()

	router.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	httpSrv := http.Server{
		Addr:              net.JoinHostPort(cfg.HTTPServer.Host, cfg.HTTPServer.Port),
		Handler:           router,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	eg, egCtx := errgroup.WithContext(context.Background())
	ctx, stop := signal.NotifyContext(egCtx, os.Interrupt)
	defer stop()

	eg.Go(func() error {
		slog.Default().InfoContext(ctx, "starting http server", "address", httpSrv.Addr)
		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("failed to start http server: %w", err)
		}
		slog.InfoContext(ctx, "http server shutdown gracefully")
		return nil
	})

	eg.Go(func() error {
		<-ctx.Done()
		timeoutCtx, cancel := context.WithTimeout(context.Background(), interruptTimeout)
		defer cancel()
		err := httpSrv.Shutdown(timeoutCtx)
		if err != nil {
			return err
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("server exited with error: %w", err)
	}

	return nil
}

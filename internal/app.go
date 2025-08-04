package internal

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	_ "net/http/pprof" //nolint:gosec // pprof port is not exposed to the internet
	"os"
	"os/signal"
	"time"

	"golang.org/x/sync/errgroup"

	v1 "github.com/klementev-io/sandbox/internal/api/v1"
)

const (
	readHeaderTimeout = 5 * time.Second
	interruptTimeout  = 2 * time.Second

	pprofHost = "127.0.0.1"
)

func Run(cfg Config) error {
	if err := SetupLogger(cfg.Logger.Level, cfg.Logger.Format); err != nil {
		return fmt.Errorf("failed to setup logger: %w", err)
	}

	eg, ctx := errgroup.WithContext(context.Background())
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	startHTTPServer(ctx, eg, cfg.HTTPServer)

	if cfg.HTTPServer.PprofPort != "" {
		startPprofServer(ctx, eg, cfg.HTTPServer.PprofPort)
	}

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("server exited with error: %w", err)
	}

	return nil
}

func startHTTPServer(ctx context.Context, eg *errgroup.Group, cfg HTTPServerConfig) {
	router := v1.SetupRouter()

	httpSrv := http.Server{
		Addr:              net.JoinHostPort(cfg.Host, cfg.Port),
		Handler:           router,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	eg.Go(func() error {
		slog.InfoContext(ctx, "starting http server", "address", httpSrv.Addr)
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
}

func startPprofServer(ctx context.Context, eg *errgroup.Group, port string) {
	pprofSrv := http.Server{
		Addr:              net.JoinHostPort(pprofHost, port),
		Handler:           http.DefaultServeMux,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	eg.Go(func() error {
		slog.InfoContext(ctx, "starting pprof server", "address", pprofSrv.Addr)
		if err := pprofSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("failed to start pprof server: %w", err)
		}
		slog.InfoContext(ctx, "pprof server shutdown gracefully")
		return nil
	})

	eg.Go(func() error {
		<-ctx.Done()
		timeoutCtx, cancel := context.WithTimeout(context.Background(), interruptTimeout)
		defer cancel()
		err := pprofSrv.Shutdown(timeoutCtx)
		if err != nil {
			return err
		}
		return nil
	})
}

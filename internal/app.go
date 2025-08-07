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

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/klementev-io/sandbox/internal/config"
	"golang.org/x/sync/errgroup"

	"github.com/klementev-io/sandbox/internal/api/middleware"
	v1 "github.com/klementev-io/sandbox/internal/api/v1"
)

const (
	readHeaderTimeout = 5 * time.Second
	interruptTimeout  = 2 * time.Second
)

func Run(cfg config.Cfg) error {
	if err := setupLogger(cfg.Logger.Level, cfg.Service.Name, cfg.Service.Version); err != nil {
		return fmt.Errorf("failed to setup logger: %w", err)
	}

	eg, ctx := errgroup.WithContext(context.Background())
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	startAPIServer(ctx, eg, cfg.APIServer)

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("server exited with error: %w", err)
	}

	return nil
}

func startAPIServer(ctx context.Context, eg *errgroup.Group, cfg config.APIServer) {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.Use(
		middleware.Recovery(),
		middleware.GinLogger(slog.Default()),
	)

	v1.RegisterHandlers(router, &v1.Handlers{})

	if cfg.Pprof {
		pprof.Register(router)
	}

	httpSrv := http.Server{
		Addr:              net.JoinHostPort(cfg.Host, cfg.Port),
		Handler:           router,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	eg.Go(func() error {
		slog.InfoContext(ctx, "starting api server", "address", httpSrv.Addr)
		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("failed to start http server: %w", err)
		}
		slog.InfoContext(ctx, "api server shutdown gracefully")
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

func setupLogger(level string, service string, version string) error {
	var lvl slog.Level
	if err := lvl.UnmarshalText([]byte(level)); err != nil {
		return err
	}

	handler := slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: lvl,
		},
	)

	l := slog.New(handler)
	l = l.With(
		slog.Group("service",
			slog.String("name", service),
			slog.String("version", version),
		),
	)

	slog.SetDefault(l)

	return nil
}

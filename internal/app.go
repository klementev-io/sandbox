package internal

import (
	"context"
	"fmt"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"

	"github.com/klementev-io/sandbox/internal/api/http"
	"github.com/klementev-io/sandbox/internal/api/http/middleware"
	v1 "github.com/klementev-io/sandbox/internal/api/http/v1"
	"github.com/klementev-io/sandbox/internal/config"
)

func Run(cfg *config.Cfg) error {
	if err := SetupLogger(cfg.Log.Level, LogWithService(cfg.Service)); err != nil {
		return fmt.Errorf("could not setup logger: %w", err)
	}

	slog.Default().InfoContext(context.Background(), "starting service")

	eg, egCtx := errgroup.WithContext(context.Background())
	ctx, stop := signal.NotifyContext(egCtx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer stop()

	startAPIServer(ctx, eg, cfg.APIServer)

	if cfg.PprofServer.Enable {
		startPprofServer(ctx, eg, cfg.PprofServer)
	}

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("service exited with error: %w", err)
	}

	slog.Default().InfoContext(context.Background(), "service successfully completed")

	return nil
}

func startAPIServer(ctx context.Context, eg *errgroup.Group, cfg config.APIServer) {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.Use(
		middleware.Recovery(),
		middleware.GinLogger(),
	)

	v1.RegisterHandlers(router, v1.NewHandlers())

	apiSrv := http.NewServer(
		cfg.Host,
		cfg.Port,
		router,
		slog.Default().With("server", "api"),
	)

	eg.Go(func() error {
		return apiSrv.Start(ctx)
	})

	eg.Go(func() error {
		<-ctx.Done()
		return apiSrv.Shutdown()
	})
}

func startPprofServer(ctx context.Context, eg *errgroup.Group, cfg config.PprofServer) {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	pprof.Register(router)

	pprofSrv := http.NewServer(
		cfg.Host,
		cfg.Port,
		router,
		slog.Default().With("server", "pprof"),
	)

	eg.Go(func() error {
		return pprofSrv.Start(ctx)
	})

	eg.Go(func() error {
		<-ctx.Done()
		return pprofSrv.Shutdown()
	})
}

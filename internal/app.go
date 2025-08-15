package internal

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"

	v1 "github.com/klementev-io/sandbox/api/gen/v1"

	handlersv1 "github.com/klementev-io/sandbox/internal/api/handlers/v1"
	"github.com/klementev-io/sandbox/internal/api/middleware"
	"github.com/klementev-io/sandbox/internal/config"
	"github.com/klementev-io/sandbox/internal/httpserver"
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

	h := handlersv1.NewHandlers()

	router := gin.New()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1.RegisterHandlersWithOptions(router, h, v1.GinServerOptions{
		BaseURL: "/api/v1",
		Middlewares: []v1.MiddlewareFunc{
			middleware.Recovery(),
			middleware.GinLogger(),
		},
	})

	apiSrv := httpserver.New("api", cfg.Host, cfg.Port, router)

	eg.Go(func() error {
		return apiSrv.Start(ctx)
	})

	eg.Go(func() error {
		<-ctx.Done()
		return apiSrv.Shutdown()
	})
}

func startPprofServer(ctx context.Context, eg *errgroup.Group, cfg config.PprofServer) {
	router := gin.New()

	pprof.Register(router)

	pprofSrv := httpserver.New("pprof", cfg.Host, cfg.Port, router)

	eg.Go(func() error {
		return pprofSrv.Start(ctx)
	})

	eg.Go(func() error {
		<-ctx.Done()
		return pprofSrv.Shutdown()
	})
}

package httpserver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"
)

const (
	readHeaderTimeout = 5 * time.Second
	interruptTimeout  = 2 * time.Second
)

type Srv struct {
	server *http.Server
	log    *slog.Logger
}

func New(name, host, port string, handler http.Handler) *Srv {
	return &Srv{
		server: &http.Server{
			Addr:        net.JoinHostPort(host, port),
			Handler:     handler,
			ReadTimeout: readHeaderTimeout,
		},
		log: slog.Default().With(slog.String("server", name)),
	}
}

func (s *Srv) Start(ctx context.Context) error {
	s.log.InfoContext(ctx, "starting http server", "address", s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to start http server: %w", err)
	}

	s.log.InfoContext(ctx, "http server shutdown gracefully")
	return nil
}

func (s *Srv) Shutdown() error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), interruptTimeout)
	defer cancel()

	err := s.server.Shutdown(timeoutCtx)
	if err != nil {
		return fmt.Errorf("failed to shutdown http server: %w", err)
	}

	return nil
}

package internal

import (
	"log/slog"
	"os"
)

func SetupLogger(level string, opt ...func(l *slog.Logger) *slog.Logger) error {
	var lvl slog.Level
	if err := lvl.UnmarshalText([]byte(level)); err != nil {
		return err
	}

	l := slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: lvl,
			},
		),
	)

	if len(opt) > 0 {
		for _, o := range opt {
			l = o(l)
		}
	}

	slog.SetDefault(l)

	return nil
}

func LogWithService(name string) func(l *slog.Logger) *slog.Logger {
	return func(l *slog.Logger) *slog.Logger {
		return l.WithGroup("service").With(slog.String("name", name))
	}
}

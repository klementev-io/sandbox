package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		status := c.Writer.Status()
		method := c.Request.Method
		errors := c.Errors.String()

		attrs := []slog.Attr{
			slog.Int("status", status),
			slog.String("method", method),
			slog.String("path", path),
			slog.String("query", query),
		}

		if errors != "" {
			attrs = append(attrs, slog.String("errors", errors))
		}

		slog.Default().LogAttrs(c.Request.Context(), getLevel(status), "request completed", attrs...)
	}
}

func getLevel(status int) slog.Level {
	switch {
	case status >= http.StatusInternalServerError:
		return slog.LevelError
	case status >= http.StatusBadRequest:
		return slog.LevelWarn
	default:
		return slog.LevelDebug
	}
}

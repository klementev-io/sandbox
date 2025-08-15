package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Recovery() func(ctx *gin.Context) {
	return gin.CustomRecoveryWithWriter(nil, func(c *gin.Context, err any) {
		slog.Default().ErrorContext(c, "panic recovered", slog.Any("error", err))
		c.AbortWithStatus(http.StatusInternalServerError)
	})
}

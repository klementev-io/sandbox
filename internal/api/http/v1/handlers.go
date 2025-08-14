package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handlers struct{}

func NewHandlers() *Handlers {
	return &Handlers{}
}

func (h Handlers) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

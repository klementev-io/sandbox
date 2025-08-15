package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	v1 "github.com/klementev-io/sandbox/api/gen/v1"
	"github.com/oapi-codegen/runtime/types"
)

var _ v1.ServerInterface = (*Handlers)(nil)

type Handlers struct {
}

func NewHandlers() *Handlers {
	return &Handlers{}
}

func (h *Handlers) PostOrders(c *gin.Context) {
	c.Status(http.StatusOK)
}

func (h *Handlers) PostUsers(c *gin.Context) {
	c.Status(http.StatusOK)
}

func (h *Handlers) GetUsersID(c *gin.Context, _ types.UUID) {
	c.Status(http.StatusOK)
}

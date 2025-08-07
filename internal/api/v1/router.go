package v1

import (
	"github.com/gin-gonic/gin"
)

func RegisterHandlers(router *gin.Engine, handlers *Handlers) {
	api := router.Group("api")
	{
		v1 := api.Group("v1")
		{
			v1.GET("health", handlers.Health)
		}
	}
}

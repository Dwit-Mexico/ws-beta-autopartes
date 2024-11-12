package http

import (
	"github.com/RomanshkVolkov/test-api/internal/adapters/handler"
	"github.com/gin-gonic/gin"
)

func InternalRoutes(r *gin.Engine) {
	internal := r.Group("/internal")
	{
		internal.PUT("/permissions-synchronization", handler.InternalSynchronization)
	}
}

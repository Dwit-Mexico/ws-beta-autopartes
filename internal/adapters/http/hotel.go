package http

import (
	"github.com/RomanshkVolkov/test-api/internal/adapters/handler"
	"github.com/RomanshkVolkov/test-api/internal/adapters/middleware"
	"github.com/gin-gonic/gin"
)

func HostingCenterRoutes(r *gin.Engine) {
	protect := middleware.Protected
	hostingCenter := r.Group("/hosting-center")
	{
		hostingCenter.GET("/current", protect(), handler.GetCurrentHostingCenter)
		hostingCenter.PUT("/:id", protect(), handler.UpdateHostingCenter)

	}
}

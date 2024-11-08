package http

import (
	"github.com/RomanshkVolkov/test-api/internal/adapters/handler"
	"github.com/RomanshkVolkov/test-api/internal/adapters/middleware"
	"github.com/gin-gonic/gin"
)

func CatalogRoutes(r *gin.Engine) {
	protect := middleware.Protected
	catalog := r.Group("/catalogs")
	{
		catalog.GET("/kitchen/:id", protect(), handler.GetKitchenByID)
		catalog.GET("/shift/:id", protect(), handler.GetShiftByID)

		catalog.PUT("/:slug/:id", protect(), handler.UpdateGenericCatalog)
		catalog.DELETE("/:slug/:id", protect(), handler.DeleteGenericCatalog)
	}
}

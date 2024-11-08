package http

import (
	"github.com/RomanshkVolkov/test-api/internal/adapters/handler"
	"github.com/RomanshkVolkov/test-api/internal/adapters/middleware"
	"github.com/gin-gonic/gin"
)

func DocumentRoutes(r *gin.Engine) {
	protect := middleware.Protected
	document := r.Group("/documents")
	{
		document.GET("/", protect(), handler.GetDocuments)
		document.GET("/:id", protect(), handler.GetDocumentByID)
		document.GET("/tables", protect(), handler.GetTables)
		document.GET("/tables/:id", protect(), handler.GetTableByID)
		document.POST("upload", protect(), handler.UploadDocument)
		document.POST("", handler.CreateDocument)
		// document.PUT("/update", UpdateDocument)
		// document.DELETE("/delete", DeleteDocument)
	}
}

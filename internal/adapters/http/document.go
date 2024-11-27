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
		document.GET("/details/:document/records/:id", protect(), handler.GetDocumentRowRecord)

		document.GET("/reports", protect(), handler.GetReports)
		document.GET("/reports/:id", protect(), handler.GetReportByID)

		document.POST("upload", protect(), handler.UploadDocument)
		document.POST("/", handler.CreateDocument)

		document.PUT("/:id", protect(), handler.UpdateDocument)
		document.PUT("/records", protect(), handler.UpdateDocumentRowRecord)

		document.DELETE("/:id", protect(), handler.DeleteDocument)
		document.DELETE("/fields/:id", protect(), handler.DeleteFieldDocument)
		document.DELETE("/details/:document/records/:id", protect(), handler.DeleteDocumentRowRecord)
	}
}

package http

import (
	"github.com/RomanshkVolkov/test-api/internal/adapters/handler"
	"github.com/gin-gonic/gin"
)

func InitWebSocketRoutes(r *gin.Engine) {
	// WebSocket routes

	r.GET("/ws", handler.HandlerWebSocket)

	go handler.PollWarehouses()
}

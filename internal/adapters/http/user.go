package http

import (
	"github.com/RomanshkVolkov/test-api/internal/adapters/handler"
	"github.com/RomanshkVolkov/test-api/internal/adapters/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	protect := middleware.Protected
	users := r.Group("/users")
	{
		users.GET("/", protect(), handler.GetAllUsers)
		users.POST("/", protect(), handler.CreateUser)
		users.PUT("/:id", protect(), handler.UpdateUser)
		users.DELETE("/:id", protect(), handler.DeleteUser)

		users.GET("/me/profile", handler.GetUserProfile)
		users.GET("/:id", protect(), handler.GetUserByID)
	}

}

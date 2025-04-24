package service

import (
	"github.com/RomanshkVolkov/test-api/internal/core/domain"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Host domain.Database
}

func GetServer(c *gin.Context) *Server {
	return &Server{
		Host: domain.DBBetaAutopartes,
	}
}

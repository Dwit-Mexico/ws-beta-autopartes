package handler

import (
	"net/http"

	"github.com/RomanshkVolkov/test-api/internal/core/domain"
	"github.com/RomanshkVolkov/test-api/internal/core/service"
	"github.com/gin-gonic/gin"
)

func InternalSynchronization(c *gin.Context) {
	request, err := ValidateRequest[domain.WebPages](c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, ServerError(err, RequestError))
		return
	}

	server := service.GetServer(c)
	res, err := server.PermissionsSynchronization(request)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, ServerError(err, domain.Message{En: "error on permissions synchronization", Es: "error en la sincronización de permisos"}))
		return
	}

	c.IndentedJSON(http.StatusOK, res)
}

package handler

import (
	"net/http"

	"github.com/RomanshkVolkov/test-api/internal/core/domain"
	"github.com/RomanshkVolkov/test-api/internal/core/service"
	"github.com/gin-gonic/gin"
)

func GetCurrentHostingCenter(c *gin.Context) {
	server := service.GetServer(c)
	hotel := server.GetCurrentHostingCenter()

	c.IndentedJSON(http.StatusOK, hotel)
}

func UpdateHostingCenter(c *gin.Context) {
	request, err := ValidateRequest[domain.HostingCenter](c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, ServerError(err, RequestError))
		return
	}
	server := service.GetServer(c)
	hotel := server.UpdateHostingCenter(request)

	c.IndentedJSON(http.StatusOK, hotel)
}

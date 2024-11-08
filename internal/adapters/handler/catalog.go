package handler

import (
	"net/http"
	"strconv"

	"github.com/RomanshkVolkov/test-api/internal/adapters/repository"
	"github.com/RomanshkVolkov/test-api/internal/core/domain"
	"github.com/RomanshkVolkov/test-api/internal/core/service"
	"github.com/gin-gonic/gin"
)

// @Summary Just Get Kitchen by ID
// @Description Get kitchen by ID
// @tags Catalogs
// @Produce json
// @Security BearerAuth
// @Param id path int true "Kitchen ID"
// @Success 200 {object} domain.Kitchen "Kitchen"
// @Failure 400 {object} string "Unhandled error (report it)"
// @Failure 500 {object} string "Server error (report it)"
// @Router /catalogs/kitchens/{id} [get]
func GetKitchenByID(c *gin.Context) {
	server := service.GetServer(c)
	stringID := c.Param("id")
	id, err := strconv.ParseUint(stringID, 10, 64)

	if err != nil {
		InvalidParamError(c, "ID", err)
		return
	}

	kitchen := server.GetKitchenByID(uint(id))

	c.IndentedJSON(http.StatusOK, kitchen)
}

// @Summary Just Get Shift by ID
// @Description Get shift by ID
// @tags Catalogs
// @Produce json
// @Security BearerAuth
// @Param id path int true "Shift ID"
// @Success 200 {object} domain.Shift "Shift"
// @Failure 400 {object} string "Unhandled error (report it)"
// @Failure 500 {object} string "Server error (report it)"
// @Router /catalogs/shifts/{id} [get]
func GetShiftByID(c *gin.Context) {
	server := service.GetServer(c)
	stringID := c.Param("id")
	id, err := strconv.ParseUint(stringID, 10, 64)

	if err != nil {
		InvalidParamError(c, "ID", err)
		return
	}

	shift := server.GetShiftByID(uint(id))

	c.IndentedJSON(http.StatusOK, shift)
}

// @Summary Just Update Generic Catalog
// @Description Update a generic catalog
// @tags Catalogs
// @Produce json
// @Param id path int true "Catalog entity ID"
// @Success 200 {object} domain.APIResponse "Return message"
// @Failure 400 {object} string "Unhandled error (report it)"
// @Failure 500 {object} string "Server error (report it)"
// @Router /catalog/{slug}/{id} [put]
func UpdateGenericCatalog(c *gin.Context) {
	request, err := ValidateRequest[domain.GenericCatalog](c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, ServerError(nil, RequestError))
		return
	}

	stringSlug := c.Param("slug")
	stringID := c.Param("id")
	id, err := strconv.ParseUint(stringID, 10, 64)

	if err != nil {
		InvalidParamError(c, "ID", err)
		return
	}

	table := repository.GetCatalogTable(stringSlug)

	if table == nil {
		c.IndentedJSON(http.StatusBadRequest, ServerError(nil, domain.Message{
			En: "Invalid catalog",
			Es: "Catálogo inválido",
		}))
		return
	}

	if request.ID != uint(id) {
		c.IndentedJSON(http.StatusBadRequest, ServerError(nil, domain.Message{
			En: "ID does not match",
			Es: "El ID no coincide",
		}))
		return
	}

	server := service.GetServer(c)
	catalog := server.UpdateGenericCatalog(request, table)

	c.IndentedJSON(http.StatusOK, catalog)
}

// @Summary Just Delete Generic Catalog
// @Description Delete a generic catalog
// @tags Catalogs
// @Produce json
// @Param id path int true "Catalog entity ID"
// @Success 200 {object} domain.APIResponse "Return message"
// @Failure 400 {object} string "Unhandled error (report it)"
// @Failure 500 {object} string "Server error (report it)"
// @Router /catalog/{slug}/{id} [delete]
func DeleteGenericCatalog(c *gin.Context) {
	stringSlug := c.Param("slug")
	stringID := c.Param("id")
	id, err := strconv.ParseUint(stringID, 10, 64)

	if err != nil {
		InvalidParamError(c, "ID", err)
		return
	}

	table := repository.GetCatalogTable(stringSlug)

	server := service.GetServer(c)
	catalog := server.DeleteGenericCatalog(uint(id), table)

	c.IndentedJSON(http.StatusOK, catalog)
}

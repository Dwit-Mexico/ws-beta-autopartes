package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/RomanshkVolkov/ws-beta-autopartes/internal/core/domain"
	"github.com/RomanshkVolkov/ws-beta-autopartes/internal/core/service"
	"github.com/gin-gonic/gin"
)

func GetAllUsers(c *gin.Context) {
	server := service.GetServer(c)
	users := server.GetAllUsers()

	c.IndentedJSON(http.StatusOK, users)
}

func GetUserByID(c *gin.Context) {
	server := service.GetServer(c)
	stringID := c.Param("id")
	id, err := strconv.ParseUint(stringID, 10, 64)

	if err != nil {
		InvalidParamError(c, "ID", err)
		return
	}

	user := server.GetUserByID(uint(id))

	c.IndentedJSON(http.StatusOK, user)
}

func CreateUser(c *gin.Context) {
	request, err := ValidateRequest[domain.CreateUserRequest](c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, ServerError(nil, RequestError))
		return
	}

	server := service.GetServer(c)
	user := server.CreateUser(request)

	c.IndentedJSON(http.StatusOK, user)
}

func UpdateUser(c *gin.Context) {
	request, err := ValidateRequest[domain.EditableUser](c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, ServerError(nil, RequestError))
		return
	}

	stringID := c.Param("id")
	id, err := strconv.ParseUint(stringID, 10, 64)

	if err != nil {
		InvalidParamError(c, "ID", err)
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
	user := server.UpdateUser(request)

	c.IndentedJSON(http.StatusOK, user)
}

func DeleteUser(c *gin.Context) {
	stringID := c.Param("id")
	id, err := strconv.ParseUint(stringID, 10, 64)

	if err != nil {
		InvalidParamError(c, "ID", err)
		return
	}

	server := service.GetServer(c)
	user := server.DeleteUser(uint(id))

	c.IndentedJSON(http.StatusOK, user)
}

// @Summary Just User Profile by token
// @Description Get user profile by token
// @tags Users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} string "User profile"
// @Failure 400 {object} string "Unhandled error (report it)"
// @Failure 500 {object} string "Server error (report it)"
// @Router /users/me/profile [get]
func GetUserProfile(c *gin.Context) {

	user, exists := c.Get("userID")
	fmt.Println("user profile")
	fmt.Printf("%v", user)
	fmt.Println(exists)
	c.IndentedJSON(http.StatusOK, "user profile")
}

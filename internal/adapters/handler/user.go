package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/RomanshkVolkov/test-api/internal/core/domain"
	"github.com/RomanshkVolkov/test-api/internal/core/service"
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

// @Summary Just Users Profiles List
// @Description Get users profiles list
// @tags Users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} domain.UserProfiles "Users profiles"
// @Failure 400 {object} string "Unhandled error (report it)"
// @Failure 500 {object} string "Server error (report it)"
// @Router /users/profiles [get]
func GetUsersProfiles(c *gin.Context) {
	server := service.GetServer(c)
	users := server.GetUsersProfiles()

	c.IndentedJSON(http.StatusOK, users)
}

func CreateProfile(c *gin.Context) {
	request, err := ValidateRequest[domain.CreateProfile](c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, ServerError(nil, RequestError))
		return
	}

	server := service.GetServer(c)
	profile := server.CreateProfile(request)

	c.IndentedJSON(http.StatusOK, profile)
}

func UpdateProfile(c *gin.Context) {
	request, err := ValidateRequest[domain.EditableProfile](c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, ServerError(nil, RequestError))
		return
	}

	server := service.GetServer(c)
	profile := server.UpdateProfile(request)

	c.IndentedJSON(http.StatusOK, profile)
}

func DeleteProfile(c *gin.Context) {
	id, err := ExtractAndParseUintParam(c, "id")
	if err != nil {
		InvalidParamError(c, "ID", err)
		return
	}

	server := service.GetServer(c)
	profile := server.DeleteProfile(uint(id))

	c.IndentedJSON(http.StatusOK, profile)
}

func GetProfileByID(c *gin.Context) {
	server := service.GetServer(c)
	stringID := c.Param("id")
	id, err := strconv.ParseUint(stringID, 10, 64)

	if err != nil {
		InvalidParamError(c, "ID", err)
		return
	}

	user := server.GetProfileByID(uint(id))

	c.IndentedJSON(http.StatusOK, user)
}

func GetPermissions(c *gin.Context) {
	server := service.GetServer(c)
	permissions := server.GetPermissions()

	c.IndentedJSON(http.StatusOK, permissions)
}

// @Summary Just Kitchens List
// @Description Get Kitchens list
// @tags Users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} domain.Kitchen "Kitchens list"
// @Failure 400 {object} string "Unhandled error (report it)"
// @Failure 500 {object} string "Server error (report it)"
// @Router /users/shifts [get]
func Kitchens(c *gin.Context) {
	server := service.GetServer(c)
	kitchens := server.GetKitchens()

	c.IndentedJSON(http.StatusOK, kitchens)
}

// @Summary Just Create Kitchen
// @Description Create a new kitchen
// @tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param kitchen body domain.Kitchen true "Kitchen object"
// @Success 200 {object} domain.Kitchen "Kitchen created"
// @Failure 400 {object} string "Unhandled error (report it)"
// @Failure 500 {object} string "Server error (report it)"
// @Router /users/kitchens [post]
func CreateKitchen(c *gin.Context) {
	request, err := ValidateRequest[domain.GenericCatalog](c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, ServerError(err, RequestError))
		return
	}

	server := service.GetServer(c)
	kitchen := server.CreateKitchen(request)

	c.IndentedJSON(http.StatusOK, kitchen)
}

// @Summary Just Shifts List
// @Description Get shifts list
// @tags Users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} domain.Shift "Shifts list"
// @Failure 400 {object} string "Unhandled error (report it)"
// @Failure 500 {object} string "Server error (report it)"
// @Router /users/shifts [get]
func GetShifts(c *gin.Context) {
	server := service.GetServer(c)
	shifts := server.GetShifts()

	c.IndentedJSON(http.StatusOK, shifts)
}

// @Summary Just Create Shift
// @Description Create a new shift
// @tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param shift body domain.Shift true "Shift object"
// @Success 200 {object} domain.Shift "Shift created"
// @Failure 400 {object} string "Unhandled error (report it)"
// @Failure 500 {object} string "Server error (report it)"
// @Router /users/shifts [post]
func CreateShift(c *gin.Context) {
	request, err := ValidateRequest[domain.GenericCatalog](c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, ServerError(err, RequestError))
		return
	}

	server := service.GetServer(c)
	shift := server.CreateShift(request)

	c.IndentedJSON(http.StatusOK, shift)
}

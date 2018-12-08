package v1

import (
	"github.com/labstack/echo"
	"github.com/mohemohe/butimili-api/controllers/api"
	"github.com/mohemohe/butimili-api/models"
	"github.com/mohemohe/butimili-api/util"
	"net/http"
)

type (
	UserRequest struct {
		UserName    string `json:"username"`
		Password string `json:"password"`
	}

	UserResponse struct {
		models.APIBase
		Data *UserData `json:"data,omitempty"`
	}

	UserData struct {
		User  models.User `json:"user"`
		Token *string     `json:"token,omitempty"`
	}
)

func DeleteUser(c echo.Context) error {
	defer api.DeferGenericError(c)

	user := c.Get("user").(*models.User)
	if err := models.DeleteUser(user); err == nil {
		return c.JSON(http.StatusOK, &models.APIBase{api.GenericOK})
	} else {
		return c.JSON(http.StatusOK, &models.APIBase{api.GenericError})
	}
}

func PostUser(c echo.Context) error {
	defer api.DeferGenericError(c)

	userRequest := new(UserRequest)
	if err := c.Bind(userRequest); err != nil {
		panic("bind error")
	}

	if user := models.GetUserByUserName(userRequest.UserName); user == nil {
		err := models.UpsertUser(&models.User{
			UserName: userRequest.UserName,
			Password: *util.Bcrypt(userRequest.Password),
		})
		if err != nil {
			panic("cant create user")
		}
		user, token := models.AuthroizeUser(userRequest.UserName, userRequest.Password)
		if token == nil {
			panic("invalid token")
		}
		return c.JSON(http.StatusOK, &UserResponse{models.APIBase{api.GenericOK}, &UserData{
			User:  *user,
			Token: token,
		}})
	} else {
		panic("user already exists")
	}
}

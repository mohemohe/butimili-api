package v1

import (
	"github.com/labstack/echo"
	"github.com/mohemohe/butimili-api/controllers/api"
	"github.com/mohemohe/butimili-api/models"
	"net/http"
)

type (
	AuthRequest struct {
		UserName    string `json:"username"`
		Password string `json:"password"`
	}

	AuthResponse struct {
		models.APIBase
		Data *AuthData `json:"data,omitempty"`
	}

	AuthData struct {
		User  models.User `json:"user"`
		Token *string     `json:"token,omitempty"`
	}
)

func GetAuth(c echo.Context) error {
	defer api.DeferGenericError(c)

	user := c.Get("user").(*models.User)

	return c.JSON(http.StatusOK, &AuthResponse{
		models.APIBase{api.GenericOK},
		&AuthData{
			User: *user,
		},
	})
}

func PostAuth(c echo.Context) error {
	defer api.DeferGenericError(c)

	authRequest := new(AuthRequest)
	if err := c.Bind(authRequest); err != nil {
		panic("bind error")
	}

	user, token := models.AuthroizeUser(authRequest.UserName, authRequest.Password)
	if token == nil {
		panic("invalid token")
	}
	return c.JSON(http.StatusOK, &AuthResponse{models.APIBase{api.GenericOK}, &AuthData{
		User:  *user,
		Token: token,
	}})
}

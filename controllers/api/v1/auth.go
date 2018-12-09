package v1

import (
	"github.com/labstack/echo"
	"github.com/mohemohe/butimili-api/controllers/api"
	"github.com/mohemohe/butimili-api/models"
	"net/http"
)

type (
	AuthRequest struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}

	AuthResponse struct {
		models.APIBase `json:"omitempty"`
		Data           *AuthData `json:"data,omitempty"`
	}

	AuthData struct {
		User  models.User `json:"user"`
		Token *string     `json:"token,omitempty"`
	}
)

// @Summary 認証状態取得
// @Description AccessTokenが有効か調べることができます。
// @ID get-v1-auth
// @Accept json
// @Produce json
// @Security AccessToken
// @Success 200 {object} v1.AuthResponse
// @Failure 400 {object} models.APIBase
// @Router /api/v1/auth [get]
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

// @Summary 認証
// @Description ユーザー名とパスワードを使用して、AccessTokenを取得します。このAPIは認証なしでアクセスできます。
// @ID post-v1-auth
// @Accept json
// @Produce json
// @Param account body v1.AuthRequest true "認証情報"
// @Success 200 {object} v1.AuthResponse
// @Failure 400 {object} models.APIBase
// @Router /api/v1/auth [post]
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

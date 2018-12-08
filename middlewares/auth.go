package middlewares

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/mitchellh/mapstructure"
	"github.com/mohemohe/butimili-api/configs"
	"github.com/mohemohe/butimili-api/controllers/api"
	"github.com/mohemohe/butimili-api/models"
	"net/http"
	"strings"
)

var errorResult = models.APIBase{api.LoginNeeded}

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		switch {
		case c.Request().Header.Get(echo.HeaderAuthorization) != "":
			return jwtMiddleware(next, c)
		case c.QueryParam("token") != "":
			return jwtMiddleware(next, c)
		default:
			return c.JSON(http.StatusOK, errorResult)
		}
	}
}

func jwtMiddleware(next echo.HandlerFunc, c echo.Context) error {
	tokenString := c.QueryParams().Get("token")
	if tokenString == "" {
		authorization := c.Request().Header.Get(echo.HeaderAuthorization)
		if authorization == "" {
			return c.JSON(http.StatusUnauthorized, errorResult)
		}
		tokens := strings.Split(authorization, " ")
		if len(tokens) > 1 {
			tokenString = tokens[1]
		} else {
			tokenString = tokens[0]
		}
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// check signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			err := errors.New("unexpected signing method")
			return nil, err
		}
		return []byte(configs.GetEnv().Sign.Secret), nil
	})
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errorResult)
	}
	if !token.Valid {
		return c.JSON(http.StatusUnauthorized, errorResult)
	}

	claims := token.Claims.(jwt.MapClaims)
	jwtUser := &models.User{}
	metadata := &mapstructure.Metadata{}
	decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata: metadata,
		Result:   jwtUser,
		TagName:  "json",
	})
	if err := decoder.Decode(claims); err != nil {
		return c.JSON(http.StatusUnauthorized, errorResult)
	}

	user := models.GetUserByUserName(jwtUser.UserName)
	if user == nil {
		return c.JSON(http.StatusUnauthorized, errorResult)
	}

	c.Set("user", user)

	return next(c)
}

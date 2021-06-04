package controllers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func RedirectSwagger(c echo.Context) error {
	return c.Redirect(http.StatusPermanentRedirect, "/swagger/index.html")
}

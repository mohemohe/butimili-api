package controllers

import (
	"github.com/labstack/echo"
	"net/http"
)

func RedirectSwagger(c echo.Context) error {
	return c.Redirect(http.StatusPermanentRedirect, "/swagger/index.html")
}
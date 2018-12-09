package api

import (
	"github.com/labstack/echo"
	"github.com/mohemohe/butimili-api/models"
	"github.com/mohemohe/butimili-api/util"
	"github.com/sirupsen/logrus"
	"net/http"
)

func DeferGenericError(c echo.Context) error {
	err := recover()
	if err != nil {
		util.Logger().WithFields(logrus.Fields{
			"error": err,
		}).Warn("recover:", err)
		return c.JSON(http.StatusBadRequest, models.APIBase{GenericError})
	}
	return nil
}

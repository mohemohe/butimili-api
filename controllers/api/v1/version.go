package v1

import (
	"github.com/labstack/echo"
	"github.com/mohemohe/butimili-api/controllers/api"
	"github.com/mohemohe/butimili-api/models"
	"github.com/mohemohe/butimili-api/util"
	"net/http"
)

type (
	EnvironmentResponse struct {
		models.APIBase
		Data *VersionData `json:"data,omitempty"`
	}

	VersionData struct {
		Version string `json:"version"`
		Branch  string `json:"branch"`
		Hash    string `json:"hash"`
	}
)

func GetVersion(c echo.Context) error {
	defer api.DeferGenericError(c)

	return c.JSON(http.StatusOK, &EnvironmentResponse{
		models.APIBase{api.GenericOK},
		&VersionData{
			Version: util.GetVersion(),
			Branch:  util.GetBranch(),
			Hash:    util.GetHash(),
		},
	})
}

package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/mohemohe/butimili-api/controllers/api"
	"github.com/mohemohe/butimili-api/models"
	"net/http"
	"net/url"
	"strings"
)

type (
	InstanceRequest struct {
		FQDN string `json:"fqdn"`
	}

	InstanceListResponse struct {
		models.APIBase
		Data []string `json:"data"`
	}
)

// @Summary インスタンスリスト取得
// @Description インスタンスリストを取得します。
// @ID get-v1-instance-list
// @Produce json
// @Security AccessToken
// @Success 200 {object} v1.InstanceListResponse
// @Failure 400 {object} models.APIBase
// @Router /api/v1/instance/list [get]
func ListInstance(c echo.Context) error {
	defer api.DeferGenericError(c)

	user := c.Get("user").(*models.User)
	if instanceList := models.GetInstanceListByUserName(user.UserName); instanceList == nil {
		return c.JSON(http.StatusOK, &InstanceListResponse{
			models.APIBase{api.GenericOK},
			[]string{},
		})
	} else {
		return c.JSON(http.StatusOK, &InstanceListResponse{
			models.APIBase{api.GenericOK},
			instanceList.Instances,
		})
	}
}

// @Summary インスタンスリスト追加
// @Description インスタンスリストにFQDNを追加します。
// @ID get-v1-instance-list
// @Produce json
// @Security AccessToken
// @Param account body v1.InstanceRequest true "FQDN"
// @Success 200 {object} v1.InstanceListResponse
// @Failure 400 {object} models.APIBase
// @Router /api/v1/instance/list [put]
func PutInstance(c echo.Context) error {
	defer api.DeferGenericError(c)

	user := c.Get("user").(*models.User)

	instanceRequest := new(InstanceRequest)
	if err := c.Bind(instanceRequest); err != nil {
		panic("bind error")
	}
	screenName := strings.TrimPrefix(instanceRequest.FQDN, "@")

	if instanceList := models.GetInstanceListByUserName(user.UserName); instanceList == nil {
		instanceList = &models.InstanceList{
			UserName: user.UserName,
			Instances: []string{
				screenName,
			},
		}
		if err := models.UpsertInstanceList(instanceList); err != nil {
			panic("could not upsert")
		}
		return c.JSON(http.StatusOK, &InstanceListResponse{
			models.APIBase{api.GenericOK},
			instanceList.Instances,
		})
	} else {
		arr := append(instanceList.Instances, screenName)
		m := make(map[string]bool)
		uniq := []string{}

		for _, ele := range arr {
			if !m[ele] {
				m[ele] = true
				uniq = append(uniq, ele)
			}
		}

		instanceList.Instances = uniq
		if err := models.UpsertInstanceList(instanceList); err != nil {
			panic("could not upsert")
		} else {
			return c.JSON(http.StatusOK, &InstanceListResponse{
				models.APIBase{api.GenericOK},
				instanceList.Instances,
			})
		}
	}
}

// @Summary インスタンスリスト削除
// @Description インスタンスリストからFQDNを削除します。
// @ID get-v1-instance-list
// @Produce json
// @Security AccessToken
// @Param FQDN path string true "FQDN"
// @Success 200 {object} v1.InstanceListResponse
// @Failure 400 {object} models.APIBase
// @Router /api/v1/instance/list/{FQDN} [delete]
func DeleteInstance(c echo.Context) error {
	defer api.DeferGenericError(c)

	user := c.Get("user").(*models.User)

	target, err := url.PathUnescape(c.Param("FQDN"))
	if err != nil || target == "" {
		panic("missing target")
	}

	if instanceList := models.GetInstanceListByUserName(user.UserName); instanceList == nil {
		panic("instancest not found")
	} else {
		FQDN := strings.TrimPrefix(target, "@")
		arr := append(instanceList.Instances, FQDN)
		m := make(map[string]bool)
		uniq := []string{}

		for _, ele := range arr {
			if !m[ele] && ele != FQDN {
				m[ele] = true
				uniq = append(uniq, ele)
			}
		}

		instanceList.Instances = uniq
		if err := models.UpsertInstanceList(instanceList); err != nil {
			panic("could not upsert")
		} else {
			return c.JSON(http.StatusOK, &InstanceListResponse{
				models.APIBase{api.GenericOK},
				instanceList.Instances,
			})
		}
	}
}

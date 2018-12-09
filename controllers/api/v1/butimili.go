package v1

import (
	"github.com/labstack/echo"
	"github.com/mohemohe/butimili-api/controllers/api"
	"github.com/mohemohe/butimili-api/models"
	"net/http"
	"net/url"
	"strings"
)

type (
	ButimiliRequest struct {
		ScreenName string `json:"screen_name"`
	}

	ButimiListResponse struct {
		models.APIBase
		Data []string `json:"data"`
	}
)

const (
	Butimili = "うおおおおおおおおおおおおあああああああああああああああああああああああああああああああ！！！！！！！！！！！ (ﾌﾞﾘﾌﾞﾘﾌﾞﾘﾌﾞﾘｭﾘｭﾘｭﾘｭﾘｭﾘｭ！！！！！！ﾌﾞﾂﾁﾁﾌﾞﾌﾞﾌﾞﾁﾁﾁﾁﾌﾞﾘﾘｲﾘﾌﾞﾌﾞﾌﾞﾌﾞｩｩｩｩｯｯｯ！！！！！！！)"
)

func GetButimiliText(c echo.Context) error {
	defer api.DeferGenericError(c)

	return c.String(http.StatusOK, Butimili)
}

func ListButimili(c echo.Context) error {
	defer api.DeferGenericError(c)

	user := c.Get("user").(*models.User)
	if butimiList := models.GetButimiListByUserName(user.UserName); butimiList == nil {
		return c.JSON(http.StatusOK, &ButimiListResponse{
			models.APIBase{api.GenericOK},
			[]string{},
		})
	} else {
		return c.JSON(http.StatusOK, &ButimiListResponse{
			models.APIBase{api.GenericOK},
			butimiList.Targets,
		})
	}
}

func ListButimiliText(c echo.Context) error {
	defer api.DeferGenericError(c)

	user := c.Get("user").(*models.User)
	if butimiList := models.GetButimiListByUserName(user.UserName); butimiList == nil {
		return c.String(http.StatusOK, "")
	} else {
		text := ""
		for _, v := range butimiList.Targets {
			text += " @" + v
		}
		text = strings.TrimSpace(text + " " + Butimili)
		return c.String(http.StatusOK, text)
	}
}

func PutButimili(c echo.Context) error {
	defer api.DeferGenericError(c)

	user := c.Get("user").(*models.User)

	butimiliRequest := new(ButimiliRequest)
	if err := c.Bind(butimiliRequest); err != nil {
		panic("bind error")
	}

	if butimiList := models.GetButimiListByUserName(user.UserName); butimiList == nil {
		panic("butimilist not found")
	} else {
		screenName := strings.TrimPrefix(butimiliRequest.ScreenName, "@")
		arr := append(butimiList.Targets, screenName)
		m := make(map[string]bool)
		uniq := []string{}

		for _, ele := range arr {
			if !m[ele] {
				m[ele] = true
				uniq = append(uniq, ele)
			}
		}

		butimiList.Targets = uniq
		if err := models.UpsertButimiList(butimiList); err != nil {
			panic("could not upsert")
		} else {
			return c.JSON(http.StatusOK, &ButimiListResponse{
				models.APIBase{api.GenericOK},
				butimiList.Targets,
			})
		}
	}
}

func DeleteButimili(c echo.Context) error {
	defer api.DeferGenericError(c)

	user := c.Get("user").(*models.User)

	target, err := url.PathUnescape(c.Param("screenName"))
	if err != nil || target == "" {
		panic("missing target")
	}

	if butimiList := models.GetButimiListByUserName(user.UserName); butimiList == nil {
		panic("butimilist not found")
	} else {
		screenName := strings.TrimPrefix(target, "@")
		arr := append(butimiList.Targets, screenName)
		m := make(map[string]bool)
		uniq := []string{}

		for _, ele := range arr {
			if !m[ele] && ele != screenName {
				m[ele] = true
				uniq = append(uniq, ele)
			}
		}

		butimiList.Targets = uniq
		if err := models.UpsertButimiList(butimiList); err != nil {
			panic("could not upsert")
		} else {
			return c.JSON(http.StatusOK, &ButimiListResponse{
				models.APIBase{api.GenericOK},
				butimiList.Targets,
			})
		}
	}
}

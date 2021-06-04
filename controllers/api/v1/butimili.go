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

// @Summary 生ブチミリ
// @Description ブチミリ本文のみを取得します。このAPIは認証なしでアクセスできます。
// @ID get-v1-butimili-raw
// @Produce plain
// @Success 200 {string} Butimili "うおおおおおおおおおおおおあああああああああああああああああああああああああああああああ！！！！！！！！！！！ (ﾌﾞﾘﾌﾞﾘﾌﾞﾘﾌﾞﾘｭﾘｭﾘｭﾘｭﾘｭﾘｭ！！！！！！ﾌﾞﾂﾁﾁﾌﾞﾌﾞﾌﾞﾁﾁﾁﾁﾌﾞﾘﾘｲﾘﾌﾞﾌﾞﾌﾞﾌﾞｩｩｩｩｯｯｯ！！！！！！！)"
// @Failure 400 {object} models.APIBase
// @Router /api/v1/butimili/raw [get]
func GetButimiliText(c echo.Context) error {
	defer api.DeferGenericError(c)

	return c.String(http.StatusOK, Butimili)
}

// @Summary ブチミリスト取得
// @Description ブチミリストを取得します。
// @ID get-v1-butimili-list
// @Produce json
// @Security AccessToken
// @Success 200 {object} v1.ButimiListResponse
// @Failure 400 {object} models.APIBase
// @Router /api/v1/butimili/list [get]
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

// @Summary ブチミリ
// @Description ブチミリストと生ブチミリを結合した文字列を取得します。
// @ID get-v1-butimili
// @Produce json
// @Security AccessToken
// @Success 200 {string} Butimili
// @Failure 400 {object} models.APIBase
// @Router /api/v1/butimili [get]
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

// @Summary ブチミリスト追加
// @Description ブチミリストにScreenNameを追加します。
// @ID get-v1-butimili-list
// @Produce json
// @Security AccessToken
// @Param account body v1.ButimiliRequest true "アカウント情報"
// @Success 200 {object} v1.ButimiListResponse
// @Failure 400 {object} models.APIBase
// @Router /api/v1/butimili/list [put]
func PutButimili(c echo.Context) error {
	defer api.DeferGenericError(c)

	user := c.Get("user").(*models.User)

	butimiliRequest := new(ButimiliRequest)
	if err := c.Bind(butimiliRequest); err != nil {
		panic("bind error")
	}
	screenName := strings.TrimPrefix(butimiliRequest.ScreenName, "@")

	if butimiList := models.GetButimiListByUserName(user.UserName); butimiList == nil {
		butimiList = &models.ButimiList{
			UserName: user.UserName,
			Targets: []string{
				screenName,
			},
		}
		if err := models.UpsertButimiList(butimiList); err != nil {
			panic("could not upsert")
		}
		return c.JSON(http.StatusOK, &ButimiListResponse{
			models.APIBase{api.GenericOK},
			butimiList.Targets,
		})
	} else {
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

// @Summary ブチミリスト削除
// @Description ブチミリストからScreenNameを削除します。
// @ID get-v1-butimili-list
// @Produce json
// @Security AccessToken
// @Param screen_name path string true "ScreenName"
// @Success 200 {object} v1.ButimiListResponse
// @Failure 400 {object} models.APIBase
// @Router /api/v1/butimili/list/{screen_name} [delete]
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

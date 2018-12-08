package api

import (
	"github.com/ahl5esoft/golang-underscore"
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo"
	"github.com/mohemohe/butimili-api/models"
	"github.com/mohemohe/butimili-api/util"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
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

func ParseListQuery(c echo.Context) *models.ListQuery {
	raw := models.RawListQuery{
		Start:     c.QueryParam("_start"),
		End:       c.QueryParam("_end"),
		Sort:      c.QueryParam("_sort"),
		Order:     c.QueryParam("_order"),
		IDLike:    c.QueryParam("id_like"),
		StartTime: c.QueryParam("startTime"),
		EndTime:   c.QueryParam("endTime"),
	}
	query := models.ListQuery{
		Sort:  "_id",
		Limit: 10,
		Skip:  0,
		Query: bson.M{},
		Raw:   &raw,
	}

	if raw.Start != "" && raw.End != "" {
		startNum, errStart := strconv.Atoi(raw.Start)
		endNum, errEnd := strconv.Atoi(raw.End)
		if errStart == nil && errEnd == nil {
			query.Skip = startNum
			query.Limit = endNum - startNum
		}
	}
	if raw.Sort != "" {
		query.Sort = raw.Sort
	}
	if raw.Order != "" && strings.ToUpper(raw.Order) == "DESC" {
		query.Sort = "-" + query.Sort
	}
	switch {
	case raw.IDLike != "":
		queryParams, err := url.QueryUnescape(raw.IDLike)
		if err != nil {
			params := strings.Split(queryParams, "|")
			v := underscore.Map(params, func(s string, _ int) bson.M {
				return bson.M{
					"_id": bson.ObjectIdHex(s),
				}
			})
			query.Query = bson.M{
				"$or": v.([]bson.M),
			}
		}
	case raw.StartTime != "" && raw.EndTime != "":
		layout := "2006-01-02"
		startTime, err := time.Parse(layout, raw.StartTime)
		if err != nil {
			break
		}
		endTime, err := time.Parse(layout, raw.EndTime)
		if err != nil {
			break
		}

		query.Query = bson.M{
			"$and": []bson.M{
				{"start": bson.M{
					"$gte": startTime.Unix(),
				}},
				{"start": bson.M{
					"$lt": endTime.Unix(),
				}},
			},
		}
	case raw.StartTime != "":
		layout := "2006-01-02"
		startTime, err := time.Parse(layout, raw.StartTime)
		if err != nil {
			break
		}
		query.Query = bson.M{
			"start": bson.M{
				"$gte": startTime.Unix(),
			},
		}
	case raw.EndTime != "":
		layout := "2006-01-02"
		endTime, err := time.Parse(layout, raw.EndTime)
		if err != nil {
			break
		}
		query.Query = bson.M{
			"start": bson.M{
				"$lt": endTime.Unix(),
			},
		}
	}

	return &query
}

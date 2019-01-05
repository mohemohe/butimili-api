package models

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/mohemohe/butimili-api/models/connections"
	"github.com/mohemohe/butimili-api/util"
	"github.com/sirupsen/logrus"
)

var collections = struct {
	Users         string
	ButimiLists   string
	InstanceLists string
}{
	Users:         "users",
	ButimiLists:   "butimilists",
	InstanceLists: "instancelists",
}

type (
	APIBase struct {
		// エラーコード
		Error int `bson:"-" json:"error" example:"1"`
	}

	ListQuery struct {
		Limit int
		Skip  int
		Sort  string
		Query bson.M
		Raw   *RawListQuery
	}

	RawListQuery struct {
		Start     string
		End       string
		Sort      string
		Order     string
		IDLike    string
		StartTime string
		EndTime   string
	}
)

func InitDB() {
	createIndex()
	util.Logger().Info("DB initialized")
}

func createIndex() {
	ensureIndex(collections.Users, getIndex([]string{"username"}, true))
}

func getIndex(key []string, unique bool) mgo.Index {
	return mgo.Index{
		Key:        key,
		Unique:     unique,
		Background: true,
	}
}

func ensureIndex(collection string, index mgo.Index) {
	util.Logger().WithFields(logrus.Fields{
		"collection": collection,
		"index":      index,
	}).Debug("create index")
	if err := connections.Mongo().Collection(collection).Collection().EnsureIndex(index); err != nil {
		util.Logger().WithFields(logrus.Fields{
			"collection": collection,
			"index":      index,
			"error":      err,
		}).Warn("index create error")
	}
}

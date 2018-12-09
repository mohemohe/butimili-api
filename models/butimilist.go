package models

import (
	"github.com/globalsign/mgo/bson"
	"github.com/go-bongo/bongo"
	"github.com/mohemohe/butimili-api/models/connections"
)

type (
	ButimiList struct {
		bongo.DocumentBase `bson:",inline"`
		UserName           string   `bson:"username" json:"username"`
		Targets            []string `bson:"targets" json:"targets"`
	}
)

func GetButimiListByUserName(username string) *ButimiList {
	conn := connections.Mongo()

	butimiList := &ButimiList{}
	err := conn.Collection(collections.ButimiLists).FindOne(bson.M{
		"username": username,
	}, butimiList)
	if err != nil {
		return nil
	}

	return butimiList
}

func UpsertButimiList(butimiList *ButimiList) error {
	return connections.Mongo().Collection(collections.ButimiLists).Save(butimiList)
}

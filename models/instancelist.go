package models

import (
	"github.com/globalsign/mgo/bson"
	"github.com/go-bongo/bongo"
	"github.com/mohemohe/butimili-api/models/connections"
)

type (
	InstanceList struct {
		bongo.DocumentBase `bson:",inline"`
		UserName           string   `bson:"username" json:"username"`
		Instances          []string `bson:"targets" json:"targets"`
	}
)

func GetInstanceListByUserName(username string) *InstanceList {
	conn := connections.Mongo()

	instanceList := &InstanceList{}
	err := conn.Collection(collections.InstanceLists).FindOne(bson.M{
		"username": username,
	}, instanceList)
	if err != nil {
		return nil
	}

	return instanceList
}

func UpsertInstanceList(instanceList *InstanceList) error {
	return connections.Mongo().Collection(collections.InstanceLists).Save(instanceList)
}

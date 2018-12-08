package connections

import (
	"crypto/tls"
	"github.com/globalsign/mgo"
	"github.com/go-bongo/bongo"
	"github.com/mohemohe/butimili-api/configs"
	"github.com/mohemohe/butimili-api/util"
	"github.com/sirupsen/logrus"
	"net"
	"time"
)

var mongoConn *bongo.Connection

func Mongo() *bongo.Connection {
	if mongoConn == nil {
		mongoConn = NewMongo()

	}
	return mongoConn
}

func NewMongo() *bongo.Connection {
	util.Logger().WithFields(logrus.Fields{
		"address":  configs.GetEnv().Mongo.Address,
		"database": configs.GetEnv().Mongo.Database,
		"ssl":      configs.GetEnv().Mongo.SSL,
	}).Info("create mongo connection")

	config := &bongo.Config{
		ConnectionString: configs.GetEnv().Mongo.Address,
		Database:         configs.GetEnv().Mongo.Database,
	}

	if configs.GetEnv().Mongo.SSL {
		// REF: https://github.com/go-bongo/bongo/pull/11
		if dialInfo, err := mgo.ParseURL(config.ConnectionString); err != nil {
			util.Logger().Fatal(err)
		} else {
			config.DialInfo = dialInfo
		}

		tlsConfig := &tls.Config{}
		config.DialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
			return conn, err
		}
		config.DialInfo.Timeout = time.Second * 3
	}

	conn, err := bongo.Connect(config)
	if err != nil {
		panic(err)
	}
	util.Logger().Info("mongo connection created")

	return conn
}

package models

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/globalsign/mgo/bson"
	"github.com/go-bongo/bongo"
	"github.com/mohemohe/butimili-api/configs"
	"github.com/mohemohe/butimili-api/models/connections"
	"github.com/mohemohe/butimili-api/util"
	"time"
)

type (
	User struct {
		bongo.DocumentBase `bson:",inline"`
		UserName           string `bson:"username" json:"username"`
		Password           string `bson:"password" json:"-"`
	}

	JwtClaims struct {
		ID       string `json:"id"`
		UserName string `json:"username"`
		jwt.StandardClaims
	}
)

func GetUserByUserName(username string) *User {
	conn := connections.Mongo()

	user := &User{}
	err := conn.Collection(collections.Users).FindOne(bson.M{
		"username": username,
	}, user)
	if err != nil {
		return nil
	}

	return user
}

func GetUserByHexID(hexID string) *User {
	conn := connections.Mongo()

	user := &User{}
	err := conn.Collection(collections.Users).FindById(bson.ObjectIdHex(hexID), user)
	if err != nil {
		return nil
	}

	return user
}

func UpsertUser(user *User) error {
	if !util.IsBcrypt(user.Password) {
		user.Password = *util.Bcrypt(user.Password)
	}
	return connections.Mongo().Collection(collections.Users).Save(user)
}

func DeleteUser(user *User) error {
	return connections.Mongo().Collection(collections.Users).DeleteDocument(user)
}

func AuthroizeUser(username string, password string) (*User, *string) {
	user := GetUserByUserName(username)
	if user == nil {
		panic("user not found")
	}

	if !util.CompareHash(password, user.Password) {
		panic("wrong password")
	}

	claims := &JwtClaims{
		user.GetId().Hex(),
		user.UserName,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ts, err := token.SignedString([]byte(configs.GetEnv().Sign.Secret))
	if err != nil {
		panic("couldnt create token")
	}

	return user, &ts
}

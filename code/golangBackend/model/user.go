package model

import "github.com/globalsign/mgo/bson"

type UserRole int

const (
	AdminRole   UserRole = 0
	AuthorRole  UserRole = 1
	ReaderRole  UserRole = 2
	VisitorRole UserRole = 3
)

type User struct {
	Id          bson.ObjectId `bson:"_id,omitempty" json:"_id,omitempty"`
	UserName    string        `bson:"username" json:"username"`
	Password    string        `bson:"password,omitempty" json:"password,omitempty"`
	Avatar      string        `bson:"avatar" json:"avatar"`
	Email       string        `bson:"email,omitempty" json:"email,omitempty"`
	Role        UserRole      `bson:"role" json:"role"`
	BangumiList []string      `bson:"bangumi_list" json:"bangumi_list"`
}

type Login struct {
	UserName string `bson:"username" json:"username"`
	Password string `bson:"password,omitempty" json:"password,omitempty"`
}

func (u *User) GenerateID() {
	if u.Id == "" {
		u.Id = bson.NewObjectId()
	}
}

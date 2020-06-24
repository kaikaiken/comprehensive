package model

import (
	"github.com/globalsign/mgo/bson"
)

type Article struct {
	Id             bson.ObjectId   `bson:"_id,omitempty" json:"_id,omitempty"`
	Title          string          `bson:"title" json:"title"`
	Desc           string          `bson:"desc" json:"desc"`
	Cover          string          `bson:"cover" json:"cover"`
	PubTime        string          `bson:"pub_time" json:"pub_time"`
	LastModifyTime string          `bson:"last_modify_time" json:"last_modify_time"`
	Author         string          `bson:"author" json:"author"`
	Content        string          `bson:"content" json:"content"`
	Type           string          `bson:"type" json:"type"`
	Tags           string          `bson:"tags" json:"tags"`
	Kind           string          `bson:"kind" json:"kind"`
	Comment        []bson.ObjectId `bson:"comment" json:"comment"`
	ReadCount      int             `bson:"read_count" json:"read_count"`
	LikeCount      int             `bson:"like_count" json:"like_count"`
}

type ArticleTag struct {
	Id         bson.ObjectId `bson:"_id,omitempty" json:"_id,omitempty"`
	ArticleTag string        `bson:"articleTag" json:"articleTag"`
}

type ArticleKind struct {
	Id          bson.ObjectId `bson:"_id,omitempty" json:"_id,omitempty"`
	ArticleKind string        `bson:"articleKind" json:"articleKind"`
}

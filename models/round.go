package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Round struct {
	Id    bson.ObjectId `json:"id"       bson:"_id"`
	Date  time.Time     `json:"date"     bson:"date"`
	Buyer User          `json:"buyer"    bson:"buyer"`
	Peons []User        `json:"peons"    bson:"peons"`
}

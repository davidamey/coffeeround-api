package models

import "gopkg.in/mgo.v2/bson"

type User struct {
	Id        bson.ObjectId `json:"id"        bson:"_id"`
	Name      string        `json:"name"      bson:"name"`
	TotalDay  int           `json:"totalDay"  bson:"totalDay"`
	TotalWeek int           `json:"totalWeek" bson:"totalWeek"`
	Total     int           `json:"total"     bson:"total`
}

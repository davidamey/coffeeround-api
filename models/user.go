package models

import "gopkg.in/mgo.v2/bson"

type User struct {
	ID       bson.ObjectId `json:"id"         bson:"_id,omitempty"`
	GoogleID string        `json:"gid"        bson:"gid,omitempty"`
	Name     string        `json:"name"       bson:"name"`
	Picture  string        `json:"picture"    bson:"picture"`
}

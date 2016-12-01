package models

import (
	"log"

	"github.com/davecgh/go-spew/spew"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DataStore interface {
	GetUsers() *[]User
	GetUser(string) (*User, bool)
	UpsertUser(string, string, string) (bson.ObjectId, error)
	Close()
}

type dataStore struct {
	session *mgo.Session
}

var session *mgo.Session

func GetDataStore() DataStore {
	// Create one 'master' session
	if session == nil {
		var err error
		session, err = mgo.Dial("mongodb://localhost")
		if err != nil {
			panic(err)
		}
	}

	// Callees get a copy of the master session and are expected to close it when they're done.
	return &dataStore{session.Copy()}
}

func (ds *dataStore) db() *mgo.Database {
	return ds.session.DB("coffeeround")
}

func (ds *dataStore) userCol() *mgo.Collection {
	return ds.db().C("user")
}

func (ds *dataStore) roundCol() *mgo.Collection {
	return ds.db().C("round")
}

// Interface methods

func (ds *dataStore) Close() {
	ds.session.Close()
}

func (ds *dataStore) GetUsers() *[]User {
	var users []User
	ds.userCol().Find(nil).All(&users)

	return &users
}

func (ds *dataStore) GetUser(id string) (*User, bool) {
	var u User
	if err := ds.userCol().FindId(id).One(&u); err != nil {
		return nil, false
	}

	return &u, true
}

func (ds *dataStore) UpsertUser(gID, name, picture string) (id bson.ObjectId, err error) {
	var u User
	if err := ds.userCol().Find(bson.M{"gid": gID}).One(&u); err != nil {
		spew.Dump(err)
		if err == mgo.ErrNotFound {
			log.Println("User not found - creating new one")
			u = User{ID: bson.NewObjectId(), GoogleID: gID}
		} else {
			return "", err
		}
	}

	// If we're here then we have a complete User and can now update it.
	u.Name = name
	u.Picture = picture

	spew.Dump(u)

	if _, err := ds.userCol().UpsertId(u.ID, u); err != nil {
		return "", err
	} else {
		return u.ID, nil
	}
}

func (ds *dataStore) GetRounds() *[]Round {
	return nil
}

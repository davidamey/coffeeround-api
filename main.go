package main

import (
	"os"

	mgo "gopkg.in/mgo.v2"

	"github.com/davidamey/coffeeround-api/controllers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	uc := controllers.NewUserController(getDB())
	r.GET("/user/:id", uc.GetUser)
	r.POST("/user/", uc.CreateUser)

	r.Run(":" + os.Getenv("PORT"))
}

func getDB() *mgo.Database {
	return getSession().DB("coffeeapp")
}

func getSession() *mgo.Session {
	s, err := mgo.Dial("mongodb://localhost")

	if err != nil {
		panic(err)
	}
	return s
}

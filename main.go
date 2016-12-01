package main

import (
	"os"

	mgo "gopkg.in/mgo.v2"

	"github.com/davidamey/coffeeround-api/controllers"
	"github.com/davidamey/coffeeround-api/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// db := getDB()

	sc := controllers.NewSecurityController()
	r.POST("/login", sc.Login)

	authed := r.Group("/", handlers.Secure())

	uc := controllers.NewUserController()
	authed.GET("/user", uc.GetUsers)
	authed.GET("/user/:id", uc.GetUser)

	rc := contollers.NewRoundController()

	r.Run(":" + os.Getenv("PORT"))
}

func getDB() *mgo.Database {
	return nil
	// return getSession().DB("coffeeround")
}

func getSession() *mgo.Session {
	s, err := mgo.Dial("mongodb://localhost")

	if err != nil {
		panic(err)
	}
	return s
}

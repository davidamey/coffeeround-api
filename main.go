package main

import (
	"os"

	mgo "gopkg.in/mgo.v2"

	"github.com/davidamey/coffeeround-api/controllers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	db := getDB()

	// security
	sc := controllers.NewSecurityController(db)
	// r.GET("/login", sc.Login)
	// r.POST("/tokenSignIn", sc.TokenSignIn)

	authed := r.Group("/", sc.SecureHandler)

	uc := controllers.NewUserController(db)
	authed.GET("/user", uc.GetUsers)
	authed.GET("/user/:id", uc.GetUser)
	authed.POST("/user", uc.CreateUser)

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

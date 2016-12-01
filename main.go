package main

import (
	"os"

	"github.com/davidamey/coffeeround-api/controllers"
	"github.com/davidamey/coffeeround-api/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	sc := controllers.NewSecurityController()
	r.POST("/login", sc.Login)

	authed := r.Group("/", handlers.Secure())

	uc := controllers.NewUserController()
	authed.GET("/user", uc.GetUsers)
	authed.GET("/user/:id", uc.GetUser)

	rc := controllers.NewRoundController()
	authed.GET("/round", rc.GetRounds)
	authed.GET("/round/:id", rc.GetRound)

	r.Run(":" + os.Getenv("PORT"))
}

package main

import (
	"os"

	"github.com/davidamey/coffeeround-api/controllers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	uc := controllers.NewUserController()
	r.GET("/user/:id", uc.GetUser)

	r.Run(":" + os.Getenv("PORT"))
}

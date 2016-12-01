package handlers

import (
	"log"

	"github.com/davidamey/coffeeround-api/models"
	"github.com/gin-gonic/gin"
)

const stubbed = false

func Secure() gin.HandlerFunc {
	if stubbed {
		return secureStub()
	} else {
		return secure()
	}
}

func secureStub() gin.HandlerFunc {
	log.Println("*** stubbed jwt ***")
	return func(c *gin.Context) {
		// c.Set("userID", &models.User{
		// 	Id:      "117594728625162368089",
		// 	Name:    "Mr Bob",
		// 	Picture: "",
		// })
		c.Next()
	}
}

func secure() gin.HandlerFunc {
	jm := models.NewJWTMiddleware()

	return func(c *gin.Context) {
		err := jm.CheckJWT(c.Writer, c.Request)

		if err != nil {
			c.AbortWithStatus(401)
		} else {
			c.Next()
		}
	}
}

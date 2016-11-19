package controllers

import (
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	mgo "gopkg.in/mgo.v2"
)

var signingKey = []byte("secret")
var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	},
	SigningMethod: jwt.SigningMethodHS256,
})

type SecurityController interface {
	Login(*gin.Context)
	SecureHandler(*gin.Context)
}

type securityController struct {
	db *mgo.Database
}

func NewSecurityController(db *mgo.Database) SecurityController {
	return &securityController{db}
}

func (sc *securityController) Login(c *gin.Context) {
	token := jwt.New(jwt.SigningMethodHS256)

	// token.Claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	ts, _ := token.SignedString(signingKey)

	c.JSON(200, gin.H{"token": ts})
}

func (sc *securityController) SecureHandler(c *gin.Context) {
	c.Next()

	// err := jwtMiddleware.CheckJWT(c.Writer, c.Request)

	// if err != nil {
	// 	c.AbortWithStatus(401)
	// } else {
	// 	c.Next()
	// }
}

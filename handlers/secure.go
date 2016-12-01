package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"gopkg.in/mgo.v2/bson"

	jwt "github.com/dgrijalva/jwt-go"
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
		c.Set("userID", bson.ObjectId("584013ec099c9847fc1d6f1a"))
		c.Next()
	}
}

func secure() gin.HandlerFunc {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte("beans"), nil
	}

	return func(c *gin.Context) {
		token, err := tokenFromHeader(c.Request)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		if token == "" {
			c.String(http.StatusBadRequest, "Empty token")
			c.Abort()
			return
		}

		parsedToken, err := jwt.Parse(token, keyFunc)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if parsedToken.Header["alg"] != jwt.SigningMethodHS256.Alg() {
			log.Printf(" Error: token algorithm %q doesn't match expected %q", parsedToken.Header["alg"], jwt.SigningMethodHS256.Alg())
			c.String(http.StatusBadRequest, "Invalid algorithm")
			c.Abort()
			return
		}

		if !parsedToken.Valid {
			c.String(http.StatusUnauthorized, "Invalid token")
			c.Abort()
			return
		}

		// We have a valid token here - set the userID on the context
		idHex := parsedToken.Claims.(jwt.MapClaims)["id"].(string)
		id := bson.ObjectIdHex(idHex)

		c.Set("userID", id)
		c.Next()
	}
}

func tokenFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("Missing token")
	}

	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", fmt.Errorf("Authorization header format must be Bearer <token>")
	}

	return authHeaderParts[1], nil
}

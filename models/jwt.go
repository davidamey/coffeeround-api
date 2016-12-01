package models

import (
	"os"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
)

func NewJWT(id string) string {
	signingKey := []byte(os.Getenv("SIGNING_KEY"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"id":  id,
	})

	ts, _ := token.SignedString(signingKey)

	return ts
}

func NewJWTMiddleware() *jwtmiddleware.JWTMiddleware {
	signingKey := []byte(os.Getenv("SIGNING_KEY"))
	return jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return signingKey, nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})
}

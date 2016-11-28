package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ericchiang/oidc"
	"github.com/gin-gonic/gin"
	mgo "gopkg.in/mgo.v2"
)

var signingKey = os.Getenv("GOOGLE_CLIENT_KEY")

type Claims struct {
	Name       string `json:"name"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Locale     string `json:"locale"`
	IssuedAt   int64  `json:"iat"`
	// Audience      string `json:"aud"`
	UserID        string `json:"sub"`
	EmailVerified bool   `json:"email_verified"`
	// AZP           string `json:"azp"`
	Email      string `json:"email"`
	PictureURL string `json:"picture"`
	// ISS           string `json:"iss"`
	Expires int64 `json:"exp"`
}

type SecurityController interface {
	Info(*gin.Context)
	SecureHandler(*gin.Context)
}

type securityController struct {
	verifier *oidc.IDTokenVerifier
	db       *mgo.Database
}

func NewSecurityController(db *mgo.Database) SecurityController {
	provider, _ := oidc.NewProvider(context.Background(), "https://accounts.google.com")
	verifier := provider.NewVerifier(context.Background())
	return &securityController{verifier, db}
}

func (sc *securityController) Info(c *gin.Context) {
	if userID, exists := c.Get("userID"); exists {
		claims, _ := c.Get("claims")
		c.JSON(200, gin.H{
			"userID": userID,
			"claims": claims,
		})
	} else {
		c.String(200, "missing userID")
	}
}

func getTokenFromRequest(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", nil // No error, just no token
	}

	// TODO: Make this a bit more robust, parsing-wise
	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", fmt.Errorf("Authorization header format must be Bearer {token}")
	}

	return authHeaderParts[1], nil
}

func (sc *securityController) SecureHandler(c *gin.Context) {
	idToken, err := getTokenFromRequest(c.Request)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(401)
	}

	if len(idToken) == 0 {
		log.Println("No token")
		c.AbortWithStatus(401)
	}

	token, err := sc.verifier.Verify(idToken)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(401)
	}

	var claims Claims
	if err := token.Claims(&claims); err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.Set("userID", claims.UserID)
	c.Set("claims", claims)
	c.Next()
}

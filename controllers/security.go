package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/davidamey/coffeeround-api/models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/ericchiang/oidc"
	"github.com/gin-gonic/gin"
)

type claims struct {
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Locale        string `json:"locale"`
	IssuedAt      int64  `json:"iat"`
	Audience      string `json:"aud"`
	Subject       string `json:"sub"`
	EmailVerified bool   `json:"email_verified"`
	AZP           string `json:"azp"`
	Email         string `json:"email"`
	Picture       string `json:"picture"`
	Issuer        string `json:"iss"`
	Expires       int64  `json:"exp"`
}

type SecurityController interface {
	Login(*gin.Context)
}

type securityController struct {
	verifier *oidc.IDTokenVerifier
}

func NewSecurityController() SecurityController {
	ctx := context.Background()
	provider, _ := oidc.NewProvider(ctx, "https://accounts.google.com")
	verifier := provider.NewVerifier(ctx) //, oidc.VerifyAudience(android_client_key), oidc.VerifyExpiry())
	return &securityController{verifier}
}

func (sc *securityController) Login(c *gin.Context) {
	ds := models.GetDataStore()
	defer ds.Close()

	id, err := ds.UpsertUser("117594728625162368089", "David Amey", "https://lh6.googleusercontent.com/-lek52PddU30/AAAAAAAAAAI/AAAAAAAAl-E/tSL0W-QcesA/s96-c/photo.jpg")
	if err != nil {
		log.Printf("Error upserting user: %q\n", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	jwt := newJWT(id)

	log.Printf("UserId: %q, JWT: %q\n", id, jwt)
	c.String(200, jwt)
}

func (sc *securityController) LoginReal(c *gin.Context) {
	idToken := c.PostForm("idToken")

	if len(idToken) == 0 {
		log.Println("No token")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	token, err := sc.verifier.Verify(idToken)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(401)
		return
	}

	var cl claims
	if err := token.Claims(&cl); err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// If we're here, we trust them and have details
	// "Upsert" their details into our DB and send them a JWT
	// u := &models.User{GoogleID: cl.Subject, Name: cl.Name, Picture: cl.Picture}

	ds := models.GetDataStore()
	defer ds.Close()

	id, err := ds.UpsertUser(cl.Subject, cl.Name, cl.Picture)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.String(200, newJWT(id))
}

func getTokenFromRequest(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", nil // No error, just no token
	}

	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", fmt.Errorf("Authorization header format must be Bearer {token}")
	}

	return authHeaderParts[1], nil
}

func newJWT(id bson.ObjectId) string {
	signingKey := []byte(os.Getenv("SIGNING_KEY"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"id":  id,
	})

	ts, _ := token.SignedString(signingKey)

	return ts
}

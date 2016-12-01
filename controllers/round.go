package controllers

import (
	"github.com/davidamey/coffeeround-api/models"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

type RoundController interface {
	GetRounds(*gin.Context)
	GetRound(*gin.Context)
}

type roundController struct{}

func NewRoundController() RoundController {
	return &roundController{}
}

func (rc *roundController) GetRounds(c *gin.Context) {
	id := c.MustGet("userID").(bson.ObjectId)

	ds := models.GetDataStore()
	defer ds.Close()

	u, _ := ds.GetUser(id)

	c.JSON(200, u)
}

func (rc *roundController) GetRound(c *gin.Context) {

}

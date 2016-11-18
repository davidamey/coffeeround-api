package controllers

import (
	"net/http"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/davidamey/coffeeround-api/models"
	"github.com/gin-gonic/gin"
)

type UserController interface {
	GetUser(*gin.Context)
	CreateUser(*gin.Context)
}

type userController struct {
	db *mgo.Database
}

func NewUserController(db *mgo.Database) UserController {
	return &userController{db}
}

func (uc *userController) GetUser(c *gin.Context) {
	id := c.Param("id")

	if !bson.IsObjectIdHex(id) {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	oid := bson.ObjectIdHex(id)
	u := models.User{}

	if err := uc.db.C("users").FindId(oid).One(&u); err != nil {
		c.AbortWithError(http.StatusNotFound, err)
	}

	c.JSON(200, u)
}

func (uc *userController) CreateUser(c *gin.Context) {
	u := models.User{}

	if c.Bind(&u) != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	u.Id = bson.NewObjectId()

	uc.db.C("users").Insert(u)

	c.JSON(201, u)
}

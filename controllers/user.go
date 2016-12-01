package controllers

import (
	"net/http"

	"github.com/davidamey/coffeeround-api/models"
	"github.com/gin-gonic/gin"
)

type UserController interface {
	GetUsers(*gin.Context)
	GetUser(*gin.Context)
}

type userController struct{}

func NewUserController() UserController {
	return &userController{}
}

func (uc *userController) GetUsers(c *gin.Context) {
	ds := models.GetDataStore()
	defer ds.Close()

	users := ds.GetUsers()

	c.JSON(200, users)
}

func (uc *userController) GetUser(c *gin.Context) {
	id := c.Param("id")

	ds := models.GetDataStore()
	defer ds.Close()

	if user, found := ds.GetUser(id); found {
		c.JSON(200, user)
	} else {
		c.AbortWithStatus(http.StatusNotFound)
	}
}

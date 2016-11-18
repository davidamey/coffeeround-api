package controllers

import (
	"net/http"
	"strconv"

	"github.com/davidamey/coffeeround-api/models"
	"github.com/gin-gonic/gin"
)

type UserController interface {
	GetUser(*gin.Context)
}

type userController struct{}

func NewUserController() UserController {
	return &userController{}
}

func (uc *userController) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"Error": err})
		return
	}

	u := models.User{
		Id:   id,
		Name: "A. User",
	}

	c.JSON(200, u)
}

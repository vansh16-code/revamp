package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"proj/models"
)

func GetUsers(c *gin.Context) {
	users := []models.User{
		{ID: 1, Name: "Vansh"},
		{ID: 2, Name: "Mahajan"},
	}

	c.JSON(http.StatusOK, users)
}

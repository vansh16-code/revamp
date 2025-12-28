package handlers

import (
	"net/http"
	"proj/config"
	"proj/models"
	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context){
	var users []models.User
	
	if err := config.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	
	c.JSON(http.StatusOK, users)
}

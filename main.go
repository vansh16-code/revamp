package main

import (
	"github.com/gin-gonic/gin"
	"proj/config"
	"proj/models"
	"proj/routes"
)

func main(){

	config.ConnectDB()
	

	config.DB.AutoMigrate(&models.User{})

	r := gin.Default()

	routes.Register(r)

	r.Run(":8080")
}

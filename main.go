package main

import (
	"github.com/gin-gonic/gin"
	"proj/config"
	"proj/middleware"
	"proj/models"
	"proj/routes"
)

func main(){

	config.ConnectDB()
	

	config.DB.AutoMigrate(
		&models.User{},
		&models.Vehicle{},
		&models.Availability{},
		&models.Booking{},
		&models.OBDTracker{},
		&models.OBDReading{},
	)

	r := gin.Default()
	

	r.Use(middleware.CORS())      
	r.Use(middleware.Logger())    
	

	routes.Register(r)

	r.Run(":8080")
}

package routes

import (
	"proj/handlers"
	"proj/middleware"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	api := r.Group("/api")

	api.POST("/register", handlers.Register)
	api.POST("/login", handlers.Login)

	protected := api.Group("")
	protected.Use(middleware.AuthRequired())
	{
		protected.GET("/profile", handlers.GetProfile)
		protected.GET("/users", handlers.GetUsers)

		protected.POST("/vehicles", handlers.CreateVehicle)
		protected.GET("/my-vehicles", handlers.GetMyVehicles)
		protected.PUT("/vehicles/:id", handlers.UpdateVehicle)
		protected.DELETE("/vehicles/:id", handlers.DeleteVehicle)
	}

	api.GET("/vehicles", handlers.GetVehicles)
	api.GET("/vehicles/:id", handlers.GetVehicleByID)
}
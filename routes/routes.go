package routes

import (
	"github.com/gin-gonic/gin"
	"proj/handlers"
)

func Register(r *gin.Engine) {
	api := r.Group("/api")

	api.GET("/users", handlers.GetUsers)
}
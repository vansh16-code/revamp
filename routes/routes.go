package routes

import (
	"proj/handlers"
	"proj/middleware"
	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	api := r.Group("/api")

	api.GET("/users", handlers.GetUsers)
}
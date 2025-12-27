package main

import (
	"github.com/gin-gonic/gin"
	"proj/routes"
)

func main(){

	r := gin.Default()

	routes.Register(r)

	r.Run(":8080")
}

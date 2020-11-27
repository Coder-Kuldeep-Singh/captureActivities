package main

import (
	"os"

	"github.com/gin-gonic/gin"
)

//SetupRouter sets up routes
func SetupRouter() *gin.Engine {
	router := gin.Default()
	gin.SetMode(os.Getenv("GIN_MODE"))

	router.LoadHTMLGlob("html/*")

	router.GET("/days", GetDays)
	router.GET("/daily", Daily)
	return router
}

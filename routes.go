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

	router.GET("/", GetVisual)
	router.GET("/graphs", Daily)
	return router
}

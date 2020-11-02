package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Core router
	router := gin.Default()

	// Versioned API group
	v1 := router.Group("/v1")
	v1.GET("/ping", pingV1)

	// Serve on default port (8080)
	router.Run()
}

// route handler
func pingV1(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"message": "pong"})
}

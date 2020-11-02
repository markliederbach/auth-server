package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type server struct{}

func main() {
	router := gin.Default()

	v1 := router.Group("/v1")
	v1.GET("/ping", pingV1)

	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func pingV1(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"message": "pong"})
}

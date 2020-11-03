package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	controllerv1 "auth-server/pkg/v1/controller"
	middlewarev1 "auth-server/pkg/v1/middleware"
)

const (
	defaultLogLevel log.Level = log.InfoLevel
)

func main() {

	// Setup logger
	configureLogger()

	// Core router
	router := gin.Default()

	// Versioned API group
	registerV1Routes(router)

	// Serve on default port (8080)
	router.Run()
}

func configureLogger() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	logLevelRaw := os.Getenv("LOG_LEVEL")

	var logLevel log.Level = defaultLogLevel
	var err error
	if logLevelRaw != "" {
		logLevel, err = log.ParseLevel(logLevelRaw)
		if err != nil {
			panic(err)
		}
	}

	log.SetLevel(logLevel)
}

func registerV1Routes(router *gin.Engine) {
	v1 := router.Group("/v1")

	// Add a test authorized endpoint
	testAuth := v1.Group("/test")
	testAuth.Use(middlewarev1.AuthorizeToken())
	testAuth.GET("/ping", pingV1)

	controllerv1.NewLoginController(v1)
}

// route handler
func pingV1(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"message": "pong"})
}

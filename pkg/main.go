package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"auth-server/pkg/config"
	controllerv1 "auth-server/pkg/v1/controller"
	middlewarev1 "auth-server/pkg/v1/middleware"
	tokenservicev1 "auth-server/pkg/v1/service"
)

const (
	defaultLogLevel log.Level = log.InfoLevel
)

func main() {

	appConfig := config.Load()

	// Core router
	router := gin.New()
	router.Use(middlewarev1.GinLogger(), gin.Recovery())

	// Versioned API group
	registerV1Routes(appConfig, router)

	// Serve on default port (8080)
	router.Run()
}

func registerV1Routes(config config.Config, router *gin.Engine) {
	v1 := router.Group("/v1")

	jwtServiceV1 := tokenservicev1.NewJWTService(config)

	// Add a test authorized endpoint
	testAuth := v1.Group("/test")
	testAuth.Use(middlewarev1.AuthorizeToken(jwtServiceV1))
	testAuth.GET("/ping", pingV1)

	controllerv1.NewLoginController(v1, jwtServiceV1)
	controllerv1.NewTokenController(v1, jwtServiceV1)
	controllerv1.NewLogoutController(v1, jwtServiceV1)
}

// route handler
func pingV1(context *gin.Context) {
	jwtUser, _ := context.MustGet("user").(tokenservicev1.JWTUser)
	context.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Hello %s!", jwtUser.Username)})
}

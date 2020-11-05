package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	tokenservice "auth-server/pkg/v1/service"
)

const (
	loginRoute string = "/login"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginController struct {
	log        *log.Entry
	group      *gin.RouterGroup
	jwtService tokenservice.JWTService
}

func NewLoginController(group *gin.RouterGroup, jwtService tokenservice.JWTService) *LoginController {
	loginController := &LoginController{
		log:        log.WithFields(log.Fields{"logger": "LoginControllerV1"}),
		group:      group,
		jwtService: jwtService,
	}
	loginController.registerRoutes()
	return loginController
}

func (c *LoginController) registerRoutes() {
	// c.log.Info("Registering routes")
	c.group.POST(loginRoute, c.Login)
}

func (c *LoginController) Login(context *gin.Context) {
	var request LoginRequest
	requestLogger, _ := context.MustGet("request_logger").(*log.Entry)
	requestLogger.Info("Handling request")

	if err := context.ShouldBindJSON(&request); err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Authenticate User

	// Generate JWTs
	accessToken, refreshToken, err := c.jwtService.GenerateToken(tokenservice.JWTUser{Username: request.Username}, true)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	context.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

package controller

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	tokenservice "auth-server/pkg/v1/token"

	"github.com/gin-gonic/gin"
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
	contextLog := log.WithFields(
		log.Fields{
			"logger":    "LoginController",
			"base_path": group.BasePath(),
		},
	)
	loginController := &LoginController{
		log:        contextLog,
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
	logger, _ := context.MustGet("logger").(*log.Entry)

	logger.Info("Handling login request")

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

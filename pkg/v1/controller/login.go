package controller

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	tokenservice "auth-server/pkg/v1/token"

	"github.com/gin-gonic/gin"
)

var (
	_ Controller = &LoginController{}
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

func NewLoginController(group *gin.RouterGroup) *LoginController {

	jwtService := tokenservice.NewJWTService()

	contextLog := log.WithFields(
		log.Fields{
			"controller": "LoginController",
			"base_path":  group.BasePath(),
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
	c.log.Info("Registering routes")
	c.group.POST(loginRoute, c.Login)
}

func (c *LoginController) Login(context *gin.Context) {
	var request LoginRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Authenticate User - TODO

	// Generate JWT
	token, err := c.jwtService.GenerateToken(
		tokenservice.JWTUser{
			Username: request.Username,
		},
	)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	context.JSON(http.StatusOK, gin.H{
		"access_token": token,
	})
}

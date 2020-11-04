package controller

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	tokenservice "auth-server/pkg/v1/token"

	"github.com/gin-gonic/gin"
)

const (
	logoutRoute string = "/logout"
)

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LogoutController struct {
	log        *log.Entry
	group      *gin.RouterGroup
	jwtService tokenservice.JWTService
}

func NewLogoutController(group *gin.RouterGroup, jwtService tokenservice.JWTService) *LogoutController {
	contextLog := log.WithFields(
		log.Fields{
			"logger":    "LogoutController",
			"base_path": group.BasePath(),
		},
	)
	logoutController := &LogoutController{
		log:        contextLog,
		group:      group,
		jwtService: jwtService,
	}
	logoutController.registerRoutes()
	return logoutController
}

func (c *LogoutController) registerRoutes() {
	// c.log.Info("Registering routes")
	c.group.DELETE(logoutRoute, c.Logout)
}

func (c *LogoutController) Logout(context *gin.Context) {
	var request LogoutRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate JWTs
	c.jwtService.RemoveRefreshToken(request.RefreshToken)
	context.Status(http.StatusNoContent)
}

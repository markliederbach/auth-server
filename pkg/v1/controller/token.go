package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	tokenservice "auth-server/pkg/v1/service"
)

const (
	tokenRoute string = "/token"
)

type TokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type TokenController struct {
	log        *log.Entry
	group      *gin.RouterGroup
	jwtService tokenservice.JWTService
}

func NewTokenController(group *gin.RouterGroup, jwtService tokenservice.JWTService) *TokenController {
	loginController := &TokenController{
		log:        log.WithFields(log.Fields{"logger": "TokenControllerV1"}),
		group:      group,
		jwtService: jwtService,
	}
	loginController.registerRoutes()
	return loginController
}

func (c *TokenController) registerRoutes() {
	// c.log.Info("Registering routes")
	c.group.POST(tokenRoute, c.Token)
}

func (c *TokenController) Token(context *gin.Context) {
	var request TokenRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Lookup refresh token to make sure it's valid
	refreshToken, authClaims, err := c.jwtService.ValidateRefreshToken(request.RefreshToken)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	if !refreshToken.Valid {
		context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Generate new JWT
	accessToken, _, err := c.jwtService.GenerateToken(authClaims.User, false)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
	})
}

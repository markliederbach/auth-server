package controller

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"

	tokenservice "auth-server/pkg/v1/token"

	"github.com/gin-gonic/gin"
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
	contextLog := log.WithFields(
		log.Fields{
			"logger":    "TokenController",
			"base_path": group.BasePath(),
		},
	)
	loginController := &TokenController{
		log:        contextLog,
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
	refreshToken, err := c.jwtService.ValidateRefreshToken(request.RefreshToken)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	if !refreshToken.Valid {
		context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Load claims to find the username
	claims, ok := refreshToken.Claims.(jwt.MapClaims)
	if !ok {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Unparsable claims"})
		return
	}

	username, ok := claims["sub"].(string)
	if !ok {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Missing subject claim"})
		return
	}

	// Generate new JWT
	accessToken, _, err := c.jwtService.GenerateToken(tokenservice.JWTUser{Username: username}, false)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
	})
}

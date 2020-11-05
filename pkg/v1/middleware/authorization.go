package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	tokenservice "auth-server/pkg/v1/service"
)

const (
	authorizationHeader string = "Authorization"
)

// AuthorizeToken checks that a JWT token is valid and attaches the corresponding JWTUser to the context
func AuthorizeToken(jwtService tokenservice.JWTService) gin.HandlerFunc {
	return func(context *gin.Context) {
		authHeader := context.GetHeader(authorizationHeader)
		if authHeader == "" {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Missing %s header", authorizationHeader)})
			return
		}

		splits := strings.Split(authHeader, " ")
		if len(splits) < 2 {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			return
		}

		tokenString := splits[1]
		token, authClaims, err := jwtService.ValidateAccessToken(tokenString)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		if !token.Valid {
			context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid access token"})
			return
		}
		context.Set("user", authClaims.User)

		// TODO: Optionally check other fields on a user, like roles
	}
}

package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"

	"auth-server/pkg/config"
	"auth-server/pkg/utils"
)

type JWTService interface {
	GenerateToken(user JWTUser, generateRefreshToken bool) (string, string, error)
	ValidateAccessToken(encodedToken string) (*jwt.Token, *AuthCustomClaims, error)
	ValidateRefreshToken(encodedToken string) (*jwt.Token, *AuthCustomClaims, error)
	RemoveRefreshToken(encodedToken string)
}

type JWTUser struct {
	Username string
}

type AuthCustomClaims struct {
	jwt.StandardClaims
	User JWTUser
}

type jwtService struct {
	config             config.Config
	validRefreshTokens []string
}

func NewJWTService(config config.Config) JWTService {
	return &jwtService{
		config: config,
		// TODO: move list to DB
		validRefreshTokens: []string{},
	}
}

func (s *jwtService) GenerateToken(user JWTUser, generateRefreshToken bool) (string, string, error) {
	now := time.Now()

	// Access token, including expiration date
	accessClaims := &AuthCustomClaims{
		StandardClaims: jwt.StandardClaims{
			Subject: user.Username,

			ExpiresAt: now.Add(s.config.AccessTokenExpire).Unix(),
			Issuer:    s.config.Issuer,
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
		},
		User: user,
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.config.AccessTokenSecret))
	if err != nil {
		return "", "", err
	}

	if !generateRefreshToken {
		// If we don't care about refresh tokens, we're done
		return accessTokenString, "", nil
	}

	// Refresh token, including extended expiration date
	refreshClaims := &AuthCustomClaims{
		StandardClaims: jwt.StandardClaims{
			Subject: user.Username,

			ExpiresAt: now.Add(s.config.RefreshTokenExpire).Unix(),
			Issuer:    s.config.Issuer,
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
		},
		User: user,
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.config.RefreshTokenSecret))
	if err != nil {
		return "", "", err
	}

	// Store refresh token
	s.validRefreshTokens = append(s.validRefreshTokens, refreshTokenString)

	return accessTokenString, refreshTokenString, nil
}

func (s *jwtService) ValidateAccessToken(encodedToken string) (*jwt.Token, *AuthCustomClaims, error) {
	token, authClaims, err := validateToken(encodedToken, s.config.AccessTokenSecret)
	if err != nil {
		return nil, nil, err
	}

	return token, authClaims, err
}

func (s *jwtService) ValidateRefreshToken(encodedToken string) (*jwt.Token, *AuthCustomClaims, error) {
	tokenIndex := utils.IndexOf(s.validRefreshTokens, encodedToken)
	if utils.IndexOf(s.validRefreshTokens, encodedToken) == -1 {
		return nil, nil, errors.New("Invalid refresh token")
	}
	token, authClaims, err := validateToken(encodedToken, s.config.RefreshTokenSecret)
	if err != nil {
		// Cleanup after ourselves. This is not thread-safe, FYI
		s.validRefreshTokens = utils.RemoveIndex(s.validRefreshTokens, tokenIndex)
		return nil, nil, err
	}
	return token, authClaims, nil
}

func (s *jwtService) RemoveRefreshToken(encodedToken string) {
	tokenIndex := utils.IndexOf(s.validRefreshTokens, encodedToken)
	if utils.IndexOf(s.validRefreshTokens, encodedToken) == -1 {
		return
	}
	s.validRefreshTokens = utils.RemoveIndex(s.validRefreshTokens, tokenIndex)
}

func validateToken(encodedToken string, tokenSecret string) (*jwt.Token, *AuthCustomClaims, error) {
	authClaims := &AuthCustomClaims{}
	token, err := jwt.ParseWithClaims(encodedToken, authClaims, func(token *jwt.Token) (interface{}, error) {
		if _, valid := token.Method.(*jwt.SigningMethodHMAC); !valid {
			return nil, fmt.Errorf("Invalid token algorithm %v", token.Header["alg"])
		}
		return []byte(tokenSecret), nil
	})
	return token, authClaims, err
}

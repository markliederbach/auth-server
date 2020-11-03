package token

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	accessTokenVariable        string = "ACCESS_TOKEN_SECRET"
	refreshTokenVariable       string = "REFRESH_TOKEN_SECRET"
	accessTokenExpireVariable  string = "ACCESS_TOKEN_EXPIRE"
	refreshTokenExpireVariable string = "REFRESH_TOKEN_EXPIRE"

	issuer string = "markliederbach/auth-service"

	defaultAccessTokenExpire  time.Duration = time.Second * 15
	defaultRefreshTokenExpire time.Duration = time.Second * 1
)

type JWTService interface {
	GenerateToken(user JWTUser, generateRefreshToken bool) (string, string, error)
	ValidateAccessToken(encodedToken string) (*jwt.Token, error)
	ValidateRefreshToken(encodedToken string) (*jwt.Token, error)
}

type JWTUser struct {
	Username string
}

type authCustomClaims struct {
	jwt.StandardClaims
}

type jwtService struct {
	accessTokenSecret  string
	refreshTokenSecret string

	accessTokenExpire  time.Duration
	refreshTokenExpire time.Duration

	issuer             string
	validRefreshTokens []string
}

func NewJWTService() JWTService {
	var err error
	accessTokenSecret := os.Getenv(accessTokenVariable)
	if accessTokenSecret == "" {
		panic(fmt.Errorf("Please set %s", accessTokenVariable))
	}

	refreshTokenSecret := os.Getenv(refreshTokenVariable)
	if refreshTokenSecret == "" {
		panic(fmt.Errorf("Please set %s", refreshTokenVariable))
	}

	var accessTokenExpire time.Duration
	accessTokenExpireRaw := os.Getenv(accessTokenExpireVariable)
	if accessTokenExpireRaw == "" {
		accessTokenExpire = defaultAccessTokenExpire
	} else {
		accessTokenExpire, err = time.ParseDuration(accessTokenExpireRaw)
		if err != nil {
			panic(err)
		}
	}

	var refreshTokenExpire time.Duration
	refreshTokenExpireRaw := os.Getenv(accessTokenExpireVariable)
	if refreshTokenExpireRaw == "" {
		refreshTokenExpire = defaultRefreshTokenExpire
	} else {
		refreshTokenExpire, err = time.ParseDuration(refreshTokenExpireRaw)
		if err != nil {
			panic(err)
		}
	}

	return &jwtService{
		accessTokenSecret:  accessTokenSecret,
		refreshTokenSecret: refreshTokenSecret,
		accessTokenExpire:  accessTokenExpire,
		refreshTokenExpire: refreshTokenExpire,

		issuer: issuer,
		// TODO: move list to DB
		validRefreshTokens: []string{},
	}
}

func (s *jwtService) GenerateToken(user JWTUser, generateRefreshToken bool) (string, string, error) {
	now := time.Now()

	// Access token, including expiration date
	accessClaims := &authCustomClaims{
		jwt.StandardClaims{
			Subject: user.Username,

			ExpiresAt: now.Add(s.accessTokenExpire).Unix(),
			Issuer:    s.issuer,
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.accessTokenSecret))
	if err != nil {
		return "", "", err
	}

	if !generateRefreshToken {
		// If we don't care about refresh tokens, we're done
		return accessTokenString, "", nil
	}

	// Refresh token, including extended expiration date
	refreshClaims := &authCustomClaims{
		jwt.StandardClaims{
			Subject: user.Username,

			ExpiresAt: now.Add(s.refreshTokenExpire).Unix(),
			Issuer:    s.issuer,
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.refreshTokenSecret))
	if err != nil {
		return "", "", err
	}

	// Store refresh token
	s.validRefreshTokens = append(s.validRefreshTokens, refreshTokenString)

	return accessTokenString, refreshTokenString, nil
}

func (s *jwtService) ValidateAccessToken(encodedToken string) (*jwt.Token, error) {
	return validateToken(encodedToken, s.accessTokenSecret)
}

func (s *jwtService) ValidateRefreshToken(encodedToken string) (*jwt.Token, error) {
	tokenIndex := indexOf(s.validRefreshTokens, encodedToken)
	if indexOf(s.validRefreshTokens, encodedToken) == -1 {
		return nil, errors.New("Invalid refresh token")
	}
	token, err := validateToken(encodedToken, s.refreshTokenSecret)
	if err != nil {
		// Cleanup after ourselves. This is not thread-safe, FYI
		s.validRefreshTokens = removeIndex(s.validRefreshTokens, tokenIndex)
		return nil, err
	}
	return token, nil
}

func validateToken(encodedToken string, tokenSecret string) (*jwt.Token, error) {
	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, valid := token.Method.(*jwt.SigningMethodHMAC); !valid {
			return nil, fmt.Errorf("Invalid token algorithm %v", token.Header["alg"])
		}
		return []byte(tokenSecret), nil
	})
}

func indexOf(items []string, value string) int {
	for index, item := range items {
		if value == item {
			return index
		}
	}
	return -1
}

func removeIndex(items []string, index int) []string {
	return append(items[:index], items[index+1:]...)
}

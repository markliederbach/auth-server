package token

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	accessTokenVariable  string = "ACCESS_TOKEN_SECRET"
	refreshTokenVariable string = "REFRESH_TOKEN_SECRET"

	issuer string = "markliederbach/auth-service"

	accessTokenExpire time.Duration = time.Second * 15
)

type JWTService interface {
	GenerateToken(user JWTUser) (string, error)
	ValidateToken(encodedToken string) (*jwt.Token, error)
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
	issuer             string
}

func NewJWTService() JWTService {
	accessTokenSecret := os.Getenv(accessTokenVariable)
	if accessTokenSecret == "" {
		panic(fmt.Errorf("Please set %s", accessTokenVariable))
	}

	refreshTokenSecret := os.Getenv(refreshTokenVariable)
	if refreshTokenSecret == "" {
		panic(fmt.Errorf("Please set %s", refreshTokenVariable))
	}

	return &jwtService{
		accessTokenSecret:  accessTokenSecret,
		refreshTokenSecret: refreshTokenSecret,
		issuer:             issuer,
	}
}

func (s *jwtService) GenerateToken(user JWTUser) (string, error) {
	now := time.Now()
	claims := &authCustomClaims{
		jwt.StandardClaims{
			Subject: user.Username,

			ExpiresAt: now.Add(accessTokenExpire).Unix(),
			Issuer:    s.issuer,
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.accessTokenSecret))
}

func (s *jwtService) ValidateToken(encodedToken string) (*jwt.Token, error) {
	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, valid := token.Method.(*jwt.SigningMethodHMAC); !valid {
			return nil, fmt.Errorf("Invalid token algorithm %v", token.Header["alg"])
		}
		return []byte(s.accessTokenSecret), nil
	})
}

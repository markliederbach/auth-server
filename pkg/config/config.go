package config

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	logLevelVariable           string = "LOG_LEVEL"
	accessTokenVariable        string = "ACCESS_TOKEN_SECRET"
	refreshTokenVariable       string = "REFRESH_TOKEN_SECRET"
	accessTokenExpireVariable  string = "ACCESS_TOKEN_EXPIRE"
	refreshTokenExpireVariable string = "REFRESH_TOKEN_EXPIRE"

	// Default values
	defaultAccessTokenExpire  time.Duration = time.Second * 15
	defaultRefreshTokenExpire time.Duration = time.Minute * 1
	defaultLogLevel           log.Level     = log.InfoLevel
	defaultIssuer             string        = "markliederbach/auth-service"
)

// Config holds all configuration data about the currently-running service
type Config struct {
	// Required variables
	AccessTokenSecret  string
	RefreshTokenSecret string

	// Optional variables
	LogLevel           log.Level
	AccessTokenExpire  time.Duration
	RefreshTokenExpire time.Duration
	Issuer             string
}

// Load creates a new instance of Config, using all available
// defaults and overrides.
func Load() Config {
	config := Config{
		AccessTokenSecret:  fromEnvString(accessTokenVariable, true, ""),
		RefreshTokenSecret: fromEnvString(refreshTokenVariable, true, ""),

		LogLevel:           fromEnvLogLevel(logLevelVariable, false, defaultLogLevel),
		AccessTokenExpire:  fromEnvDuration(accessTokenExpireVariable, false, defaultAccessTokenExpire),
		RefreshTokenExpire: fromEnvDuration(refreshTokenExpireVariable, false, defaultRefreshTokenExpire),
		Issuer:             fromEnvString(refreshTokenExpireVariable, false, defaultIssuer),
	}

	config.configureLogger()

	return config
}

func (c *Config) configureLogger() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	gin.SetMode(gin.ReleaseMode)

	switch c.LogLevel {
	case log.TraceLevel:
		gin.SetMode(gin.DebugMode)
		log.SetReportCaller(true)
	}

	log.SetLevel(c.LogLevel)
}

func fromEnvString(variable string, required bool, defaultValue string) string {
	rawValue, exists := fromEnv(variable, required)
	if !exists {
		rawValue = defaultValue
	}
	return rawValue
}

func fromEnvDuration(variable string, required bool, defaultValue time.Duration) time.Duration {
	var err error
	value := defaultValue
	rawValue, exists := fromEnv(variable, required)
	if exists {
		value, err = time.ParseDuration(rawValue)
		if err != nil {
			panic(err)
		}
	}
	return value
}

func fromEnvLogLevel(variable string, required bool, defaultValue log.Level) log.Level {
	var err error
	value := defaultValue
	rawValue, exists := fromEnv(variable, required)
	if exists {
		value, err = log.ParseLevel(rawValue)
		if err != nil {
			panic(err)
		}
	}
	return value
}

func fromEnv(variable string, required bool) (string, bool) {
	value, exists := os.LookupEnv(variable)
	if !exists && required {
		panic(fmt.Errorf("Missing required environment variable %s", variable))
	}
	return value, exists
}

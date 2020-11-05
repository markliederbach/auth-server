package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// GinLogger creates a contextual logger used for request/response information. The logger is
// also added to the request context, should any downstream controllers wish to use it.
func GinLogger() gin.HandlerFunc {
	logger := log.StandardLogger()

	return func(context *gin.Context) {
		// Create a contextual logger for downstream controllers
		contextLogger := logger.WithFields(
			log.Fields{
				"method":     context.Request.Method,
				"uri":        context.Request.RequestURI,
				"referer":    context.Request.Referer(),
				"source_ip":  context.ClientIP(),
				"user_agent": context.Request.UserAgent(),
			},
		)

		context.Set("request_logger", contextLogger)

		startTime := time.Now()

		// Fulfill request
		context.Next()

		dataLength := context.Writer.Size()
		if dataLength < 0 {
			dataLength = 0
		}
		statusCode := context.Writer.Status()

		// Add additional fields about the response
		contextLogger = contextLogger.WithFields(
			log.Fields{
				"logger":              "RequestLogger",
				"status":              statusCode,
				"response_data_bytes": dataLength,
				"latency_ns":          time.Since(startTime),
			},
		)

		if len(context.Errors) > 0 {
			contextLogger.Error(context.Errors.ByType(gin.ErrorTypePrivate).String())
			return
		}

		message := "Request finished"
		if statusCode >= http.StatusInternalServerError {
			contextLogger.Error(message)
		} else if statusCode >= http.StatusBadRequest {
			contextLogger.Warn(message)
		} else {
			contextLogger.Info(message)
		}
	}
}

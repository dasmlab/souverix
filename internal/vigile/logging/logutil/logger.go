package logging

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// InitLogger initializes a logrus logger with structured formatting
func InitLogger(component string) *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
	})
	log.SetLevel(logrus.InfoLevel)

	// Set log level from environment
	if level := getLogLevel(); level != "" {
		if parsed, err := logrus.ParseLevel(level); err == nil {
			log.SetLevel(parsed)
		}
	}

	log.WithField("component", component).Info("logger initialized")
	return log
}

// GinLogger returns a Gin middleware that logs requests using logrus
func GinLogger(log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log request
		latency := time.Since(start)
		status := c.Writer.Status()

		entry := log.WithFields(logrus.Fields{
			"status":     status,
			"method":     c.Request.Method,
			"path":       path,
			"query":      raw,
			"ip":         c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
			"latency":    latency,
		})

		if status >= 500 {
			entry.Error("HTTP request")
		} else if status >= 400 {
			entry.Warn("HTTP request")
		} else {
			entry.Info("HTTP request")
		}
	}
}

func getLogLevel() string {
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		return "info"
	}
	return level
}

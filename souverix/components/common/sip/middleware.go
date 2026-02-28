package sip

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// LoggingMiddleware logs SIP messages
func LoggingMiddleware(logger *logrus.Logger) MiddlewareFunc {
	return func(ctx context.Context, msg *Message, c *gin.Context) (context.Context, error) {
		if msg.IsRequest() {
			logger.Infof("SIP Request: %s %s (Call-ID: %s)", msg.Method, msg.URI, msg.GetCallID())
		} else {
			logger.Infof("SIP Response: %d %s (Call-ID: %s)", msg.StatusCode, msg.StatusText, msg.GetCallID())
		}
		return ctx, nil
	}
}

// ValidateHeadersMiddleware validates required SIP headers
func ValidateHeadersMiddleware(requiredHeaders []string) MiddlewareFunc {
	return func(ctx context.Context, msg *Message, c *gin.Context) (context.Context, error) {
		for _, header := range requiredHeaders {
			if msg.GetHeader(header) == "" {
				return ctx, fmt.Errorf("missing required header: %s", header)
			}
		}
		return ctx, nil
	}
}

// RecordRouteMiddleware adds Record-Route header for proxies
func RecordRouteMiddleware(componentName, host string, port int) MiddlewareFunc {
	return func(ctx context.Context, msg *Message, c *gin.Context) (context.Context, error) {
		if msg.IsRequest() {
			recordRoute := fmt.Sprintf("<sip:%s:%d;lr>", host, port)
			msg.AddHeader("Record-Route", recordRoute)
		}
		return ctx, nil
	}
}

// ViaMiddleware adds Via header
func ViaMiddleware(componentName, host string, port int) MiddlewareFunc {
	return func(ctx context.Context, msg *Message, c *gin.Context) (context.Context, error) {
		if msg.IsRequest() {
			via := fmt.Sprintf("SIP/2.0/UDP %s:%d;branch=z9hG4bK-%s", host, port, componentName)
			msg.AddHeader("Via", via)
		}
		return ctx, nil
	}
}

// ContactMiddleware sets Contact header
func ContactMiddleware(contact string) MiddlewareFunc {
	return func(ctx context.Context, msg *Message, c *gin.Context) (context.Context, error) {
		if msg.IsRequest() {
			msg.SetHeader("Contact", contact)
		}
		return ctx, nil
	}
}

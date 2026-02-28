package sip

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// HandlerFunc is a function that handles a SIP message
type HandlerFunc func(ctx context.Context, msg *Message, c *gin.Context) (*Message, error)

// SIPHandler handles SIP messages for a component
type SIPHandler struct {
	logger      *logrus.Logger
	handlers    map[string]HandlerFunc // Method -> handler
	middleware  []MiddlewareFunc
	errorHandler ErrorHandlerFunc
}

// MiddlewareFunc is a middleware function for SIP processing
type MiddlewareFunc func(ctx context.Context, msg *Message, c *gin.Context) (context.Context, error)

// ErrorHandlerFunc handles errors during SIP processing
type ErrorHandlerFunc func(ctx context.Context, err error, c *gin.Context)

// NewSIPHandler creates a new SIP handler
func NewSIPHandler(logger *logrus.Logger) *SIPHandler {
	return &SIPHandler{
		logger:      logger,
		handlers:    make(map[string]HandlerFunc),
		middleware:  []MiddlewareFunc{},
		errorHandler: defaultErrorHandler,
	}
}

// RegisterHandler registers a handler for a SIP method
func (h *SIPHandler) RegisterHandler(method string, handler HandlerFunc) {
	h.handlers[strings.ToUpper(method)] = handler
}

// Use adds middleware to the handler chain
func (h *SIPHandler) Use(middleware MiddlewareFunc) {
	h.middleware = append(h.middleware, middleware)
}

// SetErrorHandler sets a custom error handler
func (h *SIPHandler) SetErrorHandler(handler ErrorHandlerFunc) {
	h.errorHandler = handler
}

// HandleRequest handles an incoming SIP request
func (h *SIPHandler) HandleRequest(c *gin.Context) {
	// Parse SIP message from request body
	rawBody, err := c.GetRawData()
	if err != nil {
		h.errorHandler(c.Request.Context(), fmt.Errorf("failed to read request body: %w", err), c)
		return
	}

	msg, err := ParseMessage(string(rawBody))
	if err != nil {
		h.errorHandler(c.Request.Context(), fmt.Errorf("failed to parse SIP message: %w", err), c)
		return
	}

	if !msg.IsRequest() {
		h.errorHandler(c.Request.Context(), fmt.Errorf("expected SIP request, got response"), c)
		return
	}

	ctx := c.Request.Context()

	// Apply middleware
	for _, mw := range h.middleware {
		var err error
		ctx, err = mw(ctx, msg, c)
		if err != nil {
			h.errorHandler(ctx, err, c)
			return
		}
	}

	// Find and execute handler
	handler, exists := h.handlers[strings.ToUpper(msg.Method)]
	if !exists {
		h.logger.Warnf("No handler registered for method: %s", msg.Method)
		resp := CreateResponse(405, "Method Not Allowed", msg)
		c.Data(http.StatusMethodNotAllowed, "application/sip", []byte(resp.String()))
		return
	}

	// Execute handler
	response, err := handler(ctx, msg, c)
	if err != nil {
		h.errorHandler(ctx, err, c)
		return
	}

	// Send response
	if response != nil {
		c.Data(http.StatusOK, "application/sip", []byte(response.String()))
	} else {
		// No response means handler will send it manually
	}
}

// RegisterRoutes registers SIP routes on a Gin router
func (h *SIPHandler) RegisterRoutes(router *gin.Engine, path string) {
	router.POST(path, h.HandleRequest)
	router.ANY(path, h.HandleRequest) // Also handle other methods
}

// defaultErrorHandler is the default error handler
func defaultErrorHandler(ctx context.Context, err error, c *gin.Context) {
	// Try to create a 500 error response
	rawBody, _ := c.GetRawData()
	if len(rawBody) > 0 {
		if msg, parseErr := ParseMessage(string(rawBody)); parseErr == nil {
			resp := CreateResponse(500, "Internal Server Error", msg)
			c.Data(http.StatusInternalServerError, "application/sip", []byte(resp.String()))
			return
		}
	}

	// Fallback to plain text error
	c.String(http.StatusInternalServerError, "SIP Error: %s", err.Error())
}

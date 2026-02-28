package sip

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Component represents a SIP component with handlers
type Component struct {
	Name        string
	Host        string
	Port        int
	Handler     *SIPHandler
	Router      *gin.Engine
	Logger      *logrus.Logger
	Middlewares []MiddlewareFunc
}

// NewComponent creates a new SIP component
func NewComponent(name, host string, port int, logger *logrus.Logger) *Component {
	handler := NewSIPHandler(logger)
	router := gin.New()
	router.Use(gin.Recovery())

	return &Component{
		Name:        name,
		Host:        host,
		Port:        port,
		Handler:     handler,
		Router:      router,
		Logger:      logger,
		Middlewares: []MiddlewareFunc{},
	}
}

// RegisterMethod registers a handler for a SIP method
func (c *Component) RegisterMethod(method string, handler HandlerFunc) {
	c.Handler.RegisterHandler(method, handler)
}

// Use adds middleware to the component
func (c *Component) Use(middleware MiddlewareFunc) {
	c.Middlewares = append(c.Middlewares, middleware)
	c.Handler.Use(middleware)
}

// SetupDefaultMiddleware sets up default middleware for the component
func (c *Component) SetupDefaultMiddleware() {
	// Logging
	c.Use(LoggingMiddleware(c.Logger))

	// Validate required headers
	c.Use(ValidateHeadersMiddleware([]string{"Call-ID", "CSeq", "From", "To"}))

	// Add Via header
	c.Use(ViaMiddleware(c.Name, c.Host, c.Port))

	// Add Record-Route for proxies
	if c.Name == "pcscf" || c.Name == "icscf" || c.Name == "scscf" {
		c.Use(RecordRouteMiddleware(c.Name, c.Host, c.Port))
	}
}

// RegisterRoutes registers SIP routes on the router
func (c *Component) RegisterRoutes(path string) {
	c.Handler.RegisterRoutes(c.Router, path)
}

// Forward forwards a SIP message to another component
func (c *Component) Forward(ctx context.Context, msg *Message, targetHost string, targetPort int) (*Message, error) {
	url := fmt.Sprintf("http://%s:%d/sip", targetHost, targetPort)
	
	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(msg.String()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/sip")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response
	responseMsg, err := ParseMessage(string(body))
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return responseMsg, nil
}

// ProxyRequest proxies a SIP request to another component
func (c *Component) ProxyRequest(ctx context.Context, msg *Message, targetHost string, targetPort int) (*Message, error) {
	// Add Record-Route if not present
	if msg.GetHeader("Record-Route") == "" {
		recordRoute := fmt.Sprintf("<sip:%s:%d;lr>", c.Host, c.Port)
		msg.AddHeader("Record-Route", recordRoute)
	}

	return c.Forward(ctx, msg, targetHost, targetPort)
}

// Start starts the SIP component server
func (c *Component) Start() error {
	addr := fmt.Sprintf(":%d", c.Port)
	c.Logger.Infof("Starting SIP component %s on %s", c.Name, addr)
	return http.ListenAndServe(addr, c.Router)
}

// StartAsync starts the SIP component server in a goroutine
func (c *Component) StartAsync() error {
	go func() {
		if err := c.Start(); err != nil {
			c.Logger.WithError(err).Fatal("Failed to start SIP component")
		}
	}()
	return nil
}

package diagnostics

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// FauxComponentServer simulates other components for unit testing
type FauxComponentServer struct {
	registry      *CallFlowRegistry
	logger        *logrus.Logger
	servers       map[int]*http.Server // port -> server
	componentMap  map[string]int       // component name -> port
	portMap       map[int]string       // port -> component name
	mu            sync.RWMutex
	basePort      int // Starting port (e.g., 19000)
	callerIP      string
	activeFlowID  string
	activeContext *FauxCallContext
}

// FauxCallContext tracks the current call flow execution
type FauxCallContext struct {
	FlowID      string
	Component   string // Component under test
	Step        int    // Current step
	CallID      string
	From        string
	To          string
	State       map[string]interface{}
}

// NewFauxComponentServer creates a new faux component server
func NewFauxComponentServer(registry *CallFlowRegistry, logger *logrus.Logger, basePort int, callerIP string) *FauxComponentServer {
	return &FauxComponentServer{
		registry:     registry,
		logger:       logger,
		servers:      make(map[int]*http.Server),
		componentMap: make(map[string]int),
		portMap:      make(map[int]string),
		basePort:     basePort,
		callerIP:     callerIP,
	}
}

// Start starts the faux component server on specified ports
func (f *FauxComponentServer) Start(ctx context.Context, flowID string, componentUnderTest string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.activeFlowID = flowID
	f.activeContext = &FauxCallContext{
		FlowID:    flowID,
		Component: componentUnderTest,
		State:     make(map[string]interface{}),
	}

	// Get all components in the flow
	flow, exists := f.registry.GetFlow(flowID)
	if !exists {
		return fmt.Errorf("flow %s not found", flowID)
	}

	// Map components to ports
	components := f.getComponentsInFlow(flow)
	port := f.basePort

	for _, compName := range components {
		// Skip the component under test
		if compName == componentUnderTest {
			continue
		}

		// Assign port
		f.componentMap[compName] = port
		f.portMap[port] = compName

		// Start server for this component
		if err := f.startComponentServer(port, compName, ctx); err != nil {
			return fmt.Errorf("failed to start server for %s on port %d: %w", compName, port, err)
		}

		f.logger.Debugf("Started faux server for %s on port %d", compName, port)
		port++
	}

	return nil
}

// Stop stops all faux component servers
func (f *FauxComponentServer) Stop(ctx context.Context) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	for port, server := range f.servers {
		if err := server.Shutdown(ctx); err != nil {
			f.logger.WithError(err).Warnf("Error shutting down server on port %d", port)
		}
		delete(f.servers, port)
	}

	f.componentMap = make(map[string]int)
	f.portMap = make(map[int]string)
	f.activeContext = nil

	return nil
}

// GetComponentPort returns the port for a component
func (f *FauxComponentServer) GetComponentPort(componentName string) (int, bool) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	port, exists := f.componentMap[componentName]
	return port, exists
}

// GetComponentURL returns the URL for a component
func (f *FauxComponentServer) GetComponentURL(componentName string) (string, bool) {
	port, exists := f.GetComponentPort(componentName)
	if !exists {
		return "", false
	}
	return fmt.Sprintf("http://%s:%d", f.callerIP, port), true
}

// startComponentServer starts an HTTP server for a faux component
func (f *FauxComponentServer) startComponentServer(port int, componentName string, ctx context.Context) error {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	// Register handler for SIP messages
	router.POST("/sip", f.handleSIPRequest(componentName))
	router.Any("/sip", f.handleSIPRequest(componentName))

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	f.servers[port] = server

	// Start server in goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			f.logger.WithError(err).Errorf("Faux server on port %d failed", port)
		}
	}()

	return nil
}

// handleSIPRequest handles incoming SIP requests for a faux component
func (f *FauxComponentServer) handleSIPRequest(componentName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		f.mu.RLock()
		flowID := f.activeFlowID
		context := f.activeContext
		f.mu.RUnlock()

		if flowID == "" || context == nil {
			c.String(http.StatusInternalServerError, "No active call flow")
			return
		}

		// Read SIP message
		rawBody, err := c.GetRawData()
		if err != nil {
			c.String(http.StatusBadRequest, "Failed to read body")
			return
		}

		// Parse SIP message (simplified - just extract method and headers)
		sipMsg := string(rawBody)
		method := f.extractSIPMethod(sipMsg)
		callID := f.extractHeader(sipMsg, "Call-ID")
		from := f.extractHeader(sipMsg, "From")
		to := f.extractHeader(sipMsg, "To")
		cseq := f.extractHeader(sipMsg, "CSeq")

		f.logger.Debugf("Faux %s received: %s (Call-ID: %s)", componentName, method, callID)

		// Generate response based on call flow
		response := f.generateFauxResponse(componentName, method, flowID, context, callID, from, to, cseq)

		c.Data(http.StatusOK, "application/sip", []byte(response))
	}
}

// generateFauxResponse generates a SIP response for a faux component
func (f *FauxComponentServer) generateFauxResponse(componentName, method, flowID string, context *FauxCallContext, callID, from, to, cseq string) string {
	// Validate flow exists
	_, exists := f.registry.GetFlow(flowID)
	if !exists {
		return f.createSIPResponse(500, "Internal Server Error", callID, from, to, cseq)
	}

	// Find the next step for this component in the flow
	compSteps := f.registry.GetComponentSteps(componentName, flowID)
	if len(compSteps) == 0 {
		return f.createSIPResponse(404, "Not Found", callID, from, to, cseq)
	}

	// Get the first step for this component (simplified - in real implementation, track state)
	step := compSteps[0]

	// Generate response based on step
	switch step.Message {
	case "200 OK":
		return f.createSIPResponse(200, "OK", callID, from, to, cseq)
	case "180 Ringing":
		return f.createSIPResponse(180, "Ringing", callID, from, to, cseq)
	case "401 Unauthorized":
		return f.createSIPResponse(401, "Unauthorized", callID, from, to, cseq)
	case "UAA", "SAA":
		// HSS responses - return JSON
		return f.createHSSResponse(step.Message, callID)
	default:
		// Default to 200 OK
		return f.createSIPResponse(200, "OK", callID, from, to, cseq)
	}
}

// createSIPResponse creates a SIP response message
func (f *FauxComponentServer) createSIPResponse(statusCode int, statusText, callID, from, to, cseq string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("SIP/2.0 %d %s\r\n", statusCode, statusText))
	if via := f.extractHeader("", "Via"); via != "" {
		sb.WriteString(fmt.Sprintf("Via: %s\r\n", via))
	}
	if from != "" {
		sb.WriteString(fmt.Sprintf("From: %s\r\n", from))
	}
	if to != "" {
		sb.WriteString(fmt.Sprintf("To: %s\r\n", to))
	}
	if callID != "" {
		sb.WriteString(fmt.Sprintf("Call-ID: %s\r\n", callID))
	}
	if cseq != "" {
		sb.WriteString(fmt.Sprintf("CSeq: %s\r\n", cseq))
	}
	sb.WriteString("Content-Length: 0\r\n")
	sb.WriteString("\r\n")
	return sb.String()
}

// createHSSResponse creates an HSS response (JSON format)
func (f *FauxComponentServer) createHSSResponse(messageType, callID string) string {
	// Simplified HSS response
	return fmt.Sprintf(`{"result_code": 2001, "scscf": "scscf.example.com", "call_id": "%s"}`, callID)
}

// extractSIPMethod extracts the SIP method from a message
func (f *FauxComponentServer) extractSIPMethod(msg string) string {
	lines := strings.Split(msg, "\r\n")
	if len(lines) > 0 {
		parts := strings.Fields(lines[0])
		if len(parts) > 0 {
			return parts[0]
		}
	}
	return ""
}

// extractHeader extracts a header value from a SIP message
func (f *FauxComponentServer) extractHeader(msg, headerName string) string {
	lines := strings.Split(msg, "\r\n")
	for _, line := range lines {
		if strings.HasPrefix(strings.ToLower(line), strings.ToLower(headerName)+":") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return ""
}

// getComponentsInFlow gets all unique components in a call flow
func (f *FauxComponentServer) getComponentsInFlow(flow *CallFlow) []string {
	components := make(map[string]bool)
	for _, step := range flow.Steps {
		if step.From != "UE" && step.From != "PSTN" && step.From != "Destination" {
			components[step.From] = true
		}
		if step.To != "UE" && step.To != "PSTN" && step.To != "Destination" {
			components[step.To] = true
		}
	}

	var result []string
	for comp := range components {
		// Map to short names
		shortName := f.mapToShortName(comp)
		if shortName != "" {
			result = append(result, shortName)
		}
	}

	return result
}

// mapToShortName maps component names to short names
func (f *FauxComponentServer) mapToShortName(name string) string {
	nameMap := map[string]string{
		"P-CSCF":  "pcscf",
		"I-CSCF":  "icscf",
		"S-CSCF":  "scscf",
		"BGCF":    "bgcf",
		"MGCF":    "mgcf",
		"HSS":     "hss",
		"IC-SCF":  "icscf",
		"SC-SCF":  "scscf",
	}

	if short, exists := nameMap[name]; exists {
		return short
	}

	// Fallback: convert to lowercase and remove hyphens
	return strings.ToLower(strings.ReplaceAll(name, "-", ""))
}

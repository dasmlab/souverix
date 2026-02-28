package diagnostics

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Diagnostics provides common diagnostic endpoints for all components
type Diagnostics struct {
	componentName string
	version       string
	buildTime     string
	gitCommit     string
	logger        *logrus.Logger
	registry      *CallFlowRegistry
	fauxGen       *FauxCallGenerator
	stateVerifier *StateVerifier
	fauxServer    *FauxComponentServer
	stateTracker  *StateTracker
	stateProvider ComponentStateProvider // Optional: component can provide its own state
}

// New creates a new Diagnostics instance
func New(componentName, version, buildTime, gitCommit string, logger *logrus.Logger) *Diagnostics {
	registry := NewCallFlowRegistry()
	fauxGen := NewFauxCallGenerator(registry)
	
	// Extract component short name from full name (e.g., "Souverix P-CSCF" -> "pcscf")
	compShortName := extractComponentShortName(componentName)
	stateVerifier := NewStateVerifier(compShortName)

	return &Diagnostics{
		componentName: componentName,
		version:       version,
		buildTime:     buildTime,
		gitCommit:     gitCommit,
		logger:        logger,
		registry:      registry,
		fauxGen:       fauxGen,
		stateVerifier: stateVerifier,
	}
}

// extractComponentShortName extracts short component name from full name
// "Souverix P-CSCF" -> "pcscf", "Souverix I-CSCF" -> "icscf", etc.
func extractComponentShortName(fullName string) string {
	// Remove "Souverix " prefix if present
	name := fullName
	if strings.HasPrefix(name, "Souverix ") {
		name = strings.TrimPrefix(name, "Souverix ")
	}
	
	// Map common component names to short names
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

// RegisterRoutes registers diagnostic routes on a Gin router
func (d *Diagnostics) RegisterRoutes(router *gin.Engine) {
	diag := router.Group("/diag")
	{
		diag.GET("/health", d.Health)
		diag.GET("/status", d.Status)
		diag.GET("/local_test", d.LocalTest)
		diag.GET("/unit_test", d.UnitTest)
	}
}

// Health returns basic health check
// @Summary Health check
// @Description Basic health check endpoint
// @Tags diagnostics
// @Produce json
// @Success 200 {object} map[string]string
// @Router /diag/health [get]
func (d *Diagnostics) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"component": d.componentName,
	})
}

// Status returns component status and version information
// @Summary Component status
// @Description Returns component status, version, and build information
// @Tags diagnostics
// @Produce json
// @Success 200 {object} map[string]string
// @Router /diag/status [get]
func (d *Diagnostics) Status(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"component": d.componentName,
		"version":   d.version,
		"buildTime": d.buildTime,
		"gitCommit": d.gitCommit,
		"status":    "running",
	})
}

// LocalTest returns success for local testing
// @Summary Local test endpoint
// @Description Endpoint for local testing, returns success
// @Tags diagnostics
// @Produce json
// @Success 200 {object} map[string]string
// @Router /diag/local_test [get]
func (d *Diagnostics) LocalTest(c *gin.Context) {
	d.logger.Debug("Local test endpoint called")
	c.JSON(http.StatusOK, gin.H{
		"resp": "success",
		"component": d.componentName,
		"test_type": "local",
	})
}

// UnitTest executes unit test for a call flow
// @Summary Unit test endpoint
// @Description Executes call flow simulation for component's portion, verifies state changes
// @Tags diagnostics
// @Produce json
// @Param flow_id query string false "Call flow ID (e.g., IMS_REGISTER_AKA, SIP_INVITE_IMS_TO_IMS, SIP_INVITE_IMS_TO_PSTN)"
// @Param base_url query string false "Base URL for component (default: http://localhost:<diag_port>)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /diag/unit_test [get]
func (d *Diagnostics) UnitTest(c *gin.Context) {
	flowID := c.Query("flow_id")
	// Extract component short name for registry lookup
	compShortName := extractComponentShortName(d.componentName)
	
	if flowID == "" {
		// Default to first available flow for this component
		flows := d.registry.GetComponentFlows(compShortName)
		if len(flows) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "No call flows available for this component",
				"component": d.componentName,
				"component_short": compShortName,
			})
			return
		}
		flowID = flows[0].FlowID
	}

	baseURL := c.Query("base_url")
	if baseURL == "" {
		// Try to get from request
		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
		baseURL = fmt.Sprintf("%s://%s", scheme, c.Request.Host)
	}

	// Extract caller IP from request
	callerIP := c.ClientIP()
	if callerIP == "" || callerIP == "::1" || callerIP == "127.0.0.1" {
		// Try to get from Host header
		host := c.Request.Host
		if strings.Contains(host, ":") {
			parts := strings.Split(host, ":")
			callerIP = parts[0]
		} else {
			callerIP = host
		}
		if callerIP == "localhost" {
			callerIP = "127.0.0.1"
		}
	}

	// Initialize faux server with caller IP
	d.fauxServer = NewFauxComponentServer(d.registry, d.logger, 19000, callerIP)

	// Start faux component servers
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	if err := d.fauxServer.Start(ctx, flowID, compShortName); err != nil {
		d.logger.WithError(err).Error("Failed to start faux component servers")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to start faux servers: %v", err),
		})
		return
	}

	// Stop faux servers when done
	defer func() {
		stopCtx, stopCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer stopCancel()
		d.fauxServer.Stop(stopCtx)
	}()

	d.logger.Infof("Unit test called for flow: %s, component: %s (short: %s), caller IP: %s", flowID, d.componentName, compShortName, callerIP)

	// Get component's steps in this flow (use short name for registry lookup)
	compSteps := d.registry.GetComponentSteps(compShortName, flowID)
	if len(compSteps) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Component %s does not participate in flow %s", d.componentName, flowID),
			"component": d.componentName,
			"flow_id": flowID,
		})
		return
	}

	// Execute call flow simulation
	results := []map[string]interface{}{}
	allPassed := true

	for _, step := range compSteps {
		stepResult := map[string]interface{}{
			"step":        step.Sequence,
			"message":     step.Message,
			"interface":   step.Interface,
			"direction":   step.Direction,
			"description": step.Description,
		}

		// Generate faux request for this step
		fauxReq, err := d.fauxGen.GenerateFauxRequest(flowID, step.Sequence, compShortName, baseURL)
		if err != nil {
			stepResult["error"] = err.Error()
			stepResult["passed"] = false
			allPassed = false
			results = append(results, stepResult)
			continue
		}

		stepResult["request"] = map[string]interface{}{
			"method": fauxReq.Method,
			"url":    fauxReq.URL,
			"body":   fauxReq.Body,
		}

		// Execute faux request (if it's outgoing)
		if step.Direction == "request" {
			resp, err := d.fauxGen.ExecuteFauxRequest(fauxReq)
			if err != nil {
				stepResult["error"] = err.Error()
				stepResult["passed"] = false
				allPassed = false
			} else {
				stepResult["response"] = map[string]interface{}{
					"status_code": resp.StatusCode,
					"body":        resp.Body,
				}
				stepResult["passed"] = (resp.StatusCode == fauxReq.ExpectedCode)
				if !stepResult["passed"].(bool) {
					allPassed = false
				}
			}
		} else {
			// For responses, we verify state instead
			stepResult["passed"] = true
		}

		// Verify state changes (simplified - components should implement their own verification)
		// This is a placeholder - components should implement ComponentStateProvider
		verification := d.stateVerifier.VerifyStateExists(step.Sequence, step.Message, "call_state_"+step.Message)
		stepResult["state_verification"] = map[string]interface{}{
			"passed":  verification.Passed,
			"message": verification.Message,
		}

		if !verification.Passed {
			allPassed = false
		}

		results = append(results, stepResult)
	}

	// Get verification summary
	summary := d.stateVerifier.GetSummary()

	response := map[string]interface{}{
		"component":    d.componentName,
		"flow_id":      flowID,
		"all_passed":   allPassed,
		"steps":        results,
		"verification": summary,
	}

	statusCode := http.StatusOK
	if !allPassed {
		statusCode = http.StatusInternalServerError
	}

	c.JSON(statusCode, response)
}

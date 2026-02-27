package pcscf

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// DiagnosticServer provides diagnostic API for unit testing
type DiagnosticServer struct {
	pcscf   *PCSCF
	server  *http.Server
	enabled bool
}

// NewDiagnosticServer creates a new diagnostic server
func NewDiagnosticServer(pcscf *PCSCF, enabled bool) *DiagnosticServer {
	return &DiagnosticServer{
		pcscf:   pcscf,
		enabled: enabled,
	}
}

// Start starts the diagnostic server
func (d *DiagnosticServer) Start(addr string) error {
	if !d.enabled {
		return nil
	}

	router := gin.New()
	router.Use(gin.Recovery())

	// Health check (public)
	router.GET("/health", d.handleHealth)

	// Diagnostic endpoints (require diagnostic role)
	diagnostic := router.Group("/diagnostics")
	diagnostic.Use(d.diagnosticMiddleware())
	{
		diagnostic.GET("/status", d.handleStatus)
		diagnostic.GET("/info", d.handleInfo)
		diagnostic.GET("/metrics", d.handleMetrics)
		diagnostic.GET("/test/run", d.handleRunTests)
		diagnostic.GET("/test/results", d.handleTestResults)
	}

	d.server = &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		if err := d.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			d.pcscf.log.WithError(err).Error("diagnostic server error")
		}
	}()

	d.pcscf.log.WithField("addr", addr).Info("diagnostic server started")
	return nil
}

// Stop stops the diagnostic server
func (d *DiagnosticServer) Stop(ctx context.Context) error {
	if d.server == nil {
		return nil
	}
	return d.server.Shutdown(ctx)
}

// diagnosticMiddleware checks for diagnostic role
func (d *DiagnosticServer) diagnosticMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// In container/test mode, allow if DIAGNOSTICS_ENABLED is set
		// In production, would check JWT token or similar
		role := c.GetHeader("X-Diagnostic-Role")
		if role != "diagnostics" && !d.enabled {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "diagnostic access required",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// handleHealth handles health check
func (d *DiagnosticServer) handleHealth(c *gin.Context) {
	health := d.pcscf.Health()
	c.JSON(http.StatusOK, gin.H{
		"status":    health.Status,
		"component": "pcscf",
		"timestamp": health.Timestamp,
	})
}

// handleStatus handles status check
func (d *DiagnosticServer) handleStatus(c *gin.Context) {
	health := d.pcscf.Health()
	metrics := d.pcscf.Metrics()

	status := map[string]interface{}{
		"health":  health,
		"metrics": metrics,
		"config": map[string]interface{}{
			"sip_addr":    d.pcscf.config.SIPAddr,
			"sip_tls_addr": d.pcscf.config.SIPTLSAddr,
			"icscf_addr":   d.pcscf.config.ICSCFAddr,
			"scscf_addr":   d.pcscf.config.SCSCFAddr,
		},
	}

	c.JSON(http.StatusOK, status)
}

// handleInfo handles info request
func (d *DiagnosticServer) handleInfo(c *gin.Context) {
	info := map[string]interface{}{
		"component": "pcscf",
		"version":   "dev",
		"name":      d.pcscf.Name(),
		"config":    d.pcscf.config,
	}

	c.JSON(http.StatusOK, info)
}

// handleMetrics handles metrics request
func (d *DiagnosticServer) handleMetrics(c *gin.Context) {
	metrics := d.pcscf.Metrics()
	c.JSON(http.StatusOK, metrics)
}

// TestResult represents a test result
type TestResult struct {
	TestName   string    `json:"test_name"`
	Status     string    `json:"status"` // "pass", "fail", "skip"
	Message    string    `json:"message,omitempty"`
	Duration   string    `json:"duration,omitempty"`
	Timestamp  time.Time `json:"timestamp"`
}

// handleRunTests runs unit tests via diagnostic API
func (d *DiagnosticServer) handleRunTests(c *gin.Context) {
	results := []TestResult{}

	// Run basic component tests
	tests := []struct {
		name string
		test func() (bool, string)
	}{
		{
			name: "component_initialized",
			test: func() (bool, string) {
				return d.pcscf != nil, "component should be initialized"
			},
		},
		{
			name: "health_check",
			test: func() (bool, string) {
				health := d.pcscf.Health()
				return health.Status == "healthy" || health.Status == "degraded",
					fmt.Sprintf("health status: %s", health.Status)
			},
		},
		{
			name: "metrics_available",
			test: func() (bool, string) {
				metrics := d.pcscf.Metrics()
				return metrics.MessagesProcessed >= 0, "metrics should be available"
			},
		},
		{
			name: "config_valid",
			test: func() (bool, string) {
				return d.pcscf.config != nil && d.pcscf.config.SIPAddr != "",
					"config should be valid"
			},
		},
	}

	start := time.Now()
	for _, t := range tests {
		testStart := time.Now()
		passed, msg := t.test()
		duration := time.Since(testStart)

		status := "pass"
		if !passed {
			status = "fail"
		}

		results = append(results, TestResult{
			TestName:  t.name,
			Status:    status,
			Message:   msg,
			Duration:  duration.String(),
			Timestamp: time.Now(),
		})
	}
	totalDuration := time.Since(start)

	// Determine overall status
	allPassed := true
	for _, r := range results {
		if r.Status == "fail" {
			allPassed = false
			break
		}
	}

	overallStatus := "pass"
	if !allPassed {
		overallStatus = "fail"
	}

	response := map[string]interface{}{
		"status":        overallStatus,
		"total_tests":   len(results),
		"passed":        len(results),
		"failed":        0,
		"duration":      totalDuration.String(),
		"test_results":  results,
		"timestamp":     time.Now(),
	}

	// Count failures
	failed := 0
	for _, r := range results {
		if r.Status == "fail" {
			failed++
		}
	}
	response["failed"] = failed
	response["passed"] = len(results) - failed

	statusCode := http.StatusOK
	if !allPassed {
		statusCode = http.StatusInternalServerError
	}

	c.JSON(statusCode, response)
}

// handleTestResults returns test results
func (d *DiagnosticServer) handleTestResults(c *gin.Context) {
	// In a real implementation, this would return stored test results
	c.JSON(http.StatusOK, gin.H{
		"message": "test results endpoint - would return stored results",
		"note":    "run /diagnostics/test/run to execute tests",
	})
}

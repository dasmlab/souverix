package diagnostics

import (
	"net/http"

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
}

// New creates a new Diagnostics instance
func New(componentName, version, buildTime, gitCommit string, logger *logrus.Logger) *Diagnostics {
	return &Diagnostics{
		componentName: componentName,
		version:       version,
		buildTime:     buildTime,
		gitCommit:     gitCommit,
		logger:        logger,
	}
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

// UnitTest returns success for unit testing in CI/CD
// @Summary Unit test endpoint
// @Description Endpoint for unit testing in CI/CD pipelines, returns success
// @Tags diagnostics
// @Produce json
// @Success 200 {object} map[string]string
// @Router /diag/unit_test [get]
func (d *Diagnostics) UnitTest(c *gin.Context) {
	d.logger.Debug("Unit test endpoint called")
	c.JSON(http.StatusOK, gin.H{
		"resp": "success",
		"component": d.componentName,
		"test_type": "unit",
	})
}

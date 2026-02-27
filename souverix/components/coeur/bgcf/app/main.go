package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Depado/ginprom"
	"github.com/dasmlab/souverix/common/diagnostics"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"
)

// @title Souverix BGCF Diagnostic API
// @version 1.0
// @description Diagnostic endpoints for Souverix BGCF
// @host localhost:8084
// @BasePath /
func main() {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	logger.WithFields(logrus.Fields{
		"component": "Souverix BGCF",
		"version":   version,
		"build":     gitCommit,
	}).Info("Souverix - Souverix BGCF - Version: " + version + " Build: " + gitCommit)

	// Get ports from environment or use defaults
	mainPort := os.Getenv("PORT")
	if mainPort == "" {
		mainPort = "8084"
	}
	metricsPort := os.Getenv("METRICS_PORT")
	if metricsPort == "" {
		metricsPort = "9094"
	}
	diagPort := os.Getenv("DIAG_PORT")
	if diagPort == "" {
		diagPort = "9084"
	}
	testPort := os.Getenv("TEST_PORT")
	if testPort == "" {
		testPort = "9184"
	}

	// Initialize main Gin router (r1) - main application server
	gin.SetMode(gin.ReleaseMode)
	r1 := gin.New()
	r1.Use(gin.LoggerWithWriter(logger.Writer()))
	r1.Use(gin.Recovery())

	// Initialize metrics router (r2) - Prometheus metrics, out of band
	r2 := gin.New()
	r2.Use(gin.LoggerWithWriter(logger.Writer()))
	r2.Use(gin.Recovery())

	// Initialize diagnostics router (r3) - diagnostics endpoints, out of band
	r3 := gin.New()
	r3.Use(gin.LoggerWithWriter(logger.Writer()))
	r3.Use(gin.Recovery())

	// Initialize test router (r4) - test endpoints, out of band
	r4 := gin.New()
	r4.Use(gin.LoggerWithWriter(logger.Writer()))
	r4.Use(gin.Recovery())

	// Setup ginprom - wrap r1 with instrumentation, r2 serves metrics
	p := ginprom.New(
		ginprom.Engine(r2),
		ginprom.Subsystem("gin"),
		ginprom.Path("/metrics"),
	)

	// Wrap main router (r1) with ginprom instrumentation
	r1.Use(p.Instrument())

	// Register diagnostic endpoints on r3 (diagnostics server)
	diag := diagnostics.New("Souverix BGCF", version, buildTime, gitCommit, logger)
	diag.RegisterRoutes(r3)

	// Register test endpoints on r4 (test server)
	// TODO: Add testing framework when ready
	r4.GET("/test/local", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"resp": "success",
			"component": "Souverix BGCF",
			"test_type": "local",
		})
	})
	r4.GET("/test/unit", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"resp": "success",
			"component": "Souverix BGCF",
			"test_type": "unit",
		})
	})

	// Create HTTP servers
	srv1 := &http.Server{
		Addr:    ":" + mainPort,
		Handler: r1,
	}

	srv2 := &http.Server{
		Addr:    ":" + metricsPort,
		Handler: r2,
	}

	srv3 := &http.Server{
		Addr:    ":" + diagPort,
		Handler: r3,
	}

	srv4 := &http.Server{
		Addr:    ":" + testPort,
		Handler: r4,
	}

	// Start all servers as goroutines (out of band)
	// Start metrics server (r2) first
	go func() {
		logger.Infof("Starting metrics server (Prometheus) on :%s", metricsPort)
		if err := srv2.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatal("failed to start metrics server")
		}
	}()

	// Start diagnostics server (r3)
	go func() {
		logger.Infof("Starting diagnostics server on :%s", diagPort)
		if err := srv3.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatal("failed to start diagnostics server")
		}
	}()

	// Start test server (r4)
	go func() {
		logger.Infof("Starting test server on :%s", testPort)
		if err := srv4.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatal("failed to start test server")
		}
	}()

	// Start main server (r1) last
	go func() {
		logger.Infof("Starting main server on :%s", mainPort)
		if err := srv1.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatal("failed to start main server")
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down Souverix BGCF...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown all servers
	if err := srv4.Shutdown(shutdownCtx); err != nil {
		logger.WithError(err).Error("error during test server shutdown")
	}
	if err := srv3.Shutdown(shutdownCtx); err != nil {
		logger.WithError(err).Error("error during diagnostics server shutdown")
	}
	if err := srv2.Shutdown(shutdownCtx); err != nil {
		logger.WithError(err).Error("error during metrics server shutdown")
	}
	if err := srv1.Shutdown(shutdownCtx); err != nil {
		logger.WithError(err).Error("error during main server shutdown")
	}

	logger.Info("Souverix BGCF stopped")
}

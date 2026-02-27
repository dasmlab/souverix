package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	// Initialize main Gin router (r1)
	gin.SetMode(gin.ReleaseMode)
	r1 := gin.New()
	r1.Use(gin.LoggerWithWriter(logger.Writer()))
	r1.Use(gin.Recovery())

	// Register diagnostic endpoints
	diag := diagnostics.New("Souverix BGCF", version, buildTime, gitCommit, logger)
	diag.RegisterRoutes(r1)

	// Initialize metrics router (r2) - out of band
	// TODO: Add ginprom wrapper when metrics package is ready
	r2 := gin.New()
	r2.Use(gin.LoggerWithWriter(logger.Writer()))
	r2.Use(gin.Recovery())
	
	// Metrics endpoint placeholder
	r2.GET("/metrics", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "metrics_endpoint_ready",
			"component": "Souverix BGCF",
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

	// Start metrics server (r2) first as goroutine
	go func() {
		logger.Infof("Starting metrics server on :%s", metricsPort)
		if err := srv2.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatal("failed to start metrics server")
		}
	}()

	// Start main server (r1) as goroutine
	go func() {
		logger.Infof("Starting diagnostic server on :%s", mainPort)
		if err := srv1.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatal("failed to start diagnostic server")
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down Souverix BGCF...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown both servers
	if err := srv2.Shutdown(shutdownCtx); err != nil {
		logger.WithError(err).Error("error during metrics server shutdown")
	}
	if err := srv1.Shutdown(shutdownCtx); err != nil {
		logger.WithError(err).Error("error during diagnostic server shutdown")
	}

	logger.Info("Souverix BGCF stopped")
}

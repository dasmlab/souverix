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

// @title Souverix HSS Diagnostic API
// @version 1.0
// @description Diagnostic endpoints for Souverix HSS
// @host localhost:8086
// @BasePath /
func main() {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	logger.WithFields(logrus.Fields{
		"component": "Souverix HSS",
		"version":   version,
		"build":     gitCommit,
	}).Info("Souverix - Souverix HSS - Version: " + version + " Build: " + gitCommit)

	// Initialize Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.LoggerWithWriter(logger.Writer()))
	router.Use(gin.Recovery())

	// Register diagnostic endpoints
	diag := diagnostics.New("Souverix HSS", version, buildTime, gitCommit, logger)
	diag.RegisterRoutes(router)

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":8086",
		Handler: router,
	}

	// Start server in goroutine
	go func() {
		logger.Info("Starting diagnostic server on :8086")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatal("failed to start diagnostic server")
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down Souverix HSS...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.WithError(err).Error("error during shutdown")
	}

	logger.Info("Souverix HSS stopped")
}

package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	{{COMPONENT_PKG}} "github.com/dasmlab/ims/internal/{{COMPONENT_DIR}}"
	gouverneConfig "github.com/dasmlab/ims/internal/gouverne/config"
	"github.com/sirupsen/logrus"
)

var (
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"
)

func main() {
	// Initialize logger
	log := logrus.New()
	log.SetLevel(logrus.InfoLevel)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Standard component startup log - FIRST LINE
	log.WithFields(logrus.Fields{
		"component": "{{COMPONENT_NAME}}",
		"version":   version,
		"build":     gitCommit,
	}).Info("Souverix - {{COMPONENT_NAME}} - Version: " + version + " Build: " + gitCommit)

	// Load configuration
	cfg := gouverneConfig.Load()
	{{COMPONENT_INIT}}

	// Start component
	ctx := context.Background()
	if err := {{COMPONENT_START}}; err != nil {
		log.WithError(err).Fatal("failed to start {{COMPONENT_NAME}}")
	}

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down {{COMPONENT_NAME}}...")

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := {{COMPONENT_STOP}}; err != nil {
		log.WithError(err).Error("error during shutdown")
	}

	log.Info("{{COMPONENT_NAME}} stopped")
}

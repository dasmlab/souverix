package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	
	"github.com/dasmlab/ims/components/coeur/mgcf/app/internal/mgcf"
	"github.com/sirupsen/logrus"
)

var (
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"
)

// @title Souverix MGCF Diagnostic API
// @version 1.0
// @description Diagnostic endpoints for Souverix MGCF
// @host localhost:8081
// @BasePath /
func main() {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	
	logger.WithFields(logrus.Fields{
		"component": "Souverix MGCF",
		"version":   version,
		"build":     gitCommit,
	}).Info("Souverix - Souverix MGCF - Version: " + version + " Build: " + gitCommit)
	
	// Create MGCF handler
	stdLogger := log.New(os.Stdout, "[MGCF] ", log.LstdFlags)
	handler := mgcf.NewHandler(stdLogger)

	// Start MGCF
	logger.Info("MGCF component started")
	_ = handler

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	logger.Info("shutting down Souverix MGCF...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = shutdownCtx
	logger.Info("Souverix MGCF stopped")
}

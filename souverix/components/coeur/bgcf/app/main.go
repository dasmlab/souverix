package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dasmlab/ims/components/coeur/bgcf/app/internal/bgcf"
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
// @host localhost:8081
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

	// Create BGCF handler
	stdLogger := log.New(os.Stdout, "[BGCF] ", log.LstdFlags)
	handler := bgcf.NewHandler(stdLogger)

	// Start BGCF
	logger.Info("BGCF component started")
	_ = handler

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down Souverix BGCF...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = shutdownCtx
	logger.Info("Souverix BGCF stopped")
}

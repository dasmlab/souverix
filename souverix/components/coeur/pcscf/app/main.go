package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	
	"github.com/dasmlab/ims/components/coeur/pcscf/app/internal/pcscf"
	"github.com/sirupsen/logrus"
)

var (
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"
)

// @title Souverix P-CSCF Diagnostic API
// @version 1.0
// @description Diagnostic endpoints for Souverix P-CSCF
// @host localhost:8081
// @BasePath /
func main() {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	
	logger.WithFields(logrus.Fields{
		"component": "Souverix P-CSCF",
		"version":   version,
		"build":     gitCommit,
	}).Info("Souverix - Souverix P-CSCF - Version: " + version + " Build: " + gitCommit)
	
	// Create P-CSCF handler
	stdLogger := log.New(os.Stdout, "[P-CSCF] ", log.LstdFlags)
	handler := pcscf.NewHandler("icscf.example.com:5060", stdLogger)

	// Start P-CSCF
	logger.Info("P-CSCF component started")
	logger.Info("Handler initialized with I-CSCF: icscf.example.com:5060")
	
	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	logger.Info("shutting down Souverix P-CSCF...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = shutdownCtx
	logger.Info("Souverix P-CSCF stopped")
}

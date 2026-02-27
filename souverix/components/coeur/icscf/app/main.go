package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	
	"github.com/dasmlab/ims/components/common/hss"
	"github.com/dasmlab/ims/components/coeur/icscf/app/internal/icscf"
	"github.com/sirupsen/logrus"
)

var (
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"
)

// @title Souverix I-CSCF Diagnostic API
// @version 1.0
// @description Diagnostic endpoints for Souverix I-CSCF
// @host localhost:8081
// @BasePath /
func main() {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	
	logger.WithFields(logrus.Fields{
		"component": "Souverix I-CSCF",
		"version":   version,
		"build":     gitCommit,
	}).Info("Souverix - Souverix I-CSCF - Version: " + version + " Build: " + gitCommit)
	
	// Create HSS client
	hssClient := hss.NewHSSClient()

	// Create I-CSCF handler
	stdLogger := log.New(os.Stdout, "[I-CSCF] ", log.LstdFlags)
	handler := icscf.NewHandler(hssClient, stdLogger)

	// Start I-CSCF
	logger.Info("I-CSCF component started")
	logger.Info("HSS client initialized")
	_ = handler

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	logger.Info("shutting down Souverix I-CSCF...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = shutdownCtx
	logger.Info("Souverix I-CSCF stopped")
}

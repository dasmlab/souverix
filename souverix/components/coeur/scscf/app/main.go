package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	
	"github.com/dasmlab/ims/components/common/hss"
	"github.com/dasmlab/ims/components/coeur/scscf/app/internal/scscf"
	"github.com/sirupsen/logrus"
)

var (
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"
)

// @title Souverix S-CSCF Diagnostic API
// @version 1.0
// @description Diagnostic endpoints for Souverix S-CSCF
// @host localhost:8081
// @BasePath /
func main() {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	
	logger.WithFields(logrus.Fields{
		"component": "Souverix S-CSCF",
		"version":   version,
		"build":     gitCommit,
	}).Info("Souverix - Souverix S-CSCF - Version: " + version + " Build: " + gitCommit)
	
	// Create HSS client
	hssClient := hss.NewHSSClient()

	// Create S-CSCF handler
	stdLogger := log.New(os.Stdout, "[S-CSCF] ", log.LstdFlags)
	handler := scscf.NewHandler(hssClient, "bgcf.example.com:5060", stdLogger)

	// Start S-CSCF
	logger.Info("S-CSCF component started")
	logger.Info("HSS client initialized")
	logger.Info("BGCF address: bgcf.example.com:5060")
	_ = handler

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	logger.Info("shutting down Souverix S-CSCF...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = shutdownCtx
	logger.Info("Souverix S-CSCF stopped")
}

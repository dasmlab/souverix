package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	coeur "github.com/dasmlab/ims/internal/coeur"
	gouverneConfig "github.com/dasmlab/ims/internal/gouverne/config"
	"github.com/sirupsen/logrus"
)

var (
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"
)

// @title Souverix Coeur Diagnostic API
// @version 1.0
// @description Diagnostic and health check endpoints for Souverix Coeur (IMS Core)
// @termsOfService http://swagger.io/terms/
// @contact.name Souverix Support
// @license.name Apache 2.0
// @host localhost:8081
// @BasePath /
func main() {
	// Initialize logger
	log := logrus.New()
	log.SetLevel(logrus.InfoLevel)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Standard component startup log - FIRST LINE
	log.WithFields(logrus.Fields{
		"component": "Souverix Coeur",
		"version":   version,
		"build":     gitCommit,
	}).Info("Souverix - Souverix Coeur - Version: " + version + " Build: " + gitCommit)

	// Load configuration
	cfg := gouverneConfig.Load()
	
	// Initialize Coeur component
	componentCfg := coeur.ConfigFromGouverne(cfg)
	component := coeur.New(componentCfg, log)

	// Start component
	ctx := context.Background()
	if err := component.Start(ctx); err != nil {
		log.WithError(err).Fatal("failed to start Souverix Coeur")
	}

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down Souverix Coeur...")

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := component.Stop(shutdownCtx); err != nil {
		log.WithError(err).Error("error during shutdown")
	}

	log.Info("Souverix Coeur stopped")
}

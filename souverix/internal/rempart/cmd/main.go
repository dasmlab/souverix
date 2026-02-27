package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	rempart "github.com/dasmlab/ims/internal/rempart"
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
		"component": "Souverix Rempart",
		"version":   version,
		"build":     gitCommit,
	}).Info("Souverix - Souverix Rempart - Version: " + version + " Build: " + gitCommit)

	// Load configuration
	cfg := gouverneConfig.Load()
	// Initialize Rempart component
	component, err := rempart.NewSBC(cfg, log)
	if err != nil {
		log.WithError(err).Fatal("failed to create Rempart")
	}

	// Start component
	ctx := context.Background()
	if err := component.Start(ctx); err != nil {
		log.WithError(err).Fatal("failed to start Souverix Rempart")
	}

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down Souverix Rempart...")

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := component.Stop(shutdownCtx); err != nil {
		log.WithError(err).Error("error during shutdown")
	}

	log.Info("Souverix Rempart stopped")
}

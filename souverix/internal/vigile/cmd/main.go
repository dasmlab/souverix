package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	vigile "github.com/dasmlab/ims/internal/vigile"
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
		"component": "Souverix Vigile",
		"version":   version,
		"build":     gitCommit,
	}).Info("Souverix - Souverix Vigile - Version: " + version + " Build: " + gitCommit)

	// Load configuration
	cfg := gouverneConfig.Load()
	// Initialize Souverix Vigile (stub)
	component := vigile.New(cfg, log)

	// Start component
	ctx := context.Background()
	if err := component.Start(ctx); err != nil {
		log.WithError(err).Fatal("failed to start Souverix Vigile")
	}

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down Souverix Vigile...")

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := component.Stop(shutdownCtx); err != nil {
		log.WithError(err).Error("error during shutdown")
	}

	log.Info("Souverix Vigile stopped")
}

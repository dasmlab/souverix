package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	gouverne "github.com/dasmlab/ims/components/gouverne"
	gouverneConfig "github.com/dasmlab/ims/components/gouverne/config"
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
		"component": "Souverix Gouverne",
		"version":   version,
		"build":     gitCommit,
	}).Info("Souverix - Souverix Gouverne - Version: " + version + " Build: " + gitCommit)

	// Load configuration
	cfg := gouverneConfig.Load()
	// Initialize Souverix Gouverne (stub)
		log.WithError(err).Fatal("failed to start Souverix Gouverne")
	}

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down Souverix Gouverne...")

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

		log.WithError(err).Error("error during shutdown")
	}

	log.Info("Souverix Gouverne stopped")
}

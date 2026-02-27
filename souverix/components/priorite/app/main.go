package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	priorite "github.com/dasmlab/ims/components/priorite"
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
		"component": "Souverix Priorite",
		"version":   version,
		"build":     gitCommit,
	}).Info("Souverix - Souverix Priorite - Version: " + version + " Build: " + gitCommit)

	// Load configuration
	cfg := gouverneConfig.Load()
	// Initialize Souverix Priorite (stub)
		log.WithError(err).Fatal("failed to start Souverix Priorite")
	}

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down Souverix Priorite...")

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

		log.WithError(err).Error("error during shutdown")
	}

	log.Info("Souverix Priorite stopped")
}

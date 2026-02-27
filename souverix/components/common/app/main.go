package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
	
	"github.com/sirupsen/logrus"
)

var (
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"
)

// @title Souverix Common Diagnostic API
// @version 1.0
// @description Common utilities and shared code
// @host localhost:8081
// @BasePath /
func main() {
	log := logrus.New()
	log.SetLevel(logrus.InfoLevel)
	log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	
	log.WithFields(logrus.Fields{
		"component": "Souverix Common",
		"version":   version,
		"build":     gitCommit,
	}).Info("Souverix - Souverix Common - Version: " + version + " Build: " + gitCommit)
	
	// Common is a library component - minimal runtime
	log.Info("Common utilities available")
	
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Info("shutting down Souverix Common...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = shutdownCtx
	log.Info("Souverix Common stopped")
}

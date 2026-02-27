package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	hssPkg "github.com/dasmlab/ims/internal/coeur/hss"
	gouverneConfig "github.com/dasmlab/ims/internal/gouverne/config"
	"github.com/sirupsen/logrus"
)

var (
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"
)

// @title Souverix Hss Diagnostic API
// @version 1.0
// @description Diagnostic endpoints for Souverix Hss
// @host localhost:8081
// @BasePath /
func main() {
	log := logrus.New()
	log.SetLevel(logrus.InfoLevel)
	log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	
	log.WithFields(logrus.Fields{
		"component": "Souverix Hss",
		"version":   version,
		"build":     gitCommit,
	}).Info("Souverix - Souverix Hss - Version: " + version + " Build: " + gitCommit)
	
	cfg := gouverneConfig.Load()
	component := hssPkg.New(cfg, log)
	
	ctx := context.Background()
	if err := component.Start(ctx); err != nil {
		log.WithError(err).Fatal("failed to start Souverix Hss")
	}
	
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Info("shutting down Souverix Hss...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	component.Stop(shutdownCtx)
	log.Info("Souverix Hss stopped")
}

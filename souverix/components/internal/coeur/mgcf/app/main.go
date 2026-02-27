package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	mgcfPkg "github.com/dasmlab/ims/internal/coeur/mgcf"
	gouverneConfig "github.com/dasmlab/ims/internal/gouverne/config"
	"github.com/sirupsen/logrus"
)

var (
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"
)

// @title Souverix Mgcf Diagnostic API
// @version 1.0
// @description Diagnostic endpoints for Souverix Mgcf
// @host localhost:8081
// @BasePath /
func main() {
	log := logrus.New()
	log.SetLevel(logrus.InfoLevel)
	log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	
	log.WithFields(logrus.Fields{
		"component": "Souverix Mgcf",
		"version":   version,
		"build":     gitCommit,
	}).Info("Souverix - Souverix Mgcf - Version: " + version + " Build: " + gitCommit)
	
	cfg := gouverneConfig.Load()
	component := mgcfPkg.New(cfg, log)
	
	ctx := context.Background()
	if err := component.Start(ctx); err != nil {
		log.WithError(err).Fatal("failed to start Souverix Mgcf")
	}
	
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Info("shutting down Souverix Mgcf...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	component.Stop(shutdownCtx)
	log.Info("Souverix Mgcf stopped")
}

package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	pcscfPkg "github.com/dasmlab/ims/internal/coeur/pcscf"
	gouverneConfig "github.com/dasmlab/ims/internal/gouverne/config"
	"github.com/sirupsen/logrus"
)

func main() {
	// Initialize logger
	log := logrus.New()
	log.SetLevel(logrus.InfoLevel)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	log.Info("starting P-CSCF component")

	// Load configuration
	cfg := gouverneConfig.Load()
	pcscfCfg := pcscfPkg.ConfigFromGouverne(cfg)

	// Create P-CSCF instance
	pcscf := pcscfPkg.New(pcscfCfg, log)

	// Start P-CSCF
	ctx := context.Background()
	if err := pcscf.Start(ctx); err != nil {
		log.WithError(err).Fatal("failed to start P-CSCF")
	}

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down P-CSCF...")

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := pcscf.Stop(shutdownCtx); err != nil {
		log.WithError(err).Error("error during shutdown")
	}

	log.Info("P-CSCF stopped")
}

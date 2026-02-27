package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/dasmlab/ims/internal/config"
	"github.com/dasmlab/ims/internal/logutil"
	"github.com/dasmlab/ims/internal/metrics"
	"github.com/dasmlab/ims/internal/sbc"
	"github.com/dasmlab/ims/internal/store"
	"github.com/sirupsen/logrus"
)

var (
	version = "dev"
	buildTime = "unknown"
)

func main() {
	// Load configuration
	cfg := config.Load()
	cfg.Version = version

	// Initialize logger
	log := logutil.InitLogger("ims-core")
	log.WithFields(logrus.Fields{
		"version":    version,
		"build_time": buildTime,
	}).Info("starting IMS core")

	// Initialize OpenTelemetry
	otelShutdown := metrics.InitOTel(log)
	defer func() {
		if err := otelShutdown(context.Background()); err != nil {
			log.WithError(err).Error("failed to shutdown OpenTelemetry")
		}
	}()

	// Initialize HSS store
	hssStore, err := store.NewMemHSSStore(log)
	if err != nil {
		log.WithError(err).Fatal("failed to initialize HSS store")
	}

	// Initialize SBC if enabled
	var sbcInstance *sbc.SBC
	if cfg.IMS.EnableSBC {
		sbcInstance, err = sbc.NewSBC(cfg, log)
		if err != nil {
			log.WithError(err).Fatal("failed to initialize SBC")
		}

		if err := sbcInstance.Start(); err != nil {
			log.WithError(err).Fatal("failed to start SBC")
		}
		defer sbcInstance.Stop()
	}

	// Setup HTTP routers
	apiRouter := gin.New()
	metricsRouter := gin.New()

	apiRouter.Use(gin.Recovery())
	metricsRouter.Use(gin.Recovery())

	// Logging middleware
	apiRouter.Use(logutil.GinLogger(log))
	metricsRouter.Use(logutil.GinLogger(log))

	// Prometheus metrics
	apiRouter.Use(metrics.GinPromMiddleware())
	metricsRouter.GET("/metrics", metrics.MetricsHandler())

	// Health check
	apiRouter.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"version": version,
		})
	})

	// API routes
	setupAPIRoutes(apiRouter, hssStore, log)

	// Setup servers
	apiAddr := cfg.Server.APIAddr
	metricsAddr := cfg.Server.MetricsAddr

	apiSrv := &http.Server{
		Addr:         apiAddr,
		Handler:      apiRouter,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	metricsSrv := &http.Server{
		Addr:         metricsAddr,
		Handler:      metricsRouter,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start metrics server
	go func() {
		log.WithField("addr", metricsAddr).Info("metrics server listening")
		if err := metricsSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("metrics server error")
		}
	}()

	// Start API server
	go func() {
		log.WithField("addr", apiAddr).Info("API server listening")
		if err := apiSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("API server error")
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down IMS core...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := apiSrv.Shutdown(ctx); err != nil {
		log.WithError(err).Error("API server shutdown error")
	}

	if err := metricsSrv.Shutdown(ctx); err != nil {
		log.WithError(err).Error("metrics server shutdown error")
	}

	log.Info("IMS core stopped")
}

// setupAPIRoutes sets up API routes
func setupAPIRoutes(router *gin.Engine, hssStore store.HSSStore, log *logrus.Logger) {
	v1 := router.Group("/api/v1")

	// Subscriber endpoints
	subscribers := v1.Group("/subscribers")
	{
		subscribers.GET("", func(c *gin.Context) {
			subs, err := hssStore.ListSubscribers()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, subs)
		})

		subscribers.GET("/:impi", func(c *gin.Context) {
			impi := c.Param("impi")
			sub, err := hssStore.GetSubscriber(impi)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, sub)
		})
	}

	// Registration endpoints
	registrations := v1.Group("/registrations")
	{
		registrations.GET("/:impi", func(c *gin.Context) {
			impi := c.Param("impi")
			reg, err := hssStore.GetRegistration(impi)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, reg)
		})
	}
}

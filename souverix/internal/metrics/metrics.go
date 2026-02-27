package metrics

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

var (
	// HTTP metrics
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ims_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "ims_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	// SIP metrics
	sipMessagesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ims_sip_messages_total",
			Help: "Total number of SIP messages",
		},
		[]string{"method", "direction", "status"},
	)

	sipMessageDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "ims_sip_message_duration_seconds",
			Help:    "SIP message processing duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	// IMS component metrics
	activeSessions = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ims_active_sessions",
			Help: "Number of active IMS sessions",
		},
		[]string{"component"},
	)

	subscriberCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "ims_subscriber_count",
			Help: "Total number of subscribers in HSS",
		},
	)
)

func init() {
	prometheus.MustRegister(
		httpRequestsTotal,
		httpRequestDuration,
		sipMessagesTotal,
		sipMessageDuration,
		activeSessions,
		subscriberCount,
	)
}

// InitOTel initializes OpenTelemetry tracing
func InitOTel(log *logrus.Logger) func(context.Context) error {
	exporter, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
	)
	if err != nil {
		log.WithError(err).Warn("failed to initialize OpenTelemetry exporter")
		return func(context.Context) error { return nil }
	}

	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceName("ims-core"),
			semconv.ServiceVersion("dev"),
		),
	)
	if err != nil {
		log.WithError(err).Warn("failed to create OpenTelemetry resource")
		return func(context.Context) error { return nil }
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)

	otel.SetTracerProvider(tp)

	return func(ctx context.Context) error {
		return tp.Shutdown(ctx)
	}
}

// GinPromMiddleware returns a Gin middleware for Prometheus metrics
func GinPromMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		c.Next()

		status := c.Writer.Status()
		duration := time.Since(start).Seconds()

		httpRequestsTotal.WithLabelValues(
			c.Request.Method,
			path,
			http.StatusText(status),
		).Inc()

		httpRequestDuration.WithLabelValues(
			c.Request.Method,
			path,
		).Observe(duration)
	}
}

// MetricsHandler returns the Prometheus metrics handler
func MetricsHandler() gin.HandlerFunc {
	return gin.WrapH(promhttp.Handler())
}

// RecordSIPMessage records a SIP message metric
func RecordSIPMessage(method, direction, status string, duration time.Duration) {
	sipMessagesTotal.WithLabelValues(method, direction, status).Inc()
	sipMessageDuration.WithLabelValues(method).Observe(duration.Seconds())
}

// SetActiveSessions sets the active sessions gauge
func SetActiveSessions(component string, count int) {
	activeSessions.WithLabelValues(component).Set(float64(count))
}

// SetSubscriberCount sets the subscriber count gauge
func SetSubscriberCount(count int) {
	subscriberCount.Set(float64(count))
}

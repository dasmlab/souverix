# Multi-stage build for Souverix Gouverne component
FROM golang:1.23-alpine AS builder

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY internal/common ./internal/common
COPY internal/gouverne ./internal/gouverne
COPY pkg ./pkg

# Build arguments
ARG COMPONENT=gouverne
ARG VERSION=dev
ARG BUILD_TIME
ARG GIT_COMMIT

# Build a simple main that uses Souverix Gouverne
RUN cat > /build/main.go << 'EOF'
package main
import (
    "context"
    "os"
    "os/signal"
    "syscall"
    "time"
    gouverne "github.com/dasmlab/ims/internal/gouverne"
    gouverneConfig "github.com/dasmlab/ims/internal/gouverne/config"
    "github.com/sirupsen/logrus"
)
func main() {
    log := logrus.New()
    log.SetLevel(logrus.InfoLevel)
    log.SetFormatter(&logrus.TextFormatter{
        FullTimestamp: true,
    })
    
    // Standard component startup log
    log.WithFields(logrus.Fields{
        "component": "Souverix Gouverne",
        "version": "${VERSION}",
        "build": "${GIT_COMMIT}",
    }).Info("Souverix - Souverix Gouverne - Version: ${VERSION} Build: ${GIT_COMMIT}")
    
    cfg := gouverneConfig.Load()
    // Initialize Souverix Gouverne (stub)
	component := gouverne.New(cfg, log)
    
    ctx := context.Background()
    if err := component.Start(ctx); err != nil {
        log.WithError(err).Fatal("failed to start Souverix Gouverne")
    }
    
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    log.Info("shutting down Souverix Gouverne...")
    shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    component.Stop(shutdownCtx)
    log.Info("Souverix Gouverne stopped")
}
EOF

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo \
    -ldflags "-X main.version=${VERSION} -X main.buildTime=${BUILD_TIME} -X main.gitCommit=${GIT_COMMIT}" \
    -o /build/gouverne \
    /build/main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/gouverne /app/gouverne

# Expose ports (component-specific)
EXPOSE 8081/tcp
EXPOSE 8081/tcp


# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8081/health || exit 1

# Run the component
CMD ["/app/gouverne"]

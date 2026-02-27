# Multi-stage build for Souverix Vigile component
FROM golang:1.26-alpine AS builder

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY internal/common ./internal/common
COPY internal/vigile ./internal/vigile
COPY pkg ./pkg

# Build arguments
ARG COMPONENT=vigile
ARG VERSION=dev
ARG BUILD_TIME
ARG GIT_COMMIT

# Build a simple main that uses Souverix Vigile
RUN cat > /build/main.go << 'EOF'
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
func main() {
    log := logrus.New()
    log.SetLevel(logrus.InfoLevel)
    log.SetFormatter(&logrus.TextFormatter{
        FullTimestamp: true,
    })
    
    // Standard component startup log
    log.WithFields(logrus.Fields{
        "component": "Souverix Vigile",
        "version": "${VERSION}",
        "build": "${GIT_COMMIT}",
    }).Info("Souverix - Souverix Vigile - Version: ${VERSION} Build: ${GIT_COMMIT}")
    
    cfg := gouverneConfig.Load()
    // Initialize Souverix Vigile (stub)
	component := vigile.New(cfg, log)
    
    ctx := context.Background()
    if err := component.Start(ctx); err != nil {
        log.WithError(err).Fatal("failed to start Souverix Vigile")
    }
    
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    log.Info("shutting down Souverix Vigile...")
    shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    component.Stop(shutdownCtx)
    log.Info("Souverix Vigile stopped")
}
EOF

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo \
    -ldflags "-X main.version=${VERSION} -X main.buildTime=${BUILD_TIME} -X main.gitCommit=${GIT_COMMIT}" \
    -o /build/vigile \
    /build/main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/vigile /app/vigile

# Expose ports (component-specific)
EXPOSE 8081/tcp
EXPOSE 8081/tcp


# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8081/health || exit 1

# Run the component
CMD ["/app/vigile"]

# Multi-stage Dockerfile for IMS Core
# Supports both Docker and Podman with staged builds

FROM golang:1.26 AS builder

WORKDIR /app

# Allow override of goproxy
ARG goproxy=https://proxy.golang.org
ENV GOPROXY=${goproxy}

# Install build dependencies
RUN apt-get update && apt-get install -y git && rm -rf /var/lib/apt/lists/*

# Copy Go module files
COPY go.mod go.sum* ./

# Cache dependencies
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Copy source code
COPY . .

# Build arguments
ARG package=./cmd/ims
ARG ldflags="-s -w -X main.version=dev -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)"
ARG version=dev

# Generate go.sum if needed
RUN go mod tidy

# Build binary
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags "${ldflags} -X main.version=${version} -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    -o /out/ims-core ${package}

# Final stage: distroless
FROM gcr.io/distroless/base-debian12

WORKDIR /

# Copy binary
COPY --from=builder /out/ims-core /ims-core

# Expose ports
# SIP: 5060 (UDP/TCP), 5061 (TLS)
# HTTP API: 8080
# Metrics: 9443
EXPOSE 5060/udp 5060/tcp 5061/tcp 8080 9443

USER 65532:65532

ENTRYPOINT ["/ims-core"]

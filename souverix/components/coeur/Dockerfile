# Multi-stage build for Souverix Coeur component
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install git and swag CLI
RUN apk add --no-cache git && \
    go install github.com/swaggo/swag/cmd/swag@latest

# Allow overload of goproxy
ARG goproxy=https://proxy.golang.org
ENV GOPROXY=$goproxy

# Copy Go mod files
COPY app/go.mod app/go.sum ./
RUN go mod download

# Copy source code
COPY app/ ./

# Generate Swagger docs for diagnostic endpoints (if main.go has swagger annotations)
RUN if grep -q "@title\|@swagger" main.go 2>/dev/null; then \
    /go/bin/swag init --generalInfo main.go --output docs --parseDependency --parseInternal || echo "Swagger generation skipped"; \
    else echo "No Swagger annotations found, skipping"; \
    fi

# Build arguments
ARG VERSION=dev
ARG BUILD_TIME
ARG GIT_COMMIT

# Build flags
ARG package=.
ARG ldflags="-X main.version=${VERSION} -X main.buildTime=${BUILD_TIME} -X main.gitCommit=${GIT_COMMIT}"

# Build with cache mounts for faster rebuilds
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags "${ldflags} -extldflags '-static'" \
    -o /out/coeur ${package}

# Runtime stage - distroless for security
FROM gcr.io/distroless/base-debian12:nonroot

WORKDIR /

# Copy binary from builder
COPY --from=builder /out/coeur /coeur

# Copy Swagger docs
COPY --from=builder /app/docs /docs

# Expose ports
EXPOSE 5060/tcp 5060/udp 5061/tcp 8081/tcp

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD ["/coeur", "--health-check"] || exit 1

# Run as non-root user (distroless default)
USER nonroot:nonroot

# Entrypoint
ENTRYPOINT ["/coeur"]

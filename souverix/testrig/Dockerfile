# Test Rig Container
# Contains all testing tools and orchestration

FROM golang:1.23 AS builder

WORKDIR /app

# Install test tools
RUN apt-get update && apt-get install -y \
    sipgrep \
    sipsak \
    sipp \
    curl \
    jq \
    && rm -rf /var/lib/apt/lists/*

# Copy test code
COPY testrig/ ./testrig/
COPY go.mod go.sum ./
COPY internal/ ./internal/
COPY pkg/ ./pkg/

# Build test rig
RUN go build -o /out/testrig ./testrig

FROM gcr.io/distroless/base-debian12

WORKDIR /

COPY --from=builder /out/testrig /testrig
COPY --from=builder /usr/bin/sipp /usr/bin/sipp
COPY --from=builder /usr/bin/sipsak /usr/bin/sipsak

ENTRYPOINT ["/testrig"]

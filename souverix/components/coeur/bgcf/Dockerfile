# Multi-stage build for Souverix bgcf subcomponent
FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git &&     go install github.com/swaggo/swag/cmd/swag@latest

ARG goproxy=https://proxy.golang.org
ENV GOPROXY=$goproxy

# Copy common module for local replace directive (development)
COPY common/ ./common/

# Copy app files
COPY components/coeur/bgcf/app/go.mod components/coeur/bgcf/app/go.sum* ./
RUN go mod download && go mod tidy

COPY components/coeur/bgcf/app/ ./

RUN if grep -q "@title\|@swagger" main.go 2>/dev/null; then     /go/bin/swag init --generalInfo main.go --output docs --parseDependency --parseInternal || echo "Swagger skipped";     else echo "No Swagger annotations";     fi

ARG VERSION=dev
ARG BUILD_TIME
ARG GIT_COMMIT
ARG package=.
ARG ldflags="-X main.version=${VERSION} -X main.buildTime=${BUILD_TIME} -X main.gitCommit=${GIT_COMMIT}"

RUN --mount=type=cache,target=/root/.cache/go-build     --mount=type=cache,target=/go/pkg/mod     CGO_ENABLED=0 GOOS=linux GOARCH=amd64     go build -ldflags "${ldflags} -extldflags '-static'"     -o /out/bgcf ${package}

FROM gcr.io/distroless/base-debian12:nonroot

WORKDIR /
COPY --from=builder /out/bgcf /bgcf
COPY --from=builder /app/docs /docs

EXPOSE 8081/tcp

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3   CMD ["/bgcf", "--health-check"] || exit 1

USER nonroot:nonroot
ENTRYPOINT ["/bgcf"]

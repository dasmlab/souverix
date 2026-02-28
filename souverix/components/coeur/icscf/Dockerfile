# Multi-stage build for Souverix icscf subcomponent
FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git &&     go install github.com/swaggo/swag/cmd/swag@latest

ARG goproxy=https://proxy.golang.org
ENV GOPROXY=$goproxy

# Copy app files
COPY components/coeur/icscf/app/go.mod components/coeur/icscf/app/go.sum* ./

# Copy common module - will be used via replace directive
COPY common/ /common/

# Update replace directive to point to /common in container
RUN sed -i 's|replace github.com/dasmlab/souverix/common => ../../../../common|replace github.com/dasmlab/souverix/common => /common|g' go.mod

RUN go mod download && go mod tidy

COPY components/coeur/icscf/app/ ./

RUN if grep -q "@title\|@swagger" main.go 2>/dev/null; then     /go/bin/swag init --generalInfo main.go --output docs --parseDependency --parseInternal || echo "Swagger skipped";     else echo "No Swagger annotations";     fi

ARG VERSION=dev
ARG BUILD_TIME
ARG GIT_COMMIT
ARG package=.
ARG ldflags="-X main.version=${VERSION} -X main.buildTime=${BUILD_TIME} -X main.gitCommit=${GIT_COMMIT}"

RUN --mount=type=cache,target=/root/.cache/go-build     --mount=type=cache,target=/go/pkg/mod     CGO_ENABLED=0 GOOS=linux GOARCH=amd64     go build -ldflags "${ldflags} -extldflags '-static'"     -o /out/icscf ${package}

FROM gcr.io/distroless/base-debian12:nonroot

WORKDIR /
COPY --from=builder /out/icscf /icscf
COPY --from=builder /app/docs /docs

EXPOSE 8082/tcp 9092/tcp 9082/tcp 9182/tcp

USER nonroot:nonroot
ENTRYPOINT ["/icscf"]

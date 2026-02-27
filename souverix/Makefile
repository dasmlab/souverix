.PHONY: test test-coverage test-unit test-integration build

# Test targets
test:
	go test ./... -v

test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

test-unit:
	go test ./... -short -v

test-integration:
	go test ./... -tags=integration -v

# Build targets
build:
	go build ./cmd/ims/...

# Lint targets
lint:
	golangci-lint run

# Format targets
fmt:
	go fmt ./...

# Vet targets
vet:
	go vet ./...

# All checks
check: fmt vet lint test

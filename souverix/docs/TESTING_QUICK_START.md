# Testing Quick Start Guide

## Quick Reference

### Run All Tests
```bash
make test
```

### Run with Coverage
```bash
make test-coverage
# View: open coverage.html
```

### Run Specific Component
```bash
go test ./internal/stir/... -v
go test ./internal/sbc/... -v
go test ./internal/ibcf/... -v
```

## Container Diagnostics

### Start Container
```bash
docker run -d \
  -p 8080:8080 \
  -e DIAGNOSTICS_ENABLED=true \
  -e X-Diagnostic-Role=diagnostics \
  ims-core:local
```

### Run Diagnostics
```bash
# Health check
curl http://localhost:8080/health

# System status (requires role)
curl -H "X-Diagnostic-Role: diagnostics" \
  http://localhost:8080/diagnostics/status

# STIR/SHAKEN tests
curl -H "X-Diagnostic-Role: diagnostics" \
  http://localhost:8080/diagnostics/test/stir/sign

curl -H "X-Diagnostic-Role: diagnostics" \
  http://localhost:8080/diagnostics/test/stir/verify
```

## Test Catalogs

- **IBCF/SIG-GW**: 50+ test cases (`docs/TEST_CATALOG.md`)
- **STIR/SHAKEN**: 30 test cases (`docs/STIR_SHAKEN_TEST_CATALOG.md`)

## Test Areas

### IBCF Tests
- SIG-*: SIP signaling
- TOP-*: Topology hiding
- SEC-*: Security
- TLS-*: TLS/mTLS
- ROU-*: Routing
- And more...

### STIR/SHAKEN Tests
- ORG-*: Originating signing
- TER-*: Terminating verification
- ATT-*: Attestation policy
- CRT-*: Certificate handling
- KEY-*: Key management
- And more...

## CI/CD

Tests run automatically on:
- Push to main/develop
- Pull requests
- Manual trigger

See `.github/workflows/test.yml` for details.

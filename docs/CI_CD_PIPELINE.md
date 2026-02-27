# CI/CD Pipeline Implementation

## Overview

This document describes the full CI/CD pipeline implementation following our 6-phase SDLC methodology (see [DESIGN_PHILOSOPHY.md](DESIGN_PHILOSOPHY.md)).

## Pipeline Flow

```
Build → Unit Test (Diagnostic API) → Publish → System Test → Publish Stable
```

## Phase Breakdown

### Phase 1: Build
- **Workflow**: `.github/workflows/pcscf-full-pipeline.yml` - `build` job
- **Actions**:
  - Checkout code
  - Set up Docker Buildx
  - Build container image
  - Push to registry (GHCR)
  - Generate metadata and tags
- **Outputs**: Image tag, image digest

### Phase 2: Unit Test (via Diagnostic API)
- **Workflow**: `unit-test` job
- **Actions**:
  - Pull built container
  - Start container with `DIAGNOSTICS_ENABLED=true`
  - Wait for health check
  - Call `/diagnostics/test/run` endpoint
  - Parse test results
  - Upload test artifacts
  - Cleanup container
- **Diagnostic Endpoints**:
  - `GET /health` - Public health check
  - `GET /diagnostics/test/run` - Run unit tests (requires `X-Diagnostic-Role: diagnostics`)
- **Outputs**: Test status, test results

### Phase 3: Publish (Tag for System Test)
- **Workflow**: `publish` job
- **Actions**:
  - Create test tag: `test-YYYYMMDD-HHMMSS-SHA`
  - Tag and push image
  - Prepare for system testing
- **Outputs**: Published tag

### Phase 4: System Test (Phantom/Trigger)
- **Workflow**: `system-test` job
- **Actions**:
  - Trigger system test deployment (phantom/placeholder)
  - In production: Deploy to test cluster (K8s/OCP)
  - In production: Trigger AAP/EDA workflow
  - In production: Run system integration tests
  - Upload system test results
- **Note**: Currently a phantom trigger - ready for integration with actual test infrastructure

### Phase 5: Publish Stable
- **Workflow**: `publish-stable` job
- **Actions**:
  - Tag as `stable-YYYYMMDD`
  - Tag as `released-SHA`
  - Push stable tags
  - Create GitHub release
- **Condition**: Only runs on `main` branch after all tests pass

## Component Structure

### P-CSCF Diagnostic Server

**File**: `internal/coeur/pcscf/diagnostics.go`

**Endpoints**:
- `GET /health` - Health check (public)
- `GET /diagnostics/status` - Component status (requires diagnostic role)
- `GET /diagnostics/info` - Component information
- `GET /diagnostics/metrics` - Metrics
- `GET /diagnostics/test/run` - **Run unit tests** (main endpoint for CI/CD)
- `GET /diagnostics/test/results` - Test results

**Configuration**:
- `DIAGNOSTICS_ENABLED=true` - Enable diagnostic server
- `DIAGNOSTICS_ADDR=:8081` - Diagnostic server address (default: :8081)
- `X-Diagnostic-Role: diagnostics` - Header required for diagnostic endpoints

### Dockerfile

**File**: `Dockerfile.pcscf`

Multi-stage build:
1. Builder stage: Compile Go binary
2. Runtime stage: Alpine-based container with health checks

**Ports**:
- `5060/udp` - SIP UDP
- `5060/tcp` - SIP TCP
- `5061/tcp` - SIP TLS
- `8081/tcp` - Diagnostic API

## Workflow Triggers

The pipeline triggers on:
- Push to paths:
  - `internal/coeur/pcscf/**`
  - `internal/common/**`
  - `.github/workflows/pcscf-full-pipeline.yml`
- Pull requests (same paths)
- Manual dispatch (`workflow_dispatch`)

## Self-Hosted Runners

All jobs run on `self-hosted` runners. Ensure:
- Docker is installed and running
- Docker Buildx is available
- Access to GHCR (GitHub Container Registry)
- `GITHUB_TOKEN` secret is configured

## Test Results

Test results are uploaded as artifacts:
- `unit-test-results` - JSON test results from diagnostic API
- `system-test-results` - System test results (when implemented)

## Badges

README badges point to the full pipeline:
- Build badge: Shows overall pipeline status
- Test badge: Shows test phase status

## Template for Other Components

This pipeline serves as a **boilerplate template** for all components:

1. Copy `.github/workflows/pcscf-full-pipeline.yml`
2. Update `COMPONENT` env var
3. Update `IMAGE_NAME` to match component
4. Create component-specific Dockerfile
5. Implement diagnostic server in component
6. Update path triggers

## Future Enhancements

- [ ] Integrate with actual K8s/OCP test cluster
- [ ] Integrate with AAP (Ansible Automation Platform)
- [ ] Integrate with EDA (Event-Driven Ansible)
- [ ] Add code coverage reporting
- [ ] Add security scanning
- [ ] Add performance benchmarks
- [ ] Add stability tests
- [ ] Add release validation tests

## References

- [Design Philosophy](DESIGN_PHILOSOPHY.md) - SDLC methodology
- [Testing Framework](TESTING_SUMMARY.md) - Test infrastructure
- [Component Architecture](ARCHITECTURE.md) - Component design

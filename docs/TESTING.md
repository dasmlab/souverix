# Testing Framework Documentation

## Overview

Comprehensive testing framework for IMS Core with unit tests, container tests, system tests, and integration with CI/CD pipelines.

## Test Structure

### Unit Tests

Every `.go` file has a corresponding `*_test.go` file with:
- Code coverage targets (aim for >80%)
- Positive and negative test cases
- Edge case handling
- Performance benchmarks

**Example:**
```bash
go test ./internal/sip/... -v -cover
```

### Container Tests

Built-in diagnostic APIs in every container:
- Role-based access control
- Component health checks
- Self-test capabilities
- Configuration validation

**Endpoints:**
- `/health` - Public health check
- `/diagnostics/status` - System status (role required)
- `/diagnostics/test/*` - Component tests (role required)

### System Tests

Orchestrated test containers that:
- Deploy test scenarios
- Validate integration
- Perform load testing
- Execute chaos experiments

**Test Rig Container:**
- Contains all testing tools
- Orchestrates test execution
- Collects results
- Reports metrics

## Test Catalogs

### IBCF/SIG-GW Test Catalog

See [TEST_CATALOG.md](TEST_CATALOG.md) for comprehensive test cases covering:
- SIP signaling (SIG-*)
- Topology hiding (TOP-*)
- Security (SEC-*)
- TLS/mTLS (TLS-*)
- Routing (ROU-*)
- NAT/SDP (NAT-*)
- Media (MED-*)
- Interoperability (INT-*)
- Observability (OAM-*)
- Resilience (HAZ-*)
- Chaos (CHA-*)

### STIR/SHAKEN Test Catalog

See [STIR_SHAKEN_TEST_CATALOG.md](STIR_SHAKEN_TEST_CATALOG.md) for STIR/SHAKEN specific tests covering:
- ORG-*: Originating Signing (6 tests)
- TER-*: Terminating Verification (4 tests)
- ATT-*: Attestation Policy (3 tests)
- CRT-*: Certificate Handling (5 tests)
- KEY-*: Key Management (2 tests)
- POL-*: Policy Enforcement (1 test)
- NET-*: Cross-Network / Peering (4 tests)
- RES-*: Resilience & Failure (1 test)
- SEC-*: Security Hardening (3 tests)
- OBS-*: Observability (1 test)

**Total: 30 STIR/SHAKEN test cases**

### Lawful Intercept & Emergency Services Test Catalog

See [LI_EMERGENCY_TEST_CATALOG.md](LI_EMERGENCY_TEST_CATALOG.md) for LI and Emergency tests covering:
- LI-C-*: Lawful Intercept – Control Plane (6 tests)
- LI-M-*: Lawful Intercept – Media Plane (4 tests)
- LI-A-*: Lawful Intercept – Audit & Logging (2 tests)
- EMR-R-*: Emergency Routing (3 tests)
- EMR-P-*: Emergency Policy (3 tests)
- EMR-M-*: Emergency Media (1 test)
- EMR-L-*: Emergency Location Handling (2 tests)
- EMR-F-*: Emergency Failure Handling (1 test)

**Total: 22 LI/Emergency test cases**

### PIXIT (Protocol Implementation eXtra Information for Testing)

See [PIXIT.md](PIXIT.md) for execution parameters and environment control:
- Timer settings (T1, T2, Timer B/F, Session Timer)
- Codec policies (Audio/Video, SRTP, DTMF)
- TLS/Security profiles (TLS version, mTLS, OCSP)
- Peer profiles (Transport, Auth, STIR trust)
- STIR controls (Attestation, iat skew, cert cache)
- LI controls (Mode, Mediation IP, Overload policy)
- Emergency controls (Numbers, Priority, Overrides)
- Chaos injection (Packet loss, Jitter, Service crashes)

**PIXIT Configuration Files**:
- `testrig/pixit/default.yaml` - Default baseline
- `testrig/pixit/str-012.yaml` - STIR OCSP test example
- `testrig/pixit/lie-107.yaml` - Emergency failover test example
- `testrig/pixit/sec-002.yaml` - Rate limit flood test example

Each test includes:
- Positive case
- Negative case
- Load case
- Scale/ramp case
- Chaos case

## Running Tests

### Local Development

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run unit tests only
make test-unit

# Run integration tests
make test-integration
```

### Container Tests

```bash
# Start container with diagnostics enabled
docker run -e DIAGNOSTICS_ENABLED=true \
  -e X-Diagnostic-Role=diagnostics \
  ims-core:local

# Run diagnostics
curl -H "X-Diagnostic-Role: diagnostics" \
  http://localhost:8080/diagnostics/test/sip
```

### System Tests

```bash
# Run test rig container
docker run --network host \
  -v $(pwd)/testrig/config:/config \
  ims-testrig:latest \
  -suite ibcf \
  -config /config/test-config.yaml
```

## Production Readiness Documentation

### KPI/SLO Matrix
See [KPI_SLO_MATRIX.md](KPI_SLO_MATRIX.md) for service-level objectives and acceptance criteria:
- Core Signaling KPIs (8 metrics)
- STIR/SHAKEN KPIs (8 metrics)
- Lawful Intercept KPIs (8 metrics)
- Emergency KPIs (8 metrics)
- Resilience KPIs (8 metrics)
- Observability KPIs (7 metrics)
- Security KPIs (5 metrics)
- SLO Gate Criteria

### Interconnect Certification
See [INTERCONNECT_CERTIFICATION.md](INTERCONNECT_CERTIFICATION.md) for end-to-end certification plan:
- 5 Certification Phases
- Functional Validation
- Interconnect Scenarios
- Negative & Abuse Testing
- Load & Soak Testing
- Failover Testing

### OpenShift CNF Test Harness
See [OPENSHIFT_CNF_HARNESS.md](OPENSHIFT_CNF_HARNESS.md) for test harness architecture:
- Traffic Generation (SIPp, Custom)
- Chaos Engineering (Litmus)
- Metrics Stack (Prometheus, Grafana)
- Scaling Strategy (HPA, SR-IOV)
- CI/CD Integration

### Compliance Checklist
See [COMPLIANCE_CHECKLIST.md](COMPLIANCE_CHECKLIST.md) for regulatory gate review:
- Lawful Intercept Compliance (11 requirements)
- Emergency Compliance (11 requirements)
- STIR/SHAKEN Compliance (11 requirements)
- Security Compliance (11 requirements)
- Operational Compliance (10 requirements)
- Resilience Compliance (10 requirements)

## CI/CD Integration

### GitHub Actions

Workflow runs on:
- Push to main/develop
- Pull requests
- Manual trigger

**Steps:**
1. Checkout code
2. Run tests
3. Generate coverage
4. Upload to codecov
5. Build artifacts

### GitLab CI

Integrated with OCP cluster:
- Runs on self-hosted runners
- Executes Ansible playbooks
- Deploys to test environment
- Runs system tests

## Observability in Tests

### Prometheus Metrics

All tests export metrics:
- Test execution time
- Pass/fail counts
- Component health
- Resource usage

### OpenTelemetry Tracing

Distributed tracing for:
- Test execution flow
- Component interactions
- Performance analysis
- Debugging

### Loki Logging

Structured logs:
- Test results
- Component logs
- Error traces
- Performance data

## Test Tools

### Built-in Tools

- SIPp - SIP traffic generation
- sipsak - SIP testing utility
- Custom test harness
- Load generators

### External Tools (Open Source Only)

- Prometheus - Metrics
- Grafana - Visualization
- Loki - Log aggregation
- Jaeger - Tracing

## Coverage Goals

- **Unit Tests**: >80% code coverage
- **Integration Tests**: All critical paths
- **System Tests**: All test catalog cases
- **Chaos Tests**: Resilience scenarios

## Security in Testing

- Role-based access to diagnostics
- No PII in test logs
- Secure test data handling
- Certificate validation in tests

## Future Enhancements

1. AI-powered test generation
2. Automated test case discovery
3. Predictive failure analysis
4. Continuous test optimization

# PIXIT Usage Guide

## Overview

PIXIT (Protocol Implementation eXtra Information for Testing) provides execution parameters for test cases. Each test can use default PIXIT values or override them with test-specific configurations.

## Quick Start

### Using Default PIXIT

```bash
# Run test with default PIXIT
./testrig -test STR-001
```

### Using Test-Specific PIXIT

```bash
# Run test with PIXIT file
./testrig -test STR-012 -pixit-file testrig/pixit/str-012.yaml

# Override specific parameters
./testrig -test STR-012 \
  -pixit ocsp_mode=hard-fail \
  -pixit chaos.ocsp_responder_disabled=true
```

## PIXIT File Structure

PIXIT files are YAML format:

```yaml
pixit:
  timers:
    t1: 500ms
    max_cps: 2000
  
  tls:
    ocsp_mode: "hard-fail"
  
  chaos:
    ocsp_responder_disabled: true
```

## Test-Specific PIXIT Files

### STR-012: OCSP Validation Test

**File**: `testrig/pixit/str-012.yaml`

**Key Overrides**:
- `tls.ocsp_mode: "hard-fail"` - Force hard-fail mode
- `chaos.ocsp_responder_disabled: true` - Disable OCSP responder

**Expected Behavior**:
- Calls fail verification
- Hard-fail blocks call
- Metrics increment verification_fail counter

### LIE-107: Emergency During Failover

**File**: `testrig/pixit/lie-107.yaml`

**Key Overrides**:
- `timers.max_cps: 5000` - Background load
- `emergency.stir_override: true` - Bypass STIR
- `chaos.kill_primary_instance: true` - Kill primary
- `chaos.active_active_nodes: 2` - Active/Active setup

**Expected Behavior**:
- Emergency call continues
- No PDD violation
- No STIR enforcement block

### SEC-002: Rate Limit Flood Protection

**File**: `testrig/pixit/sec-002.yaml`

**Key Overrides**:
- `timers.max_cps: 2000` - Legitimate max
- `peer.max_cps: 1000` - Peer max
- `chaos.burst_cps: 20000` - Attack burst
- `chaos.rate_limit_threshold: 2500` - Rate limit

**Expected Behavior**:
- Legit calls pass
- Excess rejected with 503
- No crash
- CPU stable <85%

## Programmatic Usage

### Load PIXIT in Go

```go
import "github.com/dasmlab/ims/internal/testrig"

// Load default
pixit := testrig.DefaultPIXIT()

// Load from file
pixit, err := testrig.LoadPIXIT("testrig/pixit/str-012.yaml")
if err != nil {
    log.Fatal(err)
}

// Validate
if err := pixit.ValidatePIXIT(); err != nil {
    log.Fatal(err)
}

// Use in test
test := NewTest("STR-012", pixit)
result := test.Run()
```

### Merge PIXIT Configs

```go
// Start with default
pixit := testrig.DefaultPIXIT()

// Load test-specific overrides
overrides, _ := testrig.LoadPIXIT("testrig/pixit/str-012.yaml")

// Merge
pixit.MergePIXIT(overrides)
```

## PIXIT Parameters Reference

### Timers

- `t1`: SIP base retransmission (100-1000ms)
- `t2`: Non-INVITE retrans max (2-8s)
- `timer_b`: INVITE client timeout (16-64s)
- `timer_f`: Non-INVITE timeout (16-64s)
- `session_timer`: re-INVITE refresh (90-3600s)
- `max_dialogs`: Capacity test (1k-100k)
- `max_cps`: Load envelope (100-20k)

### TLS/Security

- `tls.version`: TLS versions (["1.2", "1.3"])
- `tls.mtls_required`: Require mTLS (true/false)
- `tls.ocsp_mode`: OCSP mode ("hard-fail" / "soft-fail")
- `tls.stir_enforcement`: STIR enforcement ("soft" / "hard")

### STIR

- `stir.attestation_policy`: Attestation ("auto", "A", "B", "C")
- `stir.iat_skew`: iat skew tolerance (0-300s)
- `stir.cert_cache_ttl`: Cert cache TTL (1h-72h)

### Emergency

- `emergency.numbers`: Emergency numbers (["911", "112"])
- `emergency.stir_override`: Bypass STIR (true/false)
- `emergency.fraud_override`: Bypass fraud (true/false)

### Chaos

- `chaos.packet_loss`: Packet loss (0-20%)
- `chaos.jitter`: Jitter (0-200ms)
- `chaos.dns_failure`: DNS failure (true/false)
- `chaos.signing_service_crash`: Kill signing service (true/false)

## Best Practices

1. **Use test-specific PIXIT files** for complex test scenarios
2. **Document overrides** in test catalog comments
3. **Validate PIXIT** before test execution
4. **Log PIXIT values** in test results for traceability
5. **Version control** PIXIT files alongside test code

## Traceability

Every test execution logs:
- PIXIT file used
- All parameter values
- Overrides applied
- Validation results

This ensures full traceability from test catalog → PIXIT → execution → results.

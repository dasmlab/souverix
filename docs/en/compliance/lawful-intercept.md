# Lawful Intercept & Emergency Services Implementation

## Overview

Implementation of Lawful Intercept (LI) and Emergency Services for carrier-grade IMS deployments, compliant with 3GPP standards and regulatory requirements.

## Components

### Lawful Intercept (`internal/li/`)

#### `intercept.go`
- **InterceptController**: Manages active warrants and interception
- **MediationDevice**: Interface for LI Mediation Device (MD)
- **AuditLogger**: Tamper-evident audit logging

**Features**:
- Warrant activation/deactivation
- Signaling interception
- Media interception (if warrant type includes media)
- Multi-target support
- Audit trail

### Emergency Services (`internal/emergency/`)

#### `emergency.go`
- **EmergencyDetector**: Detects emergency numbers (911/112/etc.)
- **EmergencyRouter**: Routes emergency calls to PSAP
- **EmergencyPolicy**: Enforces emergency bypass policies
- **EmergencyLocationHandler**: Handles location information

**Features**:
- Emergency number detection
- PSAP routing
- Priority handling
- Location preservation
- Bypass all restrictions

## Integration

### SBC Integration

Emergency calls are processed **before** any restrictions:
1. Emergency detection (highest priority)
2. Bypass rate limiting
3. Bypass STIR verification (if configured)
4. Bypass fraud detection
5. Route to PSAP
6. Preserve location

### IBCF Integration

Lawful Intercept is integrated at IBCF boundary:
1. Check if call involves intercept target
2. Mirror signaling to LI MD
3. Duplicate media (if warrant includes media)
4. Maintain intercept across interconnect

## Configuration

### Lawful Intercept

```bash
# Enable LI
LI_ENABLED=true
LI_MEDIATION_DEVICE=https://li-md.example.com
LI_HANDOVER_INTERFACE=https://li-hi.example.com
LI_AUDIT_LOGGING=true
```

### Emergency Services

```bash
# Enable Emergency (enabled by default)
EMERGENCY_ENABLED=true
EMERGENCY_BYPASS_STIR=true
EMERGENCY_BYPASS_FRAUD=true
EMERGENCY_BYPASS_RATE_LIMIT=true
```

## Test Coverage

### Lawful Intercept Tests (12 test cases)
- LIE-001 to LIE-012
- Control plane, media plane, audit logging
- Multi-target, interconnect, overload scenarios

### Emergency Services Tests (10 test cases)
- LIE-101 to LIE-110
- Routing, policy, location, failure handling

## Critical Requirements

### Emergency Calls

**MUST NEVER be blocked by**:
- ✅ STIR failure
- ✅ Fraud detection
- ✅ Billing restrictions
- ✅ Rate limiting
- ✅ Policy enforcement

### Lawful Intercept

**MUST**:
- ✅ Be silent to user
- ✅ Not alter signaling behavior
- ✅ Not degrade call quality
- ✅ Be auditable and tamper-evident
- ✅ Persist across re-INVITE, transfer, interconnect

## Regulatory Compliance

- **3GPP TS 33.107**: IMS Lawful Intercept
- **3GPP TS 23.167**: Emergency Sessions in IMS
- **ETSI TS 102 232**: Handover Interfaces
- **FCC Requirements**: 911/E911 (US)
- **CRTC Requirements**: 911 (Canada)
- **EU Requirements**: 112

## Operational KPIs

- Emergency PDD ≤ regulatory target
- Intercept latency overhead ≤ X ms
- 0% emergency drop during failover
- 100% intercept coverage for active warrants
- No topology leakage in emergency routing

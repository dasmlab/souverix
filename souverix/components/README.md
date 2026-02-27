# Internal Code Structure

## Overview

This directory contains all internal (non-exported) code for the Souverix platform, organized by:
- **IMS Core Nodes** (3GPP standard nodes)
- **Supporting Features** (cross-cutting concerns)
- **Supporting Infrastructure** (platform services)
- **Shared Utilities** (common code)

---

## Directory Structure

### IMS Core Nodes

#### `coeur/` - Souverix Coeur (IMS Core)
Contains all IMS core signaling nodes:

- **`pcscf/`** - Proxy CSCF (first contact point for UE)
- **`icscf/`** - Interrogating CSCF (inter-domain routing, S-CSCF selection)
- **`scscf/`** - Serving CSCF (core session control, service logic)
- **`bgcf/`** - Breakout Gateway Control Function (PSTN routing)
- **`mgcf/`** - Media Gateway Control Function (SIP to ISUP conversion)
- **`hss/`** - Home Subscriber Server / UDM (subscriber data store)

#### `relais/` - Souverix Relais (Media Gateway)
- **`mgw/`** - Media Gateway (RTP to TDM conversion)

### Border Control

#### `rempart/` - Souverix Rempart (SIG-GW/IBCF)
Border control and SIP gateway functionality:

- **`ibcf/`** - Interconnection Border Control Function
- **`sbc/`** - Session Border Controller
- **`gateway/`** - Gateway functions

**Note**: Rempart is not an IMS node itself, but acts on behalf of nodes.

### Supporting Features

#### `features/` - Cross-Cutting Features
Features that work across multiple nodes:

- **`stir/`** - STIR/SHAKEN (caller ID authentication)
- **`emergency/`** - Emergency services (911/112 routing)
- **`li/`** - Lawful Intercept (signaling/media duplication)

### Supporting Infrastructure

#### `autorite/` - Souverix Autorite (PKI/CA)
- **`ca/`** - Certificate Authority
- **`vault/`** - Vault integration
- **`hsm/`** - HSM integration

#### `vigie/` - Souverix Vigie (AI Intelligence)
- **`ai/`** - AI hooks and integration

#### `vigile/` - Souverix Vigile (Observability)
- **`metrics/`** - Prometheus metrics
- **`logging/`** - Structured logging
- **`tracing/`** - OpenTelemetry tracing

#### `gouverne/` - Souverix Gouverne (Policy Control)
- **`policy/`** - Policy management

#### `federation/` - Souverix Federation (Inter-domain)
- **`peering/`** - Peering control

### Shared Utilities

#### `common/` - Shared Utilities
Code used by multiple components:

- **`node/`** - Base node interface and boilerplate
- **`sip/`** - SIP message handling
- **`config/`** - Configuration management
- **`logging/` - Logging utilities
- **`metrics/`** - Metrics utilities

---

## Node Boilerplate

All IMS core nodes implement the `node.Node` interface defined in `common/node/node.go`:

```go
type Node interface {
    Name() string
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
    Health() HealthStatus
    Metrics() Metrics
}
```

Each node embeds `node.BaseNode` which provides:
- Health status tracking
- Metrics collection
- Common lifecycle management

---

## Migration Status

### ✅ Completed
- Directory structure created
- Base node interface defined
- Node boilerplate created for all IMS core nodes:
  - P-CSCF
  - I-CSCF
  - S-CSCF
  - BGCF
  - MGCF
  - HSS
  - MGW

### 🔄 Pending Migration
- `internal/store/hss.go` → `internal/coeur/hss/hss.go`
- `internal/sbc/` → `internal/rempart/sbc/`
- `internal/ibcf/` → `internal/rempart/ibcf/`
- `internal/stir/` → `internal/features/stir/`
- `internal/emergency/` → `internal/features/emergency/`
- `internal/li/` → `internal/features/li/`
- `internal/sip/` → `internal/common/sip/`
- `internal/config/` → `internal/common/config/`
- `internal/zta/` → `internal/autorite/ca/`
- `internal/ai/` → `internal/vigie/ai/`
- `internal/metrics/` → `internal/vigile/metrics/`
- `internal/logutil/` → `internal/vigile/logging/`

---

## Development Guidelines

1. **All IMS core nodes** should be in `coeur/` or `relais/`
2. **All supporting features** should be in `features/`
3. **All shared utilities** should be in `common/`
4. **All nodes** must implement the `node.Node` interface
5. **All nodes** should embed `node.BaseNode` for common functionality

---

## Next Steps

1. Migrate existing code to proper locations
2. Update all imports across the codebase
3. Implement actual node functionality (currently just boilerplate)
4. Add unit tests for each node
5. Integrate nodes with supporting features (STIR, Emergency, LI)

# Internal Code Structure - IMS Core Nodes vs Supporting Features

## Overview

This document defines the proper organization of `internal/` to separate:
- **IMS Core Nodes** (actual 3GPP IMS nodes)
- **Supporting Features** (cross-cutting concerns, not nodes)
- **Supporting Infrastructure** (platform services)

---

## IMS Core Nodes (3GPP Standard Nodes)

These are the actual IMS nodes defined by 3GPP standards:

### Souverix Coeur (IMS Core)
**Location**: `internal/coeur/`

**Nodes:**
- **P-CSCF** (Proxy CSCF) - `internal/coeur/pcscf/`
- **I-CSCF** (Interrogating CSCF) - `internal/coeur/icscf/`
- **S-CSCF** (Serving CSCF) - `internal/coeur/scscf/`
- **BGCF** (Breakout Gateway Control Function) - `internal/coeur/bgcf/`
- **MGCF** (Media Gateway Control Function) - `internal/coeur/mgcf/`
- **HSS/UDM** (Home Subscriber Server / Unified Data Management) - `internal/coeur/hss/`

### Souverix Relais (Media Gateway)
**Location**: `internal/relais/`

**Nodes:**
- **MGW** (Media Gateway) - `internal/relais/mgw/`

### Souverix Rempart (SIG-GW/IBCF)
**Location**: `internal/rempart/`

**Note**: Rempart acts on behalf of nodes but is not itself an IMS node. It's a border control function.

---

## Supporting Features (Not IMS Nodes)

These are cross-cutting features that work across nodes:

### STIR/SHAKEN
**Location**: `internal/features/stir/`
- Not a node, but a feature used by Rempart and Coeur
- PASSporT generation/verification
- Certificate management (via Autorite)

### Emergency Services
**Location**: `internal/features/emergency/`
- Not a node, but a feature used by all nodes
- Emergency number detection
- PSAP routing
- Priority handling

### Lawful Intercept
**Location**: `internal/features/li/`
- Not a node, but a feature used by all nodes
- Signaling/media duplication
- Warrant management

---

## Supporting Infrastructure (Platform Services)

These are platform-level services, not IMS nodes:

### Souverix Autorite (PKI/CA)
**Location**: `internal/autorite/`
- Certificate authority
- HSM integration
- Certificate lifecycle

### Souverix Vigie (AI Intelligence)
**Location**: `internal/vigie/`
- AI hooks and integration
- Fraud detection
- Anomaly detection

### Souverix Vigile (Observability)
**Location**: `internal/vigile/`
- Metrics (Prometheus)
- Logging (Logrus)
- Tracing (OpenTelemetry)

### Souverix Gouverne (Policy Control)
**Location**: `internal/gouverne/`
- Policy management
- Configuration orchestration

### Souverix Federation (Inter-domain)
**Location**: `internal/federation/`
- Cross-domain peering
- Trust mapping

---

## Shared Utilities

**Location**: `internal/common/`

- **SIP Handling** - `internal/common/sip/` (used by all nodes)
- **Config** - `internal/common/config/` (shared configuration)
- **Logging** - `internal/common/logging/` (shared logging)
- **Metrics** - `internal/common/metrics/` (shared metrics)

---

## Proposed Structure

```
internal/
├── coeur/                    # Souverix Coeur (IMS Core)
│   ├── pcscf/               # P-CSCF node
│   ├── icscf/               # I-CSCF node
│   ├── scscf/               # S-CSCF node
│   ├── bgcf/                # BGCF node
│   ├── mgcf/                # MGCF node
│   └── hss/                 # HSS/UDM node
│
├── relais/                   # Souverix Relais (Media Gateway)
│   └── mgw/                 # MGW node
│
├── rempart/                  # Souverix Rempart (SIG-GW/IBCF)
│   ├── ibcf/                # IBCF functionality
│   ├── sbc/                 # SBC functionality
│   └── gateway/             # Gateway functions
│
├── features/                  # Cross-cutting features (not nodes)
│   ├── stir/                # STIR/SHAKEN
│   ├── emergency/           # Emergency services
│   └── li/                  # Lawful intercept
│
├── autorite/                  # Souverix Autorite (PKI/CA)
│   ├── ca/                  # Certificate authority
│   ├── vault/               # Vault integration
│   └── hsm/                 # HSM integration
│
├── vigie/                     # Souverix Vigie (AI Intelligence)
│   └── ai/                  # AI hooks and integration
│
├── vigile/                    # Souverix Vigile (Observability)
│   ├── metrics/             # Prometheus metrics
│   ├── logging/             # Structured logging
│   └── tracing/             # OpenTelemetry tracing
│
├── gouverne/                  # Souverix Gouverne (Policy Control)
│   └── policy/              # Policy management
│
├── federation/                # Souverix Federation (Inter-domain)
│   └── peering/             # Peering control
│
└── common/                    # Shared utilities
    ├── sip/                 # SIP message handling
    ├── config/              # Configuration
    ├── logging/             # Logging utilities
    └── metrics/             # Metrics utilities
```

---

## Migration Plan

1. **Create new structure** with proper node directories
2. **Move existing code** to appropriate locations:
   - `internal/store/hss.go` → `internal/coeur/hss/`
   - `internal/sip/` → `internal/common/sip/`
   - `internal/stir/` → `internal/features/stir/`
   - `internal/emergency/` → `internal/features/emergency/`
   - `internal/li/` → `internal/features/li/`
   - `internal/sbc/` → `internal/rempart/sbc/`
   - `internal/ibcf/` → `internal/rempart/ibcf/`
   - `internal/zta/` → `internal/autorite/ca/`
   - `internal/ai/` → `internal/vigie/ai/`
   - `internal/metrics/` → `internal/vigile/metrics/`
   - `internal/logutil/` → `internal/vigile/logging/`
   - `internal/config/` → `internal/common/config/`
3. **Create node boilerplate** for all IMS core nodes
4. **Update imports** across the codebase

---

## Node Boilerplate Structure

Each IMS core node should follow this structure:

```
internal/coeur/pcscf/
├── pcscf.go              # Main node implementation
├── pcscf_test.go         # Unit tests
├── handler.go            # Request handlers
├── handler_test.go       # Handler tests
└── config.go             # Node-specific config
```

Each node should implement:
- Node interface (start, stop, process message)
- Health checks
- Metrics export
- Configuration loading
- Integration with common utilities

---

## Next Steps

1. Create the new directory structure
2. Migrate existing code to proper locations
3. Create node boilerplate templates
4. Implement base node interface
5. Start building actual IMS nodes (P-CSCF, I-CSCF, S-CSCF, etc.)

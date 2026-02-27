# Internal Directory Reorganization

## Structure Overview

### IMS Core Nodes (`internal/coeur/`)
All actual 3GPP IMS core nodes:
- `node.go` - Base node foundation (Node interface, BaseNode, NodeStatus)
- `pcscf.go` - Proxy CSCF (first contact point for UE)
- `icscf.go` - Interrogating CSCF (inter-domain routing, HSS query)
- `scscf.go` - Serving CSCF (core session control, service logic)
- `bgcf.go` - Breakout Gateway Control Function (PSTN breakout routing)
- `mgcf.go` - Media Gateway Control Function (SIP to ISUP conversion)
- `mgw.go` - Media Gateway (RTP to TDM conversion)

### Souverix Components (Supporting Infrastructure)

#### `internal/rempart/` - SIG-GW / IBCF / Border Control
- `sbc.go` - Session Border Controller
- `ibcf.go` - Interconnection Border Control Function
- `stir.go` - STIR/SHAKEN integration
- `ratelimiter.go` - DoS protection
- `emergency.go` - Emergency handling

#### `internal/priorite/` - Emergency & Priority Services
- `emergency/emergency.go` - Emergency routing and detection

#### `internal/mandat/` - Lawful Intercept
- `li/intercept.go` - Intercept controller

#### `internal/autorite/` - PKI / HSM / Vault
- `zta/ca.go` - Certificate Authority management

#### `internal/vigie/` - AI Intelligence Engine
- `ai/hooks.go` - AI agent hooks and MCP integration

#### `internal/vigile/` - Observability & Audit
- `metrics/metrics.go` - Prometheus metrics
- `logging/logutil/logger.go` - Structured logging
- `diagnostics/diagnostics.go` - Diagnostic APIs

#### `internal/gouverne/` - Policy & Control Plane
- `config/config.go` - Configuration management

#### `internal/relais/` - Media Plane (Future)
- Media relay and RTP anchoring

#### `internal/federation/` - Inter-domain Control (Future)
- Cross-border peering control

### Shared Packages

#### `internal/sip/` - SIP Protocol
- SIP message handling (used by all nodes)

#### `internal/store/` - Data Storage
- `hss.go` - HSS/UDM store

#### `internal/stir/` - STIR/SHAKEN
- `passport.go` - PASSporT token handling
- `acme_cert.go` - ACME certificate management

#### `internal/testrig/` - Test Infrastructure
- Test rigs and PIXIT configuration

## Package Naming

All packages follow the Souverix component naming:
- `package coeur` - IMS core nodes
- `package rempart` - SIG-GW/IBCF
- `package priorite` - Emergency services
- `package mandat` - Lawful intercept
- `package autorite` - PKI/CA
- `package vigie` - AI intelligence
- `package vigile` - Observability
- `package gouverne` - Policy/control

## Import Paths

All imports should use the new structure:
```go
import (
    "github.com/dasmlab/ims/internal/coeur"
    "github.com/dasmlab/ims/internal/rempart"
    "github.com/dasmlab/ims/internal/priorite"
    "github.com/dasmlab/ims/internal/mandat"
    "github.com/dasmlab/ims/internal/autorite"
    "github.com/dasmlab/ims/internal/vigie"
    "github.com/dasmlab/ims/internal/vigile"
    "github.com/dasmlab/ims/internal/gouverne"
)
```

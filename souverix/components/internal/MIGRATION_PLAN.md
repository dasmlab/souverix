# Internal Structure Migration Plan

## Current → Target Mapping

### IMS Core Nodes (Move to coeur/)
- `internal/store/hss.go` → `internal/coeur/hss/hss.go`
- Create new: `internal/coeur/pcscf/`, `icscf/`, `scscf/`, `bgcf/`, `mgcf/`

### Media Gateway (Move to relais/)
- Create new: `internal/relais/mgw/`

### SIG-GW/IBCF (Move to rempart/)
- `internal/sbc/` → `internal/rempart/sbc/`
- `internal/ibcf/` → `internal/rempart/ibcf/`

### Supporting Features (Move to features/)
- `internal/stir/` → `internal/features/stir/`
- `internal/emergency/` → `internal/features/emergency/`
- `internal/li/` → `internal/features/li/`

### Supporting Infrastructure
- `internal/zta/` → `internal/autorite/ca/`
- `internal/ai/` → `internal/vigie/ai/`
- `internal/metrics/` → `internal/vigile/metrics/`
- `internal/logutil/` → `internal/vigile/logging/`

### Shared Utilities (Move to common/)
- `internal/sip/` → `internal/common/sip/`
- `internal/config/` → `internal/common/config/`

### Keep As-Is
- `internal/diagnostics/` → Keep for now
- `internal/testrig/` → Keep for now

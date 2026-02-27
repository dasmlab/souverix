# Souverix Digital Twin Service

**Stateful Graph Model for Fleet Operations**

## Overview

The Digital Twin Service maintains event-sourced state models for all Souverix fleet nodes, enabling graph-based fleet visualization, state machine tracking, and drift analysis.

## Digital Twin Model

### Core Structure

```yaml
DigitalTwin:
  identity:
    immutable_id: "sx-rempart-2026-prod-1-node-42"
    component_type: "souverix-rempart"
    region: "region-a"
    zone: "edge-zone-a1"
    trust_chain: [...]
  
  connectivity:
    state: "live"  # live | delayed | edge-cached | autonomous | disconnected | returning | replaying
    last_sync: "2026-02-27T03:45:00Z"
    sync_debt: "0m"
    message_backlog: 0
    latency_ms: 45
  
  lifecycle:
    phase: "active"  # provisioning | active | maintenance | decommissioning
    
    hardware:
      serial: "HW-2024-Q3-0042"
      manufacturing_origin: "Factory-A"
      batch_id: "BATCH-2024-Q3-001"
      firmware_version: "2.1.3"
      firmware_signature: "sha256:abc..."
    
    software:
      desired_version: "souverix-rempart-1.4.7"
      actual_version: "souverix-rempart-1.4.6"
      drift_score: 1
      sbom_hash: "sha256:def..."
      image_provenance: "signed-by-build-system"
  
  state:
    desired: {...}
    actual: {...}
    drift: {...}
    last_applied: "2026-02-27T03:30:00Z"
  
  policy:
    desired_policy: "policy-v2.1"
    actual_policy: "policy-v2.0"
    compliance_score: 0.95
    stale_by: "15m"
  
  trust:
    certificate_chain: [...]
    attestation_state: "verified"
    last_attestation: "2026-02-27T03:45:00Z"
    sbom_lineage: [...]
  
  supply_chain:
    hardware_metadata: {...}
    software_provenance: {...}
    manufacturing_batch: {...}
    trust_chain: [...]
```

## Event Store

### Event Sourcing

All state changes are event-sourced:

```yaml
Event:
  id: "evt-2026-02-27-034500-001"
  timestamp: "2026-02-27T03:45:00Z"
  node_id: "sx-rempart-2026-prod-1-node-42"
  event_type: "state_transition"
  payload:
    from: "delayed"
    to: "live"
    reason: "connectivity_restored"
```

### Event Types

1. **State Transitions**
   - Connectivity state changes
   - Lifecycle phase changes
   - Operational state changes

2. **Drift Events**
   - Version drift detected
   - Configuration drift detected
   - Policy drift detected

3. **Lifecycle Events**
   - Provisioning started/completed
   - Update initiated/completed
   - Maintenance started/completed
   - Decommissioning started/completed

4. **Connectivity Events**
   - Connection established
   - Connection lost
   - Sync completed
   - Replay started/completed

5. **Trust Events**
   - Attestation performed
   - Certificate rotated
   - SBOM updated
   - Trust chain validated

## State Projection Engine

### Projections

UI builds from projections, not raw events:

1. **Current State Projection**
   - Latest state per node
   - Real-time updates
   - Optimized for queries

2. **Drift Projection**
   - Aggregated drift metrics
   - Per-component drift
   - Cluster-level drift

3. **Connectivity Projection**
   - Connectivity spectrum distribution
   - Sync debt aggregation
   - Backlog statistics

4. **Lifecycle Projection**
   - Lifecycle phase distribution
   - Update progress
   - Maintenance windows

### Projection Updates

```go
type ProjectionEngine struct {
    eventStore    *EventStore
    projections   map[string]*Projection
    updateStream  *UpdateStream
}

type Projection struct {
    name         string
    state        map[string]interface{}
    lastEventID  string
    updateFunc   func(event Event) error
}
```

## Fleet State Graph

### Graph Model

```
Fleet Graph
├── Nodes (Digital Twins)
│   ├── sx-rempart-2026-prod-1-node-42
│   ├── sx-coeur-2026-prod-1-node-43
│   └── sx-relais-2026-prod-1-node-44
│
├── Edges (Relationships)
│   ├── sx-rempart-42 → sx-coeur-43 (signaling)
│   ├── sx-coeur-43 → sx-relais-44 (media)
│   └── sx-rempart-42 → sx-relais-44 (bypass)
│
├── Clusters (Logical Groupings)
│   ├── region-a-cluster
│   ├── edge-zone-a1-cluster
│   └── iBCF-interconnect-cluster
│
└── Zones (Geographic/Logical)
    ├── region-a
    ├── edge-zone-a1
    └── edge-zone-a2
```

### Graph Properties

- **Nodes**: Digital twins with full state
- **Edges**: Relationships (signaling, media, control)
- **Clusters**: Logical groupings (region, zone, function)
- **Zones**: Geographic/logical boundaries
- **Trust Domains**: Security boundaries

## Drift Analyzer

### Drift Detection

```yaml
DriftAnalysis:
  node: "sx-rempart-2026-prod-1-node-42"
  timestamp: "2026-02-27T03:45:00Z"
  
  version_drift:
    component: "souverix-rempart"
    desired: "1.4.7"
    actual: "1.4.6"
    severity: "medium"
    impact: "missing_security_patch"
  
  config_drift:
    objects: 2
    last_sync: "2026-02-27T03:30:00Z"
    stale_by: "15m"
    details:
      - field: "stir.enforcement_mode"
        desired: "hard"
        actual: "soft"
      - field: "rate_limit.max_cps"
        desired: 2000
        actual: 1500
  
  policy_drift:
    policy_version: "policy-v2.1"
    actual_version: "policy-v2.0"
    compliance_score: 0.95
    violations: 1
```

### Drift Aggregation

- Per-node drift
- Per-cluster drift
- Per-zone drift
- Global drift metrics

## Connectivity Ledger

### Connectivity History

Tracks connectivity state transitions:

```yaml
ConnectivityLedger:
  node: "sx-rempart-edge-zone-c-node-15"
  
  history:
    - timestamp: "2026-02-27T00:00:00Z"
      state: "live"
      latency_ms: 45
    
    - timestamp: "2026-02-27T00:15:00Z"
      state: "delayed"
      latency_ms: 5000
    
    - timestamp: "2026-02-27T00:30:00Z"
      state: "edge-cached"
      sync_interval: "15m"
    
    - timestamp: "2026-02-27T01:00:00Z"
      state: "autonomous"
      reason: "connection_lost"
    
    - timestamp: "2026-02-27T04:00:00Z"
      state: "returning"
      reason: "connection_detected"
    
    - timestamp: "2026-02-27T04:05:00Z"
      state: "replaying"
      backlog: 1247
    
    - timestamp: "2026-02-27T04:20:00Z"
      state: "live"
      latency_ms: 50
```

## Trust / SBOM Index

### Trust Chain Tracking

```yaml
TrustIndex:
  node: "sx-rempart-2026-prod-1-node-42"
  
  certificate_chain:
    - issuer: "souverix-autorite-root-ca"
      subject: "sx-rempart-2026-prod-1-node-42"
      valid_until: "2026-08-27T00:00:00Z"
      signature: "sha256:..."
  
  sbom_lineage:
    - component: "souverix-rempart"
      version: "1.4.6"
      sbom_hash: "sha256:def..."
      provenance: "signed-by-build-system"
      dependencies: [...]
  
  attestation:
    last_verified: "2026-02-27T03:45:00Z"
    integrity_state: "verified"
    trust_score: 1.0
```

## Implementation

### Event Store

```go
type EventStore interface {
    Append(event Event) error
    GetEvents(nodeID string, since time.Time) ([]Event, error)
    GetLatestState(nodeID string) (*DigitalTwin, error)
}

type Event struct {
    ID        string
    Timestamp time.Time
    NodeID    string
    Type      string
    Payload   map[string]interface{}
}
```

### State Projection

```go
type ProjectionEngine interface {
    UpdateProjection(event Event) error
    GetProjection(name string) (*Projection, error)
    Subscribe(projectionName string, callback func(*Projection))
}
```

### Digital Twin Service

```go
type DigitalTwinService struct {
    eventStore      EventStore
    projectionEngine ProjectionEngine
    driftAnalyzer   *DriftAnalyzer
    connectivityLedger *ConnectivityLedger
    trustIndex      *TrustIndex
}
```

---

## Integration Points

### With Lifecycle Manager
- Lifecycle state updates
- Hardware/software metadata
- Supply chain tracking

### With Connectivity Manager
- Connectivity state updates
- Sync debt tracking
- Replay progress

### With Multi-Master Controller
- Authority metadata
- Federation relationships
- Regional sync status

---

## Next Steps

1. Design event schema
2. Implement event store
3. Build projection engine
4. Create drift analyzer
5. Design UI for fleet state graph

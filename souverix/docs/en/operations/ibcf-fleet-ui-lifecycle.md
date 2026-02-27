# IBCF Fleet Monitoring + Lifecycle Core Design (v0)

This document defines:
1) Canonical Fleet Object Model (schema-level)
2) UI layering model (global / region / edge) and the projections to render it cleanly
3) Node lifecycle state machine (including autonomous + return replay)

---

## 1) Canonical IBCF Fleet Object Model (Schema-Level)

### Design goals
- Stable identity, even when connectivity is absent
- Represent *desired vs actual*, drift, and authority boundaries
- Support multi-master control planes (regional autonomy) and eventual consistency
- Support supply-chain / provenance tracking as first-class
- Support store-and-forward + replay workflows

### Canonical objects
- **Node**: physical/virtual compute element (edge node, DC node, gateway, appliance)
- **Relay**: comm/control intermediaries (brokers, gateways, sat links, store-and-forward hops)
- **Domain**: trust/authority scope (region, cluster, tenant, air-gap island, etc.)
- **PolicyBundle**: desired-state definition (config + version + constraints)
- **Observation**: telemetry event or snapshot
- **Command**: control intent with delivery semantics (at-most-once / at-least-once)
- **Artifact**: signed software/firmware payload (SBOM + provenance)
- **Timeline**: lifecycle phases and transitions for the node

### Minimal JSON-style schema (v0)
> Treat as a conceptual schema; implement as Go structs + protobuf later.

```json
{
  "Node": {
    "id": "uuid/ulid",
    "name": "string",
    "kind": "EDGE|DC|GATEWAY|RELAY|VIRTUAL",
    "tags": ["string"],
    "domain_id": "uuid/ulid",
    "authority": {
      "master_of_record": "GLOBAL|REGION|LOCAL",
      "controller_ids": ["string"],
      "last_authoritative_update_ts": "rfc3339"
    },

    "connectivity": {
      "mode": "LIVE|DELAYED|EDGE_CACHED|AUTONOMOUS|DISCONNECTED|RETURNING|REPLAYING",
      "last_seen_ts": "rfc3339",
      "link_quality": {
        "rtt_ms": 0,
        "loss_pct": 0,
        "bandwidth_kbps": 0
      },
      "sync_debt": {
        "seconds_behind": 0,
        "backlog_messages": 0,
        "replay_required": false
      }
    },

    "desired_state": {
      "policy_bundle_id": "uuid/ulid",
      "target_version": "semver/string",
      "constraints": {
        "maintenance_window": "string",
        "requires_online": false,
        "min_battery_pct": 0
      }
    },

    "actual_state": {
      "reported_version": "semver/string",
      "runtime": {
        "health": "OK|DEGRADED|FAILED|UNKNOWN",
        "last_heartbeat_ts": "rfc3339",
        "uptime_s": 0
      },
      "inventory": {
        "cpu": "string",
        "ram_gb": 0,
        "storage_gb": 0,
        "gpus": ["string"]
      }
    },

    "drift": {
      "score": 0.0,
      "objects_out_of_spec": 0,
      "details_ref": "uri/string",
      "last_evaluated_ts": "rfc3339"
    },

    "supply_chain": {
      "hardware": {
        "serial": "string",
        "manufacturer": "string",
        "model": "string",
        "bom_ref": "uri/string"
      },
      "firmware": {
        "version": "string",
        "signature": "string",
        "provenance_ref": "uri/string"
      },
      "software": {
        "artifact_id": "uuid/ulid",
        "sbom_ref": "uri/string",
        "attestations": ["uri/string"],
        "signature": "string"
      }
    },

    "lifecycle": {
      "phase": "PROVISIONING|ENROLLED|ACTIVE|UPDATING|DEGRADED|QUARANTINED|DECOMMISSIONING|RETIRED",
      "state": "string (state machine)",
      "last_transition_ts": "rfc3339",
      "reason": "string"
    },

    "observations": {
      "last_snapshot_ref": "uri/string",
      "last_event_ref": "uri/string"
    }
  }
}
```

### Key invariants (non-negotiable)

- **Node.id is immutable** across all lifetimes.
- **Authority must be explicit**: who is allowed to declare "desired" for this node right now.
- **Connectivity mode is a spectrum**, not binary.
- **Drift is a computed projection** (do not store drift as the primary truth; store inputs).

### 1A) Object relationship diagram (Mermaid)

```mermaid
flowchart TB
    Domain[Domain] -->|contains| Node[Node]
    Domain -->|contains| Relay[Relay]
    PolicyBundle[PolicyBundle] -->|applies desired| Node
    Artifact[Artifact] -->|installs/updates| Node
    Node -->|emits| Observation[Observation]
    Node --> SupplyChain[Supply Chain Projection]
    Node --> Drift[Drift Projection]
    Relay -->|transports| Observation
    Relay -->|transports| Command[Command]
    Controller[Controller/Master] -->|issues| Command
    Command -->|applies| Node
```

---

## 2) UI Layering Model (Global / Region / Edge)

### UI goals

- "Simple, clean" at top level, but drill-down depth when needed
- Support massive scale (avoid rendering raw node lists as the primary UX)
- Make autonomy + replay visible without scaring operators
- Make authority boundaries visible (multi-master reality)

### UI layers

#### Layer 0: Global Fleet Overview (executive + SRE)

**Primary questions:**
- How many are live vs delayed vs autonomous vs replaying?
- Which domains are drifting the most?
- Where is the control plane authoritative vs degraded?

**Widgets:**
- Connectivity spectrum distribution
- Drift heat map by domain
- Upgrade rollout status by policy bundle
- "Attention queue" (top N critical anomalies)

#### Layer 1: Domain / Region / Cluster View (operator)

**Primary questions:**
- Which sub-fleets are behind?
- What relay paths are failing?
- Which policy bundles are safe to advance?

**Widgets:**
- Domain topology overlay (relays + edges)
- Sync-debt / backlog timeline
- Domain health summary + anomaly list
- Batch actions scoped by authority and constraints

#### Layer 2: Node Detail View (engineering)

**Primary questions:**
- What is desired vs actual, exactly?
- What drift objects differ?
- What's the replay status and last known good snapshot?
- Is this node trusted and supply-chain clean?

**Widgets:**
- Desired vs actual diff view (structured)
- Lifecycle timeline (state transitions)
- Replay progress + backlog
- Supply-chain panel (SBOM / signatures / provenance)
- Observations timeline (events / snapshots)

### Critical UI technique: projections (not raw metrics)

The UI should render projections from an event store:

- **ConnectivityProjection**
- **DriftProjection**
- **AuthorityProjection**
- **LifecycleProjection**
- **SupplyChainProjection**

Each projection is queryable by:
- `domain_id`
- `tags`
- `policy_bundle_id`
- `time window`

### 2A) UI layering architecture diagram (Mermaid)

```mermaid
flowchart TB
    subgraph Ingest["Ingest + Normalize Layer"]
        EventStream[Event Stream] --> Normalizer[Normalizer]
        Snapshots[Snapshots] --> Normalizer
        CommandAPI[Command API] --> Normalizer
    end
    
    subgraph State["State Backbone Layer"]
        EventStore[(Event Store)]
        SnapshotStore[(Snapshot Store)]
    end
    
    subgraph Projections["Projection Layer (Read)"]
        ConnProj[Connectivity Projection]
        LifeProj[Lifecycle Projection]
        AuthProj[Authority Projection]
        DriftProj[Drift Projection]
        SCProj[Supply Chain Projection]
    end
    
    subgraph APIs["Query + Control APIs"]
        QueryAPI[Query API]
        CommandAPI
    end
    
    subgraph UI["UI Layers"]
        Layer0[Layer 0: Global Overview]
        Layer1[Layer 1: Domain/Region View]
        Layer2[Layer 2: Node Detail View]
    end
    
    Normalizer --> EventStore
    Normalizer --> SnapshotStore
    
    EventStore --> ConnProj
    EventStore --> LifeProj
    EventStore --> AuthProj
    EventStore --> DriftProj
    SnapshotStore --> DriftProj
    EventStore --> SCProj
    SnapshotStore --> SCProj
    
    ConnProj --> QueryAPI
    LifeProj --> QueryAPI
    AuthProj --> QueryAPI
    DriftProj --> QueryAPI
    SCProj --> QueryAPI
    
    QueryAPI --> Layer0
    QueryAPI --> Layer1
    QueryAPI --> Layer2
```

### 2B) Operator "attention queue" model (practical)

A scalable UI needs an inbox. Example attention scoring inputs:

- Drift score
- Sync debt (seconds behind + backlog)
- Authority mismatch (who thinks they're master)
- Health degradation (OK → DEGRADED)
- Supply chain violation (signature missing)

**Output**: top N items per domain + global top N.

---

## 3) Node Lifecycle State Machine (Autonomous + Return Replay)

### Lifecycle phases vs states

- **Phase** is coarse (PROVISIONING, ACTIVE, etc.)
- **State** is precise (LIVE_OK, AUTONOMOUS_BUFFERING, REPLAYING_APPLYING, …)

### Core connectivity states (spectrum)

```
LIVE → DELAYED → EDGE_CACHED → AUTONOMOUS → DISCONNECTED → RETURNING → REPLAYING → LIVE
```

### Key behaviors

- **In AUTONOMOUS**: node continues local policies; records observations + local decisions
- **In RETURNING**: node re-establishes a control channel; does not immediately accept new commands until reconciliation
- **In REPLAYING**: node uploads backlog; control plane applies reconciliation policy
- **After REPLAYING**: node becomes LIVE or QUARANTINED depending on trust/drift results

### 3A) Lifecycle state diagram (Mermaid State Diagram)

```mermaid
stateDiagram-v2
    [*] --> PROVISIONING
    PROVISIONING --> ENROLLED: Enrollment complete
    ENROLLED --> ACTIVE: Activation
    
    ACTIVE --> LIVE_OK: Connectivity live
    ACTIVE --> DELAYED: Latency increases
    ACTIVE --> EDGE_CACHED: Periodic sync
    ACTIVE --> AUTONOMOUS: Connection lost
    
    LIVE_OK --> DELAYED: Latency threshold
    DELAYED --> EDGE_CACHED: Sync interval extends
    EDGE_CACHED --> AUTONOMOUS: Connection lost
    AUTONOMOUS --> RETURNING: Connection detected
    RETURNING --> REPLAYING: Backlog upload
    REPLAYING --> LIVE_OK: Replay complete
    REPLAYING --> QUARANTINED: Drift/trust violation
    
    ACTIVE --> UPDATING: Update initiated
    UPDATING --> ACTIVE: Update complete
    UPDATING --> DEGRADED: Update failed
    
    ACTIVE --> DEGRADED: Health degradation
    DEGRADED --> ACTIVE: Recovery
    DEGRADED --> QUARANTINED: Critical failure
    
    QUARANTINED --> ACTIVE: Remediation complete
    QUARANTINED --> DECOMMISSIONING: Unrecoverable
    
    ACTIVE --> DECOMMISSIONING: Planned decommission
    DECOMMISSIONING --> RETIRED: Decommission complete
    RETIRED --> [*]
```

### 3B) Return + Replay sequence (Mermaid Sequence Diagram)

```mermaid
sequenceDiagram
    participant Node
    participant Relay as Relay/Gateway
    participant Control as Control Plane
    participant Proj as Projection Engine
    
    Note over Node: Node has been autonomous<br/>(backlog queued)
    
    Node->>Relay: Reconnect attempt + identity proof
    Relay->>Control: Forward auth + reconnect notice
    Control->>Relay: Issue reconciliation lock + session token
    Relay->>Node: Session token + replay instruction
    
    Node->>Relay: Upload backlog (events/snapshots) in chunks
    Relay->>Control: Forward backlog
    
    Control->>Control: Validate signatures + order + dedupe
    Control->>Control: Reconcile desired vs actual + compute drift deltas
    Control->>Proj: Emit normalized events for projections
    
    alt Clean reconciliation
        Control->>Relay: Release lock + accept new commands
        Note over Node: Node returns to LIVE_OK
    else Drift/trust violation
        Control->>Relay: Quarantine instruction (limited ops)
        Note over Node: Node enters QUARANTINED
    end
```

---

## Implementation notes (non-rabbit-hole, but important)

- Keep the event store **append-only**; projections are rebuildable.
- Commands require delivery semantics (at-least-once + idempotency keys).
- **AuthorityProjection** prevents split-brain in multi-master (UI must show it).
- Drift is computed from snapshots + desired bundles, not raw metrics.

---

## Appendix: Mermaid support in Markdown

Mermaid supports:

- **flowcharts**: `flowchart LR|TB`
- **state machines**: `stateDiagram-v2`
- **sequence diagrams**: `sequenceDiagram`
- **gantt timelines**: `gantt` (useful for rollout timelines)

All Mermaid diagrams render directly in MkDocs Material when using the `pymdownx.superfences` extension with Mermaid support.

---

## Next Steps

1. Implement Go structs matching the canonical schema
2. Build event store foundation
3. Create projection engine
4. Design UI mockups for each layer
5. Implement state machine engine
6. Build reconciliation logic

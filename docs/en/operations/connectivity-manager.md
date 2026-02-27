# Souverix Connectivity Manager

**Connectivity Spectrum Management for Fleet Operations**

## Overview

The Connectivity Manager handles the spectrum of connectivity states for Souverix fleet nodes, moving beyond binary online/offline to model real-world intermittent connectivity, autonomous operation, and edge scenarios.

## Connectivity Spectrum

### State Definitions

| State | Description | Characteristics | UI Indicator |
|-------|-------------|-----------------|--------------|
| **Live** | Real-time bidirectional | < 1s latency, full sync | Green pulse |
| **Delayed** | Sync with latency | 1s-5m latency, periodic sync | Yellow bar (latency shown) |
| **Edge-cached** | Local cache, periodic sync | 5m-1h sync interval | Blue with sync timer |
| **Autonomous** | Operating independently | No connectivity, store-and-forward | Orange with return ETA |
| **Disconnected** | No connectivity | Unknown return time | Red with last-seen |
| **Returning** | Reconnecting | Actively establishing connection | Yellow pulse |
| **Replaying** | Catching up on backlog | Processing queued messages | Purple progress bar |

### State Machine

```
┌─────────┐
│  Live   │◄─────────────────┐
└────┬────┘                   │
     │ latency increases      │ sync complete
     ▼                        │
┌─────────┐                   │
│ Delayed │                   │
└────┬────┘                   │
     │ sync interval extends  │
     ▼                        │
┌──────────────┐              │
│ Edge-cached  │              │
└────┬─────────┘              │
     │ connection lost        │
     ▼                        │
┌──────────────┐              │
│ Autonomous   │──────────────┘
└────┬─────────┘
     │ connection detected
     ▼
┌──────────┐
│Returning │
└────┬─────┘
     │ connection established
     ▼
┌──────────┐
│Replaying │
└────┬─────┘
     │ backlog processed
     ▼
┌─────────┐
│  Live   │
└─────────┘
```

## Sync Debt

### Definition

Sync debt represents the accumulated drift between a node's last known state and the current desired state.

### Calculation

```yaml
SyncDebt:
  node: "sx-rempart-edge-zone-c-node-15"
  last_sync: "2026-02-27T00:00:00Z"
  current_time: "2026-02-27T03:45:00Z"
  debt_hours: 3.75
  
  components:
    config:
      updates_missed: 3
      last_known: "config-v2.0"
      current_desired: "config-v2.1"
    
    policy:
      updates_missed: 1
      last_known: "policy-v2.0"
      current_desired: "policy-v2.1"
    
    version:
      updates_missed: 0
      last_known: "souverix-1.4.6"
      current_desired: "souverix-1.4.7"
```

### UI Representation

- **Temporal Bar**: Visual representation of sync debt
- **Debt Meter**: Hours/days of accumulated debt
- **Backlog Depth**: Number of queued messages
- **Replay ETA**: Estimated time to catch up

## Message Backlog

### Tracking

When nodes are disconnected or operating autonomously:

1. **Message Queuing**
   - Configuration updates
   - Policy updates
   - Version updates
   - Control plane messages

2. **Backlog Management**
   - Priority ordering
   - Expiration handling
   - Compression
   - Deduplication

3. **Replay Coordination**
   - Sequential replay
   - Conflict resolution
   - State validation
   - Progress tracking

### Backlog Structure

```yaml
MessageBacklog:
  node: "sx-rempart-edge-zone-c-node-15"
  total_messages: 1247
  
  by_type:
    config_updates: 3
    policy_updates: 1
    version_updates: 0
    control_messages: 1243
  
  priority:
    high: 5
    medium: 12
    low: 1230
  
  oldest_message: "2026-02-27T00:00:00Z"
  newest_message: "2026-02-27T03:45:00Z"
```

## Autonomous Operation

### Characteristics

Nodes operating autonomously:
- Continue normal operation
- Store updates locally
- Maintain local state
- Queue messages for replay

### Store-and-Forward

```yaml
AutonomousNode:
  node: "sx-rempart-edge-zone-c-node-15"
  mode: "autonomous"
  entered_at: "2026-02-27T00:00:00Z"
  
  local_state:
    version: "souverix-1.4.6"
    config: "config-v2.0"
    policy: "policy-v2.0"
  
  queued_updates:
    count: 1247
    storage_used: "45MB"
    oldest: "2026-02-27T00:00:00Z"
  
  return_eta: "2026-02-27T04:00:00Z"
  estimated_replay_time: "15m"
```

## Replay State

### Replay Process

When connectivity returns:

1. **Connection Detection**
   - Health check success
   - Authentication established
   - State transition: Returning → Replaying

2. **Backlog Processing**
   - Sequential replay
   - State validation
   - Conflict resolution
   - Progress tracking

3. **Completion**
   - Backlog cleared
   - State synchronized
   - State transition: Replaying → Live

### Replay Progress

```yaml
ReplayProgress:
  node: "sx-rempart-edge-zone-c-node-15"
  state: "replaying"
  started_at: "2026-02-27T04:00:00Z"
  
  progress:
    total: 1247
    processed: 623
    remaining: 624
    percentage: 50
  
  eta: "2026-02-27T04:15:00Z"
  estimated_duration: "15m"
  
  current_operation: "Applying config update #3"
```

## UI Patterns

### Connectivity Spectrum Visualization

**Not:**
- ❌ Red/Green binary indicators
- ❌ Simple "online/offline" status

**Instead:**
- ✅ Temporal bars showing sync debt
- ✅ Color-coded spectrum (green → yellow → orange → red)
- ✅ Sync timer displays
- ✅ Message backlog depth indicators
- ✅ Replay progress bars
- ✅ Return ETA displays

### Fleet View

```
Fleet Connectivity Overview
├── Region A (Live: 45/50, Delayed: 3/50, Autonomous: 2/50)
│   ├── Zone A1: Live (100%)
│   └── Zone A2: Live (48/50), Delayed (2/50)
├── Region B (Live: 30/40, Edge-cached: 8/40, Autonomous: 2/40)
│   └── Zone B1: Edge-cached (100%)
└── Edge Zone C (Autonomous: 15/15)
    └── Sync Debt: 3.75h average
    └── Return ETA: 4h
```

## Implementation

### State Machine Engine

```go
type ConnectivityState string

const (
    StateLive        ConnectivityState = "live"
    StateDelayed     ConnectivityState = "delayed"
    StateEdgeCached  ConnectivityState = "edge-cached"
    StateAutonomous  ConnectivityState = "autonomous"
    StateDisconnected ConnectivityState = "disconnected"
    StateReturning   ConnectivityState = "returning"
    StateReplaying   ConnectivityState = "replaying"
)

type ConnectivityManager struct {
    stateMachine *StateMachine
    syncDebtCalc *SyncDebtCalculator
    backlogMgr   *MessageBacklogManager
    replayCoord  *ReplayCoordinator
}
```

### Sync Debt Calculator

- Tracks time since last sync
- Calculates missed updates
- Estimates replay time
- Updates UI projections

### Message Backlog Manager

- Queues messages during disconnection
- Prioritizes updates
- Manages expiration
- Coordinates replay

### Replay Coordinator

- Sequences replay operations
- Validates state consistency
- Resolves conflicts
- Tracks progress

---

## Integration Points

### With Digital Twin Service
- Connectivity state updates
- Sync debt tracking
- Replay progress

### With Lifecycle Manager
- Autonomous mode transitions
- Update coordination
- State synchronization

### With Multi-Master Controller
- Regional connectivity
- Federation sync
- Authority boundaries

---

## Next Steps

1. Implement connectivity state machine
2. Build sync debt calculator
3. Create message backlog manager
4. Design replay coordinator
5. Build UI for connectivity spectrum visualization

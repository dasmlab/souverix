# Souverix Frontend Architecture

**Modern Frontend Techniques for Fleet Monitoring & Lifecycle Management**

## Overview

This document defines the frontend architecture for Souverix fleet monitoring, leveraging modern techniques (2023-2026) to handle massive-scale fleet visualization, lifecycle management, and multi-master control.

## Current Stack

- **Vite**: Fast dev + modern bundling
- **Quasar**: Component system + layout control
- **CRDs**: Backend truth (Kubernetes Custom Resources)
- **Command UI Pattern**: Frontend → Adapter → K8s
- **Pinia**: State management

## Top 5 Modern Frontend Techniques (2023-2026)

### 1️⃣ Server-Driven UI + Typed Contracts (tRPC / RPC-first models)

#### Why This Matters for Fleet Systems

**Traditional Pattern:**
- REST endpoints
- Frontend fetches resources
- Frontend assembles projections
- Frontend computes drift/sync debt

**Modern Pattern:**
- Typed RPC contracts
- Server returns projection-ready view models
- Frontend renders directly
- No frontend logic duplication

#### Hot Enablers

- **tRPC**: End-to-end type safety
- **GraphQL**: Selective use for complex queries
- **Zod**: Schema inference + validation
- **TypeScript**: Full type safety from backend to frontend

#### Why It's Hot Now

The shift is toward:
- Server computes projections (drift, connectivity spectrum, authority state)
- UI consumes strongly typed view models
- No frontend logic duplication

#### For Your Scale

**Don't:**
- Compute drift in browser
- Compute sync debt in UI
- Duplicate authority logic in frontend

**Do:**
- Render projections directly
- Use typed contracts
- Trust server-side computation

#### Implementation Pattern

```typescript
// Backend adapter exposes typed projection APIs
interface FleetProjectionAPI {
  getGlobalProjection(): Promise<GlobalProjection>
  getDomainProjection(domainId: string): Promise<DomainProjection>
  getNodeProjection(nodeId: string): Promise<NodeProjection>
  
  // Typed subscriptions
  subscribeConnectivity(callback: (update: ConnectivityUpdate) => void): Subscription
  subscribeLifecycle(callback: (update: LifecycleUpdate) => void): Subscription
  subscribeDrift(callback: (update: DriftUpdate) => void): Subscription
}

// Frontend consumes typed projections
const globalProjection = await api.getGlobalProjection()
// TypeScript knows the exact shape
renderConnectivitySpectrum(globalProjection.connectivity)
```

---

### 2️⃣ Event-Stream Native UI (SSE-first / WebTransport / Partial Reactivity)

#### Why This Matters

You're building something where:
- Nodes reconnect dynamically
- Replay happens in real-time
- Drift updates continuously
- Authority shifts occur
- Sync debt changes

**Polling REST every 10 seconds is dead.**

#### Modern Pattern

- **Server-Sent Events (SSE)** as baseline
- **WebSockets** only when needed (bidirectional)
- **WebTransport** (emerging, but not mandatory)
- **Fine-grained reactive stores** (Pinia / Vue reactivity)

#### Why SSE is Hot Again

- Simpler than WebSockets
- Works well behind reverse proxies
- Perfect for projection updates
- Automatic reconnection
- Lower overhead than polling

#### UI Pattern

**Subscribe to:**
- `ConnectivityProjection` stream
- `DriftProjection` stream
- `LifecycleProjection` stream
- `AuthorityProjection` stream

**Update only affected nodes** (fine-grained reactivity)

#### You Avoid

- Full table refresh
- Heavy re-renders
- Jank under scale
- Polling overhead

#### Implementation Pattern

```typescript
// SSE-based projection streaming
class ProjectionStream {
  private eventSource: EventSource
  
  subscribeConnectivity(callback: (update: ConnectivityUpdate) => void) {
    this.eventSource = new EventSource('/api/v1/projections/connectivity/stream')
    
    this.eventSource.addEventListener('update', (event) => {
      const update = JSON.parse(event.data) as ConnectivityUpdate
      callback(update) // Fine-grained reactivity updates only affected nodes
    })
  }
  
  subscribeDrift(callback: (update: DriftUpdate) => void) {
    // Separate stream for drift updates
  }
}

// Pinia store with fine-grained reactivity
export const useFleetStore = defineStore('fleet', () => {
  const nodes = ref<Map<string, NodeProjection>>(new Map())
  
  // Only affected nodes update
  function updateNode(nodeId: string, update: Partial<NodeProjection>) {
    const node = nodes.value.get(nodeId)
    if (node) {
      Object.assign(node, update) // Vue reactivity handles the rest
    }
  }
  
  return { nodes, updateNode }
})
```

---

### 3️⃣ Virtualized + Hierarchical Rendering (at Scale)

#### The Problem

Rendering 10k+ nodes is not trivial.

**Traditional approaches fail:**
- Full DOM rendering = performance death
- Pure SVG for 50k nodes = jank
- Flat lists = unusable

#### What's Hot Now

- **Windowed rendering** (TanStack Virtual)
- **Dynamic tree virtualization**
- **Canvas/WebGL rendering** for graph layers
- **Hybrid DOM + Canvas**

#### Libraries Worth Watching

- **Cytoscape.js**: Graph visualization
- **PixiJS**: WebGL rendering
- **D3.js**: Used surgically for specific visualizations
- **TanStack Virtual**: Windowed list rendering

#### Modern Fleet UIs Are

**NOT:**
- Full DOM graph renderers
- Pure SVG for 50k nodes
- Flat node lists

**Instead:**
- Canvas/WebGL for topology
- DOM for interaction panels
- Virtualized lists for node tables

#### For IBCF

- **Domain view** = Canvas graph overlay (Cytoscape.js)
- **Node detail** = Standard Quasar components
- **Global overview** = Aggregated projection cards
- **Node list** = Virtualized table (TanStack Virtual)

#### Implementation Pattern

```typescript
// Canvas-based topology layer
import cytoscape from 'cytoscape'

class FleetTopologyView {
  private cy: cytoscape.Core
  
  renderDomain(domain: DomainProjection) {
    this.cy = cytoscape({
      container: document.getElementById('topology-canvas'),
      elements: this.buildGraphElements(domain),
      style: this.getTopologyStyle(),
      layout: { name: 'cose' }
    })
    
    // Only render visible nodes (viewport culling)
    this.cy.on('viewport', () => {
      this.updateVisibleNodes()
    })
  }
}

// Virtualized node list
import { useVirtualizer } from '@tanstack/vue-virtual'

const virtualizer = useVirtualizer({
  count: nodes.value.length,
  getScrollElement: () => scrollElement.value,
  estimateSize: () => 50, // Node row height
  overscan: 10 // Render 10 extra rows for smooth scrolling
})
```

---

### 4️⃣ State Machines in Frontend (XState / Formalized UI Logic)

#### Why This Matters

When your node lifecycle is:
```
LIVE → DELAYED → AUTONOMOUS → RETURNING → REPLAYING → QUARANTINED
```

Your UI logic should **not** be:
```typescript
if (status === 'AUTONOMOUS') {
  // ...
} else if (status === 'REPLAYING') {
  // ...
} else if (status === 'QUARANTINED') {
  // ...
}
```

#### Modern Approach

- Model lifecycle in a **state machine**
- UI reacts to machine transitions
- Formalized UI logic

#### Enabler: XState

XState provides:
- Visual state machine definition
- Type-safe transitions
- Guard conditions
- Actions on transitions
- React/Vue integration

#### Why This is Powerful for You

- **Replay state** can disable certain buttons
- **Quarantined state** changes UI theme
- **Updating state** locks command issuance
- **Authority mismatch** disables controls

**Instead of spaghetti conditionals:**
You bind UI to a machine.

#### Implementation Pattern

```typescript
import { createMachine, interpret } from 'xstate'
import { useMachine } from '@xstate/vue'

// Define node lifecycle machine
const nodeLifecycleMachine = createMachine({
  id: 'nodeLifecycle',
  initial: 'live',
  states: {
    live: {
      on: {
        LATENCY_INCREASES: 'delayed',
        CONNECTION_LOST: 'autonomous'
      },
      entry: 'enableCommands',
      exit: 'disableCommands'
    },
    delayed: {
      on: {
        SYNC_INTERVAL_EXTENDS: 'edgeCached',
        CONNECTION_LOST: 'autonomous',
        LATENCY_DECREASES: 'live'
      }
    },
    autonomous: {
      on: {
        CONNECTION_DETECTED: 'returning'
      },
      entry: 'showAutonomousWarning',
      exit: 'hideAutonomousWarning'
    },
    returning: {
      on: {
        BACKLOG_UPLOADED: 'replaying'
      }
    },
    replaying: {
      on: {
        REPLAY_COMPLETE: 'live',
        DRIFT_VIOLATION: 'quarantined'
      },
      entry: 'showReplayProgress',
      exit: 'hideReplayProgress'
    },
    quarantined: {
      entry: 'disableAllCommands',
      exit: 'enableCommands'
    }
  }
})

// Use in Vue component
export default defineComponent({
  setup() {
    const [state, send] = useMachine(nodeLifecycleMachine)
    
    // UI automatically reacts to state changes
    const canIssueCommands = computed(() => 
      state.value.matches('live') || state.value.matches('delayed')
    )
    
    const uiTheme = computed(() => {
      if (state.value.matches('quarantined')) return 'negative'
      if (state.value.matches('autonomous')) return 'warning'
      return 'primary'
    })
    
    return { state, send, canIssueCommands, uiTheme }
  }
})
```

---

### 5️⃣ Edge-Aware Offline-Capable Frontend (PWA + IndexedDB Replay)

#### Why This Matters

If nodes can go autonomous…  
Why shouldn't your UI?

#### Modern Fleet Systems

- Cache projections locally
- Allow filtered offline browsing
- Queue user commands until control plane reachable
- Survive control plane restarts

#### Hot Tech

- **Service Workers**: Mature, now widely adopted properly
- **IndexedDB wrappers**: Dexie.js
- **Background Sync**: Queue commands offline
- **Workbox**: Service worker management

#### Pattern

1. UI loads last known projections
2. Marks them "stale"
3. Syncs when control plane reconnects
4. Queues commands when offline
5. Replays commands when online

#### For Edge-Aware IBCF

This is extremely coherent with your system philosophy:
- Nodes operate autonomously
- UI should too
- Store-and-forward for commands
- Replay when connectivity returns

#### Implementation Pattern

```typescript
// Service Worker for offline support
// sw.js
self.addEventListener('fetch', (event) => {
  if (event.request.url.includes('/api/v1/projections')) {
    event.respondWith(
      caches.match(event.request).then((response) => {
        return response || fetch(event.request).then((fetchResponse) => {
          const cache = caches.open('projections-v1')
          cache.put(event.request, fetchResponse.clone())
          return fetchResponse
        })
      })
    )
  }
})

// IndexedDB for command queue
import Dexie from 'dexie'

class CommandQueue extends Dexie {
  commands!: Table<QueuedCommand>
  
  constructor() {
    super('CommandQueue')
    this.version(1).stores({
      commands: '++id, nodeId, command, timestamp, status'
    })
  }
  
  async queueCommand(nodeId: string, command: Command) {
    await this.commands.add({
      nodeId,
      command,
      timestamp: Date.now(),
      status: 'queued'
    })
  }
  
  async replayCommands() {
    const queued = await this.commands
      .where('status')
      .equals('queued')
      .toArray()
    
    for (const cmd of queued) {
      try {
        await api.sendCommand(cmd.nodeId, cmd.command)
        await this.commands.update(cmd.id!, { status: 'sent' })
      } catch (error) {
        // Keep queued, retry later
      }
    }
  }
}

// Vue composable for offline-aware commands
export function useOfflineCommands() {
  const queue = new CommandQueue()
  const isOnline = ref(navigator.onLine)
  
  watch(isOnline, (online) => {
    if (online) {
      queue.replayCommands()
    }
  })
  
  async function sendCommand(nodeId: string, command: Command) {
    if (isOnline.value) {
      try {
        await api.sendCommand(nodeId, command)
      } catch (error) {
        // Queue if network error
        await queue.queueCommand(nodeId, command)
      }
    } else {
      // Queue when offline
      await queue.queueCommand(nodeId, command)
    }
  }
  
  return { sendCommand, isOnline }
}
```

---

## What's NOT Hot Anymore

❌ **Massive Redux-style monolith stores**  
❌ **Poll-based dashboards**  
❌ **"Single pane of glass" mega dashboards**  
❌ **Raw metric walls**  
❌ **SVG-only topology for large fleets**

---

## Implementation Roadmap

### Phase 1: Typed Projection APIs
- [ ] Define TypeScript interfaces for projections
- [ ] Implement tRPC or typed REST endpoints
- [ ] Generate frontend types from backend schema

### Phase 2: SSE Streaming
- [ ] Implement SSE endpoints for projections
- [ ] Create Pinia stores with fine-grained reactivity
- [ ] Build subscription management

### Phase 3: Virtualized Rendering
- [ ] Integrate Cytoscape.js for topology
- [ ] Add TanStack Virtual for node lists
- [ ] Implement canvas/WebGL hybrid rendering

### Phase 4: State Machines
- [ ] Define XState machines for lifecycle
- [ ] Integrate with Vue components
- [ ] Build UI bindings to machine states

### Phase 5: Offline Support
- [ ] Implement Service Worker
- [ ] Add IndexedDB command queue
- [ ] Build offline UI indicators

---

## Mapping to IBCF Monitoring

| Need | Modern Frontend Technique |
|------|--------------------------|
| Massive fleet | Virtualized + Canvas graph |
| Drift visibility | Projection-driven UI |
| Multi-master | Authority projection + state machines |
| Autonomous nodes | Lifecycle state machine in UI |
| Replay workflows | SSE stream + machine transitions |
| Supply chain panel | Typed projection view model |
| Offline operation | PWA + IndexedDB replay |

---

## The Real Architectural Upgrade

**The big shift in the last 2–3 years:**

Frontend is no longer:
> A REST client that builds views.

It is:
> A projection renderer that subscribes to authoritative state.

**That distinction matters massively for scale.**

---

## Recommended Stack Additions

Given your existing stack (Vite + Quasar + Pinia):

1. **tRPC** or **typed REST** for projection APIs
2. **SSE** for projection streaming
3. **Cytoscape.js** for topology visualization
4. **XState** for lifecycle state machines
5. **Dexie.js** + **Workbox** for offline support

---

## Next Steps

1. Design typed projection API contracts
2. Implement SSE streaming infrastructure
3. Build canvas-based topology component
4. Define XState machines for lifecycle
5. Add PWA offline support

---

## References

- [tRPC](https://trpc.io/) - End-to-end typesafe APIs
- [XState](https://xstate.js.org/) - State machines
- [Cytoscape.js](https://js.cytoscape.org/) - Graph visualization
- [TanStack Virtual](https://tanstack.com/virtual/latest) - Virtual scrolling
- [Dexie.js](https://dexie.org/) - IndexedDB wrapper
- [Workbox](https://developers.google.com/web/tools/workbox) - Service worker toolkit

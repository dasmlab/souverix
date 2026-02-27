# Souverix Platform UI Architecture

**Multi-Component Fleet Monitoring & Lifecycle Management Frontend**

## Overview

This document defines the Souverix Platform UI architecture: a multi-component frontend platform where IBCF, Lifecycle Manager, Supply Chain, and other IMS components are modules within a unified shell.

**Key Principle:** Souverix is a **Platform UI**, not an App UI.

- One shell (navigation, auth, streaming substrate)
- Many component modules (IBCF, LCM, Supply Chain, Policy, etc.)
- Shared projections (fleet, domains, identity, authority)
- Per-component projections + commands

---

## Architecture Principles

### Platform-Level Frontend Techniques

1. **Micro-frontend modules** (pragmatic approach)
2. **Cross-component common projections** (fleet substrate)
3. **Topic routing for streams** (SSE multiplexer)
4. **Contract registry** (schemas as product interfaces)
5. **State machine composition** (global + module)

---

## 1) Frontend Projection Contract Layer (Typed)

### 1.1 Principles

- **Projection-first:** UI consumes read-models designed for rendering (not raw CRDs)
- **E2E types:** Contract schemas are shared/validated (runtime + compile-time)
- **Stable identities:** Everything keyed by immutable IDs (nodeId, domainId, policyBundleId)
- **Partial updates:** Streaming updates patch only what changed
- **Platform + Module separation:** Core substrate vs. component-specific projections

### 1.2 Contract Structure

```
contracts/
  core/              # Souverix substrate (shared)
    types.ts
    api.ts
    schemas.ts
  modules/
    ibcf/
      types.ts
      api.ts
    lcm/
      types.ts
      api.ts
    supply-chain/
      types.ts
      api.ts
```

### 1.3 Core Substrate Types (Shared)

```typescript
// src/contracts/core/types.ts
export type ID = string;

export type ConnectivityMode =
  | "LIVE"
  | "DELAYED"
  | "EDGE_CACHED"
  | "AUTONOMOUS"
  | "DISCONNECTED"
  | "RETURNING"
  | "REPLAYING";

export type Health = "OK" | "DEGRADED" | "FAILED" | "UNKNOWN";

export type LifecyclePhase =
  | "PROVISIONING"
  | "ENROLLED"
  | "ACTIVE"
  | "UPDATING"
  | "DEGRADED"
  | "QUARANTINED"
  | "DECOMMISSIONING"
  | "RETIRED";

// Core substrate view models (shared across all modules)
export interface SxNodeCard {
  nodeId: ID;
  name: string;
  kind: "EDGE" | "DC" | "GATEWAY" | "RELAY" | "VIRTUAL";
  domainId: ID;

  connectivity: {
    mode: ConnectivityMode;
    lastSeenTs: string; // RFC3339
    syncDebtSeconds: number;
    backlogMessages: number;
    rttMs?: number;
    lossPct?: number;
  };

  health: Health;

  desired: {
    policyBundleId: ID;
    targetVersion: string;
  };

  actual: {
    reportedVersion: string;
  };

  drift: {
    score: number;         // 0..1
    objectsOutOfSpec: number;
    lastEvaluatedTs: string;
  };

  authority: {
    masterOfRecord: "GLOBAL" | "REGION" | "LOCAL";
    controllerIds: string[];
  };

  flags: {
    supplyChainAlert: boolean;
    authorityMismatch: boolean;
    quarantined: boolean;
    updateInProgress: boolean;
  };
}

export interface SxGlobalProjection {
  generatedAtTs: string;
  totals: {
    nodes: number;
    domains: number;
  };

  connectivityCounts: Record<ConnectivityMode, number>;
  healthCounts: Record<Health, number>;

  topDriftDomains: Array<{ domainId: ID; domainName: string; avgDriftScore: number }>;
  attentionQueue: Array<{
    itemId: ID;
    severity: "P0" | "P1" | "P2";
    domainId: ID;
    nodeId?: ID;
    reason: string;
    ts: string;
  }>;
}

export interface SxDomainProjection {
  generatedAtTs: string;
  domainId: ID;
  domainName: string;

  // For topology overlays
  relays: Array<{ relayId: ID; name: string; kind: string }>;
  edges: Array<{ fromId: ID; toId: ID; linkHealth: Health; rttMs?: number }>;

  // For lists/cards (virtualized)
  nodeCards: SxNodeCard[];

  rolloutStatus: Array<{
    policyBundleId: ID;
    targetVersion: string;
    percentComplete: number;
    blockedBy?: string;
  }>;

  attentionQueue: SxGlobalProjection["attentionQueue"];
}

export interface SxNodeProjection {
  generatedAtTs: string;
  node: SxNodeCard;

  // Render-ready "diff"
  desiredVsActual: Array<{
    path: string;         // e.g. "spec.runtime.foo"
    desired: unknown;
    actual: unknown;
    status: "MATCH" | "DIFF" | "MISSING";
  }>;

  lifecycleTimeline: Array<{
    ts: string;
    fromState: string;
    toState: string;
    reason?: string;
  }>;

  supplyChain: {
    hardwareSerial?: string;
    firmwareVersion?: string;
    artifactId?: ID;
    sbomRef?: string;
    signatureOk?: boolean;
    attestations?: string[];
  };

  recentEvents: Array<{
    ts: string;
    type: string;
    summary: string;
  }>;
}

// Core substrate projections
export interface SxConnectivityProjection {
  nodeId: ID;
  mode: ConnectivityMode;
  lastSeenTs: string;
  syncDebtSeconds: number;
  backlogMessages: number;
}

export interface SxAuthorityProjection {
  domainId: ID;
  masterOfRecord: "GLOBAL" | "REGION" | "LOCAL";
  controllerIds: string[];
  lastAuthoritativeUpdateTs: string;
  conflicts?: Array<{ controllerId: ID; claim: string }>;
}
```

### 1.4 Module-Specific Projections (Example: IBCF)

```typescript
// src/contracts/modules/ibcf/types.ts
import type { ID, SxNodeCard, SxDomainProjection } from "src/contracts/core/types";

export interface IbcfDomainProjection extends SxDomainProjection {
  // IBCF-specific extensions
  relayGraph: {
    nodes: Array<{ relayId: ID; name: string; health: Health }>;
    edges: Array<{ from: ID; to: ID; pathType: "PRIMARY" | "BACKUP" }>;
  };

  interconnectPaths: Array<{
    pathId: ID;
    fromDomain: ID;
    toDomain: ID;
    status: "ACTIVE" | "DEGRADED" | "DOWN";
    replayInProgress?: boolean;
  }>;

  topologyHiding: {
    enabled: boolean;
    lastVerifiedTs: string;
    violations?: Array<{ nodeId: ID; leakType: string }>;
  };
}

export interface IbcfNodeProjection {
  nodeId: ID;
  base: SxNodeCard;

  // IBCF-specific
  peeringStatus: {
    activePeers: number;
    blockedPeers: number;
    lastPolicyUpdateTs: string;
  };

  stirShaken: {
    signingEnabled: boolean;
    verificationEnabled: boolean;
    certExpiryTs?: string;
  };
}
```

### 1.5 API Contracts

```typescript
// src/contracts/core/api.ts
import type { SxGlobalProjection, SxDomainProjection, SxNodeProjection, ID } from "./types";

export interface SouverixCoreQueryAPI {
  getGlobalProjection(): Promise<SxGlobalProjection>;
  getDomainProjection(domainId: ID): Promise<SxDomainProjection>;
  getNodeProjection(nodeId: ID): Promise<SxNodeProjection>;
}

export interface SouverixCoreCommandAPI {
  issueCommand(cmd: {
    commandId: ID;        // idempotency key
    target: { nodeId?: ID; domainId?: ID; tags?: string[] };
    type: "ROLL_OUT" | "QUARANTINE" | "CLEAR_QUARANTINE" | "RECONCILE" | "REBOOT" | "SET_POLICY";
    payload?: unknown;
  }): Promise<{ accepted: boolean; trackingId: ID }>;
}

// Module-specific API pattern
export interface SouverixModuleQueryAPI<ModuleName extends string> {
  getModuleProjection(moduleId: ModuleName, scope: { domainId?: ID; nodeId?: ID }): Promise<unknown>;
}

export interface SouverixModuleCommandAPI<ModuleName extends string> {
  issueModuleCommand(
    moduleId: ModuleName,
    cmd: { commandId: ID; target: { nodeId?: ID; domainId?: ID }; type: string; payload?: unknown }
  ): Promise<{ accepted: boolean; trackingId: ID }>;
}
```

### 1.6 Module Registry

```typescript
// src/contracts/modules/registry.ts
import type { RouteRecordRaw } from "vue-router";

export interface ModuleDescriptor {
  id: "ibcf" | "lcm" | "supply-chain" | "policy" | string;
  name: string;
  version: string;
  
  routes: RouteRecordRaw[];
  requiredTopics: string[];  // SSE topics this module needs
  
  // Module initialization
  init(store: RootStore): void;
  
  // Optional: module-specific stores
  stores?: Record<string, any>;
}

export interface ModuleRegistry {
  register(module: ModuleDescriptor): void;
  get(id: string): ModuleDescriptor | undefined;
  getAll(): ModuleDescriptor[];
}
```

---

## 2) Vite/Quasar Architecture Layout

### 2.1 Folder Structure

```
src/
  app/
    boot/
      api.ts                 # instantiate API adapters
      sse.ts                 # SSE client + topic subscriptions
      modules.ts             # module registry initialization
      auth.ts                # optional
    router/
      routes.ts              # core routes + module routes
    layouts/
      MainLayout.vue         # shell layout
  contracts/
    core/
      types.ts
      api.ts
      schemas.ts             # zod schemas
    modules/
      ibcf/
        types.ts
        api.ts
      lcm/
        types.ts
        api.ts
      supply-chain/
        types.ts
        api.ts
      registry.ts            # module registry
  adapters/
    core/
      queryHttp.ts           # implements SouverixCoreQueryAPI
      commandHttp.ts          # implements SouverixCoreCommandAPI
      sseClient.ts           # typed SSE subscribe + reconnection
    modules/
      ibcf/
        queryHttp.ts
        commandHttp.ts
  stores/
    core/
      projections.ts         # Pinia store for core projections
      attention.ts
      authority.ts
      connectivity.ts
    modules/
      ibcf/
        store.ts
      lcm/
        store.ts
  pages/
    core/
      GlobalOverviewPage.vue
      DomainPage.vue
      NodeDetailPage.vue
    modules/
      ibcf/
        IbcfDomainPage.vue
        IbcfNodePage.vue
      lcm/
        LcmNodePage.vue
  components/
    core/
      cards/
        NodeCard.vue
        DomainSummaryCard.vue
      topology/
        TopologyCanvas.vue     # Canvas/WebGL graph layer
      tables/
        VirtualizedNodeTable.vue
      panels/
        DesiredVsActualDiff.vue
        SupplyChainPanel.vue
        LifecycleTimeline.vue
    modules/
      ibcf/
        RelayGraph.vue
        InterconnectPaths.vue
      lcm/
        LifecycleStateMachine.vue
  utils/
    normalize.ts             # entity normalization
    time.ts
    ids.ts
```

### 2.2 Runtime Model

**Adapters:**
- Fetch initial projection snapshots (core + modules)
- Handle command issuance

**SSE:**
- Delivers incremental updates (patches)
- Topic-routed to appropriate stores

**Pinia Stores:**
- Core stores: hold normalized entities + indexes
- Module stores: hold module-specific state

**Pages:**
- Render from stores (not from API calls directly)
- Compose core + module components

**Components:**
- Core components: reusable across modules
- Module components: module-specific UI

### 2.3 Store Design (Normalized for Scale)

```typescript
// src/stores/core/projections.ts
import { defineStore } from "pinia";
import type { ID, SxNodeCard, SxDomainProjection, SxGlobalProjection } from "src/contracts/core/types";

export const useCoreProjectionStore = defineStore("core/projections", {
  state: () => ({
    global: null as SxGlobalProjection | null,

    nodesById: {} as Record<ID, SxNodeCard>,
    domainsById: {} as Record<ID, { domainId: ID; domainName: string }>,

    nodeIdsByDomain: {} as Record<ID, ID[]>,
    lastGeneratedAtByDomain: {} as Record<ID, string>,
  }),

  actions: {
    applyGlobalSnapshot(g: SxGlobalProjection) {
      this.global = g;
    },

    applyDomainSnapshot(d: SxDomainProjection) {
      this.domainsById[d.domainId] = { domainId: d.domainId, domainName: d.domainName };
      this.nodeIdsByDomain[d.domainId] = d.nodeCards.map(n => n.nodeId);
      this.lastGeneratedAtByDomain[d.domainId] = d.generatedAtTs;

      for (const n of d.nodeCards) this.nodesById[n.nodeId] = n;
    },

    patchNodeCard(node: Partial<SxNodeCard> & { nodeId: ID }) {
      this.nodesById[node.nodeId] = { ...this.nodesById[node.nodeId], ...node } as SxNodeCard;
    }
  },

  getters: {
    nodesByDomain: (state) => (domainId: ID) => {
      return (state.nodeIdsByDomain[domainId] || []).map(id => state.nodesById[id]);
    }
  }
});
```

### 2.4 Module Store Pattern

```typescript
// src/stores/modules/ibcf/store.ts
import { defineStore } from "pinia";
import { useCoreProjectionStore } from "src/stores/core/projections";
import type { IbcfDomainProjection, IbcfNodeProjection } from "src/contracts/modules/ibcf/types";

export const useIbcfStore = defineStore("modules/ibcf", {
  state: () => ({
    domainProjections: {} as Record<string, IbcfDomainProjection>,
    nodeProjections: {} as Record<string, IbcfNodeProjection>,
  }),

  actions: {
    applyDomainProjection(domain: IbcfDomainProjection) {
      this.domainProjections[domain.domainId] = domain;
      
      // Also update core store with base domain data
      const coreStore = useCoreProjectionStore();
      coreStore.applyDomainSnapshot(domain);
    },

    patchDomainProjection(domainId: string, patch: Partial<IbcfDomainProjection>) {
      this.domainProjections[domainId] = {
        ...this.domainProjections[domainId],
        ...patch
      } as IbcfDomainProjection;
    }
  }
});
```

---

## 3) SSE Event Schema (Fleet Projections)

### 3.1 Goals

- Topic-based subscription: global, domain, node, module-specific
- Support replay/resume after disconnect (Last-Event-ID)
- Incremental patches to avoid heavy payloads
- Idempotent processing (eventId + revision)
- Clear "snapshot vs patch" semantics
- Multiplexed streams (single connection, multiple topics)

### 3.2 SSE Stream Endpoints

**Recommended:**
```
/api/stream?topics=sx:global,sx:domain:abc,ibcf:domain:abc,lcm:node:xyz
```

Where:
- `sx:` = Souverix core substrate topics
- `ibcf:`, `lcm:`, etc. = module topics

**Alternative (separate endpoints):**
```
/api/stream/core/global
/api/stream/core/domains/{domainId}
/api/stream/core/nodes/{nodeId}
/api/stream/modules/{moduleId}/domains/{domainId}
```

### 3.3 Event Envelope

```typescript
interface SSEEnvelope {
  eventId: string;        // ULID
  ts: string;            // RFC3339
  topic: string;         // "sx:global" | "sx:domain:<id>" | "ibcf:domain:<id>" | "lcm:node:<id>"
  type: "SNAPSHOT" | "PATCH" | "TOMBSTONE" | "HEARTBEAT";
  projection: "GLOBAL" | "DOMAIN" | "NODE" | "CONNECTIVITY" | "LIFECYCLE" | "DRIFT" | "AUTHORITY" | "SUPPLYCHAIN" | string;
  revision: number;
  payload: unknown;
}
```

**SSE Fields:**
- `id:` = eventId
- `event:` = type or projection type
- `data:` = JSON envelope

### 3.4 Patch Format

Use JSON Patch (RFC 6902) or merge-patch.

**Example PATCH payload (node card update):**
```json
{
  "op": "patch",
  "entity": "SxNodeCard",
  "entityId": "node-123",
  "patch": [
    { "op": "replace", "path": "/connectivity/mode", "value": "REPLAYING" },
    { "op": "replace", "path": "/connectivity/backlogMessages", "value": 922 },
    { "op": "replace", "path": "/flags/updateInProgress", "value": false }
  ]
}
```

### 3.5 Snapshot Semantics

- On initial connect: server emits a SNAPSHOT for that topic
- Then PATCH events
- On reconnect: client sends Last-Event-ID header; server replays from event store if available, else re-sends SNAPSHOT

### 3.6 Heartbeats

Send HEARTBEAT events every N seconds:
- Prevents idle disconnects behind proxies
- Gives UI a sense of "stream alive"

### 3.7 Example SSE Frames

```
id: 01JABC...
event: PATCH
data: {"eventId":"01JABC...","ts":"...","topic":"sx:domain:abc","type":"PATCH","projection":"CONNECTIVITY","revision":554,"payload":{...}}
```

### 3.8 Client Processing Rules

**Maintain:**
- `lastEventId` per stream
- `lastRevision` per (topic, projection)

**Drop events if:**
- `revision <= lastRevision` (idempotency)

**Apply:**
- Snapshots replace
- Patches merge/patch

### 3.9 SSE Client Implementation

```typescript
// src/adapters/core/sseClient.ts
type Envelope = {
  eventId: string;
  ts: string;
  topic: string;
  type: "SNAPSHOT" | "PATCH" | "TOMBSTONE" | "HEARTBEAT";
  projection: string;
  revision: number;
  payload: any;
};

export class SSEClient {
  private eventSource: EventSource | null = null;
  private lastEventId: string | null = null;
  private topicHandlers: Map<string, (env: Envelope) => void> = new Map();

  subscribe(topics: string[], onMessage: (env: Envelope) => void): () => void {
    const topicString = topics.join(",");
    const url = `/api/stream?topics=${encodeURIComponent(topicString)}${this.lastEventId ? `&lastEventId=${this.lastEventId}` : ""}`;
    
    this.eventSource = new EventSource(url, { withCredentials: true });

    this.eventSource.onmessage = (evt) => {
      try {
        const envelope = JSON.parse(evt.data) as Envelope;
        this.lastEventId = envelope.eventId;
        
        // Route to topic handler
        const handler = this.topicHandlers.get(envelope.topic);
        if (handler) {
          handler(envelope);
        }
        
        // Also call general handler
        onMessage(envelope);
      } catch (e) {
        console.error("Failed to parse SSE message", e);
      }
    };

    this.eventSource.addEventListener("PATCH", (evt: MessageEvent) => {
      try {
        const envelope = JSON.parse(evt.data) as Envelope;
        this.lastEventId = envelope.eventId;
        onMessage(envelope);
      } catch (e) {
        console.error("Failed to parse PATCH event", e);
      }
    });

    this.eventSource.addEventListener("SNAPSHOT", (evt: MessageEvent) => {
      try {
        const envelope = JSON.parse(evt.data) as Envelope;
        this.lastEventId = envelope.eventId;
        onMessage(envelope);
      } catch (e) {
        console.error("Failed to parse SNAPSHOT event", e);
      }
    });

    this.eventSource.onerror = () => {
      // Browser handles reconnect, but you can also close + backoff
      console.warn("SSE connection error, will attempt reconnect");
    };

    return () => {
      this.eventSource?.close();
      this.eventSource = null;
    };
  }

  registerTopicHandler(topic: string, handler: (env: Envelope) => void) {
    this.topicHandlers.set(topic, handler);
  }
}
```

---

## 4) Platform-Level Architecture Patterns

### 4.1 Micro-Frontend Modules (Pragmatic)

**Not "micro-frontends because trend", but because you'll have many IMS components.**

**Best Practice:**
- Single Vite/Quasar app shell
- Modules loaded by route + feature flags
- Shared design system + shared stores/contracts
- Optionally: Module Federation later, only if you truly need independent deployments

**What You Standardize:**
- Module registers routes
- Module declares required projections + commands
- Module declares SSE topics it needs

### 4.2 Cross-Component Common Projections

**Souverix-wide read models:**
- Identity / inventory
- Connectivity spectrum
- Authority / tenancy domains
- Event timeline index
- Attention queue

**Then each module adds:**
- Its own projections (IBCF-relay graph, cert lifecycle, etc.)
- Its own commands

**This prevents every component from reinventing:**
- "node list"
- "domain view"
- "who's master"
- "what's stale"

### 4.3 Topic Routing for Streams (SSE Multiplexer)

**Single multiplexed stream:**
```
/api/stream?topics=sx:global,sx:domain:abc,ibcf:domain:abc,lcm:node:xyz
```

**Where:**
- `sx:` = Souverix core substrate topics
- `ibcf:`, `lcm:`, etc. = module topics

**Modern technique:** A single streaming substrate with topic routing + per-module handlers.

### 4.4 Contract Registry (Schemas as Product Interfaces)

**In a platform UI, the contract layer becomes a "public interface" between modules and backend adapters.**

**Structure:**
```
contracts/
  core/              # Souverix substrate
  modules/
    ibcf/            # IBCF module contracts
    lcm/             # Lifecycle Manager contracts
    supply-chain/    # Supply Chain contracts
```

**Use runtime validation (Zod or equivalent)** so modules can't break the shell silently.

**This is a big "last 2 years" pattern:** Treat schemas as an interface registry.

### 4.5 State Machine Composition (Global + Module)

**Platform has global UI states:**
- authenticated / unauthenticated
- connected / degraded / offline
- authority conflict detected
- maintenance mode
- data stale

**Each module has its own state machine, but it should compose with global state guards.**

**Example:**
If platform says "authority conflict", disable "rollout" in every module that issues commands.

**This is how you keep correctness across many components without spaghetti guards everywhere.**

---

## 5) Module Integration Pattern

### 5.1 Module Descriptor

```typescript
// src/contracts/modules/registry.ts
import type { RouteRecordRaw } from "vue-router";
import type { Store } from "pinia";

export interface ModuleDescriptor {
  id: "ibcf" | "lcm" | "supply-chain" | "policy" | string;
  name: string;
  version: string;
  
  routes: RouteRecordRaw[];
  requiredTopics: string[];  // SSE topics this module needs
  
  // Module initialization
  init(store: Store): void;
  
  // Optional: module-specific stores
  stores?: Record<string, any>;
  
  // Optional: sidebar items
  sidebarItems?: Array<{
    label: string;
    icon: string;
    route: string;
  }>;
}
```

### 5.2 Module Registration

```typescript
// src/app/boot/modules.ts
import { ModuleRegistry } from "src/contracts/modules/registry";
import { ibcfModule } from "src/modules/ibcf";
import { lcmModule } from "src/modules/lcm";

export function initializeModules(registry: ModuleRegistry, router: Router, store: Store) {
  // Register modules
  registry.register(ibcfModule);
  registry.register(lcmModule);
  
  // Add routes
  const modules = registry.getAll();
  for (const module of modules) {
    router.addRoute(module.routes);
    module.init(store);
  }
}
```

### 5.3 Example Module (IBCF)

```typescript
// src/modules/ibcf/index.ts
import type { ModuleDescriptor } from "src/contracts/modules/registry";
import { useIbcfStore } from "src/stores/modules/ibcf/store";

export const ibcfModule: ModuleDescriptor = {
  id: "ibcf",
  name: "IBCF",
  version: "1.0.0",
  
  routes: [
    {
      path: "/ibcf/domains/:domainId",
      component: () => import("./pages/IbcfDomainPage.vue"),
    },
    {
      path: "/ibcf/nodes/:nodeId",
      component: () => import("./pages/IbcfNodePage.vue"),
    },
  ],
  
  requiredTopics: [
    "sx:global",
    "sx:domain:*",
    "ibcf:domain:*",
  ],
  
  init(store) {
    // Initialize IBCF-specific stores
    const ibcfStore = useIbcfStore();
    // ... setup
  },
  
  sidebarItems: [
    {
      label: "IBCF Domains",
      icon: "network",
      route: "/ibcf/domains",
    },
  ],
};
```

---

## 6) Implementation Roadmap

### Phase 1: Core Substrate
- [ ] Define core projection contracts
- [ ] Implement core query/command APIs
- [ ] Build core Pinia stores
- [ ] Create global overview page

### Phase 2: SSE Infrastructure
- [ ] Implement SSE client with topic routing
- [ ] Build projection update handlers
- [ ] Add reconnection logic
- [ ] Test with core projections

### Phase 3: Module System
- [ ] Create module registry
- [ ] Implement module descriptor pattern
- [ ] Build module initialization
- [ ] Add route registration

### Phase 4: First Module (IBCF)
- [ ] Define IBCF projection contracts
- [ ] Implement IBCF stores
- [ ] Build IBCF pages/components
- [ ] Integrate with core substrate

### Phase 5: Advanced Features
- [ ] Virtualized rendering
- [ ] Canvas topology
- [ ] State machines (XState)
- [ ] Offline support (PWA)

---

## Next Steps

1. Implement core projection contracts
2. Build SSE client with topic routing
3. Create module registry
4. Build first module (IBCF)
5. Add virtualization and canvas rendering

---

## References

- [Frontend Architecture](./frontend-architecture.md) - Modern frontend techniques
- [IBCF Fleet UI & Lifecycle](./ibcf-fleet-ui-lifecycle.md) - IBCF-specific design
- [Fleet Monitoring](./fleet-monitoring.md) - Fleet monitoring architecture

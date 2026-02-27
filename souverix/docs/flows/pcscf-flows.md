# P-CSCF Flow Diagrams

This document illustrates the key operational flows for the Proxy Call Session Control Function (P-CSCF) from its perspective within the IMS architecture.

## P-CSCF In One Line

**P-CSCF secures and anchors UE SIP signaling into the IMS core.**

## 1. High-Level Access Edge View

P-CSCF serves as the first SIP contact point for User Equipment, providing security and policy enforcement at the IMS access edge.

### Network View

```
[ UE / Access Network ]
         |
    (Gm - SIP over IPSec/TLS)
         |
      | P-CSCF |
         |
    (Mw - SIP)
         |
  I-CSCF / S-CSCF
         |
    [ IMS Core ]
```

### What P-CSCF Does Here

- ✅ Terminates secure SIP from UE (Gm interface)
- ✅ Validates and forwards SIP messages
- ✅ Inserts Record-Route to maintain signaling path
- ✅ Triggers policy control (PCRF/PCF)
- ✅ Maintains dialog anchoring

### Sequence Diagram

```mermaid
sequenceDiagram
    participant UE as User Equipment<br/>(Access Network)
    participant PCSCF as P-CSCF
    participant ICSCF as I-CSCF
    participant SCSCF as S-CSCF
    participant IMS as IMS Core

    Note over UE,IMS: Access Edge - First Contact Point

    UE->>PCSCF: SIP REGISTER<br/>(Gm - over IPSec/TLS)
    Note right of PCSCF: Security validation<br/>Policy trigger
    
    PCSCF->>ICSCF: SIP REGISTER<br/>(Mw interface)
    ICSCF->>SCSCF: Forward to S-CSCF
    SCSCF->>IMS: Process registration
    
    IMS-->>SCSCF: Registration accepted
    SCSCF-->>ICSCF: 200 OK
    ICSCF-->>PCSCF: 200 OK
    PCSCF-->>UE: 200 OK<br/>(Gm - secure)
    
    Note over UE,IMS: UE registered, P-CSCF anchored
```

---

## 2. Primary Flow – SIP REGISTER with Security Exchange

This is the standard registration flow where P-CSCF establishes security association with UE.

### Sequence Diagram

```mermaid
sequenceDiagram
    participant UE as UE
    participant PCSCF as P-CSCF
    participant ICSCF as I-CSCF
    participant SCSCF as S-CSCF
    participant HSS as HSS

    Note over UE,HSS: Initial Registration Flow

    UE->>PCSCF: REGISTER<br/>(Gm - initial, no auth)
    Note right of PCSCF: First contact point<br/>Validate request
    
    PCSCF->>ICSCF: REGISTER<br/>(Mw - forward)
    ICSCF->>SCSCF: REGISTER<br/>(select S-CSCF)
    SCSCF->>HSS: Query subscriber data
    
    HSS-->>SCSCF: Subscriber profile
    SCSCF-->>ICSCF: 401 Unauthorized<br/>(challenge)
    ICSCF-->>PCSCF: 401 Unauthorized
    PCSCF-->>UE: 401 Unauthorized<br/>(Gm - secure)
    
    Note over UE,PCSCF: IPSec Security-Client /<br/>Security-Server negotiation
    
    UE->>PCSCF: REGISTER<br/>(Gm - with auth response)
    Note right of PCSCF: IPSec SA established<br/>Integrity protection active
    
    PCSCF->>ICSCF: REGISTER<br/>(Mw - authenticated)
    ICSCF->>SCSCF: REGISTER
    SCSCF->>HSS: Verify authentication
    
    HSS-->>SCSCF: Auth success
    SCSCF-->>ICSCF: 200 OK<br/>(Service-Route)
    ICSCF-->>PCSCF: 200 OK
    PCSCF-->>UE: 200 OK<br/>(Gm - secure)
    
    Note over UE,HSS: Registration complete,<br/>P-CSCF in path
```

### What P-CSCF Does Here

- ✅ Receives initial REGISTER from UE (Gm)
- ✅ Forwards to I-CSCF (Mw)
- ✅ Handles 401 Unauthorized response
- ✅ Establishes IPSec Security Association
- ✅ Enforces integrity protection
- ✅ Maintains Record-Route for dialog anchoring

---

## 3. Primary Flow – SIP INVITE (Call Setup)

P-CSCF forwards call setup requests and maintains dialog state.

### Sequence Diagram

```mermaid
sequenceDiagram
    participant UE1 as UE-A<br/>(Caller)
    participant PCSCF1 as P-CSCF-A
    participant SCSCF1 as S-CSCF-A
    participant SCSCF2 as S-CSCF-B
    participant PCSCF2 as P-CSCF-B
    participant UE2 as UE-B<br/>(Callee)

    Note over UE1,UE2: Call Setup Flow

    UE1->>PCSCF1: INVITE<br/>(Gm - secure SIP)
    Note right of PCSCF1: Validate headers<br/>Trigger policy control
    
    PCSCF1->>SCSCF1: INVITE<br/>(Mw - with Record-Route)
    Note right of PCSCF1: Insert Record-Route<br/>to remain in path
    
    SCSCF1->>SCSCF2: INVITE<br/>(inter-S-CSCF)
    SCSCF2->>PCSCF2: INVITE<br/>(Mw)
    PCSCF2->>UE2: INVITE<br/>(Gm - secure)
    
    UE2-->>PCSCF2: 180 Ringing<br/>(Gm)
    PCSCF2-->>SCSCF2: 180 Ringing
    SCSCF2-->>SCSCF1: 180 Ringing
    SCSCF1-->>PCSCF1: 180 Ringing
    PCSCF1-->>UE1: 180 Ringing<br/>(Gm)
    
    UE2-->>PCSCF2: 200 OK<br/>(Gm)
    PCSCF2-->>SCSCF2: 200 OK
    SCSCF2-->>SCSCF1: 200 OK
    SCSCF1-->>PCSCF1: 200 OK
    PCSCF1-->>UE1: 200 OK<br/>(Gm)
    
    UE1->>PCSCF1: ACK<br/>(Gm)
    PCSCF1->>SCSCF1: ACK<br/>(Mw - via Record-Route)
    SCSCF1->>SCSCF2: ACK
    SCSCF2->>PCSCF2: ACK
    PCSCF2->>UE2: ACK<br/>(Gm)
    
    Note over UE1,UE2: Call established,<br/>P-CSCF maintains dialog state
```

### What P-CSCF Does Here

- ✅ Receives INVITE from UE (Gm)
- ✅ Validates SIP headers (From/To/Contact)
- ✅ Triggers policy control (QoS authorization)
- ✅ Inserts Record-Route header
- ✅ Forwards to S-CSCF (Mw)
- ✅ Proxies all subsequent messages (180, 200 OK, ACK)
- ✅ Maintains dialog state for mid-dialog signaling

---

## 4. Emergency Call Handling

P-CSCF routes emergency calls to E-CSCF (Emergency CSCF) instead of normal I-CSCF.

### Sequence Diagram

```mermaid
sequenceDiagram
    participant UE as UE<br/>(Emergency)
    participant PCSCF as P-CSCF
    participant ECSCF as E-CSCF<br/>(Emergency)
    participant LRF as LRF<br/>(Location)
    participant PSAP as PSAP<br/>(Public Safety)

    Note over UE,PSAP: Emergency Call Flow

    UE->>PCSCF: INVITE<br/>(Gm - emergency indication)
    Note right of PCSCF: Detect emergency call<br/>Priority routing
    
    PCSCF->>ECSCF: INVITE<br/>(emergency routing)
    Note right of ECSCF: Emergency handling
    
    ECSCF->>LRF: Query location<br/>(for routing)
    LRF-->>ECSCF: Location + PSAP info
    
    ECSCF->>PSAP: INVITE<br/>(emergency call)
    
    PSAP-->>ECSCF: 200 OK
    ECSCF-->>PCSCF: 200 OK
    PCSCF-->>UE: 200 OK<br/>(Gm)
    
    Note over UE,PSAP: Emergency call established
```

### What P-CSCF Does Here

- ✅ Detects emergency call indication
- ✅ Routes to E-CSCF (not I-CSCF)
- ✅ Bypasses normal registration checks
- ✅ Maintains security (IPSec still required)
- ✅ Supports emergency registration if needed

---

## 5. Policy Control Interaction (Rx Interface)

P-CSCF interacts with PCRF/PCF for QoS and policy enforcement.

### Sequence Diagram

```mermaid
sequenceDiagram
    participant UE as UE
    participant PCSCF as P-CSCF
    participant PCRF as PCRF/PCF<br/>(Policy Control)
    participant SCSCF as S-CSCF
    participant PCEF as PCEF<br/>(Access Gateway)

    Note over UE,PCEF: Policy Control Flow

    UE->>PCSCF: INVITE<br/>(with SDP)
    Note right of PCSCF: Extract media info<br/>from SDP
    
    PCSCF->>PCRF: AAR (Authorization Request)<br/>(Rx interface - Diameter)
    Note right of PCRF: Policy decision<br/>QoS authorization
    
    PCRF->>PCEF: Install QoS rules
    PCRF-->>PCSCF: AAA (Authorization Answer)<br/>(QoS authorized)
    
    PCSCF->>SCSCF: INVITE<br/>(forward)
    SCSCF-->>PCSCF: 200 OK<br/>(with SDP)
    
    PCSCF->>PCRF: RAR (Re-Auth Request)<br/>(update media)
    PCRF-->>PCSCF: RAA (Re-Auth Answer)
    
    PCSCF-->>UE: 200 OK<br/>(Gm)
    
    Note over UE,PCEF: Media session authorized,<br/>QoS enforced
```

### What P-CSCF Does Here

- ✅ Extracts media information from SDP
- ✅ Requests QoS authorization from PCRF/PCF (Rx)
- ✅ Receives policy decisions
- ✅ Updates policy on media changes
- ✅ Maintains policy session state

---

## 6. Multi-Site IMS Deployment (P-CSCF Focus)

P-CSCF is deployed close to access networks, with each site having its own P-CSCF pool.

### Topology Diagram

```mermaid
graph TB
    subgraph SiteA["IMS Site A"]
        UE1[UE-A]
        P1[P-CSCF-A]
        S1[S-CSCF-A]
        Core1[Core-A]
        
        UE1 -->|Gm| P1
        P1 -->|Mw| S1
        S1 --> Core1
    end
    
    subgraph SiteB["IMS Site B"]
        UE2[UE-B]
        P2[P-CSCF-B]
        S2[S-CSCF-B]
        Core2[Core-B]
        
        UE2 -->|Gm| P2
        P2 -->|Mw| S2
        S2 --> Core2
    end
    
    Core1 -.->|Inter-site Core<br/>Connectivity| Core2
    
    style P1 fill:#4a90e2,stroke:#333,stroke-width:2px
    style P2 fill:#4a90e2,stroke:#333,stroke-width:2px
    style SiteA fill:#e8f4f8
    style SiteB fill:#e8f4f8
```

### Key Observations

- ✅ **P-CSCF is always deployed close to access edge**
- ✅ **Each site typically has its own P-CSCF pool**
- ✅ **UE attaches to geographically nearest P-CSCF**
- ✅ **Inter-site signaling does NOT traverse P-CSCF ↔ P-CSCF**
- ✅ **P-CSCF is access-bound, not inter-domain routing logic**

### What P-CSCF Does NOT Do

- ❌ Does **not** perform service logic
- ❌ Does **not** select breakout
- ❌ Does **not** anchor media (that's SBC/IMS-ALG/IMS-ATCF domain)
- ❌ Does **not** perform HSS subscriber selection
- ❌ Does **not** manage inter-IMS routing

---

## 7. NAT Traversal Support

P-CSCF helps maintain SIP signaling through NATs using keep-alives and STUN.

### Sequence Diagram

```mermaid
sequenceDiagram
    participant UE as UE<br/>(Behind NAT)
    participant NAT as NAT Gateway
    participant PCSCF as P-CSCF
    participant SCSCF as S-CSCF

    Note over UE,SCSCF: NAT Traversal Flow

    UE->>NAT: REGISTER<br/>(creates NAT binding)
    NAT->>PCSCF: REGISTER<br/>(NAT public IP)
    PCSCF->>SCSCF: REGISTER<br/>(forward)
    SCSCF-->>PCSCF: 200 OK
    PCSCF-->>NAT: 200 OK
    NAT-->>UE: 200 OK
    
    Note over UE,PCSCF: Registration complete
    
    loop Keep-Alive (every 30s)
        UE->>NAT: CRLF keep-alive<br/>(maintain binding)
        NAT->>PCSCF: CRLF
        PCSCF-->>NAT: CRLF
        NAT-->>UE: CRLF
    end
    
    Note over UE,SCSCF: NAT binding maintained
    
    UE->>NAT: INVITE<br/>(uses existing binding)
    NAT->>PCSCF: INVITE
    PCSCF->>SCSCF: INVITE
    SCSCF-->>PCSCF: 200 OK
    PCSCF-->>NAT: 200 OK
    NAT-->>UE: 200 OK
    
    Note over UE,SCSCF: Call established through NAT
```

### What P-CSCF Does Here

- ✅ Receives SIP through NAT
- ✅ Responds to keep-alive messages (CRLF)
- ✅ Supports STUN binding requests
- ✅ Maintains NAT bindings via periodic keep-alives
- ✅ Enables SIP OPTIONS for NAT traversal

---

## Interface Summary

| Interface | Direction | Purpose |
|-----------|-----------|---------|
| **Gm** | UE ↔ P-CSCF | Secure SIP signaling (IPSec/TLS) |
| **Mw** | P-CSCF ↔ I-CSCF/S-CSCF | Forward SIP to IMS core |
| **Rx** | P-CSCF ↔ PCRF/PCF | Policy control and QoS authorization |

---

## Related Documentation

- [P-CSCF Features List](../../components/coeur/pcscf/FEATURES_LIST.md)
- [ETSI TS 23.228](https://www.etsi.org/deliver/etsi_ts/123200_123299/123228/07.02.00_60/ts_123228v070200p.pdf) - IMS Stage 2
- [ETSI TS 24.229](https://www.etsi.org/deliver/etsi_ts/124200_124299/124229/16.04.00_60/ts_124229v160400p.pdf) - SIP/SDP Protocol

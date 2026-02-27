# S-CSCF Flow Diagrams

This document illustrates the key operational flows for the Serving Call Session Control Function (S-CSCF) from its perspective within the IMS architecture.

## S-CSCF In One Line

**S-CSCF is the stateful SIP service brain of the IMS domain.**

## 1. Registration Flow (Core Logic)

S-CSCF handles user authentication, service profile download, and registration state management.

### Sequence Diagram

```mermaid
sequenceDiagram
    autonumber
    participant UE as UE
    participant PCSCF as P-CSCF
    participant ICSCF as I-CSCF
    participant SCSCF as S-CSCF
    participant HSS as HSS (Cx)

    UE->>PCSCF: REGISTER (Gm)
    PCSCF->>ICSCF: REGISTER (Mw)
    ICSCF->>SCSCF: REGISTER (Mw)
    
    Note over SCSCF: S-CSCF's job:<br/>AKA challenge/verify +<br/>server assignment +<br/>service profile

    SCSCF-->>UE: 401 Unauthorized (AKA challenge)\n(via I-CSCF, P-CSCF)
    
    UE->>PCSCF: REGISTER + Authorization (AKA response)
    PCSCF->>ICSCF: REGISTER + Authorization
    ICSCF->>SCSCF: REGISTER + Authorization

    SCSCF->>HSS: Cx SAR (Server-Assignment-Req)
    HSS-->>SCSCF: Cx SAA (Service Profile / iFC)
    
    Note over SCSCF: Store service profile<br/>Initial Filter Criteria (iFC)<br/>Registration state active
    
    SCSCF-->>UE: 200 OK (Registered)\n(via I-CSCF, P-CSCF)
```

### What S-CSCF Does Here

- ✅ Receives REGISTER from I-CSCF
- ✅ Issues 401 Unauthorized (AKA challenge)
- ✅ Verifies AKA authentication response
- ✅ Performs server assignment (SAR to HSS)
- ✅ Downloads service profile and iFC from HSS
- ✅ Maintains registration state
- ✅ Returns 200 OK upon successful registration

---

## 2. Basic Session Flow (Outgoing Call)

S-CSCF applies service logic, triggers Application Servers, and determines routing.

### Sequence Diagram

```mermaid
sequenceDiagram
    autonumber
    participant UE as UE
    participant PCSCF as P-CSCF
    participant ICSCF as I-CSCF
    participant SCSCF as S-CSCF
    participant AS as Application Server<br/>(if triggered)
    participant UE2 as Called Party (IMS/URI)

    UE->>PCSCF: INVITE (SDP) (Gm)
    PCSCF->>ICSCF: INVITE (Mw)
    ICSCF->>SCSCF: INVITE (Mw)

    Note over SCSCF: Apply service profile (iFC)\nTrigger AS if required (ISC)\nSelect routing target
    
    alt AS Trigger Required
        SCSCF->>AS: INVITE (ISC interface)
        AS-->>SCSCF: Continue / Modify
    end
    
    SCSCF->>UE2: INVITE (toward terminating side)

    UE2-->>SCSCF: 180 Ringing
    SCSCF-->>ICSCF: 180 Ringing
    ICSCF-->>PCSCF: 180 Ringing
    PCSCF-->>UE: 180 Ringing

    UE2-->>SCSCF: 200 OK (SDP)
    SCSCF-->>ICSCF: 200 OK
    ICSCF-->>PCSCF: 200 OK
    PCSCF-->>UE: 200 OK

    UE->>PCSCF: ACK
    PCSCF->>ICSCF: ACK
    ICSCF->>SCSCF: ACK
    
    Note over SCSCF: Insert Record-Route<br/>Anchor dialog
    
    SCSCF->>UE2: ACK
```

### What S-CSCF Does Here

- ✅ Receives INVITE from I-CSCF
- ✅ Applies Initial Filter Criteria (iFC) from service profile
- ✅ Triggers Application Servers if required (ISC interface)
- ✅ Determines routing target (IMS, PSTN, etc.)
- ✅ Inserts Record-Route to anchor dialog
- ✅ Proxies all subsequent messages
- ✅ Maintains dialog state

---

## 3. PSTN Breakout Decision

S-CSCF analyzes destination and triggers BGCF for PSTN breakout.

### Sequence Diagram

```mermaid
sequenceDiagram
    autonumber
    participant UE as UE
    participant PCSCF as P-CSCF
    participant ICSCF as I-CSCF
    participant SCSCF as S-CSCF
    participant BGCF as BGCF
    participant MGCF as MGCF
    participant PSTN as PSTN/CS

    UE->>PCSCF: INVITE (tel:+E164) (Gm)
    PCSCF->>ICSCF: INVITE (Mw)
    ICSCF->>SCSCF: INVITE (Mw)

    Note over SCSCF: Number analysis / policy\nDecide PSTN breakout required
    
    SCSCF->>BGCF: INVITE (Mi)
    
    Note over BGCF: Select breakout network\nSelect MGCF (local) or forward to remote BGCF
    
    BGCF->>MGCF: INVITE (Mj)

    Note over MGCF: SIP↔ISUP interworking\nControl MGW via H.248
    
    MGCF->>PSTN: ISUP IAM (Nc)

    PSTN-->>MGCF: ISUP ACM (alerting)
    MGCF-->>BGCF: 180 Ringing
    BGCF-->>SCSCF: 180 Ringing
    SCSCF-->>UE: 180 Ringing (via I-CSCF, P-CSCF)

    PSTN-->>MGCF: ISUP ANM (answer)
    MGCF-->>BGCF: 200 OK
    BGCF-->>SCSCF: 200 OK
    SCSCF-->>UE: 200 OK (via I-CSCF, P-CSCF)
    
    UE->>PCSCF: ACK
    PCSCF->>ICSCF: ACK
    ICSCF->>SCSCF: ACK
    SCSCF->>BGCF: ACK
    BGCF->>MGCF: ACK
```

### What S-CSCF Does Here

- ✅ Receives INVITE with tel: URI
- ✅ Performs number analysis (ENUM, routing logic)
- ✅ Determines PSTN breakout required
- ✅ Routes to BGCF (Mi interface)
- ✅ Maintains dialog anchoring
- ✅ Proxies responses and ACK

---

## 4. Application Server Triggering (ISC Interface)

S-CSCF triggers Application Servers based on Initial Filter Criteria.

### Sequence Diagram

```mermaid
sequenceDiagram
    autonumber
    participant UE as UE
    participant SCSCF as S-CSCF
    participant HSS as HSS
    participant AS1 as AS1<br/>(Voicemail)
    participant AS2 as AS2<br/>(Call Forwarding)
    participant Dest as Destination

    UE->>SCSCF: INVITE
    SCSCF->>HSS: Query service profile
    
    Note over SCSCF: Apply iFC rules<br/>from service profile
    
    alt iFC Match: Voicemail
        SCSCF->>AS1: INVITE (ISC)
        AS1-->>SCSCF: Continue
    end
    
    alt iFC Match: Call Forwarding
        SCSCF->>AS2: INVITE (ISC)
        AS2->>Dest: Forward call
        Dest-->>AS2: 200 OK
        AS2-->>SCSCF: 200 OK
    else No AS Match
        SCSCF->>Dest: INVITE (normal routing)
        Dest-->>SCSCF: 200 OK
    end
    
    SCSCF-->>UE: 200 OK
```

### What S-CSCF Does Here

- ✅ Evaluates Initial Filter Criteria (iFC) from service profile
- ✅ Triggers Application Servers based on iFC matches
- ✅ Handles AS responses (Continue, Terminate, etc.)
- ✅ Applies service logic in correct order
- ✅ Maintains service execution state

---

## 5. Mid-Dialog Request Handling

S-CSCF handles mid-dialog requests like re-INVITE and UPDATE.

### Sequence Diagram

```mermaid
sequenceDiagram
    autonumber
    participant UE1 as UE-A
    participant SCSCF as S-CSCF
    participant UE2 as UE-B

    Note over UE1,UE2: Active call in progress

    UE1->>SCSCF: re-INVITE (modify SDP)
    Note over SCSCF: Dialog anchored<br/>Apply service logic
    
    SCSCF->>UE2: re-INVITE (modify media)
    UE2-->>SCSCF: 200 OK (SDP answer)
    SCSCF-->>UE1: 200 OK
    
    UE1->>SCSCF: ACK
    SCSCF->>UE2: ACK
    
    Note over UE1,UE2: Media modified (e.g., hold/resume)
    
    UE1->>SCSCF: UPDATE (QoS change)
    SCSCF->>UE2: UPDATE
    UE2-->>SCSCF: 200 OK
    SCSCF-->>UE1: 200 OK
```

### What S-CSCF Does Here

- ✅ Receives mid-dialog requests (re-INVITE, UPDATE)
- ✅ Validates dialog state
- ✅ Applies service logic if required
- ✅ Forwards to other dialog participant
- ✅ Maintains dialog continuity

---

## Interface Summary

| Interface | Direction | Purpose |
|-----------|-----------|---------|
| **Mw** | I-CSCF ↔ S-CSCF | SIP signaling routing |
| **Cx** | S-CSCF ↔ HSS | Registration, authentication, service profile (Diameter) |
| **ISC** | S-CSCF ↔ AS | Application Server triggering |
| **Mi** | S-CSCF → BGCF | PSTN breakout routing |
| **Ro/Rf** | S-CSCF ↔ CCF/OCF | Charging events (Diameter) |

---

## Related Documentation

- [S-CSCF Features List](../../components/coeur/scscf/FEATURES_LIST.md)
- [ETSI TS 23.228](https://www.etsi.org/deliver/etsi_ts/123200_123299/123228/07.02.00_60/ts_123228v070200p.pdf) - IMS Stage 2
- [ETSI TS 29.228](https://www.etsi.org/deliver/etsi_ts/129200_129299/129228/) - Cx and Dx interfaces

# HSS Flow Diagrams

This document illustrates the key operational flows for the Home Subscriber Server (HSS) from its perspective within the IMS architecture.

## HSS In One Line

**HSS is the IMS "source of truth" for subscriber identity, authentication vectors, and which S-CSCF + service profile applies.**

## 1. S-CSCF Assignment During REGISTER (I-CSCF ↔ HSS)

I-CSCF queries HSS to determine which S-CSCF should serve the user during registration.

### Sequence Diagram

```mermaid
sequenceDiagram
    autonumber
    participant UE as UE
    participant PCSCF as P-CSCF
    participant ICSCF as I-CSCF
    participant HSS as HSS (Cx)
    participant SCSCF as S-CSCF

    UE->>PCSCF: REGISTER (Gm)
    PCSCF->>ICSCF: REGISTER (Mw)
    
    Note over ICSCF: Query HSS for<br/>S-CSCF assignment
    
    ICSCF->>HSS: Cx UAR (User-Authorization-Request)
    
    Note over HSS: HSS's job:<br/>S-CSCF selection/authorization<br/>based on user profile
    
    HSS-->>ICSCF: Cx UAA (Assigned S-CSCF)
    
    Note over ICSCF: Forward to<br/>assigned S-CSCF
    
    ICSCF->>SCSCF: REGISTER (Mw)
```

### What HSS Does Here

- ✅ Receives UAR (User Authorization Request) from I-CSCF
- ✅ Looks up subscriber data (IMPI/IMPU)
- ✅ Determines S-CSCF assignment (static or dynamic selection)
- ✅ Returns UAA (User Authorization Answer) with S-CSCF capabilities
- ✅ Authorizes user for IMS service
- ✅ Provides S-CSCF selection criteria

---

## 2. Authentication Vectors + Service Profile Download (S-CSCF ↔ HSS)

S-CSCF requests authentication vectors and service profile from HSS during registration.

### Sequence Diagram

```mermaid
sequenceDiagram
    autonumber
    participant SCSCF as S-CSCF
    participant HSS as HSS (Cx)

    Note over SCSCF: User sent REGISTER<br/>with credentials
    
    SCSCF->>HSS: Cx MAR (Multimedia-Auth-Request)
    
    Note over HSS: HSS's job:<br/>Generate IMS AKA vectors<br/>for authentication
    
    HSS-->>SCSCF: Cx MAA (Multimedia-Auth-Answer)<n/>(AKA vectors: RAND, AUTN, XRES, CK, IK)
    
    Note over SCSCF: Verify AKA response<br/>from UE
    
    SCSCF->>HSS: Cx SAR (Server-Assignment-Request)
    
    Note over HSS: HSS's job:<br/>Store S-CSCF assignment<br/>Download service profile
    
    HSS-->>SCSCF: Cx SAA (Server-Assignment-Answer)<n/>(Service Profile / iFC)
    
    Note over SCSCF: Store service profile<br/>Apply iFC for service logic
```

### What HSS Does Here

- ✅ Receives MAR (Multimedia Auth Request) from S-CSCF
- ✅ Generates IMS AKA authentication vectors (RAND, AUTN, XRES, CK, IK)
- ✅ Returns MAA (Multimedia Auth Answer) with vectors
- ✅ Receives SAR (Server Assignment Request) after successful authentication
- ✅ Stores S-CSCF assignment for user
- ✅ Downloads service profile including Initial Filter Criteria (iFC)
- ✅ Returns SAA (Server Assignment Answer) with complete service profile
- ✅ Updates registration state (user registered, serving S-CSCF)

---

## 3. Complete Registration Flow (HSS Perspective)

Full registration flow showing all HSS interactions.

### Sequence Diagram

```mermaid
sequenceDiagram
    autonumber
    participant UE as UE
    participant PCSCF as P-CSCF
    participant ICSCF as I-CSCF
    participant HSS as HSS
    participant SCSCF as S-CSCF

    UE->>PCSCF: REGISTER (initial, no auth)
    PCSCF->>ICSCF: REGISTER (Mw)
    
    Note over ICSCF: Need to find S-CSCF
    
    ICSCF->>HSS: Cx UAR (User-Authorization-Request)
    Note right of HSS: Lookup subscriber<br/>Select S-CSCF
    HSS-->>ICSCF: Cx UAA (S-CSCF capabilities)
    
    ICSCF->>SCSCF: REGISTER (Mw)
    
    Note over SCSCF: Need auth vectors
    
    SCSCF->>HSS: Cx MAR (Multimedia-Auth-Request)
    Note right of HSS: Generate AKA vectors
    HSS-->>SCSCF: Cx MAA (AKA vectors)
    
    SCSCF-->>UE: 401 Unauthorized (AKA challenge)
    
    UE->>PCSCF: REGISTER + Authorization (AKA response)
    PCSCF->>ICSCF: REGISTER + Authorization
    ICSCF->>SCSCF: REGISTER + Authorization
    
    Note over SCSCF: Verify AKA,<br/>now assign server
    
    SCSCF->>HSS: Cx SAR (Server-Assignment-Request)
    Note right of HSS: Store S-CSCF assignment<br/>Prepare service profile
    HSS-->>SCSCF: Cx SAA (Service Profile / iFC)
    
    SCSCF-->>UE: 200 OK (Registered)
```

### What HSS Does Here

- ✅ **UAR/UAA**: S-CSCF selection and authorization
- ✅ **MAR/MAA**: Authentication vector generation
- ✅ **SAR/SAA**: Server assignment and service profile delivery
- ✅ Maintains registration state
- ✅ Tracks serving S-CSCF for user

---

## 4. Service Profile Query (Sh Interface - Application Server)

Application Servers can query HSS for subscriber data via Sh interface.

### Sequence Diagram

```mermaid
sequenceDiagram
    autonumber
    participant AS as Application Server
    participant HSS as HSS (Sh)
    participant SCSCF as S-CSCF

    Note over AS: Service logic needs<br/>subscriber data
    
    AS->>HSS: Sh UDR (User-Data-Request)
    
    Note over HSS: HSS's job:<br/>Provide subscriber data<br/>to authorized AS
    
    HSS-->>AS: Sh UDA (User-Data-Answer)<n/>(Subscriber profile, preferences)
    
    Note over AS: Apply service logic<br/>based on data
    
    AS->>SCSCF: Service decision (ISC)
```

### What HSS Does Here

- ✅ Receives UDR (User Data Request) from AS
- ✅ Validates AS authorization
- ✅ Retrieves subscriber data (profile, preferences, etc.)
- ✅ Returns UDA (User Data Answer) with requested data
- ✅ Supports subscription/notification model for data updates

---

## 5. Registration State Management

HSS tracks which S-CSCF is serving each user and their registration status.

### State Diagram

```mermaid
stateDiagram-v2
    [*] --> Unregistered: User not registered
    Unregistered --> Registering: REGISTER received
    Registering --> Registered: SAR successful
    Registered --> Registered: Re-registration
    Registered --> Unregistered: Deregistration
    Registered --> Unregistered: Registration timeout
    
    note right of Registered
        HSS stores:
        - Serving S-CSCF
        - Registration timestamp
        - Service profile active
    end note
```

### What HSS Does Here

- ✅ Tracks registration state (registered/unregistered)
- ✅ Stores serving S-CSCF reference
- ✅ Updates state on SAR (Server Assignment Request)
- ✅ Updates state on deregistration
- ✅ Provides state to I-CSCF for routing decisions

---

## 6. S-CSCF Assignment Policies

HSS implements S-CSCF assignment logic based on subscriber profile and network policies.

### Assignment Flow

```mermaid
flowchart TD
    Start[UAR Received] --> Check{User has<br/>static S-CSCF?}
    Check -->|Yes| ReturnStatic[Return Static S-CSCF]
    Check -->|No| CheckCap{Check S-CSCF<br/>Capabilities}
    CheckCap --> SelectPool[Select from S-CSCF Pool]
    SelectPool --> LoadBalance[Apply Load Balancing]
    LoadBalance --> ReturnDynamic[Return Dynamic S-CSCF]
    ReturnStatic --> UAA[Send UAA]
    ReturnDynamic --> UAA
    UAA --> End[End]
```

### What HSS Does Here

- ✅ Checks for static S-CSCF assignment
- ✅ Evaluates S-CSCF capabilities if dynamic
- ✅ Applies load balancing across S-CSCF pool
- ✅ Returns appropriate S-CSCF in UAA
- ✅ Maintains assignment consistency

---

## 7. Authentication Vector Lifecycle

HSS generates and manages IMS AKA authentication vectors.

### Vector Generation Flow

```mermaid
sequenceDiagram
    autonumber
    participant SCSCF as S-CSCF
    participant HSS as HSS
    participant UE as UE

    Note over SCSCF: User attempting registration

    SCSCF->>HSS: Cx MAR (request vectors)
    
    Note over HSS: Generate AKA vectors:<n/>- RAND (random challenge)<n/>- AUTN (authentication token)<n/>- XRES (expected response)<n/>- CK (cipher key)<n/>- IK (integrity key)
    
    HSS-->>SCSCF: Cx MAA (AKA vectors)
    
    SCSCF->>UE: 401 Unauthorized (RAND, AUTN)
    
    UE->>SCSCF: REGISTER (RES, MAC)
    
    Note over SCSCF: Verify RES matches XRES<br/>Verify MAC
    
    alt Authentication Success
        SCSCF->>HSS: Cx SAR (server assignment)
        HSS-->>SCSCF: Cx SAA (success)
    else Authentication Failure
        SCSCF->>HSS: Cx MAR (request new vectors)
        HSS-->>SCSCF: Cx MAA (new vectors)
    end
```

### What HSS Does Here

- ✅ Generates cryptographically secure AKA vectors
- ✅ Uses subscriber's K key (shared secret)
- ✅ Provides vectors for authentication challenge
- ✅ Supports vector regeneration on failure
- ✅ Manages vector lifecycle and expiration

---

## Interface Summary

| Interface | Direction | Purpose | Key Messages |
|-----------|-----------|---------|--------------|
| **Cx** | I-CSCF/S-CSCF ↔ HSS | S-CSCF assignment, auth vectors, service profile | UAR/UAA, MAR/MAA, SAR/SAA |
| **Sh** | AS ↔ HSS | Service/subscriber data access | UDR/UDA, SNR/SNA |

## Key Diameter Messages (Cx Interface)

| Message | Direction | Purpose |
|---------|-----------|---------|
| **UAR** (User-Authorization-Request) | I-CSCF → HSS | Query for S-CSCF assignment |
| **UAA** (User-Authorization-Answer) | HSS → I-CSCF | Return S-CSCF capabilities |
| **MAR** (Multimedia-Auth-Request) | S-CSCF → HSS | Request authentication vectors |
| **MAA** (Multimedia-Auth-Answer) | HSS → S-CSCF | Return AKA vectors |
| **SAR** (Server-Assignment-Request) | S-CSCF → HSS | Assign S-CSCF, request service profile |
| **SAA** (Server-Assignment-Answer) | HSS → S-CSCF | Return service profile (iFC) |

---

## Related Documentation

- [HSS Features List](../../components/coeur/hss/FEATURES_LIST.md)
- [ETSI TS 23.228](https://www.etsi.org/deliver/etsi_ts/123200_123299/123228/07.02.00_60/ts_123228v070200p.pdf) - IMS Stage 2
- [ETSI TS 29.228](https://www.etsi.org/deliver/etsi_ts/129200_129299/129228/) - Cx and Dx interfaces
- [ETSI TS 29.329](https://www.etsi.org/deliver/etsi_ts/129300_129399/129329/) - Sh interface

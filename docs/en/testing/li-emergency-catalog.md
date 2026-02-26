# Lawful Intercept (LI) + Emergency Services Behavioral Test Catalogue
SIG-GW / IBCF Context — 2026 Edition

**Author**: Daniel  
**Scope**: Behavioral testing of Lawful Intercept and Emergency call handling  
**Context**: IMS + IBCF deployment (NNI/Ic)

**Standards Anchors** (jurisdiction dependent):
- 3GPP TS 33.107 (LI for IMS)
- 3GPP TS 23.167 (Emergency Sessions in IMS)
- ETSI TS 102 232 (Handover Interfaces)
- National regulatory frameworks (FCC, CRTC, etc.)

---

## Legend

**TEST ID format**: `LIE-###`

**Columns**:
- TEST ID
- SLOGAN
- AREA
- POSITIVE CASE
- NEGATIVE CASE
- LOAD CASE
- SCALE/RAMP CASE
- CHAOS CASE

**AREA codes**:
- `LI-C` = Lawful Intercept – Control Plane
- `LI-M` = Lawful Intercept – Media Plane
- `LI-A` = Lawful Intercept – Audit & Logging
- `EMR-R` = Emergency Routing
- `EMR-P` = Emergency Policy
- `EMR-M` = Emergency Media
- `EMR-L` = Emergency Location Handling
- `EMR-F` = Emergency Failure Handling

---

## Lawful Intercept + Emergency Test Matrix

| TEST ID | SLOGAN | AREA | POSITIVE CASE | NEGATIVE CASE | LOAD CASE | SCALE/RAMP CASE | CHAOS CASE |
|---------|--------|------|---------------|---------------|-----------|-----------------|------------|
| LIE-001 | "Targeted subscriber intercepted" | LI-C | Targeted TN → SIP signaling mirrored to LI mediation device | Non-target TN not mirrored | 2k CPS w/ 5% targets | Ramp targets 0→10% | LI mediation link down mid-call |
| LIE-002 | "Signaling-only intercept" | LI-C | SIP metadata exported correctly | Media not exported if signaling-only warrant | 1k CPS | Ramp signaling export rate | Corrupt warrant config reload |
| LIE-003 | "Media interception active" | LI-M | RTP duplicated to LI mediation | Media not leaked for non-target | 5k RTP streams | Ramp media flows | Media duplication node crash |
| LIE-004 | "Interception continuity on re-INVITE" | LI-M | Mid-call SDP change still intercepted | No loss of intercept during hold/resume | 1k CPS w/ re-INVITE | Ramp re-INVITE rate | Force media relay restart |
| LIE-005 | "Intercept after call transfer" | LI-C | REFER transfer retains interception | Transfer to non-target not intercepted | 500 CPS REFER | Ramp transfers | Kill LI controller mid-transfer |
| LIE-006 | "Warrant activation/deactivation" | LI-A | Activation immediate; intercept begins next call | Deactivated warrant stops intercept | 500 CPS | Ramp warrant toggles | Toggle warrant under 2k CPS |
| LIE-007 | "Audit log completeness" | LI-A | Intercept start/stop logged with timestamp | No missing entries | 2k CPS | Ramp log volume | Logging backend down |
| LIE-008 | "Privacy enforcement" | LI-A | Only authorized operators can view intercept logs | Unauthorized access blocked | N/A | Ramp audit queries | Attempt privilege escalation |
| LIE-009 | "Multi-target scaling" | LI-M | Multiple simultaneous intercept targets handled | No cross-leak between targets | 5k CPS 10% targets | Ramp targets 1→500 | Exhaust LI mediation throughput |
| LIE-010 | "Interconnect intercept policy" | LI-C | Intercepted calls across NNI preserved | Intercept not dropped at IBCF boundary | 2k CPS cross-network | Ramp NNI traffic | Peer restart mid-intercept |
| LIE-011 | "Encrypted signaling intercept" | LI-C | TLS decrypted internally for LI export | Cannot intercept without lawful authority | 1k CPS TLS | Ramp TLS ratio | TLS key reload during intercept |
| LIE-012 | "Intercept during overload" | LI-M | Under overload, intercept preserved | No silent drop of LI | 10k CPS overload | Ramp overload | CPU throttling injected |
| LIE-101 | "Emergency call routing works" | EMR-R | Emergency dial (e.g. 911/112) routed to PSAP | Non-emergency cannot use emergency path | 200 CPS emergency | Ramp emergency % | Remove primary PSAP route |
| LIE-102 | "Emergency priority handling" | EMR-P | Emergency INVITE prioritized in queue | Normal traffic cannot starve emergency | 5k CPS mixed | Ramp emergency bursts | Flood normal traffic 20k CPS |
| LIE-103 | "Emergency location header handling" | EMR-L | P-Access-Network-Info / Geolocation preserved | Missing location flagged | 200 CPS | Ramp location data size | Corrupt location header |
| LIE-104 | "Emergency call bypass policy" | EMR-P | Emergency call allowed even if subscriber barred | Fraud block not applied to emergency | 500 CPS barred subs | Ramp barred % | Toggle subscriber state mid-call |
| LIE-105 | "Emergency callback routing" | EMR-R | Callback from PSAP routed correctly | Non-PSAP spoof rejected | 200 CPS callback | Ramp callback % | Kill routing DB |
| LIE-106 | "Emergency media integrity" | EMR-M | Media stable with low latency | No transcoding distortion | 500 RTP flows | Ramp media load | Packet loss 10% |
| LIE-107 | "Emergency during failover" | EMR-F | Active/Active IBCF failover does not drop emergency | Call preserved | 2k CPS | Ramp failovers | Kill active instance mid-emergency |
| LIE-108 | "Emergency STIR handling" | EMR-P | Emergency allowed even if STIR fails | Attestation not blocking emergency | 500 CPS mixed | Ramp invalid STIR % | Remove cert chain mid-call |
| LIE-109 | "Emergency logging compliance" | EMR-L | Logs show timestamp, route, PSAP ID | No sensitive location leakage | 200 CPS | Ramp log verbosity | Logging DB unavailable |
| LIE-110 | "International emergency handling" | EMR-R | Roaming emergency routed correctly | Foreign emergency codes not misrouted | 200 CPS roaming | Ramp roaming % | Remove roaming routing entry |

---

## Explicit Examples (Required Coverage)

**Positive Example**:  
LIE-001 — Targeted subscriber signaling mirrored correctly.

**Negative Example**:  
LIE-004 — Interception not applied to non-target after transfer.

**Load Example**:  
LIE-009 — 5k CPS with 10% intercept targets.

**Scale/Ramp Example**:  
LIE-102 — Ramp emergency call ratio during 5k CPS background load.

**Chaos Example**:  
LIE-107 — Kill active IBCF instance mid-emergency call; verify continuity.

---

## Critical Behavioral Assertions

1. **Emergency calls must NEVER be blocked due to**:
   - STIR failure
   - Fraud detection
   - Billing restriction
   - Rate limiting

2. **Lawful Intercept must**:
   - Be silent to user
   - Not alter signaling behavior
   - Not degrade call quality
   - Be auditable and tamper-evident

3. **Intercept must persist across**:
   - Re-INVITE
   - Call transfer
   - Interconnect traversal
   - Media renegotiation

---

## Operational KPIs to Validate

- Emergency PDD ≤ configured regulatory target
- Intercept latency overhead ≤ X ms
- 0% emergency drop during failover
- 100% intercept coverage for active warrants
- No topology leakage in emergency routing

---

## Regulatory Compliance

### Lawful Intercept Requirements

- **3GPP TS 33.107**: IMS Lawful Intercept
- **ETSI TS 102 232**: Handover Interface specifications
- **National Regulations**: Jurisdiction-specific requirements

### Emergency Services Requirements

- **3GPP TS 23.167**: Emergency sessions in IMS
- **FCC Requirements**: 911/E911 compliance (US)
- **CRTC Requirements**: 911 compliance (Canada)
- **EU Requirements**: 112 compliance

---

## Implementation Notes

### Lawful Intercept Architecture

```
IMS Core
  |
  v
IBCF/SIG-GW
  |
  +-- LI Controller
  |   - Warrant management
  |   - Target identification
  |   - Intercept activation
  |
  +-- LI Mediation Device (MD)
  |   - Signaling duplication
  |   - Media duplication
  |   - Handover Interface (HI)
  |
  v
Law Enforcement Agency (LEA)
```

### Emergency Services Architecture

```
UE
  |
  v
P-CSCF
  |
  v
S-CSCF
  |
  +-- Emergency Detection
  |   - Number recognition (911/112/etc.)
  |   - Priority handling
  |
  +-- Emergency Routing
  |   - PSAP selection
  |   - Location-based routing
  |
  v
PSAP (Public Safety Answering Point)
```

### Key Implementation Points

1. **LI must be transparent** - No impact on call quality or behavior
2. **Emergency must bypass all restrictions** - Highest priority
3. **Audit trails required** - Tamper-evident logging
4. **Privacy protection** - Only authorized access
5. **Regulatory compliance** - Jurisdiction-specific rules

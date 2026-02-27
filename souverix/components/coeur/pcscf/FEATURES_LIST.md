# P-CSCF Features List

This document outlines the standard responsibilities and features of the Proxy Call Session Control Function (P-CSCF) as defined by ETSI/3GPP specifications.

## Related Specifications

- **ETSI TS 123 228 V7.2.0** (3GPP TS 23.228): [IP Multimedia Subsystem (IMS); Stage 2](https://www.etsi.org/deliver/etsi_ts/123200_123299/123228/07.02.00_60/ts_123228v070200p.pdf)
- **ETSI TS 124 229 V16.4.0** (3GPP TS 24.229): [IP multimedia call control protocol based on Session Initiation Protocol (SIP) and Session Description Protocol (SDP)](https://www.etsi.org/deliver/etsi_ts/124200_124299/124229/16.04.00_60/ts_124229v160400p.pdf)

## P-CSCF Operational Definition

**P-CSCF is the first SIP contact point for the UE and the IMS security + policy enforcement edge for SIP signaling.**

### Key Functions

- ✅ Anchors UE SIP signaling
- ✅ Enforces security (IPsec/TLS)
- ✅ Applies policy control (PCRF/PCF interaction)
- ✅ Maintains SIP routing state toward S-CSCF
- ✅ Performs topology hiding toward access side

### What P-CSCF Is NOT

- ❌ It is **not** a service logic engine
- ❌ It is **not** breakout logic
- ❌ It is an **access-edge control function**

## P-CSCF Base-Set Feature Table

| Feature / Responsibility | Primary vs Secondary | ETSI/3GPP Spec | Call / Message / Interface Context (What You'll Actually See) |
|------------------------|---------------------|----------------|------------------------------------------------------------------|
| First contact point for UE SIP signaling | Primary | [TS 23.228](https://www.etsi.org/deliver/etsi_ts/123200_123299/123228/07.02.00_60/ts_123228v070200p.pdf) | **Gm interface**: UE → P-CSCF over Gm (SIP REGISTER, INVITE, MESSAGE, etc.) |
| SIP message forwarding toward S-CSCF | Primary | [TS 23.228](https://www.etsi.org/deliver/etsi_ts/123200_123299/123228/07.02.00_60/ts_123228v070200p.pdf) | **Mw interface**: P-CSCF → I-CSCF/S-CSCF over Mw |
| Security association establishment (IPsec / TLS) | Primary | [TS 24.229](https://www.etsi.org/deliver/etsi_ts/124200_124299/124229/16.04.00_60/ts_124229v160400p.pdf) | During REGISTER; IPSec SA negotiation (IMS AKA, Security-Client/Server headers) |
| Integrity and replay protection | Primary | [TS 24.229](https://www.etsi.org/deliver/etsi_ts/124200_124299/124229/16.04.00_60/ts_124229v160400p.pdf) | IPSec ESP protection of SIP signaling |
| Emergency call handling (routing toward E-CSCF) | Primary | [TS 23.228](https://www.etsi.org/deliver/etsi_ts/123200_123299/123228/07.02.00_60/ts_123228v070200p.pdf) | INVITE with emergency indication |
| SIP compression support (SigComp) | Secondary (optional) | [TS 24.229](https://www.etsi.org/deliver/etsi_ts/124200_124299/124229/16.04.00_60/ts_124229v160400p.pdf) | Gm compression negotiation |
| Topology hiding toward UE side | Secondary | [TS 23.228](https://www.etsi.org/deliver/etsi_ts/123200_123299/123228/07.02.00_60/ts_123228v070200p.pdf) | Via header manipulation / Via/Record-Route management |
| Policy Control interaction (PCRF / PCF) | Primary in EPC/5GC | [TS 23.228](https://www.etsi.org/deliver/etsi_ts/123200_123299/123228/07.02.00_60/ts_123228v070200p.pdf) | **Rx interface**: Diameter / HTTP2 in 5GC |
| QoS authorization trigger | Primary (in LTE/5G) | [TS 23.228](https://www.etsi.org/deliver/etsi_ts/123200_123299/123228/07.02.00_60/ts_123228v070200p.pdf) | Media authorization request toward PCRF/PCF |
| SIP header inspection and validation | Primary | [TS 24.229](https://www.etsi.org/deliver/etsi_ts/124200_124299/124229/16.04.00_60/ts_124229v160400p.pdf) | Validate From/To/Contact/Route integrity |
| NAT traversal support | Primary in real-world deployments | [TS 24.229](https://www.etsi.org/deliver/etsi_ts/124200_124299/124229/16.04.00_60/ts_124229v160400p.pdf) | Keep-alive (CRLF), SIP OPTIONS, STUN binding |
| Emergency registration handling | Primary | [TS 23.228](https://www.etsi.org/deliver/etsi_ts/123200_123299/123228/07.02.00_60/ts_123228v070200p.pdf) | Special REGISTER processing |
| Maintain dialog state (Record-Route) | Primary operational behavior | [TS 24.229](https://www.etsi.org/deliver/etsi_ts/124200_124299/124229/16.04.00_60/ts_124229v160400p.pdf) | Insert Record-Route to remain in path |
| Lawful intercept support (where required) | Secondary (deployment dependent) | [TS 33.106](https://www.etsi.org/deliver/etsi_ts/133100_133199/133106/) | Mirror signaling metadata |
| Charging correlation trigger | Secondary | [TS 32.240](https://www.etsi.org/deliver/etsi_ts/132200_132299/132240/) | Interaction with CCF/OCF in some deployments |

## Interface Reference

- **Gm**: Interface between UE and P-CSCF (SIP over IPsec/TLS)
- **Mw**: Interface between P-CSCF and I-CSCF/S-CSCF
- **Rx**: Interface between P-CSCF and PCRF/PCF (Policy Control)

## Notes

- **Primary** features are core responsibilities that must be implemented for P-CSCF to function correctly.
- **Secondary** features are optional or supporting capabilities that enhance functionality but may not be strictly required for basic operation.
- P-CSCF is **access-bound** and does not participate in inter-domain routing or service logic execution.
- P-CSCF deployment is typically **geographically distributed** close to access networks for low latency.

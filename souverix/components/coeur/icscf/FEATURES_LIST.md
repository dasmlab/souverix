# I-CSCF Features List

This document outlines the standard responsibilities and features of the Interrogating Call Session Control Function (I-CSCF) as defined by ETSI/3GPP specifications.

## Related Specifications

- **ETSI TS 123 228 V7.2.0** (3GPP TS 23.228): [IP Multimedia Subsystem (IMS); Stage 2](https://www.etsi.org/deliver/etsi_ts/123200_123299/123228/07.02.00_60/ts_123228v070200p.pdf)
- **ETSI TS 129 228 V16.4.0** (3GPP TS 29.228): [IP Multimedia (IM) Subsystem Cx and Dx interfaces; Signalling flows and message contents](https://www.etsi.org/deliver/etsi_ts/129200_129299/129228/)

## I-CSCF Operational Definition

**I-CSCF is the entry routing and HSS query point of an IMS domain.**

### Key Functions

- ✅ Receives SIP from external domains or P-CSCF
- ✅ Queries HSS to determine which S-CSCF serves the user
- ✅ Hides internal topology
- ✅ Routes SIP toward the correct S-CSCF

### What I-CSCF Is

- ✅ It is a **routing and lookup function**
- ✅ It is **not** a service engine

## I-CSCF Base-Set Feature Table

| Feature / Responsibility | Primary vs Secondary | ETSI/3GPP Spec | Interface / Message Context (What You'll Actually See) |
|------------------------|---------------------|----------------|--------------------------------------------------------|
| Entry point for SIP into IMS domain | Primary | [TS 23.228](https://www.etsi.org/deliver/etsi_ts/123200_123299/123228/07.02.00_60/ts_123228v070200p.pdf) | **Mw interface**: From P-CSCF or external IBCF |
| HSS query for S-CSCF assignment | Primary | [TS 29.228](https://www.etsi.org/deliver/etsi_ts/129200_129299/129228/) | **Cx interface**: Diameter UAR (User-Authorization-Request) / LIR (Location-Info-Request) |
| Select S-CSCF for registration/session | Primary | [TS 23.228](https://www.etsi.org/deliver/etsi_ts/123200_123299/123228/07.02.00_60/ts_123228v070200p.pdf) | Based on HSS response (S-CSCF capabilities) |
| Topology hiding (internal CSCF structure) | Primary | [TS 23.228](https://www.etsi.org/deliver/etsi_ts/123200_123299/123228/07.02.00_60/ts_123228v070200p.pdf) | Header manipulation / route control |
| Forward SIP toward selected S-CSCF | Primary | [TS 23.228](https://www.etsi.org/deliver/etsi_ts/123200_123299/123228/07.02.00_60/ts_123228v070200p.pdf) | **Mw interface**: Forward to S-CSCF |
| Emergency routing fallback | Secondary | [TS 23.228](https://www.etsi.org/deliver/etsi_ts/123200_123299/123228/07.02.00_60/ts_123228v070200p.pdf) | Domain policy dependent |
| Registration routing during initial attach | Primary | [TS 23.228](https://www.etsi.org/deliver/etsi_ts/123200_123299/123228/07.02.00_60/ts_123228v070200p.pdf) | REGISTER path |
| Load balancing across S-CSCF pool | Secondary (deployment) | [TS 23.228](https://www.etsi.org/deliver/etsi_ts/123200_123299/123228/07.02.00_60/ts_123228v070200p.pdf) | Policy-based selection |
| ENUM / DNS support (optional) | Secondary | [TS 23.228](https://www.etsi.org/deliver/etsi_ts/123200_123299/123228/07.02.00_60/ts_123228v070200p.pdf) | External routing decisions |

## Interface Reference

- **Mw**: Interface between P-CSCF/I-CSCF and I-CSCF/S-CSCF
- **Cx**: Interface between I-CSCF and HSS (Diameter)

## I-CSCF In One Line

**I-CSCF decides which S-CSCF handles a user and shields internal IMS structure.**

## Notes

- **Primary** features are core responsibilities that must be implemented for I-CSCF to function correctly.
- **Secondary** features are optional or supporting capabilities that enhance functionality but may not be strictly required for basic operation.
- I-CSCF is **stateless** for sessions but maintains routing state for S-CSCF selection.
- I-CSCF provides **topology hiding** by not exposing internal S-CSCF structure to external domains.

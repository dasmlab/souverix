# S-CSCF Features List

This document outlines the standard responsibilities and features of the Serving Call Session Control Function (S-CSCF) as defined by ETSI/3GPP specifications.

## Related Specifications

- **ETSI TS 123 228 V7.2.0** (3GPP TS 23.228): [IP Multimedia Subsystem (IMS); Stage 2](https://www.etsi.org/deliver/etsi_ts/123200_123299/123228/07.02.00_60/ts_123228v070200p.pdf)
- **ETSI TS 129 228 V16.4.0** (3GPP TS 29.228): [IP Multimedia (IM) Subsystem Cx and Dx interfaces; Signalling flows and message contents](https://www.etsi.org/deliver/etsi_ts/129200_129299/129228/)

## S-CSCF Operational Definition

**S-CSCF is the central SIP service engine of IMS.**

### Key Functions

- ✅ Maintains subscriber registration state
- ✅ Executes service logic
- ✅ Triggers Application Servers
- ✅ Anchors SIP dialogs
- ✅ Enforces service profiles

### What S-CSCF Is

- ✅ **This is the brain**
- ✅ **Everything intelligent happens here**

## S-CSCF Base-Set Feature Table

| Feature / Responsibility | Primary vs Secondary | ETSI/3GPP Spec | Interface / Message Context (What You'll Actually See) |
|------------------------|---------------------|----------------|--------------------------------------------------------|
| Maintain registration state | Primary | [TS 23.228](https://www.etsi.org/deliver/etsi_ts/123200_123299/123228/07.02.00_60/ts_123228v070200p.pdf) | REGISTER handling |
| User authentication (IMS AKA challenge) | Primary | [TS 23.228](https://www.etsi.org/deliver/etsi_ts/123200_123299/123228/07.02.00_60/ts_123228v070200p.pdf) | **Cx interface**: 401 challenge / AKA verification |
| Service profile download from HSS | Primary | [TS 29.228](https://www.etsi.org/deliver/etsi_ts/129200_129299/129228/) | **Cx interface**: SAR (Server-Assignment-Request) / SAA (Server-Assignment-Answer) |
| Application Server triggering (Initial Filter Criteria) | Primary | [TS 23.228](https://www.etsi.org/deliver/etsi_ts/123200_123299/123228/07.02.00_60/ts_123228v070200p.pdf) | **ISC interface**: Trigger AS based on iFC |
| SIP dialog anchoring (Record-Route) | Primary | [TS 23.228](https://www.etsi.org/deliver/etsi_ts/123200_123299/123228/07.02.00_60/ts_123228v070200p.pdf) | All INVITE dialogs |
| Service invocation (VoLTE, VoWiFi, etc.) | Primary | [TS 23.228](https://www.etsi.org/deliver/etsi_ts/123200_123299/123228/07.02.00_60/ts_123228v070200p.pdf) | INVITE processing |
| Breakout decision trigger toward BGCF | Primary | [TS 23.228](https://www.etsi.org/deliver/etsi_ts/123200_123299/123228/07.02.00_60/ts_123228v070200p.pdf) | Routing logic (tel: URI → BGCF) |
| ENUM / number analysis | Primary | [TS 23.228](https://www.etsi.org/deliver/etsi_ts/123200_123299/123228/07.02.00_60/ts_123228v070200p.pdf) | Tel URI processing |
| Call barring / ODB enforcement | Primary | [TS 23.228](https://www.etsi.org/deliver/etsi_ts/123200_123299/123228/07.02.00_60/ts_123228v070200p.pdf) | Based on subscriber profile |
| Emergency service invocation | Primary | [TS 23.228](https://www.etsi.org/deliver/etsi_ts/123200_123299/123228/07.02.00_60/ts_123228v070200p.pdf) | Routing toward E-CSCF |
| Session release management | Primary | [TS 23.228](https://www.etsi.org/deliver/etsi_ts/123200_123299/123228/07.02.00_60/ts_123228v070200p.pdf) | BYE handling |
| Charging event generation (CCF/OCF) | Primary | [TS 32.240](https://www.etsi.org/deliver/etsi_ts/132200_132299/132240/) | **Ro/Rf interface**: Diameter charging |
| Third-party registration support | Secondary | [TS 23.228](https://www.etsi.org/deliver/etsi_ts/123200_123299/123228/07.02.00_60/ts_123228v070200p.pdf) | AS integration |
| Forking logic (parallel INVITE) | Primary | [TS 23.228](https://www.etsi.org/deliver/etsi_ts/123200_123299/123228/07.02.00_60/ts_123228v070200p.pdf) | SIP routing control |
| Mid-dialog request handling (reINVITE, UPDATE) | Primary | [TS 23.228](https://www.etsi.org/deliver/etsi_ts/123200_123299/123228/07.02.00_60/ts_123228v070200p.pdf) | Session management |

## Interface Reference

- **Mw**: Interface between I-CSCF and S-CSCF
- **Cx**: Interface between S-CSCF and HSS (Diameter)
- **ISC**: Interface between S-CSCF and Application Servers
- **Ro/Rf**: Interface between S-CSCF and Charging Functions (CCF/OCF)

## S-CSCF In One Line

**S-CSCF is the stateful SIP service brain of the IMS domain.**

## Notes

- **Primary** features are core responsibilities that must be implemented for S-CSCF to function correctly.
- **Secondary** features are optional or supporting capabilities that enhance functionality but may not be strictly required for basic operation.
- S-CSCF is **stateful** and maintains registration and dialog state.
- S-CSCF is the **central service execution point** for all IMS services.

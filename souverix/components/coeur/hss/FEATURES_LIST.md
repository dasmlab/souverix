# HSS Features List

This document outlines the standard responsibilities and features of the Home Subscriber Server (HSS) as defined by ETSI/3GPP specifications.

## Related Specifications

- **ETSI TS 123 228 V7.2.0** (3GPP TS 23.228): [IP Multimedia Subsystem (IMS); Stage 2](https://www.etsi.org/deliver/etsi_ts/123200_123299/123228/07.02.00_60/ts_123228v070200p.pdf)
- **ETSI TS 129 228 V16.4.0** (3GPP TS 29.228): [IP Multimedia (IM) Subsystem Cx and Dx interfaces; Signalling flows and message contents](https://www.etsi.org/deliver/etsi_ts/129200_129299/129228/)
- **ETSI TS 129 329 V16.4.0** (3GPP TS 29.329): [Sh interface based on the Diameter protocol; Protocol details](https://www.etsi.org/deliver/etsi_ts/129300_129399/129329/)

## HSS Operational Definition

**HSS (Home Subscriber Server) is the authoritative subscriber + service-profile database for IMS, and the policy decision point for 'which S-CSCF serves this user' plus authentication vector distribution.**

### Key Functions

- ✅ Stores subscriber identities and IMS subscription data
- ✅ Issues IMS AKA authentication vectors
- ✅ Assigns / confirms S-CSCF for a user
- ✅ Provides service profiles (iFC) and registration status hooks

### What HSS Is

- ✅ It is **data + authorization + assignment**, not call control

## HSS Base-Set Feature Table

| Feature / Responsibility | Primary vs Secondary | ETSI/3GPP Spec | Interface / Message Context (What You'll Actually See) |
|------------------------|---------------------|----------------|--------------------------------------------------------|
| Store IMS subscriber data (IMPU/IMPI, public/private IDs, subscription) | Primary | [TS 23.228](https://www.etsi.org/deliver/etsi_ts/123200_123299/123228/07.02.00_60/ts_123228v070200p.pdf) | **Cx / Sh interface**: Backed data model |
| S-CSCF assignment for a user (selection/authorization) | Primary | [TS 29.228](https://www.etsi.org/deliver/etsi_ts/129200_129299/129228/) | **Cx interface**: UAR/UAA (I-CSCF queries) |
| Authentication vector generation & delivery (IMS AKA) | Primary | [TS 29.228](https://www.etsi.org/deliver/etsi_ts/129200_129299/129228/) | **Cx interface**: MAR/MAA (S-CSCF requests vectors) |
| Service profile delivery (iFC, TAS triggers, barring, etc.) | Primary | [TS 29.228](https://www.etsi.org/deliver/etsi_ts/129200_129299/129228/) | **Cx interface**: SAR/SAA (Server Assignment) |
| Registration state tracking (who is serving, registered/not) | Primary | [TS 29.228](https://www.etsi.org/deliver/etsi_ts/129200_129299/129228/) | Updated via SAR/SAA; may support deregistration |
| Location / reachability info (serving node reference, roaming hints) | Secondary (deployment) | [TS 23.228](https://www.etsi.org/deliver/etsi_ts/123200_123299/123228/07.02.00_60/ts_123228v070200p.pdf) | Used for routing decisions; varies by vendor |
| User authorization / roaming restrictions | Primary | [TS 29.228](https://www.etsi.org/deliver/etsi_ts/129200_129299/129228/) | **Cx interface**: Authorization results to I/S-CSCF |
| Support Sh queries from Application Servers | Secondary (common in MMTel) | [TS 29.329](https://www.etsi.org/deliver/etsi_ts/129300_129399/129329/) | **Sh interface**: AS → HSS data access |
| Notify service changes (profile updates) | Secondary | [TS 29.329](https://www.etsi.org/deliver/etsi_ts/129300_129399/129329/) | **Sh interface**: Notifications / subscription models |
| Support interworking evolution (HSS → UDM) | Secondary (future-proofing) | [TS 23.501](https://www.etsi.org/deliver/etsi_ts/123500_123599/123501/) | Migration path in 5GC architectures |

## Interface Reference

- **Cx**: Interface between I-CSCF/S-CSCF and HSS (Diameter)
  - **Purpose**: S-CSCF assignment, auth vectors, service profile
- **Sh**: Interface between Application Server and HSS (Diameter)
  - **Purpose**: Service/subscriber data access (optional but common)

## HSS In One Line

**HSS is the IMS "source of truth" for subscriber identity, authentication vectors, and which S-CSCF + service profile applies.**

## What HSS Does NOT Do

- ❌ Does **not** forward SIP
- ❌ Does **not** maintain SIP dialogs
- ❌ Does **not** run service logic (AS does)
- ❌ Does **not** interwork PSTN (MGCF does)
- ❌ Does **not** enforce access security (P-CSCF does)

## Souverix "Base-Set" HSS Contract (Practical)

If you're making this implementation-ready, your HSS contract for the platform should explicitly cover:

- ✅ **Identity objects**: IMPI, IMPU, service domain, aliases
- ✅ **S-CSCF assignment policy**: static vs dynamic, pool selection rules
- ✅ **Auth vectors**: AKA support + vector lifecycle rules
- ✅ **Service profile**: iFC structure + AS targets + barring flags
- ✅ **State**: registration status + serving S-CSCF binding
- ✅ **APIs**: Cx mandatory; Sh optional (but recommended for MMTel/AS ecosystems)

## Notes

- **Primary** features are core responsibilities that must be implemented for HSS to function correctly.
- **Secondary** features are optional or supporting capabilities that enhance functionality but may not be strictly required for basic operation.
- HSS is **stateless** for call control but maintains persistent subscriber and registration state.
- HSS is the **authoritative source** for all subscriber-related data in IMS.

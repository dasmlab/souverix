# MGCF Features List

This document outlines the standard responsibilities and features of the Media Gateway Control Function (MGCF) as defined by ETSI/3GPP specifications.

## Related Specifications

- **ETSI TS 123 228 V7.2.0** (3GPP TS 23.228): [IP Multimedia Subsystem (IMS); Stage 2](https://www.etsi.org/deliver/etsi_ts/123200_123299/123228/07.02.00_60/ts_123228v070200p.pdf)
- **ETSI TS 129 163 V16.4.0** (3GPP TS 29.163): [Interworking between the IP Multimedia (IM) Core Network (CN) subsystem and Circuit Switched (CS) networks](https://www.etsi.org/deliver/etsi_ts/129100_129199/129163/16.04.00_60/ts_129163v160400p.pdf)

## MGCF Operational Definition

**MGCF (Media Gateway Control Function) is the SIP ↔ circuit-switched signaling interworking controller between IMS and the CS/PSTN domain.**

### Key Functions

- ✅ Terminates SIP from IMS
- ✅ Translates SIP signaling to ISUP/BICC
- ✅ Controls Media Gateway (MGW) via H.248
- ✅ Anchors call state toward PSTN
- ✅ Handles call release, cause mapping, timers

### What MGCF Is

- ✅ It is **signaling interworking**, not media processing

## MGCF Base-Set Feature Table

| Feature / Responsibility | Primary vs Secondary | ETSI/3GPP Spec | Call / Interface Context (What You'll Actually See) |
|------------------------|---------------------|----------------|-----------------------------------------------------|
| SIP termination from IMS (from BGCF/S-CSCF) | Primary | [TS 23.228](https://www.etsi.org/deliver/etsi_ts/123200_123299/123228/07.02.00_60/ts_123228v070200p.pdf) | **Mj interface**: SIP over Mj |
| ISUP/BICC interworking toward PSTN/CS | Primary | [TS 29.163](https://www.etsi.org/deliver/etsi_ts/129100_129199/129163/16.04.00_60/ts_129163v160400p.pdf) | **Nc interface**: ISUP/BICC signaling |
| SIP ↔ ISUP message mapping | Primary | [TS 29.163](https://www.etsi.org/deliver/etsi_ts/129100_129199/129163/16.04.00_60/ts_129163v160400p.pdf) | INVITE ↔ IAM, 180 ↔ ACM, 200 OK ↔ ANM |
| H.248 control of Media Gateway (MGW) | Primary | [TS 29.163](https://www.etsi.org/deliver/etsi_ts/129100_129199/129163/16.04.00_60/ts_129163v160400p.pdf) | **Mg interface**: H.248 / Megaco |
| Cause code mapping (Q.850 ↔ SIP codes) | Primary | [TS 29.163](https://www.etsi.org/deliver/etsi_ts/129100_129199/129163/16.04.00_60/ts_129163v160400p.pdf) | 4xx/5xx ↔ ISUP REL (Release) |
| Early media handling coordination | Primary | [TS 29.163](https://www.etsi.org/deliver/etsi_ts/129100_129199/129163/16.04.00_60/ts_129163v160400p.pdf) | ACM with in-band tones |
| Call state machine control (CS domain) | Primary | [TS 29.163](https://www.etsi.org/deliver/etsi_ts/129100_129199/129163/16.04.00_60/ts_129163v160400p.pdf) | PSTN call supervision |
| Release handling and resource teardown | Primary | [TS 29.163](https://www.etsi.org/deliver/etsi_ts/129100_129199/129163/16.04.00_60/ts_129163v160400p.pdf) | BYE ↔ REL/RLC (Release Complete) |
| Timer management (T1, T7 equivalents) | Primary | [TS 29.163](https://www.etsi.org/deliver/etsi_ts/129100_129199/129163/16.04.00_60/ts_129163v160400p.pdf) | SIP ↔ ISUP timer harmonization |
| Emergency call interworking | Primary | [TS 23.228](https://www.etsi.org/deliver/etsi_ts/123200_123299/123228/07.02.00_60/ts_123228v070200p.pdf) | Priority mapping |
| TDM trunk selection | Secondary (deployment-specific) | [TS 29.163](https://www.etsi.org/deliver/etsi_ts/129100_129199/129163/16.04.00_60/ts_129163v160400p.pdf) | PSTN routing logic |
| Charging event trigger | Secondary | [TS 32.240](https://www.etsi.org/deliver/etsi_ts/132200_132299/132240/) | Correlation with CCF |
| Fax / Modem handling signaling coordination | Secondary | [TS 29.163](https://www.etsi.org/deliver/etsi_ts/129100_129199/129163/16.04.00_60/ts_129163v160400p.pdf) | T.38 coordination |
| Number normalization / E.164 formatting | Secondary | [TS 29.163](https://www.etsi.org/deliver/etsi_ts/129100_129199/129163/16.04.00_60/ts_129163v160400p.pdf) | SIP URI ↔ ISUP called number |

## Interface Reference

- **Mj**: Interface between BGCF/S-CSCF and MGCF (SIP signaling)
- **Mg**: Interface between MGCF and MGW (H.248 control)
- **Nc**: Interface between MGCF and PSTN (ISUP/BICC)

## MGCF In One Line

**MGCF translates IMS SIP signaling into legacy circuit signaling and controls the Media Gateway that carries the actual RTP↔TDM media.**

## Key Separation

- ✅ **MGCF = signaling**
- ✅ **MGW = media**
- ✅ **MGCF does not touch RTP packets directly**

## What MGCF Does NOT Do

- ❌ Does **not** perform SIP service logic
- ❌ Does **not** select breakout (BGCF does that)
- ❌ Does **not** authenticate subscribers
- ❌ Does **not** anchor RTP directly (MGW does)
- ❌ Does **not** enforce IMS policy (P-CSCF does that)

## Notes

- **Primary** features are core responsibilities that must be implemented for MGCF to function correctly.
- **Secondary** features are optional or supporting capabilities that enhance functionality but may not be strictly required for basic operation.
- MGCF is **signaling-only** - media conversion is handled by MGW.
- MGCF deployment is typically **geographically distributed** near TDM trunks.

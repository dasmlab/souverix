# BGCF Features List

This document outlines the standard responsibilities and features of the Breakout Gateway Control Function (BGCF) as defined by ETSI/3GPP specifications.

## Related Specifications

- **ETSI TS 123 228 V7.2.0** (3GPP TS 23.228): [IP Multimedia Subsystem (IMS); Stage 2](https://www.etsi.org/deliver/etsi_ts/123200_123299/123228/07.02.00_60/ts_123228v070200p.pdf)
- **ETSI TS 129 163 V16.4.0** (3GPP TS 29.163): [Interworking between the IP Multimedia (IM) Core Network (CN) subsystem and Circuit Switched (CS) networks](https://www.etsi.org/deliver/etsi_ts/129100_129199/129163/16.04.00_60/ts_129163v160400p.pdf)

## BGCF Base-Set Feature Table

| Feature / Responsibility | Primary vs Secondary | ETSI/3GPP Spec | Call / Message / Feature-Path Context (What You'll Actually See) |
|------------------------|---------------------|----------------|------------------------------------------------------------------|
| Accept breakout selection request from S-CSCF (BGCF is invoked when S-CSCF determines PSTN/CS breakout is needed) | Primary | [TS 23.228](https://www.etsi.org/deliver/etsi_ts/123200_123299/123228/07.02.00_60/ts_123228v070200p.pdf) | **Mi interface**: Typically a SIP INVITE arriving from S-CSCF → BGCF for PSTN/CS-domain termination/breakout selection. |
| Select breakout network (decide which network will perform PSTN/CS interworking) | Primary | [TS 23.228](https://www.etsi.org/deliver/etsi_ts/123200_123299/123228/07.02.00_60/ts_123228v070200p.pdf) | Internal policy/routing decision; outcome is either "local breakout" (select MGCF) or "remote breakout" (forward to another BGCF). |
| Forward signalling to BGCF in selected breakout network (when breakout is not local) | Primary | [TS 23.228](https://www.etsi.org/deliver/etsi_ts/123200_123299/123228/07.02.00_60/ts_123228v070200p.pdf) | **Mk interface**: Forward the dialog (starting with INVITE) to the BGCF of the selected network. In practice this is still SIP routing, often with Route / Request-URI steering. |
| Select MGCF in the breakout network (when breakout is local to that network) | Primary | [TS 23.228](https://www.etsi.org/deliver/etsi_ts/123200_123299/123228/07.02.00_60/ts_123228v070200p.pdf) | **Mj interface**: Choose an MGCF and forward SIP signalling (INVITE etc.) toward that MGCF. |
| Proxy/relay subsequent SIP requests & responses when BGCF remains on the signalling path (dialog continuity) | Secondary (but often operationally required) | [TS 23.228](https://www.etsi.org/deliver/etsi_ts/123200_123299/123228/07.02.00_60/ts_123228v070200p.pdf) | If BGCF stays in the Route set, it forwards mid-dialog signalling between S-CSCF ↔ MGCF, e.g. hold/resume flows show MGCF → BGCF → S-CSCF and back (BGCF forwards requests and 200 OKs). |
| Optional inter-domain "exit" steering via IBCF (when sending toward another domain) | Secondary | [TS 23.228](https://www.etsi.org/deliver/etsi_ts/123200_123299/123228/07.02.00_60/ts_123228v070200p.pdf) | Standard notes that inter-domain requests may be routed via a local network exit point (IBCF). That's typically **Mx interface** in the reference model. |
| Generate charging records (CDRs) for BGCF-handled breakout sessions | Primary | [TS 29.163](https://www.etsi.org/deliver/etsi_ts/129100_129199/129163/16.04.00_60/ts_129163v160400p.pdf) | BGCF explicitly includes "Generation of CDRs" as a performed function. (Implementation-wise: tie CDR correlation to dialog identifiers / routing decisions.) |
| Use administrative / other-protocol-derived information to make routing choice (policy + data inputs) | Secondary | [TS 23.228](https://www.etsi.org/deliver/etsi_ts/123200_123299/123228/07.02.00_60/ts_123228v070200p.pdf) | The standard allows BGCF to use administrative information and/or info from other protocols when selecting breakout network/MGCF. Think: routing tables, number ranges, carrier policy, etc. |
| Support both "local breakout" and "remote breakout" topologies (including cases where BGCF is not on-path after initial routing) | Secondary | [TS 23.228](https://www.etsi.org/deliver/etsi_ts/123200_123299/123228/07.02.00_60/ts_123228v070200p.pdf) | The procedures explicitly allow: (a) select MGCF in same network, (b) forward to another BGCF; and some flows note BGCF might or might not remain in the signalling path after first INVITE routing. |

## Interface Reference

- **Mi**: Interface between S-CSCF and BGCF
- **Mj**: Interface between BGCF and MGCF
- **Mk**: Interface between BGCFs (inter-network)
- **Mx**: Interface between BGCF and IBCF (inter-domain)

## Notes

- **Primary** features are core responsibilities that must be implemented for BGCF to function correctly.
- **Secondary** features are optional or supporting capabilities that enhance functionality but may not be strictly required for basic operation.
- Some secondary features are marked as "often operationally required" indicating they are typically needed in production deployments even if not strictly mandatory by the standard.

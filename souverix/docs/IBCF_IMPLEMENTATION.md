# IBCF Implementation Guide

## Overview

Our IBCF implementation provides 3GPP TS 23.228 compliant Interconnection Border Control Function functionality, integrated with our SBC for carrier-grade border control.

## Architecture

```
┌─────────────────────────────────────────┐
│           IBCF Component                │
│                                         │
│  ┌───────────────────────────────────┐ │
│  │  Message Validation               │ │
│  │  - SIP method validation          │ │
│  │  - Header structure validation    │ │
│  │  - Malformed message blocking     │ │
│  └───────────────────────────────────┘ │
│                                         │
│  ┌───────────────────────────────────┐ │
│  │  Policy Engine                    │ │
│  │  - Peer whitelisting              │ │
│  │  - Call acceptance policies       │ │
│  │  - Attestation requirements       │ │
│  └───────────────────────────────────┘ │
│                                         │
│  ┌───────────────────────────────────┐ │
│  │  Topology Hiding                  │ │
│  │  - Record-Route removal           │ │
│  │  - Via header rewriting           │ │
│  │  - Contact URI replacement        │ │
│  │  - Internal IP removal            │ │
│  └───────────────────────────────────┘ │
│                                         │
│  ┌───────────────────────────────────┐ │
│  │  STIR/SHAKEN Integration          │ │
│  │  - Identity header signing        │ │
│  │  - Identity header verification   │ │
│  │  - Attestation level enforcement  │ │
│  └───────────────────────────────────┘ │
│                                         │
│  ┌───────────────────────────────────┐ │
│  │  Header Normalization              │ │
│  │  - URI normalization              │ │
│  │  - Header name capitalization     │ │
│  └───────────────────────────────────┘ │
└─────────────────────────────────────────┘
```

## Core Functions

### 1. SIP Signaling Control

The IBCF validates and enforces SIP message structure:

- **Method Validation**: Only allows standard SIP methods (INVITE, ACK, BYE, etc.)
- **Header Validation**: Ensures required headers are present
- **Structure Validation**: Blocks malformed messages

### 2. Topology Hiding

Implements 3GPP-required topology hiding:

- **Record-Route Removal**: Strips internal routing information
- **Via Header Rewriting**: Replaces internal domains with border gateway
- **Contact URI Replacement**: Hides internal addresses
- **Internal IP Removal**: Removes private IP ranges

### 3. Security Enforcement

- **TLS Enforcement**: Requires SIP over TLS for interconnects
- **DoS Protection**: Rate limiting and transaction limits
- **Policy Control**: Peer whitelisting and call acceptance policies
- **SIP Method Filtering**: Restricts allowed methods

### 4. Inter-Operator Peering Control

- **Peer Whitelisting**: Only allow known peer domains
- **Call Acceptance Policies**: Enforce peering contracts
- **Attestation Requirements**: Require minimum STIR/SHAKEN attestation levels

## Configuration

### Environment Variables

```bash
# Enable IBCF
ENABLE_IBCF=true

# Topology hiding
IBCF_TOPOLOGY_HIDING=true

# Security
IBCF_REQUIRE_TLS=true
IBCF_DOS_PROTECTION=true

# Peering policy
IBCF_ALLOWED_PEERS=peer1.com,peer2.com,peer3.com
IBCF_REQUIRE_STIR=true
IBCF_MIN_ATTESTATION=A

# STIR/SHAKEN
SBC_ENABLE_STIR=true
SBC_STIR_ATTESTATION=auto
```

### Kubernetes ConfigMap

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: ims-core-config
data:
  enable-ibcf: "true"
  ibcf-topology-hiding: "true"
  ibcf-require-tls: "true"
  ibcf-allowed-peers: "peer1.com,peer2.com"
  ibcf-require-stir: "true"
  ibcf-min-attestation: "A"
```

## Integration Points

### With S-CSCF

IBCF receives messages from S-CSCF via Mw reference point:

```
S-CSCF → IBCF → External Network
```

### With BGCF

For PSTN breakout:

```
S-CSCF → BGCF → IBCF → MGCF → PSTN
```

### With STIR/SHAKEN

IBCF integrates STIR/SHAKEN for:
- Signing outbound calls
- Verifying inbound calls
- Enforcing attestation requirements

## Call Flow Example

### Inter-IMS Call

```
IMS-A (Originating)          IMS-B (Terminating)

UE-A
  |
P-CSCF
  |
S-CSCF
  |
IBCF-A
  ├─ Validates message
  ├─ Applies policy
  ├─ Hides topology
  ├─ Signs STIR/SHAKEN
  └─ Normalizes headers
  |
========= Interconnect =========
  |
IBCF-B
  ├─ Validates peer
  ├─ Verifies STIR/SHAKEN
  ├─ Applies inbound policy
  └─ Restores routing
  |
S-CSCF
  |
P-CSCF
  |
UE-B
```

## Policy Engine

The IBCF includes a policy engine for:

- **Peer Validation**: Check if peer domain is allowed
- **Call Acceptance**: Apply peering contract rules
- **Attestation Enforcement**: Require minimum STIR/SHAKEN attestation

### Policy Rules

```go
// Example policy configuration
allowedPeers := []string{
    "carrier-a.com",
    "carrier-b.com",
}

requireSTIR := true
minAttestation := stir.AttestationFull // A level
```

## High Availability

IBCF supports:

- **Active/Active Clusters**: Multiple IBCF instances
- **Stateful Failover**: Transaction replication
- **DNS-Based Failover**: Automatic peer redirection
- **Geo-Distributed**: Multi-region deployment

## Monitoring

Key metrics:

- `ims_ibcf_messages_total`: Total messages processed
- `ims_ibcf_messages_rejected`: Messages rejected by policy
- `ims_ibcf_stir_verified`: STIR/SHAKEN verifications
- `ims_ibcf_peers_active`: Active peer connections

## Troubleshooting

### Common Issues

1. **Messages Rejected**
   - Check peer whitelist configuration
   - Verify STIR/SHAKEN requirements
   - Review policy rules

2. **Topology Not Hidden**
   - Verify topology hiding enabled
   - Check internal domain configuration
   - Review header rewriting logic

3. **STIR Verification Failures**
   - Check certificate accessibility
   - Verify attestation level requirements
   - Review token expiration

## Standards Compliance

- ✅ **3GPP TS 23.228**: IMS Stage 2 Architecture
- ✅ **3GPP TS 29.165**: Inter-IMS Network to Network Interface
- ✅ **3GPP TS 29.162**: Interworking procedures
- ✅ **3GPP TS 33.203**: IMS Security

## Future Enhancements

1. **Advanced Policy Engine**: Rule-based policy configuration
2. **Lawful Intercept**: LI trigger support
3. **Emergency Routing**: Enhanced 911/E112 support
4. **Media Anchoring**: TrGW integration
5. **QUIC Transport**: Experimental QUIC-based SIP

# Souverix Rempart
SIG-GW / IBCF - Fortified Border Control

## Overview

Souverix Rempart is the carrier and military-grade Session Border Controller (SBC) and Interconnection Border Control Function (IBCF). It serves as the fortified boundary layer of the Souverix platform.

**Rempart** = fortified wall (French)

## Role

Rempart represents the fortified sovereign wall that protects the IMS core from external threats while enabling secure interconnect with other networks.

## Functions

### Border Control
- NNI (Network-to-Network Interface) border control
- IBCF behavior per 3GPP TS 23.228
- SIP normalization
- Topology hiding
- Security enforcement

### STIR/SHAKEN
- STIR/SHAKEN signing (outbound)
- STIR/SHAKEN verification (inbound)
- Certificate management integration
- Attestation level enforcement

### Security
- DoS protection
- Rate limiting
- SIP fuzzing protection
- TLS/mTLS enforcement
- Header size limits
- Method allowlisting

### Peering Policy
- Peer allowlisting
- Routing policy enforcement
- Codec policy
- SRTP enforcement
- Emergency routing priority

## Standards Compliance

- **3GPP TS 23.228** - IMS Architecture (IBCF)
- **RFC 3261** - SIP: Session Initiation Protocol
- **RFC 8224** - SIP Identity Header (STIR)
- **RFC 8225** - PASSporT Token

## Integration

Rempart integrates with:
- **Souverix Coeur** - IMS core signaling
- **Souverix Autorite** - Certificate management
- **Souverix Vigie** - AI intelligence
- **Souverix Mandat** - Lawful intercept
- **Souverix Priorite** - Emergency services

## Configuration

See [Configuration Guide](../operations/configuration.md) for Rempart configuration options.

## Testing

See [Test Catalog](../testing/catalog.md) for Rempart test cases (SIG-*, TOP-*, SEC-*, TLS-*, etc.).

---

## End of Rempart Documentation

# IMS over 5GC Implementation Profile

**VoNR (Voice over New Radio)** - IMS over 5G Core

## What Changes vs Generic IMS

- **The "policy control" world modernizes** (compared to classic EPC patterns).
- **Clean separation of IMS core from access/core policy** so you can plug in 5GC policy and charging models.
- PCF (Policy Control Function) replaces PCRF with HTTP/2 interface (vs Diameter).
- NEF (Network Exposure Function) provides new service exposure capabilities.
- 5G QoS model differs from 4G (5QI vs QCI).

## Required Adapters

### Access Adapter Layer

- **5GC Policy Adapter**: Interface with PCF (Policy Control Function)
- **HTTP/2 Interface**: N28/N7 interfaces (vs Diameter Rx)
- **NEF Integration**: Network exposure for service capabilities
- **5G QoS Model**: 5QI (5G QoS Identifier) handling

### Policy Integration

- **PCF Integration**: Modern policy control via HTTP/2
- **Service-Based Architecture**: 5GC SBA integration
- **Charging Integration**: 5GC charging models (CHF)
- **QoS Authorization**: 5QI-based QoS requests

## Test Cases

### Policy Control

- PCF interaction for QoS authorization
- HTTP/2 interface handling (N28/N7)
- Policy update during active session
- 5QI to media mapping

### Service-Based Architecture

- NEF service exposure
- Service discovery and registration
- Inter-NF communication patterns

### QoS and Media

- 5QI establishment for voice
- QoS modification during call
- Media authorization with 5GC
- Charging correlation (CHF)

### Interworking

- 5G to 4G handover (VoNR to VoLTE)
- 5G to Wi-Fi handover
- Inter-system registration

## Deployment Considerations

### Souverix Posture

**Make policy integration a module boundary (adapter layer), not a CSCF rewrite.**

- Keep IMS core unchanged
- Implement 5GC policy adapter
- Support both Diameter (4G) and HTTP/2 (5G) interfaces
- Abstract policy interface from CSCFs

### Key Differentiators

- **Modern Policy**: HTTP/2-based policy control
- **Service Exposure**: NEF integration for advanced services
- **QoS Evolution**: 5QI support alongside legacy QCI
- **Architecture Separation**: Clean IMS/policy boundary

## Related Documentation

- [P-CSCF Features](../components/coeur/pcscf/FEATURES_LIST.md)
- [Access Architecture](../ARCHITECTURE.md)
- [5G Standards](../standards/6g/README.md)

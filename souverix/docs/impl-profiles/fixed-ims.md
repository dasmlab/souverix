# Fixed IMS Implementation Profile

**Fixed IMS** - FTTH/DSL/Enterprise SIP access

## What Changes vs Generic IMS

- **Identity and access are more "network edge + enterprise" oriented**, often **SBC-heavy**.
- **Roaming/mobility is less dominant**, with **security and topology control becoming dominant**.
- Fixed endpoints are typically stationary.
- Enterprise integration requires specific handling.
- SBC (Session Border Controller) is often the first point of contact.

## Required Adapters

### Access Adapter Layer

- **SBC Integration**: Interface with SBC for enterprise/fixed access
- **Enterprise Identity**: Enterprise-specific identity handling
- **Topology Hiding**: Strong topology hiding for enterprise networks
- **Interconnect Patterns**: SBC adjacency and peering

### Policy Integration

- **Fixed Network Policy**: Different policy models for fixed access
- **Enterprise Policies**: Enterprise-specific service policies
- **QoS Models**: Fixed network QoS (not mobile bearer-based)

## Test Cases

### Enterprise Integration

- Enterprise identity handling
- Enterprise-specific routing
- Enterprise policy enforcement
- Enterprise charging models

### SBC Interworking

- SBC as first contact point
- SBC topology hiding
- SBC security enforcement
- SBC-to-IMS routing

### Topology Control

- Internal topology hiding
- External topology exposure
- Route header manipulation
- Via header management

### Security

- Enterprise security policies
- Fixed network security models
- Interconnect security
- Topology-based security

## Deployment Considerations

### Souverix Posture

**Interconnect/SBC adjacency, topology hiding, and enterprise integration.**

- Strong IBCF/SBC integration
- Enterprise identity support
- Topology hiding patterns
- Interconnect-focused architecture

### Key Differentiators

- **SBC-Centric**: SBC as primary access point
- **Topology Control**: Strong topology hiding requirements
- **Enterprise Focus**: Enterprise-specific features and policies
- **Interconnect**: Heavy focus on peering and interconnect

## Related Documentation

- [IBCF Documentation](../IBCF_DEEP_DIVE.md)
- [Interconnect Certification](../INTERCONNECT_CERTIFICATION.md)
- [Access Architecture](../ARCHITECTURE.md)

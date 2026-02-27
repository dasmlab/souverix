# Souverix Architecture

## Access-Agnostic Architecture

Souverix is designed with a **winning abstraction** that keeps the IMS core access-agnostic while supporting multiple access technologies through adapter layers.

### Core Architecture

Build Souverix as:

1. **IMS Control Plane Core**
   - P/I/S-CSCF, BGCF, MGCF control, HSS/UDM integration
   - Access-agnostic IMS signaling and service logic
   - Standard 3GPP IMS interfaces and behaviors

2. **Access & Policy Adapter Layer**
   - LTE/EPC, Wi-Fi/ePDG, 5GC, fixed/enterprise, private networks
   - Policy integration (PCRF, PCF, enterprise policies)
   - Access-specific handling without core changes

3. **Interconnect Edge**
   - IBCF/SBC adjacency patterns, topology, compliance
   - Inter-domain and inter-network connectivity
   - Security and policy enforcement at network boundaries

**That keeps you "access-agnostic" and makes "new RAN types" mostly an adapter problem, not a core rewrite.**

## Access Attachment Contract

All access implementations must conform to the **Access Attachment Contract** which defines standardized expectations for:

### Registration Lifecycle Expectations

- **Timers**: Registration refresh intervals, timeout values
- **Keepalives**: NAT keepalive mechanisms, frequency
- **NAT**: NAT traversal requirements, binding maintenance
- **Re-registration**: Re-registration triggers and patterns

### Security Association Expectations

- **IPsec/TLS**: Security protocol requirements
- **Association Lifecycle**: Establishment, maintenance, rekeying
- **Security Policies**: Trusted vs untrusted access handling
- **Certificate Management**: Certificate validation and renewal

### Policy Triggers

- **Generic "Policy Events" Interface**: Standardized policy event model
- **QoS Authorization**: Media authorization requests
- **Charging Triggers**: Charging event generation
- **Service Policies**: Service-specific policy enforcement

## Implementation Profiles

See [Implementation Profiles](./impl-profiles/README.md) for detailed profiles:

- [VoLTE](./impl-profiles/volte.md) - IMS over LTE/EPC
- [VoWiFi](./impl-profiles/vowifi.md) - EPC/ePDG flavor
- [IMS over 5GC](./impl-profiles/ims-over-5gc.md) - VoNR
- [Fixed IMS](./impl-profiles/fixed-ims.md) - FTTH/DSL/Enterprise

## Future Implementations

### Private LTE / Private 5G (Industrial/Campus)

- Same IMS patterns, but different trust boundaries + orchestration
- Private network identity and policy models
- Campus/enterprise-specific requirements

### Non-Terrestrial Networks (NTN)

- Latency/jitter constraints push survivability and timing behaviors
- Intermittent connectivity handling
- Extended delay tolerance

### Mission-Critical Services (Public Safety / High Assurance)

- Availability, deterministic behavior, auditability
- Zero-trust boundaries
- Enhanced resilience and redundancy

## Defense / High-Assurance Posture

In defense and high-assurance environments, the differentiators are usually:

- **Zero-trust boundaries**: Explicit trust zones, stricter identity, strong crypto posture
- **Auditable control plane**: Traceability, deterministic behavior, provenance of config
- **Resilience under constrained links**: NTN-ish behavior, intermittent connectivity
- **Interop conventions**: More gateways, more strict peering rules

**So: invest early in strong interface contracts + policy adapters + observability + config provenance.**

## Related Documentation

- [Component Flows](./flows/README.md)
- [Component Features](../components/coeur/)
- [6G Standards Watchlist](./standards/6g/README.md)

# VoWiFi Implementation Profile

**VoWiFi (Voice over Wi-Fi)** - EPC/ePDG flavor

## What Changes vs Generic IMS

- **UE reaches IMS through trusted/untrusted Wi-Fi paths** via ePDG (evolved Packet Data Gateway).
- **Session stability is dominated by NAT, tunnel keepalives, and security association behavior**.
- Wi-Fi network quality and availability varies significantly.
- NAT traversal is critical for maintaining connectivity.
- Security association (IPsec) must be robust for untrusted Wi-Fi.

## Required Adapters

### Access Adapter Layer

- **ePDG Integration**: Interface with ePDG for Wi-Fi access
- **NAT Traversal**: Strong NAT handling and keepalive mechanisms
- **IPsec Management**: Security association lifecycle management
- **Tunnel Keepalive**: Maintain IPsec tunnel connectivity

### Policy Integration

- **Wi-Fi Policy**: Different policy rules for trusted vs untrusted Wi-Fi
- **QoS Handling**: Wi-Fi QoS mapping to IMS media
- **Charging**: VoWiFi-specific charging models

## Test Cases

### NAT Traversal

- NAT binding maintenance via keepalives
- NAT rebinding after timeout
- Multiple NAT layers (carrier-grade NAT)
- STUN binding requests
- SIP OPTIONS for NAT keepalive

### Security Association

- IPsec tunnel establishment
- Tunnel rekeying during active session
- Tunnel failure and recovery
- Security association timeout handling

### Session Stability

- Call setup over Wi-Fi
- Mid-call Wi-Fi disconnection and recovery
- Wi-Fi to LTE handover (if supported)
- Wi-Fi quality degradation handling

### Keepalive Mechanisms

- CRLF keepalive frequency
- SIP OPTIONS keepalive
- IPsec tunnel keepalive
- Combined keepalive strategies

## Deployment Considerations

### Souverix Posture

**Strong edge survivability patterns**: keepalive handling, NAT resiliency, fast re-reg patterns.

- Implement robust NAT traversal at P-CSCF
- Support multiple keepalive mechanisms
- Fast re-registration on connectivity loss
- Graceful degradation under poor Wi-Fi conditions

### Key Differentiators

- **Edge Survivability**: P-CSCF must handle intermittent connectivity
- **NAT Resilience**: Multiple keepalive strategies
- **Security First**: Strong IPsec enforcement for untrusted Wi-Fi
- **Fast Recovery**: Quick re-registration and session recovery

## Related Documentation

- [P-CSCF Flows](../flows/pcscf-flows.md#nat-traversal-support)
- [P-CSCF Features](../components/coeur/pcscf/FEATURES_LIST.md)
- [Access Architecture](../ARCHITECTURE.md)

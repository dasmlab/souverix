# Implementation Profiles

This directory contains implementation-specific profiles for different access technologies and deployment scenarios supported by Souverix.

## Supported Implementations

- [VoLTE (IMS over LTE/EPC)](./volte.md) - Voice over LTE
- [VoWiFi (EPC/ePDG flavor)](./vowifi.md) - Voice over Wi-Fi
- [IMS over 5GC (VoNR)](./ims-over-5gc.md) - Voice over New Radio (5G)
- [Fixed IMS](./fixed-ims.md) - Fixed network IMS (FTTH/DSL/Enterprise)

## Future Implementations

- Private LTE / Private 5G (industrial/campus)
- Non-Terrestrial Networks (NTN)
- Mission-critical services (public safety / high assurance)

## Profile Structure

Each implementation profile includes:

1. **What Changes** - Differences from generic IMS
2. **Required Adapters** - Access and policy adapter requirements
3. **Test Cases** - Implementation-specific test scenarios
4. **Deployment Considerations** - Operational and architectural notes

## Access Attachment Contract

All implementations must conform to the **Access Attachment Contract** which defines:

- Registration lifecycle expectations (timers, keepalives, NAT)
- Security association expectations (IPsec/TLS)
- Policy triggers (generic "policy events" interface)

See [Access Architecture](../ARCHITECTURE.md#access-attachment-contract) for details.

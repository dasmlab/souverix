# SIP Gateway / IBCF / SBC Deep Dive

## Layered Architecture

```
+------------------------------------------------------------+
| Applications Layer                                         |
| MMTel AS | RCS AS | Voicemail | TAS | Emergency AS       |
+------------------------------------------------------------+
| Control Layer                                              |
| S-CSCF | I-CSCF | P-CSCF | BGCF | HSS/UDM | PCRF/PCF     |
+------------------------------------------------------------+
| Interconnect / Border Layer                                 |
| IBCF / SBC (SIP GW) / Security Edge                       |
+------------------------------------------------------------+
| Media & PSTN Interworking Layer                           |
| MGCF | MGW (RTP <-> TDM) | Trunks                         |
+------------------------------------------------------------+
| Access / Transport Layer                                   |
| LTE / 5G / WiFi / Fixed Broadband / Enterprise SIP        |
+------------------------------------------------------------+
```

## Functional Placement of SIP Gateway

Our SIP Gateway can operate in multiple roles:

### 1. IBCF (Interconnection Border Control Function)
- **3GPP Standardized** (TS 23.228)
- Inter-domain routing
- Topology hiding
- Security enforcement
- SIP normalization

### 2. SBC (Session Border Controller)
- **Carrier-Grade** implementation
- Enterprise SIP trunking
- PBX to IMS interworking
- Fixed Broadband voice to IMS
- Advanced security features

### 3. Enterprise SIP Gateway
- Enterprise trunk termination
- PBX integration
- Codec normalization
- NAT traversal

### 4. PSTN Gateway (Future)
- ISUP support
- TDM interworking
- Legacy switching

## IBCF vs SBC vs Generic SIP Proxy

| Feature | IBCF | SBC | Generic SIP Proxy |
|---------|------|-----|-------------------|
| Defined by 3GPP | ✅ Yes | ❌ No | ❌ No |
| Topology Hiding | ✅ Yes | ✅ Yes | ⚠️ Limited |
| NAT Traversal | ✅ Yes | ✅ Yes | ❌ No |
| DoS Protection | ⚠️ Limited | ✅ Advanced | ❌ No |
| SIP Normalization | ✅ Yes | ✅ Yes | ⚠️ Minimal |
| Media Anchoring | ⚠️ Optional | ✅ Yes | ❌ No |
| TLS/SRTP Termination | ✅ Yes | ✅ Yes | ⚠️ Rare |
| Carrier Peering | ✅ Yes | ✅ Yes | ❌ No |
| Enterprise Edge | ⚠️ Rare | ✅ Yes | ❌ No |
| STIR/SHAKEN | ✅ Yes | ✅ Yes | ❌ No |

**Key Insight**: In most real deployments, IBCF functionality is embedded inside SBC platforms. Our implementation provides both.

## Call Flow Diagrams

### Enterprise → Mobile Call (IMS Termination)

```
Enterprise PBX
    |
    | INVITE
    v
SIP GW / SBC
    |
    | INVITE (normalized, STIR signed)
    v
IBCF
    |
    | INVITE
    v
S-CSCF
    |
    | Query subscriber (Diameter)
    v
HSS/UDM
    |
    v
P-CSCF
    |
    v
UE (Mobile Device)
```

**Media Path**: RTP flows Enterprise ↔ SBC ↔ UE (or via dedicated media anchoring node)

### IMS → PSTN Breakout

```
UE
    |
    | INVITE
    v
P-CSCF
    |
    v
S-CSCF
    |
    v
BGCF (routing decision)
    |
    v
MGCF (SIP → ISUP)
    |
    v
MGW (RTP → TDM)
    |
    v
PSTN
```

## SIP Gateway in Private 5G Networks

### Why Private 5G Needs IMS

Even in private 5G SA (Standalone):
- Voice over NR (VoNR) still requires IMS
- Mission-critical voice relies on SIP signaling
- Interconnect to public networks requires SIP GW

### Deployment Model

```
Private 5G Core (AMF/SMF/UPF)
    |
    v
IMS Core (Virtualized / CNF)
    |
    v
Edge SBC / SIP GW
    |
    v
Enterprise SIP / PSTN / Public MNO
```

### Common Use Cases

- **Campus Telephony**: Internal voice services
- **Industrial Push-to-Talk**: Mission-critical communications
- **Secure Enterprise Mobility**: Encrypted voice
- **Multi-Site Private Interconnect**: Site-to-site calling

## Security Model (2026)

Modern SIP Gateway includes:

- ✅ **Mutual TLS** between carriers
- ✅ **STIR/SHAKEN** identity verification (with ACME)
- ✅ **SIP header normalization**
- ✅ **Fraud detection engines** (AI agent hooks)
- ✅ **Topology hiding**
- ✅ **SIP rate limiting**
- ✅ **TLS 1.3 enforcement**
- ✅ **DoS protection**

Carrier networks treat SIP GW as:

> **"Telecom Firewall + Protocol Translator"**

## Cloud-Native Deployment

### RHOSO (Red Hat OpenStack) Model

```
Compute Nodes (Nova)
    |
    v
IMS VNFs (S-CSCF, HSS, MGCF)
    |
    v
SR-IOV Networking
    |
    v
SIP GW VNF
```

**Benefits**:
- SR-IOV for low latency
- DPDK acceleration
- NFV MANO orchestration

### OpenShift CNF Model

```
OpenShift Cluster
    |
    +-- IMS Core Pods
    |   - S-CSCF
    |   - I-CSCF
    |   - P-CSCF
    |
    +-- SIP GW / SBC Pods
    |
    +-- Diameter Services
    |
    +-- Media Anchors
```

**Characteristics**:
- Horizontal scaling
- Kubernetes-native lifecycle
- Rolling upgrades
- Geo-redundant deployments
- GitOps-based config management

## Operational Considerations

### Scaling Strategy

- **Stateless SIP proxies** scale horizontally
- **Media anchors** scale independently
- **Separate signaling and media planes**
- **Kubernetes HPA** for signaling nodes

### High Availability

- **Active/Active SBC clusters**
- **Geo-redundant IBCF**
- **DNS-based failover**
- **Diameter redundancy**

## Standards Compliance

- **3GPP TS 23.228**: IP Multimedia Subsystem (IMS)
- **3GPP TS 29.228**: Cx and Dx interfaces
- **RFC 3261**: SIP: Session Initiation Protocol
- **RFC 8224**: SIP Identity Header (STIR)
- **RFC 8225**: PASSporT Token
- **RFC 8588**: Certificate Management
- **RFC 8555**: ACME Protocol
- **GSMA IR.92 / IR.94**: VoLTE/VoWiFi profiles

## Strategic Reality (2026)

### IMS Status

- ✅ **Stable and mandatory** for mobile voice
- ✅ **Not being replaced** in 5G
- ✅ **Becoming cloud-native**
- ✅ **Integrated into private 5G**

### SIP GW Layer Status

- 💰 **Where interconnect revenue exists**
- 🔒 **Where security threats concentrate**
- ☁️ **Where cloud-native transformation is happening**
- 🚀 **Where innovation is focused**

## Future Enhancements

1. **STIR/SHAKEN Full Implementation**: Complete RFC 8224/8225 compliance
2. **Emergency Call Routing**: Enhanced 911/E112 support
3. **SIP Normalization Examples**: Detailed header transformation
4. **Diameter vs HTTP2**: 5G SBA comparison
5. **VoNR Call Setup**: Detailed 5G voice flows
6. **Fraud Analytics**: AI-powered fraud detection

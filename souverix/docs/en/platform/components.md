# Souverix Platform Components

## Component Overview

The Souverix Platform consists of 10 core components, each serving a specific role in the sovereign signaling architecture.

## Component List

### 🧠 Souverix Coeur
**IMS Core Signaling Brain**

Core Functions:
- P-CSCF (Proxy CSCF)
- I-CSCF (Interrogating CSCF)
- S-CSCF (Serving CSCF)
- BGCF (Breakout Gateway Control Function)
- MGCF (Media Gateway Control Function)
- Policy Enforcement Hooks
- Session State Control

[Full Documentation](../components/coeur.md)

### 🛡 Souverix Rempart
**Carrier / Military-Grade SIG-GW**

The fortified boundary layer of Souverix.

Functions:
- NNI Border Control (IBCF)
- SIP normalization
- Topology hiding
- DoS mitigation
- STIR/SHAKEN enforcement
- Peering policy control

[Full Documentation](../components/rempart.md)

### 🎛 Souverix Relais
**Media Relay Engine**

Functions:
- RTP proxy
- SRTP enforcement
- NAT traversal
- QoS tagging
- Media statistics

[Full Documentation](../components/relais.md)

### 🔐 Souverix Autorite
**Cryptographic Sovereign Authority**

Functions:
- Internal CA management
- Certificate lifecycle automation
- HSM integration
- STIR signing key control
- mTLS enforcement
- Automated rotation

[Full Documentation](../components/autorite.md)

### 👁 Souverix Vigie
**AI Intelligence Layer**

Functions:
- Fraud detection
- Anomaly detection
- Dynamic routing optimization
- STIR analytics
- Traffic pattern classification
- Attack detection
- Self-healing triggers

[Full Documentation](../components/vigie.md)

### 🎯 Souverix Mandat
**Lawful Intercept Orchestration**

Functions:
- Signaling duplication
- Media duplication
- Mediation device integration
- Intercept lifecycle management
- Audit compliance tracking

[Full Documentation](../components/mandat.md)

### 🚨 Souverix Priorite
**Emergency & National Priority Services**

Functions:
- Emergency routing
- STIR override
- Fraud override
- PSAP integration
- Priority queuing
- Disaster mode

[Full Documentation](../components/priorite.md)

### 📊 Souverix Vigile
**Observability & Compliance Telemetry**

Functions:
- Structured logging
- Metrics export
- Call tracing
- Audit dashboards
- Regulatory reporting
- SLO monitoring

[Full Documentation](../components/vigile.md)

### 🌐 Souverix Federation
**Inter-Domain Sovereign Peering Control**

Functions:
- Cross-border trust mapping
- Attestation domain translation
- Peering policy orchestration
- Multi-tenant domain isolation

[Full Documentation](../components/federation.md)

### ⚙ Souverix Gouverne
**Policy & Configuration Control Plane**

Functions:
- Peer definitions
- Policy management
- Rate limits
- Enforcement toggles
- STIR mode control
- Emergency override configuration

[Full Documentation](../components/gouverne.md)

---

## Component Interactions

```
Souverix Gouverne (Policy Control)
         |
         v
Souverix Coeur (IMS Core)
         |
         +---> Souverix Rempart (Border Control)
         |           |
         |           +---> Souverix Autorite (PKI)
         |           +---> Souverix Vigie (AI Intelligence)
         |           +---> Souverix Mandat (LI)
         |           +---> Souverix Priorite (Emergency)
         |
         +---> Souverix Relais (Media)
         |
         +---> Souverix Vigile (Observability)
         |
         +---> Souverix Federation (Inter-domain)
```

---

## End of Components Overview

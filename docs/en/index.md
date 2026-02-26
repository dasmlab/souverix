# Souverix Documentation
## Sovereign Signaling Doctrine — IMS + SIG-GW

Souverix is a **sovereign, AI-native signaling platform** engineered for **carrier and defense-grade** environments. Built in modern Golang, it is a clean architectural rewrite designed for **high CPS performance**, **secure interconnect control**, and **cryptographic sovereignty**.

Souverix treats signaling as strategic infrastructure:

- **Sovereign by design**: trust roots, policy, and interconnect control are first-class concerns.
- **Resilient by default**: active/active patterns, failure isolation, and emergency continuity.
- **Intelligent at the edge**: AI-assisted detection, classification, and adaptive enforcement.
- **Regulation-ready**: STIR/SHAKEN, lawful intercept, emergency routing, and auditability.
- **Cloud-native CNF**: built to run on Kubernetes/OpenShift with predictable operations.

---

## Platform Components

### 🧠 Souverix Coeur — IMS Core
Cloud-native IMS signaling core (X-CSCF stack) responsible for core session control and policy integration.

### 🛡 Souverix Rempart — SIG-GW / IBCF
Carrier/military-grade border signaling gateway implementing interconnect control, topology hiding, SIP normalization, abuse mitigation, and STIR/SHAKEN enforcement.

### 🎛 Souverix Relais — Media Plane
Media relay and anchoring layer for RTP/SRTP policy enforcement, NAT traversal, QoS handling, and media telemetry.

### 🔐 Souverix Autorite — PKI / HSM / Vault
Sovereign cryptographic authority: CA chain management, certificate lifecycle automation, HSM integration, mTLS enforcement, and key rotation.

### 👁 Souverix Vigie — AI Intelligence
AI-driven intelligence layer for anomaly detection, fraud signals, adaptive policy, attack classification, and self-healing triggers.

### 🎯 Souverix Mandat — Lawful Intercept
Lawful intercept orchestration for signaling/media duplication, mediation integration, and audit-grade compliance tracking.

### 🚨 Souverix Priorite — Emergency & Priority
Emergency and national priority services: PSAP routing, priority queuing, override controls, and continuity under stress.

### 📊 Souverix Vigile — Observability & Audit
Metrics, logs, traces, compliance telemetry, and regulatory-grade audit reporting.

### 🌐 Souverix Federation — Inter-domain Control
Federation layer enabling controlled interoperability across sovereign domains: trust mapping, peering agreements, and multi-tenant interconnect policy.

### ⚙ Souverix Gouverne — Policy Control Plane
Policy and configuration authority: peer profiles, enforcement toggles, rate limits, runtime controls, emergency overrides, and warrant provisioning.

---

## Suggested Reading Path

1. **Platform → Overview / Components / Naming**
2. **Architecture → Layers / OpenShift CNF**
3. **Compliance → STIR/SHAKEN / Lawful Intercept / Emergency**

---

## Quick Links

- [Platform Overview](platform/doctrine.md)
- [Component Breakdown](platform/components.md)
- [Naming & Namespace](platform/naming.md)
- [Getting Started](operations/getting-started.md)
- [Architecture Hierarchy](architecture/hierarchy.md)

---

## Doctrine Statement

Souverix does not merely implement IMS.

It establishes a modern doctrine of **sovereign signaling** — where **interconnect**, **trust**, **intelligence**, and **resilience** are first-class architectural concerns.

**Designed in Canada. Built for sovereign control.**

---

## Standards Compliance

- **3GPP**: TS 23.228, TS 24.229, TS 29.228, TS 33.107, TS 23.167
- **IETF**: RFC 3261, RFC 8224, RFC 8225, RFC 8588, RFC 8555
- **Regulatory**: FCC, CRTC, ETSI

---

## End of Documentation

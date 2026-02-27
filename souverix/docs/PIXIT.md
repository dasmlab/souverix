# Implementation-Level PIXIT Extension
SIG-GW / IBCF — Execution Parameters & Environment Control

**Author**: Daniel  
**Purpose**: Define implementation-specific execution parameters (PIXIT) for all test matrices  
**Scope**: Applies to IBCF, STIR/SHAKEN, LI, Emergency

**PIXIT** = Protocol Implementation eXtra Information for Testing

This document defines:
- Timer settings
- Codec policies
- TLS profiles
- Peer profiles
- STIR enforcement modes
- LI activation states
- Emergency override rules
- Scaling thresholds
- Chaos injection controls

---

## 1. Global Timer Profile (Baseline)

| Parameter | Default | Test Override Range | Notes |
|-----------|---------|---------------------|-------|
| T1 | 500 ms | 100–1000 ms | SIP base retransmission |
| T2 | 4 s | 2–8 s | Non-INVITE retrans max |
| Timer B | 32 s | 16–64 s | INVITE client timeout |
| Timer F | 32 s | 16–64 s | Non-INVITE timeout |
| Session Timer | 1800 s | 90–3600 s | re-INVITE refresh |
| Max Dialogs | 50,000 | 1k–100k | Capacity test |
| Max CPS | 2,000 | 100–20k | Load envelope |

---

## 2. Codec Policy PIXIT

| Parameter | Default | Allowed Values |
|-----------|---------|----------------|
| Audio Codecs | AMR-WB, G.711 | Configurable allowlist |
| Video Codecs | H.264 | Add/remove |
| Transcoding | Disabled | Enabled/Disabled |
| SRTP Required | Yes (NNI) | Yes/No |
| DTMF Mode | RFC2833 | RFC2833 / SIP INFO |

---

## 3. TLS / Security Profile

| Parameter | Default | Options |
|-----------|---------|---------|
| TLS Version | 1.2/1.3 | Configurable |
| Cipher Suites | ECDHE-ECDSA-AES256-GCM | Custom set |
| mTLS Required | Yes (NNI) | Yes/No |
| Cert Reload Mode | Hot | Restart |
| OCSP Mode | Soft-fail | Hard/Soft |
| STIR Enforcement | Soft | Soft/Hard |

---

## 4. Peer Profile Definition Template

Each test must reference a peer profile.

| Field | Example Value |
|-------|---------------|
| Peer ID | PEER-A |
| Transport | SIP-TLS |
| IP Version | IPv4 |
| Auth Mode | mTLS |
| STIR Trust Level | Trusted / External |
| Emergency Routing Flag | Enabled |
| LI Cooperation Flag | Enabled |
| Max CPS | 1000 |
| Topology Hiding Mode | Full |

---

## 5. STIR PIXIT Controls

| Parameter | Default | Test Override |
|-----------|---------|---------------|
| Attestation Policy | A/B/C auto | Force A/B/C |
| iat Skew Tolerance | ±60 sec | ±0–300 sec |
| Cert Cache TTL | 24h | 1h–72h |
| Signing Key Source | HSM | File / Vault |
| Re-sign Transit Calls | No | Yes/No |
| Identity Header Max Size | 8 KB | 2–16 KB |

---

## 6. Lawful Intercept PIXIT Controls

| Parameter | Default | Options |
|-----------|---------|---------|
| LI Mode | Disabled | Signaling / Full |
| Mediation IP | 10.1.1.10 | Configurable |
| Intercept Target List | Dynamic DB | Static |
| LI Overload Policy | Preserve LI | Drop LI |
| Audit Log Retention | 180 days | Custom |
| TLS Decrypt for LI | Enabled | Disabled |

---

## 7. Emergency PIXIT Controls

| Parameter | Default | Options |
|-----------|---------|---------|
| Emergency Numbers | 911,112 | Custom list |
| Emergency Priority Queue | Enabled | Disabled |
| STIR Override | Enabled | Disabled |
| Fraud Override | Enabled | Disabled |
| Emergency Route | PSAP-A | Fallback list |
| Location Header Mandatory | Yes | Optional |

---

## 8. Chaos Injection Controls

| Parameter | Default | Range |
|-----------|---------|-------|
| Packet Loss | 0% | 0–20% |
| Jitter | 0 ms | 0–200 ms |
| DNS Failure | Off | On |
| TLS Cert Expiry Simulation | Off | On |
| Signing Service Crash | Off | On |
| State Store Partition | Off | On |
| Media Relay Kill | Off | On |

---

## 9. Example: PIXIT Overlay for Specific Test

**Example: STR-012 (OCSP outage test)**

| PIXIT Parameter | Value |
|-----------------|-------|
| OCSP Mode | Hard-fail |
| Cert TTL | 24h |
| iat Skew | ±60 sec |
| TLS Version | 1.3 |
| CPS | 2000 |
| Chaos | OCSP responder disabled |

**Expected Result:**
- Calls fail verification
- Hard-fail blocks call
- No CPU spike
- Metrics increment verification_fail counter

---

## 10. Example: Emergency During Failover (LIE-107)

| PIXIT Parameter | Value |
|-----------------|-------|
| Emergency Priority | Enabled |
| STIR Override | Enabled |
| Active/Active Nodes | 2 |
| Chaos | Kill primary instance |
| CPS Background | 5000 |
| Emergency CPS | 200 |

**Expected Result:**
- Emergency call continues
- No PDD violation
- No STIR enforcement block
- No LI drop

---

## 11. Example: IBCF Flood Protection (SEC-002)

| PIXIT Parameter | Value |
|-----------------|-------|
| Max CPS | 2000 |
| Rate Limit | 2500 CPS |
| Peer Max CPS | 1000 |
| Chaos | 20k CPS burst |
| TLS Required | Yes |

**Expected Result:**
- Legit calls pass
- Excess rejected with 503
- No crash
- CPU stable <85%

---

## 12. Execution Discipline

Every test must log:

- SIP ladder trace
- Timer values in effect
- Codec negotiation result
- STIR verification state
- LI state
- Emergency policy state
- Peer profile used
- Chaos injection status

---

## 13. Traceability Matrix (Linking Documents)

| Test Matrix | Requires PIXIT Section |
|-------------|------------------------|
| IBCF Functional | Sections 1–4, 8 |
| STIR/SHAKEN | Sections 3, 5, 8 |
| LI Tests | Sections 6, 8 |
| Emergency | Sections 7, 8 |
| Resilience | Sections 1, 8 |

---

## 14. PIXIT Configuration File Format

Tests can load PIXIT parameters from YAML/JSON:

```yaml
pixit:
  timers:
    t1: 500ms
    t2: 4s
    timer_b: 32s
    timer_f: 32s
    session_timer: 1800s
  
  codec:
    audio: ["AMR-WB", "G.711"]
    video: ["H.264"]
    srtp_required: true
  
  tls:
    version: ["1.2", "1.3"]
    mtls_required: true
    ocsp_mode: "soft-fail"
  
  stir:
    enforcement: "soft"
    iat_skew: 60s
    cert_cache_ttl: 24h
  
  li:
    mode: "disabled"
    mediation_ip: "10.1.1.10"
    overload_policy: "preserve"
  
  emergency:
    numbers: ["911", "112"]
    priority_queue: true
    stir_override: true
  
  chaos:
    packet_loss: 0%
    jitter: 0ms
    dns_failure: false
```

---

## 15. Test Execution with PIXIT

### Command Line

```bash
# Run test with PIXIT override
./testrig -test STR-012 \
  -pixit ocsp_mode=hard-fail \
  -pixit chaos=ocsp_disabled

# Load PIXIT from file
./testrig -test LIE-107 \
  -pixit-file pixit/emergency-failover.yaml
```

### Programmatic

```go
pixit := LoadPIXIT("pixit/str-012.yaml")
test := NewTest("STR-012", pixit)
result := test.Run()
```

---

## End of PIXIT Extension Document

# STIR / SHAKEN Test Catalogue (SIG-GW / IBCF Context) — 2026 Edition

**Author**: Daniel  
**Scope**: STIR/SHAKEN implementation inside SIG-GW / IBCF  
**Standards Anchors**:
- RFC 8224 (SIP Identity Header)
- RFC 8225 (PASSporT)
- RFC 8588 (SHAKEN Certificates)
- ATIS-1000074 (SHAKEN Framework)
- FCC / CRTC regulatory enforcement models

---

## Legend

**TEST ID format**: `STR-###`

**Columns**:
- TEST ID
- SLOGAN
- AREA
- POSITIVE CASE
- NEGATIVE CASE
- LOAD CASE
- SCALE/RAMP CASE
- CHAOS CASE

**AREA codes**:
- `ORG` = Originating Signing
- `TER` = Terminating Verification
- `ATT` = Attestation Policy
- `CRT` = Certificate Handling
- `KEY` = Key Management
- `POL` = Policy Enforcement
- `NET` = Cross-Network / Peering
- `RES` = Resilience & Failure
- `SEC` = Security Hardening
- `OBS` = Observability

---

## STIR/SHAKEN Test Matrix

| TEST ID | SLOGAN | AREA | POSITIVE CASE | NEGATIVE CASE | LOAD CASE | SCALE/RAMP CASE | CHAOS CASE |
|---------|--------|------|---------------|---------------|-----------|-----------------|------------|
| STR-001 | "A-level signing works" | ORG | Valid subscriber call → PASSporT generated, signed, Identity header inserted | Missing TN authorization → no signing or downgrade | 2k CPS signed calls | Ramp 100→2k CPS signing | Kill signing service pod mid-load; fail closed |
| STR-002 | "PASSporT structure valid" | ORG | JWT payload contains orig/dest/iat/attest | Malformed PASSporT rejected internally | 2k CPS | Ramp payload size variations | Inject corrupted JWT; no crash |
| STR-003 | "iat timestamp accuracy" | ORG | iat within allowed skew window | iat outside skew → reject or mark invalid | 1k CPS | Ramp skew tolerance 0→5m | Jump system clock +10m |
| STR-004 | "B-level attestation" | ATT | Enterprise trunk call → B attestation | Unauthorized enterprise TN → downgrade to C | 1k CPS enterprise | Ramp enterprise % | Remove TN ownership DB mid-run |
| STR-005 | "C-level gateway marking" | ATT | Foreign inbound call → C attestation | Attempt to force A-level from gateway blocked | 2k CPS inbound | Ramp % gateway traffic | Peer spoofs A-level claim; verification catches |
| STR-006 | "Identity header insertion" | ORG | Identity header present in outbound INVITE | Duplicate Identity header prevented | 2k CPS | Ramp header size | Header buffer overflow attempt |
| STR-007 | "Signature cryptographic validity" | ORG | ECDSA P-256 signature validates | Tampered payload fails verification | 1k CPS | Ramp verification concurrency | Corrupt signing key file |
| STR-008 | "Verification success path" | TER | Valid signature → verification pass, call flagged verified | Invalid signature → verification fail | 2k CPS verify | Ramp 100→2k verify/s | Kill cert cache service |
| STR-009 | "Identity header missing" | TER | Missing Identity header → policy: mark unverified | Require verification policy blocks call | 2k CPS no-identity | Ramp enforcement strictness | Toggle enforcement mode live |
| STR-010 | "Expired certificate" | CRT | Expired cert rejected | Soft-fail policy logs + mark unverified | 1k CPS expired | Ramp expired cert ratio | Force cert expiry at runtime |
| STR-011 | "Certificate chain validation" | CRT | Full chain validates to trusted root | Unknown root CA rejected | 500 CPS | Ramp chain depth | Replace root CA mid-traffic |
| STR-012 | "OCSP validation" | CRT | Valid OCSP → pass | OCSP revoked → fail call | 500 CPS OCSP checks | Ramp OCSP latency | OCSP responder down (hard-fail vs soft-fail) |
| STR-013 | "CRL fallback" | CRT | CRL used if OCSP unavailable | Corrupt CRL → fail safe | 500 CPS | Ramp CRL size | Corrupt CRL file injected |
| STR-014 | "Key rotation (signing)" | KEY | New private key loaded seamlessly | Invalid key rejected; old retained | 1k CPS during rotation | Ramp rotation frequency | Rotate every 60s under load |
| STR-015 | "Key compromise scenario" | KEY | Compromised key revoked & blocked | Old key cannot sign | N/A | Ramp revocation checks | Inject compromised key event mid-traffic |
| STR-016 | "Cross-border attestation downgrade" | NET | Foreign call retains C-level | Improper A-level from foreign rejected | 2k CPS cross-border | Ramp foreign % | Peer mislabels attestation |
| STR-017 | "Re-signing rules" | NET | Transit through trusted domain preserves Identity | Re-sign only if policy allows | 1k CPS transit | Ramp re-sign ratio | Remove trust relationship mid-call |
| STR-018 | "Replay attack defense" | SEC | Replay detection blocks reused PASSporT | Valid retransmission allowed | 5k replays | Ramp replay rate | Massive replay storm |
| STR-019 | "JWT tampering attempt" | SEC | Modified payload → signature mismatch | Call not marked verified | 2k CPS tampered | Ramp tamper % | Fuzz Identity header |
| STR-020 | "Attestation downgrade logic" | POL | If verification fails → downgrade to unverified | Never upgrade attestation | 2k CPS mixed | Ramp invalid ratio | Remove TN ownership DB |
| STR-021 | "Verification performance" | TER | ≤5ms verification latency target | Latency spikes trigger alert | 5k verify/s | Ramp verify concurrency | CPU throttling injected |
| STR-022 | "Burst signing resilience" | ORG | Signing stable under burst | No signature corruption | 10k CPS burst | Ramp burst window | Signing HSM slowdown |
| STR-023 | "Multi-tenant cert selection" | CRT | Correct cert chosen per TN range | Wrong cert not used | 1k CPS mixed TN | Ramp tenant count | Remove one tenant cert mid-run |
| STR-024 | "Header propagation integrity" | NET | Identity header preserved across routing | Header stripped unexpectedly flagged | 2k CPS | Ramp multi-hop | Strip header in one hop |
| STR-025 | "Soft-fail vs hard-fail policy" | POL | Configurable enforcement works | Hard-fail blocks invalid | 2k CPS | Ramp enforcement mode | Toggle mode at peak |
| STR-026 | "Observability: signature metrics" | OBS | Metrics show sign/verify counts | Metric drop triggers alert | 5k CPS | Ramp metrics volume | Metrics backend down |
| STR-027 | "Identity size limits" | SEC | Oversized Identity rejected safely | No crash | 5k CPS oversized | Ramp header sizes | Buffer overflow fuzz |
| STR-028 | "iat skew tolerance" | TER | Accept within ±X seconds | Outside window rejected | 2k CPS varied skew | Ramp skew | Time jump event |
| STR-029 | "SIP fragmentation handling" | NET | Fragmented SIP with Identity handled correctly | Truncated header rejected | 2k CPS | Ramp MTU size | Inject fragmentation storm |
| STR-030 | "Cross-cluster verification" | RES | Multiple IBCF instances share cert cache | Cache miss handled | 2k CPS | Ramp instances 1→6 | Kill cert cache pod cluster-wide |

---

## Required Explicit Examples

**Positive Example**:  
STR-001 — Valid subscriber call signed at A-level.

**Negative Example**:  
STR-010 — Expired certificate rejected.

**Load Example**:  
STR-022 — 10k CPS signing burst.

**Scale/Ramp Example**:  
STR-021 — Ramp verification from 100→5k/s.

**Chaos Example**:  
STR-012 — OCSP responder outage during peak traffic.

---

## Additional Advanced (Optional) Scenarios

1. Blockchain-backed certificate discovery (future models)
2. Quantum-resistant signature testing (forward-looking)
3. SIP over QUIC STIR header preservation
4. Multi-hop re-signing logic
5. Private 5G STIR domain separation

---

## Test Implementation Notes

### Test Harness Requirements

- SIP traffic generator (SIPp, custom)
- Certificate management simulator
- OCSP/CRL responder mock
- Time manipulation for skew tests
- Network partition simulation

### Assertions

- PASSporT token structure validation
- Signature cryptographic verification
- Attestation level correctness
- Certificate chain validation
- Performance latency measurements
- Security attack detection

### Metrics to Track

- Signing latency (p50, p95, p99)
- Verification latency
- Certificate fetch latency
- Attestation distribution (A/B/C)
- Verification success/failure rates
- Replay detection accuracy

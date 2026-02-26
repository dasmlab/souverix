# IBCF / SIG-GW Test Catalogue (IMS NNI / Ic Interconnect) — 2026-oriented

**Author**: Daniel  
**Scope**: A **SIG-GW app implementing IBCF** behaviors for IMS interconnect (NNI/Ic) per 3GPP IMS architecture, with test planning aligned to publicly available IMS NNI testing and IMS testing principles (ETSI INT series, IBCF PICS for Ic, and IMS testing guidance).

> **Notes on sources:**
> - ETSI "IMS Network Testing (INT) / IMS NNI Interoperability Test Specifications" provides formal **test purposes & test descriptions** for IMS NNI interoperability.
> - ETSI IBCF PICS work explicitly targets **IBCF requirements** on **Ic** and calls out **topology hiding** expectations.
> - ITU-T Q.3904 frames **IMS testing principles** across conformance/interoperability/functionality.
> - 3GPP TS 23.228 defines IMS architecture and IBCF's place/functionality.

---

## Legend

**Columns:**
- **TEST ID**: `AAA-###` where `AAA` = area code (3 letters).
- **SLOGAN**: quick statement of intent.
- **AREA**: functional domain.
- **POSITIVE CASE**: "should work" behavior.
- **NEGATIVE CASE**: "should fail safely" behavior.
- **LOAD CASE**: stress on throughput / CPS / concurrency.
- **SCALE/RAMP CASE**: gradual ramp with SLO tracking.
- **CHAOS CASE**: fault injection / partial failures.

**Area Codes (AAA):**
- `SIG` SIP signaling / transaction correctness
- `TOP` topology hiding & privacy
- `SEC` security posture (policy, DoS, authN/Z)
- `TLS` SIP-TLS/mTLS, cert lifecycle
- `ROU` routing / interconnect policy / peering
- `NAT` NAT traversal + SDP address manipulation
- `MED` media policy (anchor/relay decisions, SRTP expectations)
- `INT` interoperability (dialect quirks, profiles)
- `OAM` observability, logs, metrics, tracing
- `HAZ` resiliency/HA, state, failover, restart
- `CHA` explicit chaos experiments (multi-fault, partitions)

---

## Test Cases Table

| TEST ID | SLOGAN | AREA | POSITIVE CASE | NEGATIVE CASE | LOAD CASE | SCALE/RAMP CASE | CHAOS CASE |
|---|---|---|---|---|---|---|---|
| SIG-001 | "INVITE basics work" | SIP Core | INVITE→100/180/200, ACK, BYE clean teardown | Malformed INVITE → 400/488 with safe logging | 2k CPS INVITE with 20s calls | Ramp 100→2k CPS over 30m, track 95p latency | Drop upstream route mid-call; ensure BYE/cleanup |
| SIG-002 | "Dialog state sane" | SIP Core | Re-INVITE updates SDP mid-call | Out-of-dialog BYE → 481 | 10k concurrent dialogs | Ramp dialogs 1k→10k | Kill dialog-state cache node; ensure safe degradation |
| SIG-003 | "Forking rules honored" | SIP Core | Fork to 2 downstream legs; first 200 wins; cancel others | Fork loop detection; prevent recursion | 500 CPS with 2-leg forking | Ramp fork fanout 1→4 legs | Inject one leg stuck; verify timer/cancel behavior |
| SIG-004 | "CANCEL works" | SIP Core | CANCEL before 200 → 487 + 200 to CANCEL | CANCEL after final 200 ignored correctly | 1k CPS with 30% cancels | Ramp cancel ratio 0→50% | Partition to one peer: ensure cancel timers fire |
| SIG-005 | "Timer sanity" | SIP Core | Timer B/F/J respected per profile | Peer retrans storm; avoid amplification | High retrans test (loss 5–10%) | Ramp packet loss 0→10% | Drop 50% packets on one side; no resource leak |
| SIG-006 | "PRACK (100rel)" | SIP Core | 183 + PRACK supported if enabled | PRACK missing → fail per policy | 500 CPS with 100rel | Ramp 0→100% 100rel calls | PRACK path delayed; ensure timers/cleanup |
| SIG-007 | "UPDATE supported" | SIP Core | UPDATE for early media / SDP refresh | UPDATE unsupported → 501/420 per config | 500 CPS w/ UPDATE | Ramp UPDATE rate | Randomly kill UPDATE handler thread; circuit-break |
| SIG-008 | "REFER policy" | SIP Core | REFER allowed only for trusted peers | REFER from untrusted peer blocked | 200 CPS REFER attempts | Ramp REFER attempts | Peer floods REFER; rate-limit triggers |
| SIG-009 | "OPTIONS keepalive" | SIP Core | Respond to OPTIONS per peering policy | OPTIONS from unknown source blocked | 5k OPTIONS/s | Ramp OPTIONS 100→5k/s | Peer restarts; OPTIONS reflect peer down/up |
| SIG-010 | "REGISTER passthrough/deny" | SIP Core | If supported, REGISTER routed correctly | Unauthorized REGISTER rejected | 200 CPS REGISTER | Ramp REGISTER rate | Restart auth backend; fail closed |
| TOP-001 | "Hide internal Via/Record-Route" | Topology | Strip/replace internal Via + Record-Route at boundary | Leak internal IP/hostnames should be prevented | 2k CPS with header rewrite | Ramp rewrite rules 10→200 | Toggle rewrite rules store mid-flight; no crash |
| TOP-002 | "Contact rewriting" | Topology | Replace Contact host with edge FQDN | Contact with private IP → rewritten | 2k CPS | Ramp 100→2k | DNS failure for edge FQDN; safe fallback |
| TOP-003 | "P-headers policy" | Topology | Remove internal P-headers not allowed at NNI | Block P-Access-Network-Info leakage if required | 1k CPS | Ramp P-header density | Corrupt header parser input; no panic |
| TOP-004 | "SDP c= address hiding" | Topology | Rewrite SDP connection address to media relay if anchoring | If not anchoring, enforce no-private-IP policy | 1k CPS w/ SDP | Ramp calls with SDP variants | Kill media relay; ensure calls fail predictably |
| TOP-005 | "Topology hide in responses" | Topology | 200 OK strips internal Server/Allow quirks | Response containing internal route not forwarded | 1k CPS | Ramp responses per second | Peer sends oversized headers; protect memory |
| SEC-001 | "Peer allowlist" | Security | Only configured peers accepted | Unknown IP/FQDN rejected (403/488/503 per policy) | 10k connection attempts | Ramp unknown sources 0→10k | Rotate peer list while under attack; atomic update |
| SEC-002 | "Rate limit INVITE floods" | Security | Legit traffic passes at configured CPS | Flood triggers 503/429-style policy action | 20k CPS burst | Ramp 1k→20k CPS | Disable one limiter instance; system still protected |
| SEC-003 | "SIP fuzz hardening" | Security | Valid SIP parses | Fuzzed SIP never crashes; returns 4xx/5xx | 1M fuzz msgs | Ramp fuzz intensity | Crash a parser worker; supervisor restarts safely |
| SEC-004 | "Header size limits" | Security | Headers within limit accepted | Oversized headers rejected | 5k msgs/s oversized | Ramp header sizes | Out-of-memory pressure; system sheds load |
| SEC-005 | "Method allowlist" | Security | Only approved methods allowed | Unknown method → 405/501 | 2k CPS mixed methods | Ramp method diversity | Policy engine restart; default deny enforced |
| SEC-006 | "Replay / transaction abuse" | Security | Retrans handled per SIP | Replayed INVITE w/ same branch/CSeq handled safely | 5k replays/s | Ramp replay rate | Partition to state store; ensure idempotence |
| SEC-007 | "Fraud patterns (basic)" | Security | Normal calling patterns allowed | High ASR-fraud pattern triggers policy action | 2k CPS targeted | Ramp fraud patterns | Disable fraud module; fallback limits apply |
| TLS-001 | "mTLS peering handshake" | TLS | Mutual TLS succeeds with valid chain | Invalid cert/CA → fail closed | 1k new TLS/s | Ramp TLS sessions 10→1k/s | Rotate peer cert during traffic; no outage |
| TLS-002 | "TLS versions/ciphers" | TLS | Only approved ciphers negotiated | Weak cipher attempt rejected | 2k handshakes/s | Ramp handshake rate | Force cipher mismatch; clear errors, no CPU spike |
| TLS-003 | "OCSP/CRL behavior" | TLS | Revoked cert blocked (if enforced) | OCSP responder down → policy (soft/hard fail) | OCSP load 100/s | Ramp OCSP query rate | Simulate OCSP timeout; prevent thread starvation |
| TLS-004 | "SIP over TCP fallback" | TLS | If allowed, TCP accepted for trusted peers | TCP from untrusted rejected | 5k TCP connects | Ramp connects | SYN flood; kernel + app protections hold |
| TLS-005 | "Cert hot reload" | TLS | New cert applied without restarting process | Bad cert rejected; old remains active | 1k CPS during reload | Ramp reload frequency | Reload loop fault; circuit-break reload |
| ROU-001 | "Outbound routing policy" | Routing | Route to correct peer based on number plan | No route → 404/503 per policy | 2k CPS mixed destinations | Ramp destination diversity | Remove route table mid-run; fail deterministic |
| ROU-002 | "ENUM/DNS resolution" | Routing | NAPTR/SRV resolution selects correct target | DNS NXDOMAIN → fallback route | 500 DNS/s | Ramp DNS QPS | DNS server partition; cache + fallback prevents collapse |
| ROU-003 | "Loop detection" | Routing | Prevent route loops using Via/Route heuristics | Loop attempt blocked fast | 2k CPS looped | Ramp loop attempts | Peer misconfig toggles; avoid oscillation |
| ROU-004 | "Least-cost routing" | Routing | Choose cheapest peer for destination | Peer violates SLA; remove from pool | 2k CPS w/ LCR | Ramp call mix | Kill one peer; traffic rebalances |
| ROU-005 | "Emergency routing" | Routing | Emergency dialed numbers follow special route | Non-emergency cannot use emergency path | 200 CPS emergency | Ramp emergency bursts | Kill emergency peer; fallback per regulatory policy |
| NAT-001 | "NAT traversal anchoring" | NAT/SDP | Offer/answer rewrites correctly via relay | Private IP in SDP blocked if no relay | 1k CPS w/ NAT cases | Ramp NAT density 0→80% | Relay failure mid-call; tear down clean |
| NAT-002 | "Symmetric RTP" | NAT/SDP | Detect symmetric RTP and lock | Mismatched RTP source blocked | 10k RTP flows | Ramp RTP flows | Random packet reordering/loss; stable |
| NAT-003 | "ICE passthrough policy" | NAT/SDP | If supported, ICE attributes preserved/normalized | ICE abuse / oversized candidates blocked | 500 CPS ICE | Ramp ICE candidates count | Candidate explosion attack; protection holds |
| MED-001 | "SRTP enforcement" | Media | If required, reject non-SRTP offers | Downgrade attack blocked | 1k CPS SRTP | Ramp SRTP-only | Kill SRTP key service; fail closed |
| MED-002 | "Codec policy" | Media | Allowlist codecs (AMR-WB, etc.) | Disallowed codec only → 488 | 1k CPS codec mix | Ramp codec diversity | Peer sends broken SDP; parser safe |
| MED-003 | "DTMF policy" | Media | RFC2833/INFO accepted per config | Unsupported DTMF rejected | 500 CPS w/ DTMF | Ramp DTMF frequency | Mid-call media relay restart; recover or teardown |
| INT-001 | "Dialect normalization (From/To)" | Interop | Normalize URI forms across peers | Reject invalid URI schemes | 2k CPS mixed formats | Ramp "weird SIP" % | Introduce peer header bug; rule-based workaround works |
| INT-002 | "Header manipulation rules" | Interop | Apply rewrite templates (To/From/URI-host) correctly | Bad rule doesn't crash; safely rejected | 2k CPS w/ rewrites | Ramp rule count 10→500 | Corrupt ruleset store; last-known-good continues |
| INT-003 | "Topology hiding via manipulation" | Interop | Confirm multi-vendor topology hiding outcomes | Leak detection triggers alarm | 1k CPS | Ramp call volume | Turn on extra headers from peer; hiding still holds |
| INT-004 | "SIP over IPv6 peer" | Interop | IPv6 peering stable (SIP/TLS) | IPv6 malformed addresses rejected | 1k CPS IPv6 | Ramp IPv6 share | Disable IPv6 route; failover to IPv4 if allowed |
| OAM-001 | "Structured logs" | Observability | Correlate call-id, branch, peer-id in logs | No PII leakage in logs | 2k CPS with logging | Ramp log volume | Disk full / log sink down; backpressure + drop policy |
| OAM-002 | "Metrics: CPS, ASR, PDD" | Observability | Export core KPIs | Missing metrics triggers alert | 10k CPS | Ramp to saturation | Kill metrics exporter; app keeps running |
| OAM-003 | "Tracing propagation" | Observability | Trace id per dialog across modules | Bad trace header ignored | 2k CPS traced | Ramp trace sampling 0→100% | Collector down; buffering bounded |
| HAZ-001 | "Process restart recovery" | Resilience | Restart IBCF process; no config loss | Restart loop detected; safe-mode | 2k CPS during rolling restart | Ramp restarts | Kill -9 mid-call; cleanup via timers |
| HAZ-002 | "Active/active peers" | Resilience | Two instances share traffic without oscillation | Split-brain avoided | 2k CPS | Ramp instances 1→4 | Partition between instances; avoid conflicting routing |
| HAZ-003 | "State store outage" | Resilience | Continue stateless paths if possible | If state required, reject fast | 2k CPS | Ramp state dependence | Drop state store 5m; no memory leak |
| HAZ-004 | "Backpressure & shedding" | Resilience | Under overload, degrade gracefully | No unbounded queues | 20k CPS overload | Ramp to overload | CPU throttling injected; load shedding engages |
| CHA-001 | "Network partition to one peer" | Chaos | Other peers unaffected; correct failover | No crash or global lock | N/A | Ramp partition frequency | Partition A for 10m while traffic continues |
| CHA-002 | "DNS poisoning / wrong SRV" | Chaos | Detect invalid peer cert/name mismatch | Prevent routing to rogue endpoint | N/A | Ramp DNS changes | Inject bad SRV + valid IP; mTLS blocks |
| CHA-003 | "Time skew" | Chaos | TLS + token checks tolerate bounded skew | Excess skew alarms + safe fail | N/A | Ramp skew 0→10m | Jump clock +5m mid-traffic; observe failure modes |
| CHA-004 | "Packet loss + jitter storm" | Chaos | SIP retrans stable; no amplification | Protect CPU | N/A | Ramp loss 0→20% | 20% loss + 200ms jitter for 15m |
| CHA-005 | "Config reload during peak" | Chaos | Atomic config switch; no partial application | Invalid config rejected | N/A | Ramp reloads/h | Reload every 60s under 2k CPS |

---

## Examples

- **Positive example**: `SIG-001` normal call setup/teardown.
- **Negative example**: `SEC-001` unknown peer rejected.
- **Load example**: `SEC-002` flood protection at 20k CPS.
- **Scale/Ramp example**: `TLS-001` ramp TLS sessions 10→1k/s while tracking 95p handshake time.
- **Chaos example**: `CHA-001` partition one peer while others remain healthy.

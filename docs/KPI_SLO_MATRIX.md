# KPI / SLO Acceptance Criteria Matrix
SIG-GW / IBCF Platform — Carrier-Grade Acceptance Criteria

**Author**: Daniel  
**Scope**: Production-readiness SLOs for IBCF + STIR/SHAKEN + LI + Emergency  
**Applies To**: IMS Interconnect (NNI/Ic), CNF deployments

---

## 1. Objectives

Define measurable service-level objectives (SLOs) and key performance indicators (KPIs) required before:

- Production go-live
- Interconnect activation
- Regulatory approval
- Peering certification

---

## 2. Core Signaling KPIs

| KPI | Target | Critical Threshold | Measurement Window | Acceptance Criteria |
|-----|--------|-------------------|-------------------|--------------------|
| Call Setup Success Rate (CSSR) | ≥ 99.5% | < 99.0% | 15 min | No sustained drop >5 min |
| Post Dial Delay (PDD) | ≤ 250 ms (NNI) | > 400 ms | 5 min | 95th percentile |
| INVITE Processing Latency | ≤ 10 ms | > 25 ms | 1 min | 99th percentile |
| CPS Sustained | 100% rated load | N/A | 30 min | No packet loss |
| Max CPS Burst | 2× rated for 60s | Crash or degradation | Burst window | No restart |
| Dialog Establishment Time | ≤ 500 ms | > 1000 ms | 5 min | End-to-end |
| BYE Processing Latency | ≤ 5 ms | > 15 ms | 1 min | 99th percentile |
| Re-INVITE Success Rate | ≥ 99.0% | < 98.0% | 15 min | Mid-call updates |

---

## 3. STIR/SHAKEN KPIs

| KPI | Target | Critical Threshold | Measurement Window | Acceptance Criteria |
|------|--------|-------------------|-------------------|--------------------|
| Signing Latency | ≤ 5 ms | > 10 ms | 1 min | 99th percentile |
| Verification Latency | ≤ 5 ms | > 10 ms | 1 min | 99th percentile |
| Verification Failure Accuracy | 100% detection | False positives < 0.1% | 1 hour | All invalid signatures caught |
| Certificate Cache Hit Rate | ≥ 95% | < 85% | 1 hour | Reduces OCSP load |
| Attestation Integrity | 0% improper upgrades | Any violation fails | Continuous | Never upgrade A→B→C |
| Identity Header Insertion Rate | 100% (outbound) | < 99.5% | 15 min | All outbound calls signed |
| Certificate Fetch Latency | ≤ 100 ms | > 500 ms | 1 min | 95th percentile |
| OCSP Response Time | ≤ 200 ms | > 1000 ms | 1 min | 95th percentile |

---

## 4. Lawful Intercept KPIs

| KPI | Target | Critical Threshold | Measurement Window | Acceptance Criteria |
|------|--------|-------------------|-------------------|--------------------|
| Intercept Coverage | 100% of active warrants | < 99.9% | Continuous | No missed intercepts |
| Intercept Latency Overhead | ≤ 5 ms | > 15 ms | 1 min | 99th percentile |
| Media Duplication Accuracy | 100% | < 99.9% | Continuous | All media packets duplicated |
| Intercept During Failover | No loss | Any loss fails | Failover event | Zero tolerance |
| Audit Log Completeness | 100% | < 99.9% | Continuous | All events logged |
| Warrant Activation Time | ≤ 1 sec | > 5 sec | Activation event | Immediate effect |
| LI Mediation Device Availability | ≥ 99.9% | < 99.0% | 15 min | MD uptime |
| TLS Decrypt for LI Latency | ≤ 10 ms | > 30 ms | 1 min | 99th percentile |

---

## 5. Emergency KPIs

| KPI | Target | Critical Threshold | Measurement Window | Acceptance Criteria |
|------|--------|-------------------|-------------------|--------------------|
| Emergency PDD | ≤ 200 ms | > 400 ms | 5 min | 95th percentile |
| Emergency Drop Rate | 0% | > 0% | Continuous | Zero tolerance |
| Emergency Failover Recovery | ≤ 1 second | > 3 seconds | Failover event | No call loss |
| Emergency Priority Preemption | 100% success | < 99.9% | Continuous | Always highest priority |
| Location Header Preservation | 100% | < 99.9% | Continuous | All location data preserved |
| Emergency Detection Latency | ≤ 10 ms | > 50 ms | 1 min | 99th percentile |
| PSAP Routing Accuracy | 100% | < 99.5% | Continuous | Correct PSAP selection |
| Emergency Callback Routing | 100% success | < 99.0% | Continuous | PSAP callbacks work |

---

## 6. Resilience KPIs

| KPI | Target | Critical Threshold | Measurement Window | Acceptance Criteria |
|------|--------|-------------------|-------------------|--------------------|
| Failover Recovery Time | ≤ 1 sec | > 3 sec | Failover event | Active/Active switchover |
| Zero Call Loss During Rolling Upgrade | Yes | Any loss fails | Upgrade window | No active calls dropped |
| Memory Leak Rate | 0 over 24h | > 0 | 24 hours | No memory growth |
| CPU Under Load | ≤ 85% sustained | > 95% | 30 min | Headroom for bursts |
| Packet Loss Handling | No crash up to 10% loss | Crash or degradation | Test window | Graceful degradation |
| State Store Recovery | ≤ 5 sec | > 30 sec | Recovery event | Fast state restoration |
| DNS Failover Time | ≤ 2 sec | > 10 sec | Failover event | DNS-based failover |
| Service Restart Time | ≤ 10 sec | > 60 sec | Restart event | Fast recovery |

---

## 7. Observability KPIs

| KPI | Target | Critical Threshold | Measurement Window | Acceptance Criteria |
|------|--------|-------------------|-------------------|--------------------|
| Metrics Availability | ≥ 99.99% | < 99.9% | Continuous | Prometheus uptime |
| Log Correlation Accuracy | 100% call-ID trace | < 99.9% | Continuous | All logs traceable |
| Alert Detection Time | ≤ 30 sec | > 5 min | Alert event | Fast incident detection |
| False Alert Rate | < 1% | > 5% | 24 hours | Reduce noise |
| Trace Completeness | ≥ 95% | < 90% | 1 hour | OpenTelemetry coverage |
| Dashboard Load Time | ≤ 2 sec | > 10 sec | On-demand | Grafana performance |
| Log Retention Compliance | 100% | < 100% | Continuous | Regulatory requirements |

---

## 8. Security KPIs

| KPI | Target | Critical Threshold | Measurement Window | Acceptance Criteria |
|------|--------|-------------------|-------------------|--------------------|
| mTLS Handshake Success | ≥ 99.9% | < 99.0% | 15 min | Peer authentication |
| TLS Version Compliance | 100% TLS 1.2+ | Any TLS 1.0/1.1 | Continuous | No weak protocols |
| Rate Limit Effectiveness | 100% block | < 95% | Attack window | DoS protection |
| Certificate Rotation Success | 100% | < 99.9% | Rotation event | Zero downtime rotation |
| Topology Hiding Effectiveness | 100% | Any leak fails | Continuous | No internal IPs exposed |

---

## 9. SLO Gate Criteria

Production approval requires:

### 9.1 Sustained Load Test
- **Duration**: 72-hour sustained load test
- **Load**: 100% rated CPS continuously
- **All KPIs**: Within tolerance for entire duration
- **Incidents**: No severity-1 incidents

### 9.2 Regression Testing
- **Emergency Scenarios**: All emergency test cases pass
- **LI Scenarios**: All LI test cases pass
- **STIR/SHAKEN**: All STIR test cases pass
- **IBCF Functional**: All IBCF test cases pass

### 9.3 Failover Validation
- **Active/Active**: Zero call loss during failover
- **Rolling Upgrade**: Zero call loss during upgrade
- **Emergency Continuity**: Emergency calls survive failover
- **LI Continuity**: Intercepts survive failover

### 9.4 Regulatory Compliance
- **LI Compliance**: All LI requirements met
- **Emergency Compliance**: All emergency requirements met
- **STIR/SHAKEN Compliance**: All STIR requirements met
- **Security Compliance**: All security requirements met

### 9.5 Documentation
- **KPI Dashboard**: All KPIs visible and tracked
- **Test Reports**: All test catalogs executed and documented
- **PIXIT Configs**: All PIXIT parameters documented
- **Runbooks**: Operational procedures documented

---

## 10. KPI Measurement Methodology

### 10.1 Measurement Tools
- **Prometheus**: Metrics collection
- **Grafana**: KPI dashboards
- **OpenTelemetry**: Distributed tracing
- **Loki**: Log aggregation
- **Custom**: Test harness metrics

### 10.2 Sampling
- **Real-time**: Continuous monitoring
- **Aggregation**: 1-minute, 5-minute, 15-minute windows
- **Percentiles**: p50, p95, p99, p99.9
- **Rolling Windows**: 15-minute, 1-hour, 24-hour

### 10.3 Alerting
- **Critical Threshold**: Immediate alert
- **Warning Threshold**: Alert after 5 minutes
- **Escalation**: Auto-escalate if unresolved
- **On-Call**: 24/7 coverage for critical KPIs

---

## 11. KPI Dashboard Requirements

### 11.1 Real-Time Dashboard
- Current CPS
- Active dialogs
- PDD (p50, p95, p99)
- Error rates
- System health

### 11.2 Historical Dashboard
- 24-hour trends
- 7-day trends
- 30-day trends
- SLO compliance over time

### 11.3 Component Dashboards
- **IBCF**: Signaling metrics
- **STIR/SHAKEN**: Sign/verify metrics
- **LI**: Intercept metrics
- **Emergency**: Emergency call metrics
- **Infrastructure**: CPU, memory, network

---

## 12. Continuous Improvement

### 12.1 KPI Review
- **Weekly**: Review KPI trends
- **Monthly**: Deep dive on outliers
- **Quarterly**: SLO target review
- **Annually**: SLO target adjustment

### 12.2 Optimization
- Identify bottlenecks
- Optimize hot paths
- Reduce latency
- Improve reliability

---

## End of KPI / SLO Matrix

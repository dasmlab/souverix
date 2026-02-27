# OpenShift CNF Reference Test Harness
SIG-GW / IBCF Validation Architecture

**Author**: Daniel  
**Platform**: OpenShift CNF

---

## 1. Architecture Overview

```
+------------------------------+
| Traffic Generator (SIPp)     |
+------------------------------+
| Chaos Injector (Litmus)      |
+------------------------------+
| Metrics Stack (Prom + Grafana)|
+------------------------------+
| SIG-GW Pods (IBCF)           |
+------------------------------+
| STIR Service + Vault         |
+------------------------------+
| Media Relay Pods             |
+------------------------------+
```

### 1.1 Component Placement

- **Traffic Generator**: Separate namespace (`testrig`)
- **Chaos Injector**: Separate namespace (`chaos`)
- **Metrics Stack**: Monitoring namespace (`monitoring`)
- **SIG-GW Pods**: Production namespace (`ims`)
- **STIR Service**: Production namespace (`ims`)
- **Media Relay**: Production namespace (`ims`)

---

## 2. Components

### 2.1 Traffic Generation

#### SIPp
- **Purpose**: SIP traffic generation
- **Deployment**: StatefulSet or DaemonSet
- **Configuration**: Test scenario scripts
- **Metrics**: Generated CPS, call success rate

#### Custom Load Injector
- **Purpose**: Advanced load patterns
- **Features**:
  - Ramp-up/ramp-down
  - Burst patterns
  - Emergency call injection
  - STIR-aware traffic

#### STIR-Aware Traffic Scripts
- **Purpose**: Generate STIR-signed traffic
- **Features**:
  - A/B/C attestation levels
  - Certificate rotation scenarios
  - Invalid signature testing

### 2.2 Chaos Engineering

#### Litmus Chaos
- **Purpose**: Chaos injection framework
- **Experiments**:
  - Pod kill
  - Network partition
  - CPU throttling
  - Memory pressure
  - DNS corruption

#### Custom Chaos Scenarios
- **Purpose**: IMS-specific chaos
- **Scenarios**:
  - STIR service crash
  - LI mediation device failure
  - Certificate expiration
  - State store partition

### 2.3 Metrics Stack

#### Prometheus
- **Purpose**: Metrics collection
- **Scrape Interval**: 15 seconds
- **Retention**: 30 days
- **Targets**: All IMS components

#### Grafana Dashboards
- **Purpose**: Visualization
- **Dashboards**:
  - IBCF Signaling
  - STIR/SHAKEN
  - LI Intercepts
  - Emergency Calls
  - Infrastructure

#### Alertmanager
- **Purpose**: Alert routing
- **Channels**: PagerDuty, Slack, Email
- **Rules**: KPI threshold violations

---

## 3. CI/CD Integration

### 3.1 Automated Regression Pipeline

**Trigger**: Push to main/develop

**Stages**:
1. **Build**: Container images
2. **Deploy**: Test environment
3. **Functional Tests**: Unit + integration
4. **Load Tests**: Baseline performance
5. **Chaos Tests**: Resilience validation
6. **Report**: Test results

### 3.2 Performance Baseline Comparison

**Purpose**: Detect performance regressions

**Comparison**:
- Current run vs. baseline
- PDD trends
- CPU/Memory trends
- Latency trends

**Action**: Alert if > 10% degradation

### 3.3 Canary Testing

**Purpose**: Validate new releases

**Process**:
1. Deploy canary (10% traffic)
2. Monitor KPIs
3. Compare to baseline
4. Promote or rollback

---

## 4. Key Metrics Collected

### 4.1 Signaling Metrics

- **CPS**: Calls per second
- **PDD**: Post dial delay (p50, p95, p99)
- **CSSR**: Call setup success rate
- **Dialog Count**: Active dialogs
- **Error Rate**: 4xx/5xx responses

### 4.2 Infrastructure Metrics

- **CPU**: Per pod, per node
- **Memory**: Per pod, per node
- **Network**: Bandwidth, packet loss
- **Disk**: I/O, space

### 4.3 STIR/SHAKEN Metrics

- **Signing Latency**: p50, p95, p99
- **Verification Latency**: p50, p95, p99
- **Certificate Cache Hit Rate**: %
- **OCSP Response Time**: p50, p95, p99
- **Attestation Distribution**: A/B/C counts

### 4.4 LI Metrics

- **Intercept Coverage**: % of warrants
- **Intercept Latency**: Overhead
- **Media Duplication Rate**: %
- **Audit Log Completeness**: %

### 4.5 Emergency Metrics

- **Emergency PDD**: p50, p95, p99
- **Emergency Drop Rate**: %
- **Priority Preemption**: Success rate
- **Location Preservation**: %

### 4.6 TLS Metrics

- **TLS Handshake Time**: p50, p95, p99
- **mTLS Success Rate**: %
- **Certificate Rotation Time**: Seconds
- **OCSP Validation Time**: p50, p95, p99

---

## 5. Scaling Strategy

### 5.1 HPA for Signaling Pods

**Configuration**:
```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: ibcf-signaling
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: ibcf-signaling
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
```

**Triggers**:
- CPU > 70%
- Memory > 80%
- Custom metric (CPS)

### 5.2 Dedicated Media Plane

**Architecture**:
- Separate namespace for media
- Dedicated nodes (optional)
- Independent scaling
- SR-IOV enabled

### 5.3 SR-IOV Enabled Nodes

**Purpose**: Low-latency media handling

**Configuration**:
- SR-IOV network attachments
- DPDK support (optional)
- CPU pinning
- NUMA awareness

### 5.4 Node Affinity for Latency

**Purpose**: Optimize for low latency

**Configuration**:
```yaml
affinity:
  nodeAffinity:
    requiredDuringSchedulingIgnoredDuringExecution:
      nodeSelectorTerms:
      - matchExpressions:
        - key: node-role.kubernetes.io/latency-optimized
          operator: In
          values:
          - "true"
```

---

## 6. Test Harness Execution Flow

### 6.1 Deploy CNF

**Steps**:
1. Apply Kubernetes manifests
2. Wait for pods ready
3. Verify health checks
4. Confirm metrics collection

### 6.2 Deploy Test Harness

**Steps**:
1. Deploy traffic generator
2. Deploy chaos injector
3. Deploy metrics stack (if not existing)
4. Configure test scenarios

### 6.3 Inject Load

**Steps**:
1. Start baseline load (50% rated)
2. Ramp to 100% rated
3. Sustain for test duration
4. Monitor KPIs

### 6.4 Inject Chaos

**Steps**:
1. Select chaos scenario
2. Execute chaos experiment
3. Monitor system behavior
4. Validate recovery

### 6.5 Validate KPIs

**Steps**:
1. Check all KPIs within target
2. Verify no degradation
3. Confirm recovery
4. Document results

### 6.6 Generate Report

**Steps**:
1. Collect metrics
2. Generate dashboards
3. Create test report
4. Archive results

---

## 7. Test Scenarios

### 7.1 Baseline Performance

**Configuration**:
- Load: 100% rated CPS
- Duration: 1 hour
- Metrics: All KPIs

**Validation**:
- All KPIs within target
- Stable performance
- No errors

### 7.2 Burst Load

**Configuration**:
- Load: 2× rated CPS
- Duration: 60 seconds
- Metrics: CPS, PDD, CPU

**Validation**:
- No packet loss
- PDD within target
- No crash

### 7.3 Failover Test

**Configuration**:
- Load: 100% rated CPS
- Event: Kill primary pod
- Metrics: Failover time, call loss

**Validation**:
- Failover ≤ 1 second
- Zero call loss
- No emergency drop

### 7.4 Chaos Test

**Configuration**:
- Load: 100% rated CPS
- Chaos: Pod kill, network partition
- Metrics: Recovery time, call handling

**Validation**:
- Graceful degradation
- Fast recovery
- No data loss

---

## 8. Monitoring & Alerting

### 8.1 Prometheus Rules

**Example**:
```yaml
groups:
- name: ibcf_alerts
  rules:
  - alert: HighPDD
    expr: histogram_quantile(0.95, rate(ibcf_pdd_bucket[5m])) > 0.4
    for: 5m
    annotations:
      summary: "PDD above threshold"
```

### 8.2 Grafana Dashboards

**Dashboards**:
- Real-time IBCF metrics
- Historical trends
- Component health
- KPI compliance

### 8.3 Alert Routing

**Channels**:
- Critical: PagerDuty
- Warning: Slack
- Info: Email

---

## 9. Test Data Management

### 9.1 Test Scenarios

**Storage**: Git repository
**Format**: YAML/JSON
**Versioning**: Git tags

### 9.2 PIXIT Configurations

**Storage**: `testrig/pixit/`
**Format**: YAML
**Usage**: Test execution parameters

### 9.3 Test Results

**Storage**: Object storage (S3)
**Retention**: 90 days
**Format**: JSON, PCAP, logs

---

## 10. Continuous Improvement

### 10.1 Performance Optimization

- Identify bottlenecks
- Optimize hot paths
- Reduce latency
- Improve throughput

### 10.2 Test Coverage

- Add new test scenarios
- Cover edge cases
- Validate new features
- Improve chaos scenarios

---

## End of OpenShift CNF Test Harness

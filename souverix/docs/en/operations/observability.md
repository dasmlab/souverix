# Observability Guide

## Overview

Souverix Vigile provides comprehensive observability for the Souverix Platform.

## Metrics

### Prometheus

Metrics are exposed at `/metrics` endpoint on each component.

### Key Metrics

- CPS (Calls Per Second)
- PDD (Post Dial Delay)
- Call success rate
- Component health
- STIR signing/verification latency

## Logging

### Structured Logs

All components emit structured JSON logs compatible with Loki.

### Log Levels

- ERROR
- WARN
- INFO
- DEBUG

## Tracing

### OpenTelemetry

Distributed tracing via OpenTelemetry for end-to-end call flows.

## Dashboards

### Grafana

Pre-configured dashboards for:
- Platform overview
- Component health
- Performance metrics
- Compliance tracking

---

## End of Observability Guide

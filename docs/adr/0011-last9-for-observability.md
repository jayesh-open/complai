# ADR-0011: Last9 for observability

**Status:** Accepted  
**Date:** 2026-04-25  
**Deciders:** Complai Engineering

## Context

Complai runs ~28 services across EKS, communicating via SQS/SNS and calling external APIs (Adaequare, Sandbox.co.in). We need unified observability across all three signals -- metrics, logs, and traces -- with per-tenant visibility for debugging and SLO tracking.

Options evaluated:

1. **Self-hosted stack** (Prometheus + Grafana + Loki + Tempo) -- full control, but significant operational overhead for a small team.
2. **Datadog** -- excellent product, but expensive at our service count and log volume.
3. **Last9** -- OpenTelemetry-native managed platform, competitive pricing, built-in SLO tracking.
4. **AWS-native** (CloudWatch + X-Ray) -- integrated but limited query capabilities and poor cross-service trace visualization.

## Decision

Last9 as our observability platform. All services instrument with the OpenTelemetry SDK; an OTel Collector DaemonSet on every EKS node exports OTLP data to Last9. Jaeger replaces Last9 in local development.

- **OTel Collector:** DaemonSet on every node receives OTLP from every pod, enriches with Kubernetes metadata, exports to Last9's OTLP endpoint.
- **Metrics:** per-service RED (rate/errors/duration), SLI dashboards per tier, business metrics (filings/min, IRNs/min, recon match rate).
- **Logs:** structured JSON via zerolog, with `tenant_id` and correlation IDs on every line. PII is redacted before export.
- **Traces:** full request traces across services with `tenant_id`, `gstin`, `tan`, `pan` as baggage on every span.
- **Dashboards:** per-service, per-provider (Adaequare/Sandbox latency and success rates), per-tenant, and SLO burn rate.
- **Alerting:** Last9 alert rules route to PagerDuty for on-call; Slack for non-critical. Alert on filing success rate <99%, gateway error rate >1%, auth failure spikes, RLS policy violations, DLQ depth >0.

## Consequences

### Positive
- OTel-native instrumentation means zero vendor lock-in on the instrumentation side -- switching observability backends requires only changing the Collector's exporter configuration.
- Per-tenant observability is possible via `tenant_id` labels, enabling tenant-specific debugging without cross-tenant data leakage.
- Built-in SLO tracking and burn rate alerting align with our tiered availability targets (99.99% / 99.9% / 99.5%).
- Managed service -- no Prometheus/Grafana/Loki/Tempo clusters to operate.

### Negative
- SaaS dependency -- observability data is hosted externally (Last9's infrastructure).
- Must instrument every service consistently; missing instrumentation creates blind spots.
- Monthly cost scales with data ingestion volume.

### Risks
- Last9 outage would not affect application functionality but would blind us to issues during the outage. Mitigated by: local logs are always available on the node; CloudWatch basic metrics serve as a backstop.

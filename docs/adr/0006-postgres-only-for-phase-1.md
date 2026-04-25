# ADR-0006: Postgres-only for Phase 1 (OLTP + analytics)

**Status:** Accepted  
**Date:** 2026-04-25  
**Deciders:** Complai Engineering

## Context

Complai needs both transactional (OLTP) and analytical workloads:

- OLTP: filing workflows, invoice CRUD, reconciliation queries, real-time dashboards.
- Analytics: CFO dashboards, compliance score aggregations, filing summaries, ITC-at-risk calculations.

Options evaluated:

1. **Postgres (OLTP) + ClickHouse (OLAP)** -- best analytical performance, but adds operational complexity (another database to manage, CDC pipeline, schema sync).
2. **Postgres (OLTP) + MongoDB (flexible docs)** -- document flexibility, but adds another data store and breaks RLS invariant.
3. **Postgres for everything** -- single database handles both workloads; use read replica for analytics and JSONB for flexible fields.

## Decision

Single Amazon RDS PostgreSQL 16 Multi-AZ instance with a read replica serves all Phase 1 needs. No MongoDB, no Kafka, no ClickHouse in Phase 1.

- **OLTP:** primary instance (`db.r7g.2xlarge`, 8 vCPU, 64GB RAM, 12,000 IOPS) handles all transactional workloads.
- **Analytics:** asynchronous read replica (same instance class) handles all BI, reporting, and dashboard queries.
- **Flexible fields:** JSONB columns with GIN indexes replace the need for a document database.
- **Materialized views** pre-compute expensive aggregates (vendor compliance scores, filing summaries, ITC-at-risk, CFO KPIs) on scheduled refresh intervals.
- One RDS instance hosts ~20 logical databases, one per owning service. Each service owns its schema and migrations exclusively.

## Consequences

### Positive
- Operational simplicity: one database technology, one backup strategy, one monitoring setup, one set of expertise required.
- RLS works uniformly across all tables -- no second database that bypasses tenant isolation.
- Postgres handles our Phase 1 scale comfortably (500 enterprise customers, ~50K daily transactions).
- JSONB provides MongoDB-like flexibility without a separate data store.
- Materialized views handle dashboard and reporting workloads on the read replica without impacting OLTP.

### Negative
- Analytics queries on the read replica are limited by replication lag (typically <1s, but can spike under heavy write load).
- Must migrate specific high-volume tables to a dedicated OLAP store when they exceed ~500M rows (estimated Year 2 for `audit_events` and `transaction_stream`).
- No columnar storage optimizations for large analytical scans.

### Risks
- Read replica falling behind during filing peaks. Mitigated by: materialized views reduce ad-hoc analytical query load; read replica is the same instance class as primary.

**Upgrade path:** When a specific table crosses ~500M rows, add ClickHouse alongside via Debezium CDC for that table only. The rest of the system continues on Postgres.

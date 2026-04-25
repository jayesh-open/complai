# ADR-0013: Outbox pattern via polling (not Debezium) in Phase 1

**Status:** Accepted  
**Date:** 2026-04-25  
**Deciders:** Complai Engineering

## Context

Complai's architecture requires reliable delivery of state changes to external APIs (government portals via Adaequare/Sandbox) and cross-service event propagation. The core problem: a business service must atomically update its database and publish an event. Without the outbox pattern, a crash between the DB commit and the message publish causes data inconsistency.

Two implementations of the outbox pattern were evaluated:

1. **CDC-based (Debezium):** Debezium reads the Postgres WAL and publishes outbox rows to a streaming platform (Kafka/Kinesis). Sub-100ms latency, but requires a Debezium cluster, Kafka Connect, and a streaming platform.
2. **Polling-based:** A sidecar process in each service polls the outbox table at a fixed interval and publishes to SQS/SNS. Simpler infrastructure, higher latency.

## Decision

Outbox table in each service's database with a polling-based publisher. No Debezium, no Kafka Connect in Phase 1.

- **Outbox table schema:** `outbox (id UUID, aggregate_type TEXT, aggregate_id UUID, event_type TEXT, payload JSONB, target_queue TEXT, created_at TIMESTAMPTZ, published_at TIMESTAMPTZ NULL)`.
- **Atomic writes:** the business transaction writes domain changes and the outbox row in the same Postgres transaction. Both commit or neither does.
- **Sidecar publisher:** each service runs an outbox publisher goroutine that polls every 500ms for unpublished rows (`WHERE published_at IS NULL`), publishes to SQS/SNS, and marks rows as `published`.
- **Idempotency:** every outbox row carries a `request_id` UUID. External providers and downstream consumers deduplicate on this ID. If the publisher crashes after making an external call but before marking the row as published, the retry sends the same `request_id` -- the provider returns the same response.
- **Cleanup:** published outbox rows are retained for 7 days (for debugging), then deleted by a scheduled job.

## Consequences

### Positive
- Simple to implement -- no additional infrastructure beyond what we already have (Postgres + SQS).
- Atomic with the business transaction -- guaranteed consistency between domain state and published events.
- No Debezium cluster, no Kafka Connect, no streaming platform to operate in Phase 1.
- Easy to reason about and debug -- the outbox table is a plain Postgres table, queryable with SQL.

### Negative
- 500ms polling latency between commit and publish. Acceptable because external API latency (Adaequare, Sandbox) is typically 1-5 seconds, making the 500ms negligible in the total flow.
- Polling adds load to the database -- one query per service every 500ms. At Phase 1 scale (~20 services), this is ~40 queries/second, negligible for our RDS instance.
- No built-in exactly-once delivery -- relies on idempotency keys for deduplication.

### Risks
- If polling frequency becomes a bottleneck at higher scale. Mitigated by: the publisher can be tuned (shorter intervals, batch sizes) without architectural changes.

**Upgrade path:** Switch to Debezium CDC reading the same outbox table, publishing to Kinesis or MSK, when scale demands sub-100ms delivery latency. The outbox table schema does not change.

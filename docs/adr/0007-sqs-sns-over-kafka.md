# ADR-0007: SQS/SNS over Kafka for Phase 1 messaging

**Status:** Accepted  
**Date:** 2026-04-25  
**Deciders:** Complai Engineering

## Context

Complai needs asynchronous messaging for:

- Outbox delivery: reliably shipping state changes from domain services to gateway services.
- Event fan-out: broadcasting events (FilingCompleted, InvoiceCreated, VendorCreated) to multiple consumers.
- Job queues: OCR processing, notification delivery, bulk operations.

Options evaluated:

1. **Amazon MSK (Managed Kafka)** -- log-based, replay capability, high throughput, but significant operational overhead and cost at our Phase 1 scale.
2. **Amazon SQS + SNS** -- fully managed queuing and pub/sub, per-message pricing, zero ops, but no log compaction or replay.
3. **Amazon Kinesis** -- streaming, replay within retention window, but shard management complexity.

## Decision

Amazon SQS + SNS with a `MessageBus` interface abstraction in `packages/shared-kernel-go/messaging`.

- **SQS standard queues** for single-consumer flows (gateway outbound queues, OCR jobs, notifications).
- **SQS FIFO queues** where strict ordering matters (filing sagas grouped by GSTIN, TDS flows grouped by TAN).
- **SNS topics with SQS fan-out** for multi-consumer events (FilingCompleted, InvoiceCreated, VendorCreated, MasterDataChanged).
- **Dead-letter queues** on every queue with a 5-retry threshold and CloudWatch alarms on DLQ depth > 0.
- **S3 pointer pattern** for messages exceeding the 256KB SQS limit (payload stored in S3, message contains bucket+key reference).
- **MessageBus interface** abstracts SQS/SNS: publishers call `bus.Publish(topic, event)`, consumers call `bus.Subscribe(queue, handler)`. No SQS primitives leak into business code.

## Consequences

### Positive
- Zero operational overhead -- fully managed by AWS, no broker clusters to maintain.
- Per-message pricing is significantly cheaper than MSK at Phase 1 volumes (~30 USD/month vs ~800+ USD/month for a minimal MSK cluster).
- Native AWS integration (IAM, CloudWatch, X-Ray).
- DLQ support built-in with automatic alerting.
- FIFO queues provide exactly-once processing and ordering where needed.

### Negative
- No log compaction -- cannot replay historical events from the queue.
- No consumer offset tracking -- once a message is processed and deleted, it is gone.
- 256KB message size limit requires the S3 pointer pattern for larger payloads.
- Less ecosystem tooling compared to Kafka (no ksqlDB, no Kafka Streams equivalent).

### Risks
- If we need event replay (e.g., rebuilding a read model from historical events), SQS cannot provide it. Mitigated by: the outbox table in Postgres retains all events and can be replayed from there.

**Upgrade path:** Swap the SQS/SNS implementation for MSK behind the same `MessageBus` interface when scale or replay requirements demand it. No business code changes required.

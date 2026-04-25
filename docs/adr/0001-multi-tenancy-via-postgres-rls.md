# ADR-0001: Multi-tenancy via Postgres RLS

**Status:** Accepted  
**Date:** 2026-04-25  
**Deciders:** Complai Engineering

## Context

Complai is a multi-tenant compliance SaaS serving hundreds of enterprise customers on shared infrastructure. We need strict tenant isolation at the data layer to prevent cross-tenant data leakage -- a non-negotiable requirement for a platform handling GST filings, TDS returns, PAN/Aadhaar data, and financial records.

Three options were evaluated:

1. **Separate database per tenant** -- strongest isolation, but creates operational complexity (connection management, migrations, backups multiply per tenant).
2. **Schema-per-tenant** -- moderate isolation, but schema migration across hundreds of schemas is fragile and tooling support is limited.
3. **Shared schema with Row-Level Security (RLS)** -- all tenants share tables; Postgres enforces row-level access based on a session variable.

## Decision

We use a shared schema with Postgres Row-Level Security on every table.

- Every table carries `tenant_id UUID NOT NULL` as a mandatory column.
- RLS policy on every table: `USING (tenant_id = current_setting('app.tenant_id')::uuid)`.
- A shared Go middleware (`packages/shared-kernel-go/tenant`) extracts `tenant_id` from the JWT, then executes `SET LOCAL app.tenant_id = '<uuid>'` at the start of every database transaction.
- All indexes are tenant-scoped: `(tenant_id, ...)` composite indexes rather than single-column indexes, ensuring the query planner uses tenant_id for partition pruning.
- If a service method fails to set the tenant context, the query returns zero rows -- RLS acts as a compile-time-ish safety net, not just a runtime check.

## Consequences

### Positive
- Strong tenant isolation without the operational complexity of per-tenant databases or schemas.
- Works correctly with connection pooling (PgBouncer) since `SET LOCAL` is scoped to the transaction.
- Standard Postgres feature -- no custom extensions, no vendor lock-in.
- Single backup strategy, single migration pipeline, single monitoring setup for all tenants.
- Scales to thousands of tenants without multiplying operational overhead.

### Negative
- Every table and every migration must include `tenant_id` -- forgetting it is a correctness bug. Enforced via CI lint rules and code review.
- Query planner requires tenant_id-leading composite indexes; single-column indexes on business keys are insufficient.
- RLS adds a small overhead to every query (negligible at our scale, but measurable under profiling).

### Risks
- A bug in the middleware that skips `SET LOCAL` would cause queries to return zero rows (safe failure mode) or, if RLS is accidentally disabled on a table, could leak data. Mitigated by: CI check that every table has RLS enabled, integration tests that verify cross-tenant isolation, and alerting on any RLS policy violation detected in audit logs.

# ADR-0014: Canonical Invoice Schema as lingua franca

**Status:** Accepted  
**Date:** 2026-04-25  
**Deciders:** Complai Engineering

## Context

Invoices enter Complai from multiple sources, each with a different format:

- **ERP push** (SAP RFC/BAPI, Oracle SuiteTalk, Tally ODBC, Dynamics Dataverse) -- vendor-specific schemas.
- **CSV/Excel upload** -- customer-defined column layouts.
- **Email ingestion** -- unstructured attachments requiring OCR extraction.
- **e-Invoice inbound** -- government-signed JSON with IRN and QR code.
- **Manual entry** -- user-created invoices in the UI.

Downstream systems (rules engine, GST filing, reconciliation, reporting, AP matching) each need invoice data. Without a canonical schema, every downstream consumer must handle N source formats, creating an N x M integration matrix.

## Decision

A single canonical invoice schema defined in Protobuf serves as the lingua franca for all invoice data within the platform. Every integration transforms source-specific formats into the canonical schema at the point of ingestion, before data enters the platform core.

- **Schema definition:** Protobuf in `packages/events/schemas/`. Code generation produces Go structs (for backend services) and TypeScript types (for frontend/BFF).
- **Core fields:** `tenant_id`, `pan`, `gstin`, `document_number`, `document_date`, `supply_type`, `document_type`, supplier/buyer details, line items (with HSN, quantity, unit price, tax breakdowns for CGST/SGST/IGST/cess), totals, payment details, e-Invoice references (IRN, EWB number), and metadata (source system, timestamps, tags).
- **Extension points:** a `metadata JSONB` field on each invoice and line item accommodates source-specific fields that do not map to the canonical schema.
- **Transformation responsibility:** each gateway/integration service (erp-gateway, ap-service, einvoice-service) implements a transformer that converts its source format into the canonical schema. Validation runs after transformation.

## Consequences

### Positive
- All downstream processing (rules engine, filing, reconciliation, reporting, AP matching) works with a single schema -- no per-source-format branching.
- Type-safe across Go and TypeScript via Protobuf code generation -- compile-time errors catch schema mismatches.
- Schema versioning via Protobuf's backward-compatible evolution rules (adding fields is safe; removing/renaming requires migration).
- Single validation pipeline: once an invoice passes canonical validation, all downstream consumers can trust the data shape.

### Negative
- Schema evolution requires coordinated migration across services that depend on generated types. Mitigated by Protobuf's backward-compatibility rules and CI checks for breaking changes.
- Some source-specific fields do not map cleanly to the canonical schema and must use the JSONB extension point, which is less type-safe.
- Initial transformation development for each source format requires understanding both the source schema and the canonical target.

### Risks
- Schema becoming a bottleneck for iteration speed if too many teams need to modify it simultaneously. Mitigated by: clear ownership (the platform team owns the canonical schema), a proposal process for changes, and JSONB extension points for source-specific data that does not need to be in the canonical schema.

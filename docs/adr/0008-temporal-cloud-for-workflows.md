# ADR-0008: Temporal Cloud (managed) for workflow orchestration

**Status:** Accepted  
**Date:** 2026-04-25  
**Deciders:** Complai Engineering

## Context

Complai has complex, long-running, multi-step workflows that require durable execution:

- **Filing sagas:** GSTR-1 filing involves data validation, government submission, polling for status, ARN capture, and notifications -- with compensation if any step fails.
- **Reconciliation:** pulling 2A/2B data, running the 5-stage match pipeline, executing IMS actions.
- **Bulk operations:** generating thousands of IRNs, E-Way Bills, or Form 16 documents in parallel.
- **AP automation:** invoice ingestion, OCR, 3-way matching, approval chains, payment execution.
- **Onboarding:** tenant setup, KYC verification, credential provisioning.

These workflows span minutes to days, must survive service restarts, and require saga-style compensation on failure. Building this on raw queues would mean reimplementing durable execution, state management, visibility, and retry logic.

Options evaluated: Temporal (self-hosted), Temporal Cloud (managed), AWS Step Functions, custom state machine on SQS.

## Decision

Temporal Cloud (managed) hosted in AWS Mumbai for all workflow orchestration. Go SDK for workflow and activity implementations.

- **Workflow namespaces:** `complai-filings`, `complai-reconciliation`, `complai-bulk`, `complai-ap`, `complai-onboarding`.
- **Deterministic workflow code** -- workflows only call activities, never direct I/O.
- **Activities are Go functions** with typed inputs and outputs.
- **Saga compensation** -- every workflow has explicit compensation logic; if Step N fails, Steps 1..N-1 are compensated.
- **Human-task integration** -- workflows emit `human_task` events; the UI picks them up and signals the workflow on user action.
- **Visibility** -- Temporal's built-in search indexed by tenant_id, GSTIN, period, and form type.

## Consequences

### Positive
- Zero operational burden -- Temporal Cloud manages history shards, visibility backend, cluster upgrades, and scaling.
- First-class Go SDK with native support for typed workflows and activities.
- Saga/compensation pattern is a built-in primitive, not something we build from scratch.
- Built-in workflow visibility, search, and debugging UI.
- Low latency to our EKS cluster (both in AWS Mumbai).

### Negative
- SaaS dependency -- if Temporal Cloud has an outage, all workflow execution pauses (though already-running activities complete).
- Monthly cost (~300 USD in Phase 1), increasing with workflow volume.
- Team needs to learn Temporal's programming model (determinism constraints, activity patterns).

### Risks
- Temporal Cloud pricing may increase significantly at scale. Mitigated by: migration to self-hosted Temporal requires zero code changes (same SDK, same API); only infrastructure configuration changes.

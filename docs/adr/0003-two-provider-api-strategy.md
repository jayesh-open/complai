# ADR-0003: Two-provider API strategy (Adaequare + Sandbox.co.in)

**Status:** Accepted  
**Date:** 2026-04-25  
**Deciders:** Complai Engineering

## Context

Complai covers a broad compliance surface: GST returns, e-Invoicing, E-Way Bills, TDS/TCS, ITR filing, KYC verification, and tax payments. No single API aggregator covers all these domains with production-grade quality. We evaluated:

1. **Single provider for everything** -- simpler integration, but no provider covers GST + TDS + ITR + KYC end-to-end.
2. **Two best-of-breed providers** -- Adaequare for GST ecosystem, Sandbox.co.in for income tax and KYC ecosystem.
3. **Direct government API integration** -- maximum control, but requires ASP/GSP licenses, crypto implementation, and years of certification effort.

## Decision

We use two primary API providers:

- **Adaequare (uGSP)** for GST returns, e-Invoice (IRP), and E-Way Bill. Adaequare holds the required GSTN ASP/GSP licenses and has proven reliability during filing peaks.
- **Sandbox.co.in (Quicko)** for TDS computation and filing, ITR filing, KYC verification (PAN, Aadhaar, GSTIN, bank account, MCA, Udyam), and tax payment challans.

Each provider is accessed through a dedicated gateway service that normalizes responses into our internal contract. Domain services never see provider-specific details.

## Consequences

### Positive
- Best-of-breed coverage per compliance domain -- each provider specializes in their area.
- Redundancy potential: if one provider has issues, we have architectural room to add alternatives behind the same gateway interface.
- Clear ownership boundaries: GST gateways own Adaequare integration; TDS/KYC gateways own Sandbox integration.

### Negative
- Two integration contracts to maintain, two authentication flows, two sets of API credentials per tenant.
- Two vendor relationships to manage (SLAs, support escalations, billing).
- Credential rotation and monitoring must be duplicated across providers.

### Risks
- Provider outages during filing deadlines. Mitigated by: circuit breakers (sony/gobreaker) on all gateway calls, retry with backoff, and Temporal workflow compensation for failed steps.
- API contract changes from either provider. Mitigated by: gateway services absorb all provider-specific changes; domain services are insulated.

# ADR-0002: Adaequare Enriched APIs only (no pass-through/SEK)

**Status:** Accepted  
**Date:** 2026-04-25  
**Deciders:** Complai Engineering

## Context

Adaequare, our GSTN ASP/GSP provider, offers two API modes:

1. **Pass-through APIs** -- raw GSTN endpoints. The caller must handle SEK (Session Encryption Key) decryption, AES-256-ECB encryption/decryption, HMAC-SHA256 signing, and all government-mandated crypto operations. This means maintaining cryptographic code, managing encryption keys, and keeping up with GSTN's periodic cipher changes.
2. **Enriched APIs** -- Adaequare handles all government-side encryption transparently. The caller sends and receives plain JSON. No crypto code needed in our codebase.

## Decision

We use only Adaequare's enriched APIs for all GST, e-Invoice (IRP), and E-Way Bill integrations. No SEK decryption, AES encryption, or HMAC signing code exists in our repository. Adaequare handles all government-mandated cryptographic operations on their infrastructure.

## Consequences

### Positive
- Dramatically simpler gateway code -- each gateway service is a thin HTTP client that sends/receives plain JSON.
- No cryptographic key management burden for GST API interactions.
- No need to track GSTN cipher changes or encryption specification updates.
- Faster development velocity for new filing integrations.
- Reduced attack surface -- no crypto code means no crypto bugs.

### Negative
- Higher per-API cost compared to pass-through mode (Adaequare charges a premium for enrichment).
- Vendor lock-in to Adaequare's enrichment layer -- switching to another GSP would require evaluating whether they offer equivalent enrichment or implementing pass-through crypto.
- We depend on Adaequare's enrichment uptime; if their enrichment layer has issues, we cannot fall back to pass-through without code changes.

### Risks
- If Adaequare discontinues enriched APIs or significantly raises pricing, we would need to implement the crypto layer. Mitigated by: the gateway abstraction pattern means only the gateway service internals would change; all upstream domain services are insulated.

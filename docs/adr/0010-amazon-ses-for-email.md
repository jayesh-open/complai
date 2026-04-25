# ADR-0010: Amazon SES for email

**Status:** Accepted  
**Date:** 2026-04-25  
**Deciders:** Complai Engineering

## Context

Complai sends transactional and bulk email across several workflows:

- Filing notifications (submission confirmation, ARN received, deadline reminders).
- Form 16 distribution to employees.
- Vendor onboarding and compliance status notifications.
- Daily/weekly compliance digest emails to CFOs and tax managers.
- AP invoice ingestion via inbound email (vendors send invoices to per-tenant email addresses).

We need an email service that handles both outbound (transactional + bulk) and inbound (AP ingestion) at enterprise scale with proper deliverability management.

## Decision

Amazon SES for all outbound email and inbound email ingestion. Mailpit replaces SES in local development environments.

- **Outbound:** SES configured with DKIM, SPF, and DMARC for all sending domains (`complai.in`, `notifications.complai.in`, tenant-specific subdomains).
- **Inbound:** per-tenant AP email addresses (`ap-{tenant}@inbox.complai.in`). SES receives email, writes to S3, triggers Lambda, enqueues to `ap.ingestion.queue` for processing by `ap-service`.
- **Configuration sets** for delivery tracking (opens, clicks, bounces, complaints).
- **Bounce handling:** SNS event destination routes bounce/complaint notifications to `notification-service` for automatic suppression list management.

## Consequences

### Positive
- Cost-effective at scale (0.10 USD per 1,000 emails for outbound; inbound is free beyond minimal S3 costs).
- Native AWS integration -- IAM for auth, CloudWatch for metrics, SNS for event delivery.
- Good deliverability with proper DKIM/SPF/DMARC configuration.
- Supports inbound email natively, enabling the AP invoice ingestion workflow without a separate email service.
- Mailpit provides a clean local dev replacement with a UI for inspecting sent emails.

### Negative
- Production sending access requires an approval process (1-3 business days to move out of SES sandbox).
- Reputation management is our responsibility -- bounce rates above 5% or complaint rates above 0.1% risk account suspension.
- Template management in SES is limited; we handle templates in `notification-service` and send raw HTML/text.

### Risks
- SES account suspension due to reputation issues. Mitigated by: automatic bounce suppression, complaint feedback loop processing, and monitoring of reputation metrics with alerting thresholds.

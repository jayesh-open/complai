# ADR-0005: AWS as cloud provider, ap-south-1 primary

**Status:** Accepted  
**Date:** 2026-04-25  
**Deciders:** Complai Engineering

## Context

Complai handles sensitive Indian financial and tax data (PAN, Aadhaar, GST filings, TDS returns, bank details). The Digital Personal Data Protection (DPDP) Act requires data residency within India for personal data. We need a cloud provider with mature Indian regions, broad service availability, and compliance certifications.

Options evaluated: AWS, Azure, GCP.

## Decision

AWS is our cloud provider with `ap-south-1` (Mumbai) as the primary region and `ap-south-2` (Hyderabad) as the disaster recovery region. All customer data stays within India at all times.

- Primary workloads run in `ap-south-1` (Mumbai) -- the most mature AWS region in India with the broadest service availability.
- DR workloads run in `ap-south-2` (Hyderabad) as a warm standby with RDS cross-region read replica, S3 cross-region replication, and a minimal EKS cluster.
- No data replicates outside of India. S3 replication, RDS backups, and all other storage are confined to these two regions.

## Consequences

### Positive
- Full DPDP Act data-residency compliance -- all data stays within India by architecture, not just policy.
- `ap-south-1` has the broadest AWS service availability in India (EKS, RDS, ElastiCache, OpenSearch, SES, Secrets Manager, KMS all available).
- Two-region architecture provides genuine DR capability with 60-minute RTO and 5-minute RPO for Tier-0 services.
- Mature region with established AZs, peering, and interconnect infrastructure.

### Negative
- `ap-south-2` (Hyderabad) is a newer region with fewer services and potentially less capacity during large-scale outages.
- AWS pricing in Indian regions carries a premium compared to some competitors (though service maturity justifies it).
- Dual-region architecture adds cost (DR EKS cluster, cross-region replication, NAT Gateways in both regions).

### Risks
- A simultaneous outage of both Indian regions is extremely unlikely but would require manual failover to a non-Indian region, which conflicts with DPDP requirements. Mitigated by: each region uses 3 AZs; a full-region outage has never occurred in `ap-south-1`.

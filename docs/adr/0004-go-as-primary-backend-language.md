# ADR-0004: Go as primary backend language

**Status:** Accepted  
**Date:** 2026-04-25  
**Deciders:** Complai Engineering

## Context

Complai needs a backend language for ~28 services spanning domain logic, API gateways, and platform infrastructure. Requirements:

- Fast compilation and small container images (dozens of services in a monorepo).
- Excellent concurrency support (filing peaks hit 5,000 RPS with thousands of parallel workflows).
- Strong static typing for a compliance domain where correctness matters.
- Mature AWS SDK and Kubernetes ecosystem support.
- Sufficient talent pool for a startup hiring in India.

Languages evaluated: Go, Java/Kotlin (Spring Boot), Rust, TypeScript (Node/NestJS).

## Decision

Go 1.22 is the primary backend language for all domain services and gateway services. Python 3.12 is used exclusively for AI/ML services (matching-ml, llm-copilot, ocr). TypeScript is used for frontend and BFF services only. No polyglot sprawl within a layer.

Standard library choices are pinned:
- HTTP: go-chi/chi v5
- DB: jackc/pgx v5 + sqlc + goose
- DI: uber-go/fx
- Config: spf13/viper
- Logging: rs/zerolog
- Validation: go-playground/validator v10
- Tracing: OpenTelemetry Go SDK
- Testing: testify + testcontainers-go
- Money: shopspring/decimal

## Consequences

### Positive
- Fast compile times enable rapid CI across 20+ services in a monorepo.
- Small binary sizes (~15-20MB) keep container images lean, reducing ECR storage and pull times.
- Goroutines and channels handle concurrent filing workloads naturally without thread pool tuning.
- Strong static typing catches compliance logic errors at compile time.
- First-class AWS SDK (aws-sdk-go-v2), Temporal SDK, and Kubernetes client support.
- Go workspaces (`go.work`) enable monorepo-wide dependency management.

### Negative
- Verbose error handling (`if err != nil`) adds boilerplate compared to exception-based languages.
- Less expressive for complex domain modeling compared to languages with algebraic data types or pattern matching.
- Generics (Go 1.18+) are still maturing; some patterns require interface-based abstractions.

### Risks
- Team members from Java/Node backgrounds need onboarding time. Mitigated by: consistent library choices reduce decision fatigue; shared-kernel-go provides canonical patterns.

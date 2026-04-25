# ADR-0015: Monorepo with Go workspaces + pnpm + Turborepo

**Status:** Accepted  
**Date:** 2026-04-25  
**Deciders:** Complai Engineering

## Context

Complai consists of ~28 services (Go, Python, Node), 4 frontend apps (Next.js), and several shared packages. The team needs:

- Atomic cross-service changes (e.g., a Protobuf schema change + all consumers updated in one commit).
- Consistent tooling, linting, and CI configuration across all services.
- Shared libraries (shared-kernel-go, shared-kernel-node, ui-components) consumed without publishing to a registry.
- A single onboarding experience for new engineers.

Options evaluated:

1. **Monorepo** -- all code in one Git repository with workspace tooling.
2. **Polyrepo** -- one repository per service, shared libraries published to a private registry.
3. **Hybrid** -- frontend monorepo + backend polyrepo (or vice versa).

## Decision

Single Git monorepo with Go workspaces (`go.work`) for all Go modules and pnpm workspaces + Turborepo for all TypeScript/Node packages.

- **Go workspace:** `go.work` at the repo root references all Go service modules and `packages/shared-kernel-go`. Local `replace` directives are unnecessary; the workspace handles cross-module resolution.
- **pnpm workspaces:** `pnpm-workspace.yaml` covers `apps/*`, `services/node/*`, `packages/shared-kernel-node`, `packages/ui-components`, and `packages/events`.
- **Turborepo:** orchestrates TypeScript/Node build, lint, and test tasks with caching and affected-service detection.
- **CI:** GitHub Actions with a build matrix that detects which services are affected by a commit (via path filters) and runs only relevant builds and tests.
- **Standard service layout:** every Go service follows the same directory structure (`cmd/server/main.go`, `internal/{api,app,domain,infra,config}`, `migrations/`, `test/`, `Dockerfile`, `Makefile`).

## Consequences

### Positive
- Atomic cross-service changes: a breaking Protobuf schema change and all consumer updates land in a single commit and PR.
- Shared libraries are consumed directly from the workspace -- no publish-consume cycle, no version drift.
- Consistent tooling: one `golangci-lint` config, one `tsconfig.base.json`, one `commitlint.config.js` for the entire codebase.
- Easier onboarding: new engineers clone one repo and have the full system.
- CI can run integration tests that span multiple services in a single pipeline.

### Negative
- Repository size grows over time as more services and assets are added. Git operations (clone, status, diff) may slow at scale.
- CI must be smart about affected-service detection to avoid running all ~28 service builds on every commit.
- IDE performance may degrade with very large workspaces (mitigated by project-scoped IDE configurations).

### Risks
- Git performance degradation at scale. Mitigated by: sparse checkout for contributors who work on a subset of services, Git LFS for binary assets, and shallow clones in CI.
- CI pipeline complexity as service count grows. Mitigated by: path-based filtering in GitHub Actions, Turborepo caching for TypeScript, and Go's fast compilation.

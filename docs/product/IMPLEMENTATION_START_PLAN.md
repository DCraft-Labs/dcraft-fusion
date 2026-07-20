# Implementation Start Plan

Implementation start date: 2026-05-18

## Stack Guardrails

Phase 1 starts with a Go backend kernel and TypeScript UI/domain foundation. The first revenue wedge is UI-first, but the control plane itself is Go-first.

Target runtime:

- Go 1.26.x for core control-plane backend services.
- Node.js 24 LTS for TypeScript packages and React UI.
- Python remains reserved for CDC workers, AI/metadata workers, and the imported Fusion CDC Engine where it already exists.

Local note:

- This workstation currently has Node 20.15. The repo is pinned to Node 24 LTS, but the first UI code avoids Node 24-only syntax so tests can run during bootstrap.
- Go is installed separately through the terminal when the backend kernel is built.

## Phase 1 Subgoals

1. Go backend kernel: organization, tenant, project, environment, membership, deployment profile, request context, audit event.
2. Test enforcement: tenant context is mandatory, access cannot cross organization/tenant/project boundaries, audit events are emitted for writes.
3. UI foundation: organization/project switcher model, Command Center shell, Integration Hub shell.
4. Integration foundation: connector definitions and secret references, not raw secret storage.
5. Deployment foundation: pooled SaaS, dedicated DB/schema, dedicated stack, customer VPC/BYOC, future on-prem.

## Do Not Drift

- Do not make CDC the main product.
- Do not store raw data by default.
- Do not implement autonomous AI execution.
- Do not assume one user belongs to one tenant.
- Do not assume one organization runs in shared SaaS.

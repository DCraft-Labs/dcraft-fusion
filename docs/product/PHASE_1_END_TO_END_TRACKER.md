# Phase 1 End-To-End Tracker

Date: 2026-05-18

## Review Verdict

Phase 1 is now implemented as a working local product and Docker Desktop Kubernetes stack with:

- real OIDC-mode runtime
- credential-based login feeding the OIDC authorization-code flow
- routed tenant and superadmin web surfaces
- Postgres persistence
- Redis-backed async job state
- Kubernetes HA baseline
- backup verification
- metrics exposure

That said, **Phase 1 is not the same thing as enterprise production GA**. The product surface is there, but some enterprise integrations are still baseline or partial and should not be oversold.

## Current Access Model

Local runtime is no longer the old click-to-enter dev access. It now uses a seeded credential login that creates a server-side login session and then completes the OIDC code flow.

### Seeded local credentials

| Role | Actor ID | Email | Password | Scope |
| --- | --- | --- | --- | --- |
| Founder Admin | `user-founder` | `founder@dcraftlabs.com` | `(env seed)` | Tenant workspace |
| Platform Superadmin | `user-superadmin` | `superadmin@dcraftlabs.com` | `(env seed)` | Platform administration |

## Running URLs

| Surface | URL |
| --- | --- |
| Docker Compose web | `http://127.0.0.1:5174` |
| Docker Compose API | `http://127.0.0.1:8080` |
| Kubernetes web/API | `http://127.0.0.1:30173` |
| Kubernetes metrics | `http://127.0.0.1:30173/metrics` |

## Web Surface

### Public routes

| Route | Status | Notes |
| --- | --- | --- |
| `/` | Complete | Public landing page with light/dark mode, rollout framing, and secure-login CTA. |
| `/login` | Complete | Email/password login form, seeded credential helpers, and OIDC runtime details. |
| `/auth/callback` | Complete | Exchanges authorization code for signed token and establishes browser session. |

### Tenant routes

| Route | Status | Backing data |
| --- | --- | --- |
| `/workspace/command-center` | Complete | `GET /api/v1/bootstrap` |
| `/workspace/integrations` | Complete for Phase 1 | `POST /api/v1/connections`, `POST /api/v1/connections/{id}/test` |
| `/workspace/metadata` | Complete for Postgres | `GET /api/v1/datasets`, discovery results |
| `/workspace/runs` | Seeded/read-only | Seeded run visibility only |
| `/workspace/audit` | Complete baseline | `GET /api/v1/audit-events` |
| `/workspace/settings` | Complete baseline | Bootstrap policy and deployment context |
| `/workspace/ai` | Complete draft-only | `GET /api/v1/ai/recommendations` |

### Superadmin routes

| Route | Status | Backing data |
| --- | --- | --- |
| `/superadmin/overview` | Complete | `GET /api/v1/platform/overview` |
| `/superadmin/organizations` | Complete | Create org, tenant, and project flows |
| `/superadmin/identity` | Complete baseline | Identity provider inventory and posture |
| `/superadmin/secrets` | Complete baseline | Secret provider inventory and posture |
| `/superadmin/adapters` | Complete baseline | Adapter inventory and readiness |
| `/superadmin/operations` | Complete baseline | Backup, HA, jobs, and observability posture |

## Backend And API Coverage

| Area | Status | Notes |
| --- | --- | --- |
| OIDC login bootstrap | Complete | `POST /api/v1/auth/login` creates server-side login session. |
| OIDC logout | Complete | `POST /api/v1/auth/logout` clears login session cookie. |
| OIDC authorize/token/jwks | Complete | Local provider supports code flow and RS256-signed JWTs. |
| Bearer token enforcement | Complete | Middleware verifies issuer, audience, expiry, and signature. |
| Organization/tenant/project CRUD | Complete baseline | Persists to Postgres when configured. |
| Bootstrap API | Complete | Main tenant UI bootstrap. |
| Connections | Complete baseline | Secret references only; raw secrets rejected. |
| Connection test | Complete baseline | Persists job state and audit. |
| Metadata discovery | Complete for Postgres | Real `information_schema` discovery. |
| Datasets | Complete | Tenant/project scoped reads. |
| Runs | Partial | Seeded read-only run visibility only. |
| Policies | Complete baseline | Visibility only. |
| AI recommendations | Complete draft-only | No autonomous writes. |
| Platform overview | Complete | Superadmin-only overview API. |

## Infrastructure Coverage

| Area | Status | Evidence |
| --- | --- | --- |
| Docker Compose | Complete | Postgres 18, Redis 8, Go API, web. |
| Docker Desktop Kubernetes | Complete | Namespace, services, deployments, PVCs, NodePort. |
| Postgres persistence | Complete | Core repositories and Phase 1 resources persist. |
| Redis job state | Complete | Async job keys written for test/discovery flows. |
| HA baseline | Complete for local cluster | 2 API replicas, 2 web replicas, HPAs, PDBs. |
| Backup verification | Complete baseline | `postgres-backup-verify` CronJob present and verified. |
| Metrics | Complete baseline | `/metrics` exposed. |
| Alerts/dashboard scaffold | Complete baseline | Files present; customer monitoring backend still needed. |
| CI | Complete baseline | Go tests, web tests, typecheck, build, audit, k8s render. |

## What Was Verified Live

| Check | Result |
| --- | --- |
| `go test ./...` | Passing |
| `npm test` | Passing |
| `npm run typecheck` | Passing |
| `npm run build` | Passing |
| Docker Compose rebuild | Passing |
| Kubernetes rollout | Passing |
| Public landing page | Verified in browser |
| Light/dark mode | Verified in browser |
| Founder login | Verified in browser |
| Workspace route changes | Verified in browser |
| Superadmin login | Verified in browser |
| Superadmin route changes | Verified in browser |

## Must Not Be Overclaimed

- This is not a Fusion password product. Local credentials are seeded for local validation only.
- Direct enterprise IdP onboarding is not complete until a real external OIDC provider is configured and exercised.
- Direct SAML ACS/login is not implemented; the current design path is SAML-to-OIDC federation.
- Airflow, Dagster, and dbt are not yet full live ingestion adapters.
- External secret provider SDK resolution is not yet implemented end to end.
- Tenant APIs still need full membership authorization checks, not just trusted scoped claims.
- Audit listing still needs tenant-safe filtering before paid multi-customer rollout.
- Full multi-workspace switching in the web app is not implemented yet.

## Remaining Production-Level Gaps

These are the real blockers between current Phase 1 and a production-grade enterprise launch.

| Priority | Gap | Why it still matters |
| --- | --- | --- |
| P0 | Membership authorization on every tenant API | A valid token must not be enough by itself to traverse customer scope. |
| P0 | Tenant-safe audit filtering | Audit data cannot be globally listable across customers. |
| P0 | Real external OIDC provider integration | Current local OIDC provider proves the product path, not customer IdP onboarding. |
| P0 | Real Airflow read-only ingestion | Workflow Runs still rely on seeded data. |
| P1 | Real Dagster/dbt ingestion | Needed for broader orchestration and transformation visibility. |
| P1 | Real secret-provider resolution layer | Required for Vault/AWS/GCP/Azure production deployments. |
| P1 | Workspace/org/tenant switcher | Needed for actual operators managing multiple customers or projects. |
| P1 | Observability export wiring | Metrics are exposed, but customer telemetry backends still need wiring and dashboards applied. |

## Founder Read

Phase 1 is now strong enough for:

- founder and design-partner demos
- private alpha validation
- platform and tenant workflow walkthroughs
- infrastructure and deployment model conversations

It is **not yet honest to call this enterprise production complete** until the remaining P0 and P1 hardening above is closed.

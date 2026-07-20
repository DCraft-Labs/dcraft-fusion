# Phase 1 Execution Plan

Research date: 2026-05-18

## Phase 1 Goal

Build a UI-first private alpha that proves DCraft Fusion can support real customers, not just a prototype.

Phase 1 must prove four things:

1. Fusion can onboard organizations and users through the UI.
2. Fusion can connect to real external systems through UI-based integrations.
3. Fusion can support B2B, B2C, and B2B2C tenancy models from day one.
4. Fusion can run in shared SaaS, dedicated tenant infrastructure, and customer-owned environments without redesigning the product later.

## Day-One Tenancy Principle

Do not design "tenant" as a flat field.

Fusion must support this hierarchy from the first domain model:

```text
Platform
  Organization / Customer Account
    Tenant / Business Unit / Brand / Region / Client
      Project / Workspace / Environment
        Resources / Connections / Workflows / Runs / Policies
```

This covers:

- **B2B**: One company buys Fusion for employees and data teams.
- **B2C**: A business uses Fusion-powered workflows for many end customers.
- **B2B2C**: A company uses Fusion for its own customers, merchants, partners, branches, brands, or regions.
- **Internal platform**: One enterprise has many departments or business units.
- **Managed service provider**: One operator manages multiple customer accounts.

## Day-One Deployment Principle

Fusion must support multiple isolation modes through one product architecture.

| Mode | Who It Serves | Isolation | Phase 1 Requirement |
|---|---|---|---|
| Shared SaaS / pooled | Startups, SMB, lower mid-market | Shared app and shared database with tenant isolation | Must be supported |
| Shared app, dedicated database/schema | Mid-market, privacy-conscious customers | Shared control plane, stronger data isolation | Design-ready in Phase 1 |
| Dedicated tenant stack | Enterprise customers | Separate app/runtime/database per customer | Deployment template planned in Phase 1 |
| Customer VPC / BYOC | Regulated customers | Runs in customer's cloud account | Architecture-ready in Phase 1 |
| On-prem / air-gapped | Highly regulated environments | Customer-managed infrastructure | Not built in Phase 1, but no blocker in design |

AWS SaaS guidance describes pool, silo, bridge, and tier-based isolation models. Kubernetes guidance also separates softer shared-cluster isolation from harder models such as virtual control planes or dedicated clusters. Fusion should adopt that hybrid thinking immediately.

## Phase 1 Product Scope

Phase 1 is **Private Alpha: UI-first command center with real integration foundations**.

It should include:

- Web UI shell.
- Organization onboarding.
- Tenant/project/workspace model.
- Auth and membership model.
- Integration Hub.
- Read-only metadata discovery for one or two source systems.
- Read-only run ingestion from one workflow/orchestration tool.
- Basic audit log.
- Tenant-aware settings.
- Deployment profile registry.
- AI explain/summarize mode only.

Do not ship production CDC as the main wedge in Phase 1. Show CDC as a visible future/preview integration card and optionally a local sandbox demo.

## Phase 1 Workstreams

### 1. Domain And Tenancy Foundation

Tasks:

- Define canonical domain objects: `Organization`, `Tenant`, `Project`, `Environment`, `User`, `Membership`, `Role`, `ServiceAccount`.
- Define organization type: `business`, `consumer_app`, `managed_service_provider`, `internal_platform`.
- Define tenant type: `business_unit`, `brand`, `region`, `client`, `consumer_segment`, `environment`, `partner`.
- Add parent-child tenant support.
- Add active context resolution: organization, tenant, project, environment.
- Require tenant/project context on every resource.
- Add tenant-scoped feature flags.
- Add tenant-scoped quotas.
- Add tenant-scoped audit trail.

Acceptance criteria:

- A user can belong to multiple organizations.
- A user can switch active organization and project in the UI.
- A B2B customer can create departments/teams as tenants.
- A B2B2C customer can create its own customer/merchant/brand tenants.
- Backend rejects resource access without matching organization and tenant context.

### 2. Identity, Access, And Membership

Tasks:

- Implement login foundation with OIDC-compatible architecture.
- Support organization memberships.
- Support user roles at organization, tenant, and project level.
- Support service accounts for integrations.
- Support invitation flow.
- Support admin-created users for private alpha.
- Add SSO/SAML/OIDC as design-ready, even if not fully self-serve in Phase 1.
- Model B2C/CIAM users separately from B2B workforce users where needed.

Acceptance criteria:

- Every request has actor, organization, tenant, project, role, and correlation context.
- Admin can invite a user into an organization.
- Same email/user can be a member of multiple organizations.
- API tokens/service accounts are scoped to organization and project.
- UI never assumes a single tenant per user.

### 3. Deployment Profile Registry

Tasks:

- Add deployment profile model: `pooled_saas`, `dedicated_db`, `dedicated_stack`, `customer_vpc`, `on_prem`.
- Add deployment profile to organization settings.
- Add environment model: `dev`, `stage`, `prod`, `sandbox`.
- Track data residency region.
- Track control-plane location.
- Track data-plane location.
- Track whether raw data movement is allowed.
- Track customer-owned secret manager references.

Acceptance criteria:

- UI shows deployment mode and region for each organization.
- Backend can decide whether a workflow/resource belongs to shared or dedicated execution context.
- No code path assumes all tenants run in the same infrastructure.
- BYOC/on-prem are visible as target profiles even if not provisioned automatically in Phase 1.

### 4. UI Shell And Navigation

Tasks:

- Build main app shell in `apps/web`.
- Add organization switcher.
- Add tenant/project/environment switcher.
- Add sidebar navigation.
- Add command center dashboard.
- Add Integration Hub.
- Add Metadata Explorer placeholder.
- Add Workflow/Run Center placeholder.
- Add Audit Center.
- Add Settings area.

Acceptance criteria:

- User can see active organization, tenant, project, and environment at all times.
- Every page loads with loading, empty, error, and permission-denied states.
- Integration Hub is usable from day one.
- UI never exposes raw secrets.

### 5. Integration Hub

Tasks:

- Build connector catalog UI.
- Add connector cards for Postgres, Snowflake or BigQuery, dbt, Airflow or Dagster, Slack/Teams, GitHub, Superset/Metabase, and Fusion CDC Engine preview.
- Add connection creation wizard.
- Add connection test flow.
- Store credentials as secret references only.
- Add integration status and last sync timestamp.
- Add integration ownership and tenant/project scope.

Acceptance criteria:

- User can create a read-only source connection from the UI.
- User can test a connection from the UI.
- User can see integration health.
- User can connect one workflow/run source in read-only mode.
- User can see why an integration is unavailable based on deployment profile, permissions, or missing secret reference.

### 6. Metadata Discovery

Tasks:

- Discover schemas/tables/views for one warehouse/database source.
- Store metadata as Fusion resources.
- Track owner, tags, description, freshness timestamp, and source system.
- Add basic search.
- Add resource detail page.
- Record discovery events in audit log.

Acceptance criteria:

- User can run discovery from UI.
- User can search discovered datasets.
- User can view dataset source, tenant/project, and last discovered time.
- Backend enforces tenant isolation in metadata queries.

### 7. Workflow And Run Visibility

Tasks:

- Ingest read-only workflow/run data from one tool, preferably Airflow, Dagster, Prefect, or dbt artifacts.
- Create canonical run model.
- Create run timeline UI.
- Add status: queued, running, succeeded, failed, canceled, unknown.
- Add correlation IDs.
- Add artifact/log references, not raw log dumps by default.

Acceptance criteria:

- User can see recent runs in Command Center.
- User can filter runs by organization, tenant, project, environment, workflow, and status.
- User can open a run detail page.
- Failed run appears in dashboard and audit/incident timeline.

### 8. Audit And Compliance Foundation

Tasks:

- Create immutable audit event model.
- Record user, organization, tenant, project, action, resource, timestamp, source IP/context, and correlation ID.
- Record integration create/test/discover events.
- Record settings changes.
- Add audit search/filter UI.
- Add basic export design.

Acceptance criteria:

- Every write action emits audit.
- Audit queries are tenant-isolated.
- Admin can filter audit events by actor/resource/action/time.
- Audit events do not contain secrets.

### 9. Policy And Approval Foundation

Tasks:

- Model policy types: deployment policy, integration policy, data movement policy, AI action policy.
- Model approval states: draft, pending_review, approved, rejected, expired, executed.
- Add UI placeholders for approval queue.
- Require approval for any dangerous action, even if Phase 1 dangerous actions are disabled.

Acceptance criteria:

- AI-generated recommendations are drafts.
- Any future execution action has an approval gate.
- Approval model is tenant-aware.

### 10. AI Assist Read-Only Mode

Tasks:

- Add AI Operator panel.
- Allow AI to summarize metadata, explain failed runs, and draft recommendations.
- Store AI recommendations as draft resources.
- Add source/context attribution for recommendations.
- Add approve/reject UI state without executing actions in Phase 1.

Acceptance criteria:

- AI cannot execute production actions.
- AI output is visibly marked as generated.
- AI recommendation includes target organization, tenant, project, and resource context.

### 11. Fusion CDC Engine Evaluation Integration

Tasks:

- Keep `engines/fusion-cdc-engine` as a module source base.
- Add it to Integration Hub as a preview engine.
- Map its current model to Fusion concepts.
- Plan rename/migration of legacy tenancy fields into generic Fusion terms.
- Extract reusable event envelope ideas into `packages/event-envelope`.
- Extract worker health model into Fusion observability model.

Acceptance criteria:

- CDC engine is visible as an optional engine, not the core product.
- Phase 1 docs explain how it will be controlled by Fusion.
- No UI copy presents CDC as the entire product.

### 12. Local Development And Delivery

Tasks:

- Create local dev profile.
- Add seed data for one B2B org, one B2B2C org, and one internal platform org.
- Add fake/mock connectors for UI development.
- Add basic CI checks.
- Add API contract skeletons.
- Add tenant isolation tests.

Acceptance criteria:

- Developer can run the UI locally.
- Seed data demonstrates multiple tenancy models.
- Tests prove that one tenant cannot read another tenant's resources.

## Phase 1 Demo Scenarios

### Demo 1: B2B Data Team

Company: Acme Analytics

Flow:

1. Admin creates organization.
2. Admin creates projects for finance and growth.
3. Admin connects Postgres or warehouse read-only.
4. User discovers metadata.
5. User sees datasets and recent orchestration runs.
6. Admin views audit events.

### Demo 2: B2B2C Platform Business

Company: Commerce Platform Co.

Flow:

1. Organization creates tenants for Brand A and Brand B.
2. Each tenant has its own project/environment.
3. User connects separate sources per tenant.
4. Metadata and runs stay isolated.
5. Admin can see cross-tenant summary.
6. Tenant user can only see its own scope.

### Demo 3: Dedicated Deployment Buyer

Company: Regulated Enterprise

Flow:

1. Admin selects deployment profile: customer VPC.
2. UI shows data-plane location and raw-data movement policy.
3. Secret references point to customer-managed secret store.
4. Integration Hub shows which connectors are allowed.
5. Audit shows deployment/settings changes.

## Phase 1 Non-Negotiable Tests

- Tenant isolation test for every query.
- Authorization test for every write endpoint.
- Audit event test for every write endpoint.
- Secret redaction test for connection APIs.
- UI context test: every page must have active organization/project.
- Integration test: connection metadata cannot leak across tenants.
- AI safety test: generated action remains draft unless approved.

## Phase 1 Out Of Scope

- Full production CDC.
- Full BI platform.
- Full dbt/Airflow replacement.
- Full on-prem installer.
- Full self-serve SSO marketplace.
- Raw customer data storage by default.
- Autonomous AI execution.

## Founder Gate For Phase 1

Phase 1 is successful only if:

- Three design partners can connect at least one real metadata source.
- At least one design partner has a B2B2C or nested-tenant use case.
- At least one design partner needs dedicated or customer-owned deployment later.
- Users return to the UI without founder hand-holding.
- The team can demonstrate tenant isolation, audit, and deployment-profile awareness.

## Research Notes

Key external architecture inputs:

- AWS SaaS guidance distinguishes pooled, siloed, bridge, and tier-based tenant isolation models.
- AWS also notes that full-stack silos can still be SaaS when surrounded by shared identity, onboarding, metrics, deployment, analytics, and operations.
- Kubernetes documentation warns that shared clusters require careful isolation, fairness, and noisy-neighbor controls; hard multi-tenancy can require virtual control planes or dedicated clusters.
- Google Cloud GKE guidance says multi-tenant cluster planning should consider isolation at cluster, namespace, node, pod, and container layers.
- Azure SQL SaaS guidance documents multiple database tenancy patterns, including multitenant databases and database-per-tenant models.
- Auth0 Organizations, WorkOS Organizations, Clerk Organizations, and Microsoft Entra External ID all reinforce that modern SaaS identity must support organization context, memberships, roles, and external users.

Sources:

- https://docs.aws.amazon.com/whitepapers/latest/saas-tenant-isolation-strategies/core-isolation-concepts.html
- https://docs.aws.amazon.com/whitepapers/latest/saas-tenant-isolation-strategies/silo-isolation.html
- https://docs.aws.amazon.com/whitepapers/latest/saas-architecture-fundamentals/full-stack-silo-and-pool.html
- https://kubernetes.io/docs/concepts/security/multi-tenancy/
- https://cloud.google.com/kubernetes-engine/docs/concepts/multitenancy-overview
- https://learn.microsoft.com/en-us/azure/architecture/guide/multitenant/service/sql-database
- https://auth0.com/docs/design/using-auth0-with-multi-tenant-apps
- https://auth0.com/docs/organizations
- https://workos.com/blog/developers-guide-saas-multi-tenant-architecture
- https://clerk.com/docs/guides/organizations/overview
- https://learn.microsoft.com/en-us/entra/external-id/index-b2b

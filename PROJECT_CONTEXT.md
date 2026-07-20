# DCraft Fusion Project Context

DCraft Fusion is the control plane OS for modern data platforms. It gives teams one governed operating layer across ingestion, orchestration, transformation, quality, lineage, catalog, governance, cost, BI/reporting, and AI-assisted operations.

## Core Principle

Fusion is the brain/control plane. Existing tools are the muscles/execution engines.

Fusion owns:

- Metadata and resource intent.
- Tenant, project, workflow, run, policy, audit, approval, and adapter state.
- Provider and adapter contracts.
- Operational visibility and recommendations.
- Safe AI-assisted drafts.

Fusion must not own by default:

- Raw customer table data.
- Stream payloads.
- Customer compute execution.
- Customer warehouses, lakes, or orchestration engines.

## MVP Boundary

Build the control plane kernel before the full 16-layer platform:

- Tenant and project model.
- Connection registry.
- Workflow registry.
- Run registry.
- Audit service.
- Adapter registry.
- One real adapter, preferably dbt or Airflow first.
- UI command center.
- AI assist in recommendation/draft mode only.

Day-one tenancy and deployment requirements:

- Support B2B, B2C, B2B2C, internal platform, and managed-service-provider use cases.
- Model organization, tenant, project, and environment separately.
- Support shared SaaS, dedicated database/schema, dedicated stack, customer VPC/BYOC, and future on-prem deployment profiles.
- Never assume one user belongs to one tenant or one customer runs in one shared infrastructure mode.

## Non-Negotiable Engineering Rules

- Every write must include tenant and actor context.
- Do not hardcode tenant IDs or secrets.
- Do not mix tenant data.
- Do not log or return secrets.
- Enforce tenant filters in backend repositories.
- Emit audit events for significant actions.
- Use typed contracts for APIs and adapters.
- Add tests for domain logic, authorization, tenant isolation, and adapter contracts.

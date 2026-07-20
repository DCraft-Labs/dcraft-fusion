# Fusion CDC Engine Evaluation

Evaluation date: 2026-05-18

## What Was Found

The imported CDC project is a substantial implementation, not just notes. It includes:

- FastAPI control plane.
- React/Vite frontend.
- CDC workers for MySQL, PostgreSQL, MongoDB, and polling fallback.
- Redis stream transport.
- Spark consumer.
- Transform worker.
- Data quality, schema evolution, monitoring, alerting, DLQ, and UDF surfaces.
- Alembic/Flyway migrations and SQL schemas.
- Docker, Kubernetes, Helm, monitoring, and tests.

It has been moved to:

```text
engines/fusion-cdc-engine/
```

## Where It Sits In DCraft Fusion

The module should sit as a **Fusion execution/data-plane engine**, not as the main product.

Recommended ownership:

```text
engines/fusion-cdc-engine/          Existing imported CDC implementation
services/cdc-control-service/       Future thin Fusion-native service boundary
packages/event-envelope/            Shared event envelope and run/event contracts
packages/adapter-sdk/               Adapter hooks for controlling CDC engine actions
apps/web/                           Unified Fusion UI that embeds CDC workflows
infra/kubernetes|helm/              Eventually shared deployment packaging
```

Fusion should expose CDC through the same control-plane model as other engines:

- Connections.
- Workflows.
- Runs.
- Policies.
- Audit events.
- Approvals.
- Lineage.
- Incidents.
- AI recommendations.

The CDC engine should not become a separate branded product inside the repo.

## How It Adds Value

The module gives Fusion a real execution story beyond passive cataloging.

Strategic value:

- Proves Fusion can control a real data-plane engine.
- Differentiates Fusion from passive metadata catalogs.
- Supports real-time and scheduled data movement use cases.
- Creates strong demos: register source, select tables, stream changes, track lag, detect schema drift, show audit, show lineage, explain incident.
- Gives regulated/hybrid customers a data-plane option that can run in their own infrastructure.

Product value:

- Connection onboarding UI.
- Source and destination management.
- CDC stream monitoring.
- Lag/error visibility.
- Data quality checks.
- Transform configuration.
- DLQ handling.
- Schema evolution controls.
- Worker health and metrics.

Revenue value:

- Can become a paid add-on after the control plane has adoption.
- Supports enterprise/BYOC pricing.
- Creates professional-services opportunities for regulated/hybrid deployments.
- Helps sell Fusion as more than a dashboard once the core command center is working.

## Founder-Level Assessment

This module is valuable, but it should not be the day-one wedge.

Reason:

- CDC correctness is hard.
- Buyers will compare it against mature ingestion vendors.
- Running a CDC engine requires operational confidence.
- The current project still has implementation vocabulary and tenancy assumptions that must be normalized into Fusion's tenant/project/resource model.

Best path:

1. Use it as a **Phase 2/3 demo integration** after the Fusion Command Center exists.
2. Productize it as a **Phase 4/5 paid engine module** only when customers explicitly ask for Fusion-managed CDC.
3. Do not let it distract from the first revenue wedge: UI-first command center, integrations, metadata, run visibility, audit, approvals, and AI drafts.

## Required Refactoring Before Productization

Before this can become a Fusion module, normalize:

- Tenancy fields into Fusion organization/tenant/project language.
- Audit events into the global Fusion audit model.
- Run state into the global Fusion run registry.
- Connector definitions into the Fusion adapter/provider registry.
- Event envelope into `packages/event-envelope`.
- Secrets into Fusion secret references.
- Frontend routes into `apps/web`, not a separate product UI.
- Kubernetes/Helm naming into shared Fusion deployment conventions.

## Rollout Recommendation

### Now

Keep the implementation in `engines/fusion-cdc-engine/` as an engine module reference and source base.

Use it to inform:

- CDC adapter contracts.
- Event envelope design.
- Worker health model.
- Data quality and DLQ UI.
- Source/destination onboarding flows.

### Phase 1

Do not ship full CDC yet. Ship only UI surfaces that can show the concept:

- Integration Hub CDC connector card.
- Read-only source/destination metadata.
- Mock or sample CDC run timeline.
- CDC health panel.

### Phase 2

Integrate one read-only or controlled CDC demo:

- Register source.
- Discover tables.
- Start a controlled local/sandbox stream.
- Show run/lag/errors in Fusion UI.
- Emit audit events.

### Phase 3+

Turn it into a paid optional module:

- BYOC CDC runtime.
- Managed worker control.
- Enterprise monitoring.
- DLQ replay.
- Schema evolution approval.
- Data quality gates.

## Final Conclusion

Fusion CDC Engine should be treated as a high-value engine module, not the core identity of DCraft Fusion.

It strengthens the product because it gives Fusion a credible real-time data-plane capability. But the company should still sell the broader control plane first. The right message is:

> Fusion can govern and operate your existing stack. When you need real-time CDC under the same control plane, Fusion CDC Engine is one available engine.


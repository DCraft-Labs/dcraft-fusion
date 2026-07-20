# What We Are Building

DCraft Fusion is an AI-assisted control plane for modern data platforms.

It is not a warehouse, lakehouse, orchestration engine, BI engine, stream processor, or transformation runtime by default. It is the governed operating layer above those systems.

## The Core Idea

Data teams already use many tools: warehouses, lakes, dbt, Airflow, Spark, Kafka, BI tools, metadata catalogs, quality systems, and governance tools. Fusion gives those teams one place to define intent, track state, enforce policy, observe operations, and use AI safely.

The simple model:

- Fusion is the brain.
- Existing tools are the muscles.
- Adapters/providers connect the brain to the muscles.

## What Fusion Owns

- Tenants and projects.
- Connections and credential references.
- Resource definitions and desired state.
- Workflow definitions and versions.
- Run state, task state, logs, timelines, and artifact references.
- Adapter/provider registrations and capability metadata.
- Policies, approvals, and audit events.
- Metadata discovery, lineage, drift detection, and reconciliation state.
- AI-generated recommendations and drafts.

## What Fusion Does Not Own By Default

- Raw customer data.
- Table contents.
- Stream payloads.
- Compute execution.
- Customer warehouses, lakes, BI systems, or orchestration engines.

## MVP Shape

The first product should prove the control plane kernel:

1. Create tenant and project context.
2. Register a connection.
3. Discover metadata from that connection.
4. Define a simple workflow.
5. Trigger the workflow through one adapter, preferably dbt or Airflow.
6. Track run and task state.
7. Emit immutable audit events.
8. Show an operator UI with connection, workflow, run, and audit visibility.
9. Let AI propose a recommendation or change draft.
10. Require approval before any AI-generated action executes.

## First Implementation Boundaries

Start with a modular monolith or small number of services if that helps move faster. Preserve the domain boundaries in `services/`, but do not split into many deployable microservices before there is real operational pressure.

The first shared packages should likely be:

- `packages/domain`
- `packages/event-envelope`
- `packages/authz`
- `packages/adapter-sdk`
- `packages/telemetry`

The first app should be:

- `apps/web`

The first adapter should be one of:

- `engines/dbt-adapter`
- `engines/airflow-adapter`

## High-Risk Areas

- Tenant isolation.
- Secrets handling.
- Adapter contract stability.
- Run-state correctness.
- Audit completeness.
- AI action safety.
- Scope creep into becoming an execution engine.

## Build Rule

Every feature should pass this test:

Does this help Fusion coordinate, govern, observe, or safely assist existing data systems without unnecessarily taking ownership of customer data or execution?

If yes, it likely belongs in Fusion. If no, it is probably engine-specific work that should stay behind an adapter.


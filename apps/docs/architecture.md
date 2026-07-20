# Architecture

DCraft Fusion separates **coordination** from **execution**.

## Control plane

The control plane is the Fusion kernel operators interact with:

- Tenants, projects, connections, and metadata registry
- Policies, approvals, workflow and run state
- Audit history and operator UX
- Provider/adapter contracts toward external engines

It does **not** replace Spark, dbt, Airflow, Kafka, Superset, warehouses, or lakes by default. Raw data and compute remain in customer infrastructure unless explicitly configured otherwise.

## Engines and adapters

Engines run outside (or beside) the control plane and perform work Fusion orchestrates or observes. Adapters translate Fusion intent into engine-native APIs.

In v1, the **Fusion CDC engine** is a first-class included engine for change data capture. Other engines plug in through the same engine-agnostic model.

## Control plane vs CDC engine

| Layer | Responsibility |
| --- | --- |
| Control plane | Auth, tenancy, metadata, policy, run/audit state, operator UI |
| CDC engine | Capture source changes, stream/transform events, checkpoints, CDC-specific ops |

Fusion coordinates CDC (and other engines) as muscles; the control plane is the brain.

```text
┌─────────────────────────────┐
│     Fusion control plane    │
│  metadata · policy · runs   │
└──────────────┬──────────────┘
               │ adapters / APIs
     ┌─────────┴─────────┐
     ▼                   ▼
┌──────────┐      ┌─────────────┐
│ CDC eng. │      │ Other eng.  │
│ (in OSS) │      │ (adapters)  │
└──────────┘      └─────────────┘
```

## Deeper reading

In-repo architecture notes live under [`docs/architecture`](https://github.com/DCraft-Labs/dcraft-fusion/tree/main/docs/architecture). For CDC internals, see [CDC](/cdc).

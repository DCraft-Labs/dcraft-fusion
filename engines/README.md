# Engines

Workers and execution-engine **adapters** live here (public, Apache-2.0).

Fusion orchestrates and observes engines through typed contracts. It is not the default execution engine.

## Fusion CDC Engine (proprietary images)

CDC **source code is private** (`DCraft-Labs/fusion-cdc-engine`).

Community users run CDC via **public GHCR images** and the public Helm chart:

| Image | Registry |
| --- | --- |
| Control plane | `ghcr.io/dcraft-labs/fusion-cdc-control-plane` |
| Worker | `ghcr.io/dcraft-labs/fusion-cdc-worker` |
| Frontend | `ghcr.io/dcraft-labs/fusion-cdc-frontend` |
| Transform worker (DuckDB) | `ghcr.io/dcraft-labs/fusion-cdc-transform-worker` |
| Spark consumer (optional / legacy) | `ghcr.io/dcraft-labs/fusion-cdc-spark-consumer` — not in default compose/Helm; prefer transform-worker |
| Transform worker | `ghcr.io/dcraft-labs/fusion-cdc-transform-worker` |

Install: `oci://ghcr.io/dcraft-labs/charts/fusion-cdc` — see [`infra/helm/fusion-cdc`](../infra/helm/fusion-cdc/).

Local Compose (images only): [`infra/local-dev/docker-compose.cdc.yml`](../infra/local-dev/docker-compose.cdc.yml).

## Future public adapters (scaffolds)

- `dbt-adapter`, `airflow-adapter`, `spark-adapter`, `superset-adapter`, `openmetadata-adapter`
- Thin workers that call customer-owned engines without embedding proprietary connector IP

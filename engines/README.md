# Engines

Workers and execution-engine adapters live here.

Fusion should orchestrate and observe these engines through typed contracts. It should not become the default execution engine itself.

Initial adapter/worker areas:

- `cdc-workers`
- `transform-worker`
- `duckdb-worker`
- `spark-adapter`
- `dbt-adapter`
- `airflow-adapter`
- `superset-adapter`
- `openmetadata-adapter`
- `fusion-cdc-engine`: imported CDC implementation kept as a future Fusion data-plane engine module.

`fusion-cdc-engine` is valuable, but it should not define the core product. Treat it as an optional execution/data-plane module controlled by the Fusion kernel.

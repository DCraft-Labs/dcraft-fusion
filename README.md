# DCraft Fusion

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)
[![Repo](https://img.shields.io/badge/GitHub-DCraft--Labs%2Fdcraft--fusion-0D7377)](https://github.com/DCraft-Labs/dcraft-fusion)

**One brain. Many muscles.**

DCraft Fusion is an AI-assisted, engine-agnostic **control plane** for modern data platforms. It coordinates, governs, observes, and assists existing systems — it does not replace Spark, dbt, Airflow, Kafka, Superset, warehouses, or lakes by default.

**Fusion CDC** runs as an official **data-plane muscle**: connector/worker **source is private**; Community gets **public container images** and a public Helm chart.

## Quick start (control plane)

```bash
docker compose -f infra/local-dev/docker-compose.yml up --build -d
```

| Surface | URL |
| --- | --- |
| Web | http://127.0.0.1:5174 |
| API | http://127.0.0.1:8080 |

Password → JWT by default. See [`infra/local-dev/.env.example`](infra/local-dev/.env.example).

### Optional: CDC via public images

```bash
docker compose -f infra/local-dev/docker-compose.cdc.yml up -d
```

## Helm (OCI)

Airbyte-style charts — override only what you need; BYO Postgres/Redis.

```bash
# Control plane
helm install dcraft-fusion oci://ghcr.io/dcraft-labs/charts/dcraft-fusion \
  --version 1.0.1 -n dcraft-fusion --create-namespace -f values.yaml

# CDC (pulls ghcr.io/dcraft-labs/fusion-cdc-* images; no CDC source in this repo)
helm install fusion-cdc oci://ghcr.io/dcraft-labs/charts/fusion-cdc \
  --version 1.0.1 -n fusion-cdc --create-namespace -f cdc-values.yaml
```

See [`infra/helm/README.md`](infra/helm/README.md) and chart `examples/`.

## Community vs Enterprise

| | Community | Enterprise |
| --- | --- | --- |
| Control plane source | Apache 2.0 | Apache 2.0 + commercial add-ons |
| CDC | Public images + Helm | Same + commercial support |
| CDC / connector source | Private | Private |
| SSO / BYOC / SLA | — | Yes |

Details: **[OPEN_CORE.md](OPEN_CORE.md)**.

## Documentation

```bash
npm run dev -w @dcraft-fusion/docs
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). PRs welcome for the **control plane** and public charts.

## License

Apache License 2.0 — see [LICENSE](LICENSE). CDC engine source is not included in this repository.

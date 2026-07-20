# DCraft Fusion

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)
[![Repo](https://img.shields.io/badge/GitHub-DCraft--Labs%2Fdcraft--fusion-0D7377)](https://github.com/DCraft-Labs/dcraft-fusion)

**One brain. Many muscles.**

DCraft Fusion is an AI-assisted, engine-agnostic **control plane** for modern data platforms. It coordinates, governs, observes, and assists existing systems — it does not replace Spark, dbt, Airflow, Kafka, Superset, warehouses, or lakes by default.

Fusion stores metadata, intent, policies, workflow and run state, approvals, audit history, and operational visibility. Raw data and compute stay in your infrastructure unless you configure otherwise.

**Change Data Capture (CDC) is included in the Community edition in v1** via the Fusion CDC engine.

## Why Fusion

- One operator surface across heterogeneous engines
- Clear separation of control plane vs execution engines
- Open-source Community core (Apache 2.0) with optional Enterprise identity and commercial offerings
- CDC and control-plane foundations ship together for real platform work, not demos alone

## Quick start (Docker Compose)

```bash
docker compose -f infra/local-dev/docker-compose.yml up --build -d
```

| Surface | URL |
| --- | --- |
| Web | http://127.0.0.1:5174 |
| API | http://127.0.0.1:8080 |

Community login defaults to password → JWT. See [`infra/local-dev/.env.example`](infra/local-dev/.env.example) for seeded accounts and overrides.

## Helm / GHCR

Charts and images publish to GitHub Container Registry:

```bash
helm install dcraft-fusion \
  oci://ghcr.io/dcraft-labs/charts/dcraft-fusion \
  --namespace dcraft-fusion \
  --create-namespace
```

Related: `oci://ghcr.io/dcraft-labs/charts/fusion-cdc` and in-repo sources under `infra/` and `engines/fusion-cdc-engine/helm/`.

## Community vs Enterprise

| | Community (OSS) | Enterprise |
| --- | --- | --- |
| License | Apache 2.0 | Commercial |
| Control plane + UX | Yes | Yes |
| Fusion CDC engine | Yes (included in v1) | Yes |
| Password / JWT auth | Yes | Yes |
| OIDC / SSO / SCIM | Lab wiring | Gated product |
| BYOC, SLA, managed | — | Available |

Details: **[OPEN_CORE.md](OPEN_CORE.md)**.

## Documentation

Public guide (VitePress): [`apps/docs`](apps/docs) — run with `npm run dev -w @dcraft-fusion/docs`.

CDC engine docs: [`engines/fusion-cdc-engine/docs`](engines/fusion-cdc-engine/docs).

## Contributing

See **[CONTRIBUTING.md](CONTRIBUTING.md)** (DCO sign-off required). Security reports: **[SECURITY.md](SECURITY.md)**.

## License

Licensed under the [Apache License 2.0](LICENSE).

Repository: [github.com/DCraft-Labs/dcraft-fusion](https://github.com/DCraft-Labs/dcraft-fusion)

# Install with Docker

The recommended path for labs and early evaluation is Docker Compose under `infra/local-dev`.

## Compose file

Primary file:

```text
infra/local-dev/docker-compose.yml
```

Services typically include PostgreSQL, Redis, the control-plane kernel, and the web UI.

## Bring the stack up

From the repository root:

```bash
docker compose -f infra/local-dev/docker-compose.yml up --build -d
```

Tear down:

```bash
docker compose -f infra/local-dev/docker-compose.yml down
```

## Configuration

Environment defaults are suitable for local development. Override via `.env` next to the Compose file. Start from:

[`infra/local-dev/.env.example`](https://github.com/DCraft-Labs/dcraft-fusion/blob/main/infra/local-dev/.env.example)

Notable variables:

| Variable | Purpose |
| --- | --- |
| `POSTGRES_PASSWORD` | Postgres password for the local `fusion` database |
| `FUSION_AUTH_MODE` | `password` (Community default) or `oidc` (Enterprise-style) |
| `FUSION_SEED_*` | Seeded operator accounts for first login |

## Ports

| Port | Service |
| --- | --- |
| `5174` | Web UI |
| `8080` | Control plane API |
| `5432` | PostgreSQL (local) |
| `6379` | Redis (local) |

## Kubernetes manifests

Raw manifests live under [`infra/kubernetes`](https://github.com/DCraft-Labs/dcraft-fusion/tree/main/infra/kubernetes). For chart-based installs, see [Helm](/install/helm).

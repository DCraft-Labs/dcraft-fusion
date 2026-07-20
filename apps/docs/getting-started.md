# Getting started

Run the local Community stack with Docker Compose from the repository root.

## Prerequisites

- Docker Desktop (or compatible Docker Engine + Compose)
- Git clone of [DCraft-Labs/dcraft-fusion](https://github.com/DCraft-Labs/dcraft-fusion)

## Start the stack

```bash
docker compose -f infra/local-dev/docker-compose.yml up --build -d
```

On Windows PowerShell:

```powershell
docker compose -f infra\local-dev\docker-compose.yml up --build -d
```

## Endpoints

| Surface | URL |
| --- | --- |
| Web UI | http://127.0.0.1:5174 |
| Control plane API | http://127.0.0.1:8080 |

## Login (Community password mode)

Local Compose defaults to simple password authentication (`FUSION_AUTH_MODE=password`). Seeded accounts and secrets are documented in:

[`infra/local-dev/.env.example`](https://github.com/DCraft-Labs/dcraft-fusion/blob/main/infra/local-dev/.env.example)

Copy that file to `.env` in the same directory for local overrides. Do not commit real credentials.

Default seed (change immediately in shared environments):

- Founder: `founder@dcraftlabs.com` / `changeme-founder`
- Superadmin: `superadmin@dcraftlabs.com` / `changeme-superadmin`

Successful password login returns a JWT used by the web UI and API clients.

## Next steps

- [Docker install details](/install/docker)
- [Helm / GHCR charts](/install/helm)
- [Authentication modes](/authentication)
- [Architecture overview](/architecture)

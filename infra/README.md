# Infrastructure

Infrastructure and local environment definitions live here.

- `kubernetes`: Raw Kubernetes manifests (images: `ghcr.io/dcraft-labs/...:latest`).
- `helm`: Helm charts
  - `helm/dcraft-fusion` — control-plane-kernel + web (+ optional Bitnami Postgres/Redis)
  - `helm/fusion-cdc` — CDC engine (canonical public chart)
- `local-dev`: Local Compose profiles
  - `docker-compose.yml` — Fusion control plane + web
  - `docker-compose.cdc.yml` — CDC profile (`--profile cdc`) for all-in-one local
- `observability`: Prometheus/Grafana assets

Secrets should be represented as references or examples only. Do not commit live credentials or PEM private keys.

### Helm quick start

```bash
helm dependency update infra/helm/dcraft-fusion
helm upgrade --install fusion infra/helm/dcraft-fusion \
  --namespace dcraft-fusion --create-namespace \
  --set secrets.seedFounderPassword=... \
  --set secrets.seedSuperadminPassword=... \
  --set postgresql.auth.password=...

helm dependency update infra/helm/fusion-cdc
helm upgrade --install fusion-cdc infra/helm/fusion-cdc \
  --namespace fusion --create-namespace \
  --set global.databaseUrl='postgresql://...'
```

OCI (after tagging `v*`):

```bash
helm install fusion oci://ghcr.io/dcraft-labs/charts/dcraft-fusion
helm install fusion-cdc oci://ghcr.io/dcraft-labs/charts/fusion-cdc
```

### Local CDC alongside Fusion

```bash
docker compose -f infra/local-dev/docker-compose.yml -f infra/local-dev/docker-compose.cdc.yml --profile cdc up -d
```


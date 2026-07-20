# Infrastructure

Infrastructure and local environment definitions live here.

- `kubernetes`: Raw Kubernetes manifests (images: `ghcr.io/dcraft-labs/...:latest`).
- `helm`: Helm charts (see [`helm/README.md`](helm/README.md))
  - `helm/dcraft-fusion` — open-source control plane (kernel + web); BYO Postgres/Redis
  - `helm/fusion-cdc` — image-only CDC chart (official GHCR images; proprietary engine source)
  - `helm/dcraft-fusion-stack` — optional local umbrella / stack guide
- `local-dev`: Local Compose profiles
  - `docker-compose.yml` — Fusion control plane + web
  - `docker-compose.cdc.yml` — CDC profile (`--profile cdc`) using GHCR images only
- `observability`: Prometheus/Grafana assets

Secrets should be represented as references or placeholders only. Do not commit live credentials or PEM private keys.

### Helm quick start (OCI)

```bash
# Control plane
helm install fusion oci://ghcr.io/dcraft-labs/charts/dcraft-fusion \
  --version 1.0.1 \
  --namespace dcraft-fusion --create-namespace \
  --set global.secrets.existingSecret=fusion-secrets \
  --set externalRedis.addr=redis:6379

# CDC (image-only)
helm install fusion-cdc oci://ghcr.io/dcraft-labs/charts/fusion-cdc \
  --version 1.0.1 \
  --namespace fusion-cdc --create-namespace \
  --set global.secrets.existingSecret=fusion-cdc-secrets
```

From this repo (no chart dependencies to fetch):

```bash
helm upgrade --install fusion infra/helm/dcraft-fusion \
  --namespace dcraft-fusion --create-namespace \
  --set global.secrets.existingSecret=fusion-secrets \
  --set externalRedis.addr=redis:6379

helm upgrade --install fusion-cdc infra/helm/fusion-cdc \
  --namespace fusion-cdc --create-namespace \
  --set global.secrets.existingSecret=fusion-cdc-secrets
```

### Local CDC alongside Fusion

```bash
docker compose -f infra/local-dev/docker-compose.yml \
  -f infra/local-dev/docker-compose.cdc.yml --profile cdc up -d
```

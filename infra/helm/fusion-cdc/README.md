# Fusion CDC Helm Chart (image-only)

> **CDC engine source is proprietary.** This chart deploys official DCraft Labs images from GHCR. It does **not** build from source and does **not** assume CDC source is present in this repository.

Bring your own Postgres and Redis. No Bitnami (or other) database subcharts are included.

## Install (OCI)

```bash
helm install fusion-cdc oci://ghcr.io/dcraft-labs/charts/fusion-cdc \
  --version 1.0.1 \
  --namespace fusion-cdc --create-namespace \
  -f my-values.yaml
```

## Prerequisites

- Kubernetes 1.25+
- Helm 3.10+
- External Postgres metadata DB
- External Redis
- GHCR pull access for `ghcr.io/dcraft-labs/fusion-cdc-*` (public or via `imagePullSecrets`)

## Quick start

```bash
kubectl create namespace fusion-cdc

kubectl -n fusion-cdc create secret generic fusion-cdc-secrets \
  --from-literal=DATABASE_URL='postgresql://fusion_user:CHANGE_ME@postgres:5432/fusion_cdc_metadata' \
  --from-literal=REDIS_URL='redis://redis:6379/0' \
  --from-literal=SECRET_KEY="$(openssl rand -hex 32)" \
  --from-literal=ENCRYPTION_KEY="$(openssl rand -hex 16)" \
  --from-literal=WORKER_TOKEN_SECRET="$(openssl rand -hex 32)" \
  --from-literal=WORKER_TOKEN="$(openssl rand -hex 32)"

helm install fusion-cdc oci://ghcr.io/dcraft-labs/charts/fusion-cdc \
  --version 1.0.1 \
  --namespace fusion-cdc \
  --set global.secrets.existingSecret=fusion-cdc-secrets
```

## Images (pre-built only)

| Component | Image | Default |
|-----------|-------|---------|
| Control plane | `ghcr.io/dcraft-labs/fusion-cdc-control-plane` | enabled |
| CDC workers | `ghcr.io/dcraft-labs/fusion-cdc-worker` | enabled |
| Frontend | `ghcr.io/dcraft-labs/fusion-cdc-frontend` | enabled |
| Spark consumer | `ghcr.io/dcraft-labs/fusion-cdc-spark-consumer` | optional (`sparkConsumer.enabled`) |
| Transform worker | `ghcr.io/dcraft-labs/fusion-cdc-transform-worker` | optional (`transformWorker.enabled`) |

Pin production installs with `image.digest` (preferred) or a specific `image.tag`.

## Example values

- [`examples/values-minimal.yaml`](examples/values-minimal.yaml)
- [`examples/values-production.yaml`](examples/values-production.yaml)

## Configuration highlights

| Key | Description |
|-----|-------------|
| `global.imageRegistry` | Optional registry mirror prefix |
| `global.imagePullSecrets` | Pull secrets for private GHCR |
| `global.secrets.existingSecret` | Pre-created Secret (recommended) |
| `externalRedis.url` | Fallback Redis URL if not in Secret |
| `controlPlane.*` / `cdcWorkers.*` / `frontend.*` | Probes, HPA, resources, SA, affinity |
| `sparkConsumer.enabled` | Optional Spark consumer |
| `transformWorker.enabled` | Optional transform workers (+ optional KEDA) |

## Stack with DCraft Fusion control plane

See [`../README.md`](../README.md) for ordered OCI install of both charts.

## Uninstall

```bash
helm uninstall fusion-cdc --namespace fusion-cdc
```

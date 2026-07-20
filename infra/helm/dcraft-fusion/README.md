# DCraft Fusion Helm Chart

Deploy the open-source **DCraft Fusion** control plane (kernel API + web UI).

Bring your own Postgres and Redis. This chart does **not** bundle databases.

## Install (OCI)

```bash
helm install fusion oci://ghcr.io/dcraft-labs/charts/dcraft-fusion \
  --version 1.0.1 \
  --namespace dcraft-fusion --create-namespace \
  -f my-values.yaml
```

## Prerequisites

- Kubernetes 1.25+
- Helm 3.10+
- External Postgres (provide `POSTGRES_DSN`)
- External Redis (set `externalRedis.addr`)
- A Secret for credentials (recommended) or `global.secrets.createSecret=true` for local/dev

## Quick start

```bash
# 1. Create secrets (do not commit these)
kubectl create namespace dcraft-fusion
kubectl -n dcraft-fusion create secret generic fusion-secrets \
  --from-literal=POSTGRES_DSN='postgres://fusion:CHANGE_ME@postgres:5432/fusion?sslmode=require' \
  --from-literal=FUSION_SEED_FOUNDER_EMAIL='founder@example.com' \
  --from-literal=FUSION_SEED_FOUNDER_PASSWORD="$(openssl rand -base64 24)" \
  --from-literal=FUSION_SEED_SUPERADMIN_EMAIL='superadmin@example.com' \
  --from-literal=FUSION_SEED_SUPERADMIN_PASSWORD="$(openssl rand -base64 24)"

# 2. Install
helm install fusion oci://ghcr.io/dcraft-labs/charts/dcraft-fusion \
  --version 1.0.1 \
  --namespace dcraft-fusion \
  --set global.secrets.existingSecret=fusion-secrets \
  --set externalRedis.addr=redis:6379
```

## Images

| Component | Image |
|-----------|-------|
| Control plane kernel | `ghcr.io/dcraft-labs/dcraft-fusion-control-plane-kernel` |
| Web UI | `ghcr.io/dcraft-labs/dcraft-fusion-web` |

Pin production installs with `image.digest` (preferred) or a specific `image.tag`.

## Example values

- [`examples/values-minimal.yaml`](examples/values-minimal.yaml) — local / smoke test
- [`examples/values-production.yaml`](examples/values-production.yaml) — HPA, ingress, digests, existingSecret

## Configuration highlights

| Key | Description |
|-----|-------------|
| `global.imageRegistry` | Optional registry mirror prefix |
| `global.imagePullSecrets` | Pull secrets for private GHCR |
| `global.edition` | `community` (default) or `enterprise` |
| `global.secrets.existingSecret` | Pre-created Secret name (recommended) |
| `externalRedis.addr` | Redis `host:port` |
| `controlPlaneKernel.*` / `web.*` | Replicas, probes, HPA, resources, SA, affinity |

## Upgrade

```bash
helm upgrade fusion oci://ghcr.io/dcraft-labs/charts/dcraft-fusion \
  --version 1.0.1 \
  --namespace dcraft-fusion \
  -f my-values.yaml
```

## Uninstall

```bash
helm uninstall fusion --namespace dcraft-fusion
```

Secrets with `helm.sh/resource-policy: keep` (chart-created) are retained on uninstall.

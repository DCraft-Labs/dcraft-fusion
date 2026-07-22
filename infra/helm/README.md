# DCraft Fusion Helm charts

Production-oriented Helm charts for DCraft Fusion. Charts follow Airbyte-style structure (`global`, per-component image/resources/probes/HPA, `existingSecret`).

| Chart | Path | OCI | What it deploys |
|-------|------|-----|-----------------|
| **dcraft-fusion** | [`dcraft-fusion/`](dcraft-fusion/) | `oci://ghcr.io/dcraft-labs/charts/dcraft-fusion` | Open-source control plane (kernel + web) |
| **fusion-cdc** | [`fusion-cdc/`](fusion-cdc/) | `oci://ghcr.io/dcraft-labs/charts/fusion-cdc` | Image-only CDC (official GHCR images) |

There is **no Bitnami / bundled Postgres or Redis**. Bring your own datastores.

## Product note

- Control plane (kernel + web) source remains open in this repository.
- **CDC engine source is proprietary.** The public `fusion-cdc` chart only pulls pre-built images (`ghcr.io/dcraft-labs/fusion-cdc-*`). It does not template in-tree CDC source paths.

## Ordered stack install (recommended)

Install the control plane first, then CDC. Use separate namespaces or the same namespace as you prefer.

```bash
# â”€â”€ Secrets (placeholders â€” replace values; do not commit) â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€
kubectl create namespace dcraft-fusion
kubectl create namespace fusion-cdc

kubectl -n dcraft-fusion create secret generic fusion-secrets \
  --from-literal=POSTGRES_DSN='postgres://fusion:CHANGE_ME@postgres:5432/fusion?sslmode=require' \
  --from-literal=FUSION_SEED_FOUNDER_EMAIL='founder@example.com' \
  --from-literal=FUSION_SEED_FOUNDER_PASSWORD='CHANGE_ME' \
  --from-literal=FUSION_SEED_SUPERADMIN_EMAIL='superadmin@example.com' \
  --from-literal=FUSION_SEED_SUPERADMIN_PASSWORD='CHANGE_ME'

kubectl -n fusion-cdc create secret generic fusion-cdc-secrets \
  --from-literal=DATABASE_URL='postgresql://fusion_user:CHANGE_ME@postgres:5432/fusion_cdc_metadata' \
  --from-literal=REDIS_URL='redis://redis:6379/0' \
  --from-literal=SECRET_KEY='CHANGE_ME' \
  --from-literal=ENCRYPTION_KEY='CHANGE_ME_32_CHARS_MINIMUM!!!!' \
  --from-literal=WORKER_TOKEN_SECRET='CHANGE_ME' \
  --from-literal=WORKER_TOKEN='CHANGE_ME'

# â”€â”€ 1. Control plane â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€
helm install fusion oci://ghcr.io/dcraft-labs/charts/dcraft-fusion \
  --version 1.2.19 \
  --namespace dcraft-fusion \
  -f dcraft-fusion/examples/values-minimal.yaml \
  --set externalRedis.addr=redis:6379

# â”€â”€ 2. CDC (image-only) â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€
helm install fusion-cdc oci://ghcr.io/dcraft-labs/charts/fusion-cdc \
  --version 1.2.19 \
  --namespace fusion-cdc \
  -f fusion-cdc/examples/values-minimal.yaml
```

### Combined values sketch

Use two values files (one per chart). Example:

**`stack-fusion.yaml`** (control plane):

```yaml
global:
  secrets:
    existingSecret: fusion-secrets
externalRedis:
  addr: redis.platform.svc.cluster.local:6379
controlPlaneKernel:
  replicaCount: 2
web:
  replicaCount: 2
```

**`stack-cdc.yaml`** (CDC):

```yaml
global:
  secrets:
    existingSecret: fusion-cdc-secrets
  imagePullSecrets:
    - name: ghcr-pull
controlPlane:
  replicaCount: 2
sparkConsumer:
  enabled: false
transformWorker:
  enabled: false
```

An umbrella chart with `file://` dependencies is awkward for OCI publishing, so the stack is documented here as an ordered install of the two published charts.

## Local chart lint / template

```bash
helm lint infra/helm/dcraft-fusion
helm lint infra/helm/fusion-cdc

helm template fusion infra/helm/dcraft-fusion \
  --set global.secrets.existingSecret=fusion-secrets \
  --set externalRedis.addr=redis:6379

helm template fusion-cdc infra/helm/fusion-cdc \
  --set global.secrets.existingSecret=fusion-cdc-secrets
```

## Local Docker Compose (CDC)

For local CDC against published images (no in-tree engine build):

```bash
docker compose -f infra/local-dev/docker-compose.cdc.yml --profile cdc up -d
```

See [`../local-dev/docker-compose.cdc.yml`](../local-dev/docker-compose.cdc.yml).

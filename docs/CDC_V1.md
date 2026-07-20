# CDC

Fusion CDC Engine is the first official **execution muscle** for DCraft Fusion.

## What is public vs private

| Artifact | Access |
| --- | --- |
| Control plane (this monorepo) | Apache 2.0 source |
| CDC Helm chart (`infra/helm/fusion-cdc`) | Public (deploys images only) |
| CDC container images | Public on `ghcr.io/dcraft-labs/fusion-cdc-*` |
| CDC connector / worker / CDC UI **source** | **Private** (`DCraft-Labs/fusion-cdc-engine`) |

## Install

```bash
helm install fusion-cdc oci://ghcr.io/dcraft-labs/charts/fusion-cdc \
  --version 1.0.1 \
  --namespace fusion-cdc \
  --create-namespace \
  -f examples/values-minimal.yaml
```

Or Compose (images only):

```bash
docker compose -f infra/local-dev/docker-compose.cdc.yml up -d
```

## Product boundary

Fusion governs connections, runs, audit, and policy. CDC moves data under that governance.
The engine is not a second brand — it is a muscle of DCraft Fusion.

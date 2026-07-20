# CDC

Fusion CDC is the first official **execution muscle** for DCraft Fusion.

## Public vs private

| Artifact | Access |
| --- | --- |
| Control plane source | Apache 2.0 (this repo) |
| CDC Helm chart | Public — [`infra/helm/fusion-cdc`](https://github.com/DCraft-Labs/dcraft-fusion/tree/main/infra/helm/fusion-cdc) |
| CDC images | Public — `ghcr.io/dcraft-labs/fusion-cdc-*` |
| CDC / connector **source** | Private — not in this repository |

## How CDC fits

1. Operators register sources and pipelines through the Fusion control plane.
2. Official CDC images capture changes from supported databases.
3. Events stream to warehouses or other consumers with checkpoints and controls.
4. Run state, metadata, and audit remain visible in the control plane.

## Install

```bash
helm install fusion-cdc oci://ghcr.io/dcraft-labs/charts/fusion-cdc \
  --version 1.0.1 \
  --namespace fusion-cdc \
  --create-namespace
```

Compose (images only): `infra/local-dev/docker-compose.cdc.yml`

See [Install → Helm](/install/helm) and [Open core](/open-core).

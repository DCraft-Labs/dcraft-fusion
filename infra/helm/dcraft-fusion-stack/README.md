# DCraft Fusion stack (optional umbrella)

Thin local umbrella that depends on the sibling charts via `file://` paths.

**For OCI / production:** install the two published charts separately — see [`../README.md`](../README.md). File dependencies do not publish cleanly as a single OCI umbrella.

## Local install (from this repo)

```bash
cd infra/helm/dcraft-fusion-stack
helm dependency update
helm upgrade --install fusion-stack . \
  --namespace dcraft --create-namespace \
  -f values.yaml
```

## OCI ordered install (recommended)

```bash
helm install fusion oci://ghcr.io/dcraft-labs/charts/dcraft-fusion --version 1.0.1 ...
helm install fusion-cdc oci://ghcr.io/dcraft-labs/charts/fusion-cdc --version 1.0.1 ...
```

CDC engine source is proprietary. The `fusion-cdc` subchart / OCI chart deploys official GHCR images only.

# Install with Helm

Community container images and Helm charts are published to GitHub Container Registry (GHCR) under DCraft Labs.

## Chart registry

Install from the OCI registry:

```bash
helm install dcraft-fusion \
  oci://ghcr.io/dcraft-labs/charts/dcraft-fusion \
  --namespace dcraft-fusion \
  --create-namespace
```

CDC engine chart (when deploying CDC alongside the control plane):

```bash
helm install fusion-cdc \
  oci://ghcr.io/dcraft-labs/charts/fusion-cdc \
  --namespace dcraft-fusion
```

Exact chart names and versions follow GHCR package tags. Prefer pinning a chart version in production:

```bash
helm install dcraft-fusion \
  oci://ghcr.io/dcraft-labs/charts/dcraft-fusion \
  --version <chart-version>
```

## Prerequisites

- Kubernetes cluster with sufficient CPU/memory for control plane + dependencies
- `helm` 3.8+ (OCI support)
- Ability to pull from `ghcr.io` (public packages or authenticated pulls as required)

## Values

Override image tags, auth mode, ingress, and resource limits with a values file:

```bash
helm upgrade --install dcraft-fusion \
  oci://ghcr.io/dcraft-labs/charts/dcraft-fusion \
  -f my-values.yaml
```

Set Community password auth or Enterprise OIDC via chart values / secrets. See [Authentication](/authentication).

## Related paths in-repo

- Raw Kubernetes manifests: `infra/kubernetes`
- CDC Helm chart source: `engines/fusion-cdc-engine/helm/fusion-cdc`

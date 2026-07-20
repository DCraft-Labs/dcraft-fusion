# Open Core

DCraft Fusion uses an open-core model:

| Layer | License / access | What you get |
| --- | --- | --- |
| **Control plane** (kernel, web, contracts, public adapters) | Apache 2.0 in `dcraft-fusion` | Full source, fork, contribute |
| **Fusion CDC Engine** (connectors, workers, CDC UI source) | **Private source** | Official **public GHCR images** + public Helm chart |
| **Enterprise** | Commercial | SSO/SAML/SCIM, BYOC, SLA, managed cloud |

## Why CDC source is private

CDC connector/worker IP is the proprietary “muscle.” Shipping **images + Helm** lets the community run production CDC without forking engine internals — similar to commercial engines that publish installable artifacts while keeping core source closed.

You can still:

- Install CDC with `helm install … oci://ghcr.io/dcraft-labs/charts/fusion-cdc`
- Pin image digests for supply-chain control
- File bugs against the public control-plane repo (runtime / chart / integration)

## Community (Apache 2.0)

- Control plane kernel and operator UX (source)
- Simple password / JWT auth
- Public Helm charts for control plane **and** CDC (image-only)
- Public GHCR images for control plane **and** CDC
- Engine-agnostic coordination, metadata, policy, run state, audit
- Community support via GitHub Issues and Discussions

## Enterprise

- SSO / SAML / SCIM
- BYOC / air-gapped commercial support
- SLA and managed offerings
- Priority support

## Trademark

"DCraft", "DCraft Labs", and "DCraft Fusion" are trademarks of DCraft Labs.

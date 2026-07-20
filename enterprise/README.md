# DCraft Fusion Enterprise (private)

This directory documents the **paid** Enterprise / Cloud surface that depends on the public Apache-2.0 Community APIs. It is a skeleton for a future private repository (`dcraft-fusion-enterprise`).

## Never reverse the dependency

Enterprise packages may import public contracts and APIs. The public Community monorepo must not import this tree.

## Planned modules

| Module | Purpose |
|--------|---------|
| `sso/` | SAML / OIDC federation, SCIM provisioning |
| `byoc/` | Bring-your-own-cloud / air-gapped installers and license checks |
| `billing/` | Seat / connection metering for Cloud Team & Growth |
| `compliance/` | SIEM export packs, retention policies, audit evidence |
| `support/` | SLA runbooks and customer onboarding automation |

## Community vs Enterprise

See [`OPEN_CORE.md`](../OPEN_CORE.md) in the public repo.

## Status

Scaffold only for v1.0.0 Community launch. Implementation lives in a private GitHub repository owned by DCraft Labs.

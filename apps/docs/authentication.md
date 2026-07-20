# Authentication

DCraft Fusion supports two primary auth postures aligned with the open-core model.

## Community: password → JWT

Default for local Compose and Community deployments:

- `FUSION_AUTH_MODE=password` (or `local`)
- Operator signs in with email and password
- Control plane returns a signed **JWT** for API and UI sessions

Seeded local accounts are defined in [`infra/local-dev/.env.example`](https://github.com/DCraft-Labs/dcraft-fusion/blob/main/infra/local-dev/.env.example). Change defaults before any shared or production use.

Password mode is intended for labs, demos, and early self-hosted Community installs.

## Enterprise: OIDC / SSO

Enterprise-oriented deployments use OIDC (and related SSO capabilities such as SAML/SCIM in commercial offerings):

- `FUSION_AUTH_MODE=oidc`
- Authorization-code style login against an identity provider
- JWKS / issuer / audience configured via environment (for example `FUSION_OIDC_*`)

OIDC/SSO, stronger isolation patterns, and managed options are **Enterprise-gated**. Community can evaluate OIDC wiring in development; production SSO packages are part of the Enterprise offering. See [Open core](/open-core) and root [`OPEN_CORE.md`](https://github.com/DCraft-Labs/dcraft-fusion/blob/main/OPEN_CORE.md).

## Choosing a mode

| Mode | Edition | Typical use |
| --- | --- | --- |
| `password` | Community | Local Compose, labs, simple self-host |
| `oidc` | Enterprise (SSO) | Organizational IdP, SSO, production identity |

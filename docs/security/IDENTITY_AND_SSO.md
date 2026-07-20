# Identity And SSO

## Day-One Modes

- `dev`: local header-based context for Docker Desktop and tests.
- `oidc`: production mode requiring a Bearer JWT signed by an RS256 key from the configured JWKS endpoint.
- `saml`: enterprise SSO design target represented in platform admin; the first implementation exchanges SAML at the identity-provider layer and sends Fusion an OIDC token.

## OIDC Runtime Settings

- `FUSION_AUTH_MODE=oidc`
- `FUSION_OIDC_ISSUER`
- `FUSION_OIDC_AUDIENCE`
- `FUSION_OIDC_JWKS_URL`

When OIDC mode is enabled, the API rejects requests without a valid bearer token. Verified claims populate actor and tenant context:

- `sub` or `actor_id`
- `organization_id`
- `tenant_id`
- `project_id`

## Superadmin

Superadmin users are platform-scope owners. The seed bootstrap user is `user-superadmin`, and the Superadmin UI calls `/api/v1/platform/overview` with that actor in local mode. In OIDC mode this must come from a verified token claim.

## SAML

SAML is supported as an enterprise integration pattern by using the customer IdP as the SAML authority and the Fusion identity broker as the OIDC token issuer. Direct SAML ACS endpoints are deferred until a customer requires direct SAML into Fusion instead of SAML-to-OIDC federation.

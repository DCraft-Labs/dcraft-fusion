# Secret Management

Fusion stores secret references, not raw credentials.

## Supported Provider References

- HashiCorp Vault: `vault/path/to/secret`
- AWS Secrets Manager: `aws-secrets-manager/secret-name`
- GCP Secret Manager: `gcp-secret-manager/project/secret`
- Azure Key Vault: `azure-key-vault/vault/secret`
- Kubernetes Secret: `kubernetes/namespace/name/key`

## Rules

- Connection APIs reject raw URLs and values containing password-like material.
- API responses return `secretref:` values only.
- Provider bindings are visible in the Superadmin UI.
- Customer-owned infrastructure deployments can map provider references through `fusion-secret-provider-map`.

## Day-One Local Behavior

Local Postgres discovery can resolve `FUSION_SECRET_<SANITIZED_REF>` from environment variables. If the reference contains `postgres`, it falls back to `POSTGRES_DSN` for Docker Desktop validation.

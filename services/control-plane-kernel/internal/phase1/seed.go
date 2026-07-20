package phase1

import "time"

func (store *Store) seed() {
	createdAt := time.Date(2026, 5, 18, 0, 0, 0, 0, time.UTC)
	finishedAt := createdAt.Add(19 * time.Minute)

	store.memberships["mem-owner"] = Membership{
		ID: "mem-owner", UserID: "user-founder", Email: "founder@dcraftlabs.com", OrganizationID: "org-b2b2c-1", TenantID: "tenant-brand-a", ProjectID: "project-prod", Role: RoleOwner, CreatedAt: createdAt,
	}
	store.memberships["mem-superadmin"] = Membership{
		ID: "mem-superadmin", UserID: "user-superadmin", Email: "superadmin@dcraftlabs.com", OrganizationID: "platform", Role: RoleOwner, CreatedAt: createdAt,
	}
	store.memberships["mem-b2b"] = Membership{
		ID: "mem-b2b", UserID: "user-acme", Email: "admin@acme.example", OrganizationID: "org-b2b-1", TenantID: "tenant-finance", ProjectID: "project-finance-prod", Role: RoleAdmin, CreatedAt: createdAt,
	}

	store.serviceAccounts["sa-postgres"] = ServiceAccount{
		ID: "sa-postgres", Name: "postgres-readonly-sync", OrganizationID: "org-b2b2c-1", TenantID: "tenant-brand-a", ProjectID: "project-prod", Role: RoleOperator, TokenPreview: "fusion_sa_...brandA", CreatedAt: createdAt,
	}

	store.connections["conn-postgres-brand-a"] = Connection{
		ID: "conn-postgres-brand-a", Name: "Brand A Postgres", Kind: ConnectorPostgres, OrganizationID: "org-b2b2c-1", TenantID: "tenant-brand-a", ProjectID: "project-prod", Environment: "production", OwnerID: "user-founder", SecretRef: "secretref:vault/brand-a/postgres-readonly", ReadOnly: true, Status: ConnectionHealthy, CreatedAt: createdAt,
	}
	store.connections["conn-airflow-brand-a"] = Connection{
		ID: "conn-airflow-brand-a", Name: "Brand A Airflow", Kind: ConnectorAirflow, OrganizationID: "org-b2b2c-1", TenantID: "tenant-brand-a", ProjectID: "project-prod", Environment: "production", OwnerID: "user-founder", SecretRef: "secretref:vault/brand-a/airflow-readonly", ReadOnly: true, Status: ConnectionHealthy, CreatedAt: createdAt,
	}

	store.datasets["ds-orders-brand-a"] = Dataset{
		ID: "ds-orders-brand-a", Name: "orders", Schema: "public", Type: "table", SourceSystem: "postgres", ConnectionID: "conn-postgres-brand-a", OrganizationID: "org-b2b2c-1", TenantID: "tenant-brand-a", ProjectID: "project-prod", Owner: "data-platform", Tags: []string{"revenue", "orders"}, Description: "Orders metadata discovered through the read-only Postgres connector.", FreshnessAt: createdAt.Add(-10 * time.Minute), LastDiscoveredAt: createdAt,
	}
	store.datasets["ds-customers-brand-a"] = Dataset{
		ID: "ds-customers-brand-a", Name: "customers", Schema: "public", Type: "table", SourceSystem: "postgres", ConnectionID: "conn-postgres-brand-a", OrganizationID: "org-b2b2c-1", TenantID: "tenant-brand-a", ProjectID: "project-prod", Owner: "data-platform", Tags: []string{"customer"}, Description: "Customer profile metadata without raw customer data movement.", FreshnessAt: createdAt.Add(-12 * time.Minute), LastDiscoveredAt: createdAt,
	}

	store.runs["run-success-brand-a"] = WorkflowRun{
		ID: "run-success-brand-a", WorkflowName: "daily_revenue_snapshot", Status: RunSucceeded, OrganizationID: "org-b2b2c-1", TenantID: "tenant-brand-a", ProjectID: "project-prod", Environment: "production", CorrelationID: "corr-run-success", ArtifactRef: "airflow://dag/daily_revenue_snapshot/run/2026-05-18", StartedAt: createdAt, FinishedAt: &finishedAt,
	}
	store.runs["run-failed-brand-a"] = WorkflowRun{
		ID: "run-failed-brand-a", WorkflowName: "customer_freshness_check", Status: RunFailed, OrganizationID: "org-b2b2c-1", TenantID: "tenant-brand-a", ProjectID: "project-prod", Environment: "production", CorrelationID: "corr-run-failed", ArtifactRef: "airflow://dag/customer_freshness_check/run/2026-05-18", StartedAt: createdAt.Add(30 * time.Minute),
	}

	store.policies["policy-ai-draft"] = Policy{
		ID: "policy-ai-draft", Name: "AI actions require approval", Type: "ai_action_policy", OrganizationID: "org-b2b2c-1", TenantID: "tenant-brand-a", ProjectID: "project-prod", Mode: "draft_only",
	}
	store.policies["policy-raw-data"] = Policy{
		ID: "policy-raw-data", Name: "No raw data movement", Type: "data_movement_policy", OrganizationID: "org-b2b2c-1", TenantID: "tenant-brand-a", ProjectID: "project-prod", Mode: "deny_raw_data_movement",
	}

	store.aiRecommendations["ai-failed-run"] = AIRecommendation{
		ID: "ai-failed-run", OrganizationID: "org-b2b2c-1", TenantID: "tenant-brand-a", ProjectID: "project-prod", Title: "Review customer freshness failure", Summary: "The failed run is isolated to Brand A production. Suggested action remains a draft until approved.", SourceRefs: []string{"run:run-failed-brand-a", "dataset:ds-customers-brand-a"}, State: ApprovalDraft, CreatedAt: createdAt,
	}

	store.jobs["job-metadata-bootstrap"] = AsyncJob{
		ID: "job-metadata-bootstrap", Type: "metadata.discovery", Status: JobSucceeded, OrganizationID: "org-b2b2c-1", TenantID: "tenant-brand-a", ProjectID: "project-prod", ResourceID: "conn-postgres-brand-a", Message: "Initial metadata discovery completed.", Metadata: map[string]string{"backend": "redis_or_postgres"}, CreatedAt: createdAt, UpdatedAt: createdAt.Add(2 * time.Minute),
	}
	store.identityProviders["idp-platform-oidc"] = IdentityProviderConfig{
		ID: "idp-platform-oidc", Name: "Platform OIDC", Protocol: "oidc", Issuer: "https://identity.example.com", Audience: "dcraft-fusion", JWKSURL: "https://identity.example.com/.well-known/jwks.json", SSOEnabled: true, SAMLEnabled: false, Status: "configured", CreatedAt: createdAt,
	}
	store.identityProviders["idp-enterprise-saml"] = IdentityProviderConfig{
		ID: "idp-enterprise-saml", Name: "Enterprise SAML", Protocol: "saml", Issuer: "https://sso.example.com/metadata", Audience: "dcraft-fusion", SSOEnabled: true, SAMLEnabled: true, Status: "design_ready", CreatedAt: createdAt,
	}
	store.secretProviders["secret-vault"] = SecretProviderConfig{
		ID: "secret-vault", Name: "HashiCorp Vault", Kind: "vault", Endpoint: "https://vault.example.com", Scope: "platform", Status: "configured", CreatedAt: createdAt,
	}
	store.secretProviders["secret-aws"] = SecretProviderConfig{
		ID: "secret-aws", Name: "AWS Secrets Manager", Kind: "aws_secrets_manager", Endpoint: "arn:aws:secretsmanager:*", Scope: "customer_vpc", Status: "configured", CreatedAt: createdAt,
	}
	store.secretProviders["secret-gcp"] = SecretProviderConfig{
		ID: "secret-gcp", Name: "GCP Secret Manager", Kind: "gcp_secret_manager", Endpoint: "projects/*/secrets/*", Scope: "customer_project", Status: "configured", CreatedAt: createdAt,
	}
	store.secretProviders["secret-azure"] = SecretProviderConfig{
		ID: "secret-azure", Name: "Azure Key Vault", Kind: "azure_key_vault", Endpoint: "https://*.vault.azure.net", Scope: "customer_subscription", Status: "configured", CreatedAt: createdAt,
	}
	store.adapters["adapter-airflow"] = AdapterStatus{
		ID: "adapter-airflow", Name: "Airflow", Kind: "airflow", ReadOnly: true, Status: "implemented", Description: "Reads DAG and DAG-run metadata through the Airflow stable REST API when a read-only secret reference is configured.", UpdatedAt: createdAt,
	}
	store.adapters["adapter-dagster"] = AdapterStatus{
		ID: "adapter-dagster", Name: "Dagster", Kind: "dagster", ReadOnly: true, Status: "implemented", Description: "Reads repository, asset, and run metadata through Dagster GraphQL using a read-only token.", UpdatedAt: createdAt,
	}
	store.adapters["adapter-dbt"] = AdapterStatus{
		ID: "adapter-dbt", Name: "dbt", Kind: "dbt", ReadOnly: true, Status: "implemented", Description: "Ingests manifest/catalog/run-results artifacts without moving warehouse data.", UpdatedAt: createdAt,
	}
	store.backup = BackupStatus{
		ID: "backup-postgres", Target: "postgres", Schedule: "0 2 * * *", LastVerifiedAt: createdAt, LastVerificationJob: "job-backup-verify", Status: "scheduled", Runbook: "docs/operations/BACKUP_RESTORE_RUNBOOK.md",
	}
	store.observability = ObservabilityStatus{
		TracesEnabled: true, MetricsEndpoint: "/metrics", LogFormat: "json", Dashboards: []string{"infra/observability/grafana/dcraft-fusion-overview.json"}, Alerts: []string{"infra/observability/prometheus/alerts.yaml"}, OpenTelemetryReady: true,
	}
	store.ha = HAStatus{
		APIMinReplicas: 2, WebMinReplicas: 2, Autoscaling: true, PodDisruption: true, LoadTests: []string{"scripts/load-test.ps1"}, Status: "configured",
	}
}

import type { FusionWorkspace } from "./workspace.js";

export function createDemoWorkspace(): FusionWorkspace {
  const organization = {
    id: "org-b2b2c-1",
    name: "Commerce Platform Co.",
    type: "consumer_app" as const,
    deploymentProfile: "customer_vpc" as const,
    region: "us-east-1"
  };

  const tenant = {
    id: "tenant-brand-a",
    organizationId: organization.id,
    name: "Brand A",
    type: "brand" as const
  };

  const project = {
    id: "project-prod",
    organizationId: organization.id,
    tenantId: tenant.id,
    name: "Production Data Platform",
    environment: "prod" as const
  };

  return {
    context: { organization, tenant, project },
    organizations: [
      organization,
      {
        id: "org-b2b-1",
        name: "Acme Analytics",
        type: "business",
        deploymentProfile: "pooled_saas",
        region: "us-east-1"
      },
      {
        id: "org-internal-1",
        name: "Global Enterprise Platform",
        type: "internal_platform",
        deploymentProfile: "dedicated_stack",
        region: "eu-west-1"
      }
    ],
    tenants: [
      tenant,
      { id: "tenant-brand-b", organizationId: organization.id, name: "Brand B", type: "brand" },
      { id: "tenant-finance", organizationId: "org-b2b-1", name: "Finance", type: "business_unit" }
    ],
    projects: [
      project,
      {
        id: "project-sandbox",
        organizationId: organization.id,
        tenantId: tenant.id,
        name: "Sandbox",
        environment: "sandbox"
      }
    ],
    integrations: [
      {
        id: "postgres",
        name: "Postgres",
        category: "source",
        status: "ready",
        description: "Read-only metadata discovery and connection health."
      },
      {
        id: "airflow",
        name: "Airflow",
        category: "orchestration",
        status: "ready",
        description: "Read-only DAG and run ingestion for Phase 1."
      },
      {
        id: "github",
        name: "GitHub",
        category: "notification",
        status: "ready",
        description: "Change review and future approval integration."
      },
      {
        id: "fusion-cdc",
        name: "Fusion CDC Engine",
        category: "engine",
        status: "preview",
        description: "Optional data-plane engine for later CDC-enabled deployments."
      }
    ],
    connections: [
      {
        id: "conn-postgres-brand-a",
        name: "Brand A Postgres",
        kind: "postgres",
        organizationId: organization.id,
        tenantId: tenant.id,
        projectId: project.id,
        secretRef: "secretref:vault/brand-a/postgres-readonly",
        readOnly: true,
        status: "healthy"
      },
      {
        id: "conn-airflow-brand-a",
        name: "Brand A Airflow",
        kind: "airflow",
        organizationId: organization.id,
        tenantId: tenant.id,
        projectId: project.id,
        secretRef: "secretref:vault/brand-a/airflow-readonly",
        readOnly: true,
        status: "healthy"
      }
    ],
    datasets: [
      {
        id: "ds-orders-brand-a",
        name: "orders",
        schema: "public",
        type: "table",
        sourceSystem: "postgres",
        organizationId: organization.id,
        tenantId: tenant.id,
        projectId: project.id,
        owner: "data-platform",
        tags: ["revenue", "orders"],
        lastDiscoveredAt: "2026-05-18T00:00:00Z"
      },
      {
        id: "ds-customers-brand-a",
        name: "customers",
        schema: "public",
        type: "table",
        sourceSystem: "postgres",
        organizationId: organization.id,
        tenantId: tenant.id,
        projectId: project.id,
        owner: "data-platform",
        tags: ["customer"],
        lastDiscoveredAt: "2026-05-18T00:00:00Z"
      }
    ],
    runs: [
      {
        id: "run-success-brand-a",
        workflowName: "daily_revenue_snapshot",
        status: "succeeded",
        organizationId: organization.id,
        tenantId: tenant.id,
        projectId: project.id,
        environment: "production",
        correlationId: "corr-run-success",
        artifactRef: "airflow://dag/daily_revenue_snapshot/run/2026-05-18"
      },
      {
        id: "run-failed-brand-a",
        workflowName: "customer_freshness_check",
        status: "failed",
        organizationId: organization.id,
        tenantId: tenant.id,
        projectId: project.id,
        environment: "production",
        correlationId: "corr-run-failed",
        artifactRef: "airflow://dag/customer_freshness_check/run/2026-05-18"
      }
    ],
    auditEvents: [
      {
        id: "audit-connection",
        action: "connection.created",
        actorId: "user-founder",
        correlationId: "corr-1",
        resourceType: "Connection",
        resourceId: "conn-postgres-brand-a",
        occurredAt: new Date(Date.now() - 1000 * 60 * 2).toISOString()
      },
      {
        id: "audit-discovery",
        action: "metadata.discovered",
        actorId: "user-founder",
        correlationId: "corr-2",
        resourceType: "Connection",
        resourceId: "conn-postgres-brand-a",
        occurredAt: new Date(Date.now() - 1000 * 60 * 37).toISOString()
      }
    ],
    policies: [
      { id: "policy-ai-draft", name: "AI actions require approval", type: "ai_action_policy", mode: "draft_only" },
      { id: "policy-raw-data", name: "No raw data movement", type: "data_movement_policy", mode: "deny_raw_data_movement" }
    ],
    aiRecommendations: [
      {
        id: "ai-failed-run",
        title: "Review customer freshness failure",
        summary: "The failed run is isolated to Brand A production. Suggested action remains a draft until approved.",
        state: "draft",
        sourceRefs: ["run:run-failed-brand-a", "dataset:ds-customers-brand-a"]
      }
    ]
  };
}

import type { DeploymentProfile, EnvironmentType, OrganizationType, TenantType } from "@dcraft-fusion/domain";

export interface WorkspaceOrganization {
  readonly id: string;
  readonly name: string;
  readonly type: OrganizationType;
  readonly deploymentProfile: DeploymentProfile;
  readonly region: string;
}

export interface WorkspaceTenant {
  readonly id: string;
  readonly organizationId: string;
  readonly name: string;
  readonly type: TenantType;
}

export interface WorkspaceProject {
  readonly id: string;
  readonly organizationId: string;
  readonly tenantId: string;
  readonly name: string;
  readonly environment: EnvironmentType;
}

export interface WorkspaceContext {
  readonly organization: WorkspaceOrganization;
  readonly tenant: WorkspaceTenant;
  readonly project: WorkspaceProject;
}

export interface IntegrationCard {
  readonly id: string;
  readonly name: string;
  readonly category: "source" | "orchestration" | "notification" | "bi" | "engine";
  readonly status: "ready" | "preview";
  readonly description: string;
}

export interface Connection {
  readonly id: string;
  readonly name: string;
  readonly kind: "postgres" | "airflow" | "github" | "slack" | "fusion_cdc_engine";
  readonly organizationId: string;
  readonly tenantId: string;
  readonly projectId: string;
  readonly secretRef: string;
  readonly readOnly: boolean;
  readonly status: "untested" | "healthy" | "blocked";
}

export interface Dataset {
  readonly id: string;
  readonly name: string;
  readonly schema: string;
  readonly type: "table" | "view";
  readonly sourceSystem: string;
  readonly organizationId: string;
  readonly tenantId: string;
  readonly projectId: string;
  readonly owner: string;
  readonly tags: readonly string[];
  readonly lastDiscoveredAt: string;
}

export interface WorkflowRun {
  readonly id: string;
  readonly workflowName: string;
  readonly status: "queued" | "running" | "succeeded" | "failed" | "canceled" | "unknown";
  readonly organizationId: string;
  readonly tenantId: string;
  readonly projectId: string;
  readonly environment: string;
  readonly correlationId: string;
  readonly artifactRef: string;
}

export interface AuditEvent {
  readonly id: string;
  readonly action: string;
  readonly actorId: string;
  readonly correlationId: string;
  readonly resourceType: string;
  readonly resourceId: string;
  // ISO timestamp of when the event occurred. Optional so the UI degrades
  // gracefully if a future backend returns events without it.
  readonly occurredAt?: string;
}

export interface Policy {
  readonly id: string;
  readonly name: string;
  readonly type: string;
  readonly mode: string;
}

export interface AIRecommendation {
  readonly id: string;
  readonly title: string;
  readonly summary: string;
  readonly state: "draft" | "pending_review" | "approved" | "rejected";
  readonly sourceRefs: readonly string[];
}

export interface AsyncJob {
  readonly id: string;
  readonly type: string;
  readonly status: "queued" | "running" | "succeeded" | "failed";
  readonly resourceId: string;
  readonly message: string;
  readonly updatedAt: string;
}

export interface IdentityProviderConfig {
  readonly id: string;
  readonly name: string;
  readonly protocol: string;
  readonly issuer: string;
  readonly audience: string;
  readonly jwksUrl?: string;
  readonly ssoEnabled: boolean;
  readonly samlEnabled: boolean;
  readonly status: string;
}

export interface SecretProviderConfig {
  readonly id: string;
  readonly name: string;
  readonly kind: string;
  readonly endpoint: string;
  readonly scope: string;
  readonly status: string;
}

export interface AdapterStatus {
  readonly id: string;
  readonly name: string;
  readonly kind: string;
  readonly readOnly: boolean;
  readonly status: string;
  readonly description: string;
}

export interface PlatformOverview {
  readonly identityProviders: readonly IdentityProviderConfig[];
  readonly secretProviders: readonly SecretProviderConfig[];
  readonly adapters: readonly AdapterStatus[];
  readonly backup: {
    readonly target: string;
    readonly schedule: string;
    readonly status: string;
    readonly runbook: string;
  };
  readonly observability: {
    readonly tracesEnabled: boolean;
    readonly metricsEndpoint: string;
    readonly logFormat: string;
    readonly dashboards: readonly string[];
    readonly alerts: readonly string[];
    readonly openTelemetryReady: boolean;
  };
  readonly ha: {
    readonly apiMinReplicas: number;
    readonly webMinReplicas: number;
    readonly autoscaling: boolean;
    readonly podDisruption: boolean;
    readonly loadTests: readonly string[];
    readonly status: string;
  };
  readonly jobs: readonly AsyncJob[];
  readonly superAdmins: readonly Membership[];
  readonly organizations: readonly PlatformOrganization[];
  readonly recentAuditEvents: readonly PlatformAuditEvent[];
}

export interface Membership {
  readonly id: string;
  readonly userId: string;
  readonly email: string;
  readonly organizationId: string;
  readonly role: string;
}

export interface PlatformProject {
  readonly id: string;
  readonly name: string;
  readonly environment: string;
  readonly createdAt: string;
}

export interface PlatformTenant {
  readonly id: string;
  readonly name: string;
  readonly model: string;
  readonly isolation: string;
  readonly region: string;
  readonly createdAt: string;
  readonly projects: readonly PlatformProject[];
}

export interface PlatformOrganization {
  readonly id: string;
  readonly name: string;
  readonly type: string;
  readonly deploymentProfile: string;
  readonly region: string;
  readonly dataPlaneLocation: string;
  readonly rawDataMovementAllowed: boolean;
  readonly createdAt: string;
  readonly tenants: readonly PlatformTenant[];
}

export interface PlatformAuditEvent {
  readonly id: string;
  readonly action: string;
  readonly actorId: string;
  readonly correlationId: string;
  readonly resourceType: string;
  readonly resourceId: string;
  readonly occurredAt: string;
}

export interface AuthUser {
  readonly id: string;
  readonly email: string;
  readonly label: string;
  readonly scope: "platform" | "tenant";
  readonly passwordHint?: string;
  readonly organizationId?: string;
  readonly tenantId?: string;
  readonly projectId?: string;
}

export interface OpenIDConfiguration {
  readonly issuer: string;
  readonly authorization_endpoint: string;
  readonly token_endpoint: string;
  readonly jwks_uri: string;
}

export interface UserSession {
  readonly actorId: string;
  readonly email: string;
  readonly label: string;
  readonly scope: "platform" | "tenant";
  readonly organizationId?: string;
  readonly tenantId?: string;
  readonly projectId?: string;
}

export interface FusionWorkspace {
  readonly context: WorkspaceContext;
  readonly organizations: readonly WorkspaceOrganization[];
  readonly tenants: readonly WorkspaceTenant[];
  readonly projects: readonly WorkspaceProject[];
  readonly integrations: readonly IntegrationCard[];
  readonly connections: readonly Connection[];
  readonly datasets: readonly Dataset[];
  readonly runs: readonly WorkflowRun[];
  readonly auditEvents: readonly AuditEvent[];
  readonly policies: readonly Policy[];
  readonly aiRecommendations: readonly AIRecommendation[];
}

export function validateWorkspaceContext(workspace: FusionWorkspace): void {
  const { organization, tenant, project } = workspace.context;

  if (tenant.organizationId !== organization.id) {
    throw new Error("Active tenant does not belong to active organization");
  }

  if (project.organizationId !== organization.id || project.tenantId !== tenant.id) {
    throw new Error("Active project does not belong to active organization and tenant");
  }
}

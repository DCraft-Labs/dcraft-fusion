import type { Id, RequestContext, TenantScopedResource } from "./tenancy.js";
export type AuditAction = "organization.created" | "tenant.created" | "project.created" | "integration.created" | "integration.tested" | "metadata.discovery.started" | "metadata.discovery.completed" | "workflow.run.observed" | "policy.approval.requested" | "ai.recommendation.created";
export interface AuditEvent {
    readonly id: Id;
    readonly action: AuditAction;
    readonly actorId: Id;
    readonly organizationId: Id;
    readonly tenantId: Id;
    readonly projectId: Id;
    readonly resourceType: string;
    readonly resourceId: Id;
    readonly correlationId: Id;
    readonly occurredAt: Date;
    readonly metadata: Readonly<Record<string, string | number | boolean>>;
}
export interface AuditEventInput {
    readonly id: Id;
    readonly action: AuditAction;
    readonly resourceType: string;
    readonly resource: TenantScopedResource;
    readonly occurredAt: Date;
    readonly metadata?: Readonly<Record<string, string | number | boolean | undefined>>;
}
export declare function createAuditEvent(context: RequestContext, input: AuditEventInput): AuditEvent;
//# sourceMappingURL=audit.d.ts.map
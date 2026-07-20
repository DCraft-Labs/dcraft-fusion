const secretLikePattern = /(secret|password|token|credential|private_key|api_key)/i;
export function createAuditEvent(context, input) {
    const metadata = Object.fromEntries(Object.entries(input.metadata ?? {}).filter(([key, value]) => value !== undefined && !secretLikePattern.test(key)));
    return {
        id: input.id,
        action: input.action,
        actorId: context.actorId,
        organizationId: context.organizationId,
        tenantId: context.tenantId,
        projectId: context.projectId,
        resourceType: input.resourceType,
        resourceId: input.resource.id,
        correlationId: context.correlationId,
        occurredAt: input.occurredAt,
        metadata
    };
}
//# sourceMappingURL=audit.js.map
export function requireRequestContext(context) {
    const missing = [
        ["actorId", context.actorId],
        ["organizationId", context.organizationId],
        ["tenantId", context.tenantId],
        ["projectId", context.projectId],
        ["environment", context.environment],
        ["correlationId", context.correlationId]
    ]
        .filter(([, value]) => value === undefined || value === "")
        .map(([key]) => key);
    if (missing.length > 0) {
        throw new Error(`Missing request context: ${missing.join(", ")}`);
    }
    return {
        actorId: context.actorId,
        organizationId: context.organizationId,
        tenantId: context.tenantId,
        projectId: context.projectId,
        environment: context.environment,
        roles: context.roles ?? [],
        correlationId: context.correlationId
    };
}
export function isCustomerOwnedDeployment(profile) {
    return profile === "customer_vpc" || profile === "on_prem";
}
export function isDedicatedDeployment(profile) {
    return profile === "dedicated_database" || profile === "dedicated_stack" || isCustomerOwnedDeployment(profile);
}
export function assertSameScope(context, resource) {
    if (context.organizationId !== resource.organizationId ||
        context.tenantId !== resource.tenantId ||
        context.projectId !== resource.projectId) {
        throw new Error("Tenant scope violation");
    }
}
//# sourceMappingURL=tenancy.js.map
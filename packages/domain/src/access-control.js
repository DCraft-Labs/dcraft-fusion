const roleRank = {
    owner: 50,
    admin: 40,
    operator: 30,
    service_account: 20,
    viewer: 10
};
export function hasRole(context, required) {
    const requiredRank = roleRank[required];
    return context.roles.some((role) => roleRank[role] >= requiredRank);
}
export function assertRole(context, required) {
    if (!hasRole(context, required)) {
        throw new Error(`Insufficient role: ${required} required`);
    }
}
export function membershipAppliesToContext(membership, context) {
    if (membership.organizationId !== context.organizationId) {
        return false;
    }
    if (membership.tenantId !== undefined && membership.tenantId !== context.tenantId) {
        return false;
    }
    if (membership.projectId !== undefined && membership.projectId !== context.projectId) {
        return false;
    }
    return true;
}
export function rolesForContext(memberships, context) {
    return memberships
        .filter((membership) => membership.userId === context.actorId)
        .filter((membership) => membershipAppliesToContext(membership, context))
        .map((membership) => membership.role);
}
//# sourceMappingURL=access-control.js.map
import type { Membership, RequestContext, Role } from "./tenancy.js";

const roleRank: Record<Role, number> = {
  owner: 50,
  admin: 40,
  operator: 30,
  service_account: 20,
  viewer: 10
};

export function hasRole(context: RequestContext, required: Role): boolean {
  const requiredRank = roleRank[required];
  return context.roles.some((role) => roleRank[role] >= requiredRank);
}

export function assertRole(context: RequestContext, required: Role): void {
  if (!hasRole(context, required)) {
    throw new Error(`Insufficient role: ${required} required`);
  }
}

export function membershipAppliesToContext(membership: Membership, context: RequestContext): boolean {
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

export function rolesForContext(memberships: readonly Membership[], context: RequestContext): Role[] {
  return memberships
    .filter((membership) => membership.userId === context.actorId)
    .filter((membership) => membershipAppliesToContext(membership, context))
    .map((membership) => membership.role);
}


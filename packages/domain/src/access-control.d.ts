import type { Membership, RequestContext, Role } from "./tenancy.js";
export declare function hasRole(context: RequestContext, required: Role): boolean;
export declare function assertRole(context: RequestContext, required: Role): void;
export declare function membershipAppliesToContext(membership: Membership, context: RequestContext): boolean;
export declare function rolesForContext(memberships: readonly Membership[], context: RequestContext): Role[];
//# sourceMappingURL=access-control.d.ts.map
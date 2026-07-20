export type Id = string;
export type OrganizationType = "business" | "consumer_app" | "managed_service_provider" | "internal_platform";
export type TenantType = "business_unit" | "brand" | "region" | "client" | "consumer_segment" | "environment" | "partner";
export type EnvironmentType = "dev" | "stage" | "prod" | "sandbox";
export type DeploymentProfile = "pooled_saas" | "dedicated_database" | "dedicated_stack" | "customer_vpc" | "on_prem";
export type DataPlaneLocation = "fusion_managed" | "customer_managed" | "hybrid";
export interface Organization {
    readonly id: Id;
    readonly name: string;
    readonly type: OrganizationType;
    readonly deploymentProfile: DeploymentProfile;
    readonly region: string;
    readonly dataPlaneLocation: DataPlaneLocation;
    readonly rawDataMovementAllowed: boolean;
}
export interface Tenant {
    readonly id: Id;
    readonly organizationId: Id;
    readonly name: string;
    readonly type: TenantType;
    readonly parentTenantId?: Id;
}
export interface Project {
    readonly id: Id;
    readonly organizationId: Id;
    readonly tenantId: Id;
    readonly name: string;
    readonly environment: EnvironmentType;
}
export type Role = "owner" | "admin" | "operator" | "viewer" | "service_account";
export interface Membership {
    readonly userId: Id;
    readonly organizationId: Id;
    readonly tenantId?: Id;
    readonly projectId?: Id;
    readonly role: Role;
}
export interface RequestContext {
    readonly actorId: Id;
    readonly organizationId: Id;
    readonly tenantId: Id;
    readonly projectId: Id;
    readonly environment: EnvironmentType;
    readonly roles: readonly Role[];
    readonly correlationId: Id;
}
export interface TenantScopedResource {
    readonly id: Id;
    readonly organizationId: Id;
    readonly tenantId: Id;
    readonly projectId: Id;
}
export declare function requireRequestContext(context: Partial<RequestContext>): RequestContext;
export declare function isCustomerOwnedDeployment(profile: DeploymentProfile): boolean;
export declare function isDedicatedDeployment(profile: DeploymentProfile): boolean;
export declare function assertSameScope(context: RequestContext, resource: TenantScopedResource): void;
//# sourceMappingURL=tenancy.d.ts.map
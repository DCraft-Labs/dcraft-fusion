import { describe, expect, it } from "vitest";
import {
  assertSameScope,
  isCustomerOwnedDeployment,
  isDedicatedDeployment,
  requireRequestContext,
  type RequestContext,
  type TenantScopedResource
} from "./tenancy.js";

describe("request context", () => {
  it("requires actor, organization, tenant, project, environment, and correlation context", () => {
    expect(() => requireRequestContext({ actorId: "user-1" })).toThrow(
      "Missing request context: organizationId, tenantId, projectId, environment, correlationId"
    );
  });

  it("defaults roles to an empty list only after required context is present", () => {
    const context = requireRequestContext({
      actorId: "user-1",
      organizationId: "org-1",
      tenantId: "tenant-1",
      projectId: "project-1",
      environment: "prod",
      correlationId: "corr-1"
    });

    expect(context.roles).toEqual([]);
  });
});

describe("deployment profile helpers", () => {
  it("classifies customer-owned deployment profiles", () => {
    expect(isCustomerOwnedDeployment("customer_vpc")).toBe(true);
    expect(isCustomerOwnedDeployment("on_prem")).toBe(true);
    expect(isCustomerOwnedDeployment("pooled_saas")).toBe(false);
  });

  it("classifies dedicated deployment profiles", () => {
    expect(isDedicatedDeployment("dedicated_database")).toBe(true);
    expect(isDedicatedDeployment("dedicated_stack")).toBe(true);
    expect(isDedicatedDeployment("customer_vpc")).toBe(true);
    expect(isDedicatedDeployment("pooled_saas")).toBe(false);
  });
});

describe("scope enforcement", () => {
  const context: RequestContext = {
    actorId: "user-1",
    organizationId: "org-1",
    tenantId: "tenant-1",
    projectId: "project-1",
    environment: "prod",
    roles: ["admin"],
    correlationId: "corr-1"
  };

  it("allows same organization, tenant, and project", () => {
    const resource: TenantScopedResource = {
      id: "connection-1",
      organizationId: "org-1",
      tenantId: "tenant-1",
      projectId: "project-1"
    };

    expect(() => assertSameScope(context, resource)).not.toThrow();
  });

  it("rejects cross-tenant access", () => {
    const resource: TenantScopedResource = {
      id: "connection-1",
      organizationId: "org-1",
      tenantId: "tenant-2",
      projectId: "project-1"
    };

    expect(() => assertSameScope(context, resource)).toThrow("Tenant scope violation");
  });
});


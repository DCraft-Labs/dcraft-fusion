import { describe, expect, it } from "vitest";
import { assertRole, hasRole, rolesForContext } from "./access-control.js";
import type { Membership, RequestContext } from "./tenancy.js";

const context: RequestContext = {
  actorId: "user-1",
  organizationId: "org-1",
  tenantId: "tenant-1",
  projectId: "project-1",
  environment: "prod",
  roles: ["operator"],
  correlationId: "corr-1"
};

describe("role checks", () => {
  it("allows higher-ranked roles to satisfy lower-ranked requirements", () => {
    expect(hasRole({ ...context, roles: ["admin"] }, "viewer")).toBe(true);
    expect(hasRole({ ...context, roles: ["viewer"] }, "admin")).toBe(false);
  });

  it("throws when the actor lacks the required role", () => {
    expect(() => assertRole(context, "admin")).toThrow("Insufficient role: admin required");
  });
});

describe("membership scoping", () => {
  it("collects only memberships that apply to active organization, tenant, and project", () => {
    const memberships: Membership[] = [
      { userId: "user-1", organizationId: "org-1", role: "viewer" },
      { userId: "user-1", organizationId: "org-1", tenantId: "tenant-1", role: "operator" },
      { userId: "user-1", organizationId: "org-1", tenantId: "tenant-2", role: "admin" },
      { userId: "user-2", organizationId: "org-1", tenantId: "tenant-1", role: "owner" }
    ];

    expect(rolesForContext(memberships, context)).toEqual(["viewer", "operator"]);
  });
});


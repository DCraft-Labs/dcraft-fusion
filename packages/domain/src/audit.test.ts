import { describe, expect, it } from "vitest";
import { createAuditEvent } from "./audit.js";
import type { RequestContext, TenantScopedResource } from "./tenancy.js";

const context: RequestContext = {
  actorId: "user-1",
  organizationId: "org-1",
  tenantId: "tenant-1",
  projectId: "project-1",
  environment: "prod",
  roles: ["admin"],
  correlationId: "corr-1"
};

const resource: TenantScopedResource = {
  id: "connection-1",
  organizationId: "org-1",
  tenantId: "tenant-1",
  projectId: "project-1"
};

describe("audit events", () => {
  it("copies actor and tenant context into immutable audit events", () => {
    const event = createAuditEvent(context, {
      id: "audit-1",
      action: "integration.created",
      resourceType: "Connection",
      resource,
      occurredAt: new Date("2026-05-18T00:00:00.000Z"),
      metadata: { connectorType: "postgres" }
    });

    expect(event).toMatchObject({
      id: "audit-1",
      actorId: "user-1",
      organizationId: "org-1",
      tenantId: "tenant-1",
      projectId: "project-1",
      resourceId: "connection-1",
      correlationId: "corr-1",
      metadata: { connectorType: "postgres" }
    });
  });

  it("redacts secret-like metadata keys", () => {
    const event = createAuditEvent(context, {
      id: "audit-1",
      action: "integration.created",
      resourceType: "Connection",
      resource,
      occurredAt: new Date("2026-05-18T00:00:00.000Z"),
      metadata: {
        connectorType: "postgres",
        password: "should-not-be-kept",
        api_token: "should-not-be-kept",
        secretReference: "should-not-be-kept"
      }
    });

    expect(event.metadata).toEqual({ connectorType: "postgres" });
  });
});


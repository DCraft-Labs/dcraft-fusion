import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { BrowserRouter } from "react-router-dom";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { App } from "./App.js";
import { createDemoWorkspace } from "./demo-workspace.js";
import type { FusionWorkspace, UserSession } from "./workspace.js";

describe("Fusion UI shell", () => {
  beforeEach(() => {
    window.sessionStorage.clear();
    window.history.pushState({}, "", "/login");
  });

  afterEach(() => {
    vi.restoreAllMocks();
    window.sessionStorage.clear();
  });

  it("shows simple community password login", async () => {
    vi.stubGlobal("fetch", vi.fn((input: RequestInfo | URL) => Promise.resolve(mockResponse(input.toString()))));

    renderWithRouter();

    expect(await screen.findByRole("heading", { name: "Sign in" })).toBeInTheDocument();
    expect(screen.getByRole("textbox", { name: "Email" })).toBeInTheDocument();
    expect(screen.getByLabelText("Password")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Sign in" })).toBeInTheDocument();
    expect(screen.getByText(/environment variables/i)).toBeInTheDocument();
  });

  it("updates the browser URL when workspace navigation changes", async () => {
    seedSession({
      actorId: "user-founder",
      email: "founder@dcraftlabs.com",
      label: "Founder Admin",
      scope: "tenant",
      organizationId: "org-b2b2c-1",
      tenantId: "tenant-brand-a",
      projectId: "project-prod"
    });
    window.history.pushState({}, "", "/workspace/command-center");
    vi.stubGlobal("fetch", vi.fn((input: RequestInfo | URL) => Promise.resolve(mockResponse(input.toString()))));

    renderWithRouter();

    expect(await screen.findByRole("heading", { name: "Command Center" })).toBeInTheDocument();
    await userEvent.click(screen.getByRole("link", { name: "Integration Hub" }));

    await waitFor(() => expect(window.location.pathname).toBe("/workspace/integrations"));
    expect(screen.getByRole("heading", { name: "Integration Hub" })).toBeInTheDocument();
    expect(screen.getByRole("heading", { name: "Postgres" })).toBeInTheDocument();
  });

  it("shows routed superadmin screens instead of a single info page", async () => {
    seedSession({
      actorId: "user-superadmin",
      email: "superadmin@dcraftlabs.com",
      label: "Platform Superadmin",
      scope: "platform"
    });
    window.history.pushState({}, "", "/superadmin/overview");
    vi.stubGlobal(
      "fetch",
      vi.fn((input: RequestInfo | URL) => {
        const url = typeof input === "string" ? input : input instanceof URL ? input.href : input.url;
        return Promise.resolve(mockResponse(url));
      })
    );

    renderWithRouter();

    expect(await screen.findByRole("heading", { name: "Superadmin Overview" })).toBeInTheDocument();
    await userEvent.click(screen.getByRole("link", { name: "Organizations" }));

    await waitFor(() => expect(window.location.pathname).toBe("/superadmin/organizations"));
    expect(screen.getByRole("heading", { name: "Organizations" })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Create Organization" })).toBeInTheDocument();
    expect(screen.getByRole("heading", { name: "Commerce Platform Co." })).toBeInTheDocument();
  });

  it("rejects a workspace with mismatched tenant and project context", () => {
    const invalidWorkspace: FusionWorkspace = {
      ...createDemoWorkspace(),
      context: {
        ...createDemoWorkspace().context,
        project: {
          id: "project-bad",
          organizationId: "org-b2b2c-1",
          tenantId: "tenant-other",
          name: "Bad Project",
          environment: "prod"
        }
      }
    };

    expect(() =>
      render(
        <BrowserRouter>
          <App workspace={invalidWorkspace} />
        </BrowserRouter>
      )
    ).toThrow("Active project does not belong to active organization and tenant");
  });
});

function renderWithRouter() {
  render(
    <BrowserRouter>
      <App workspace={createDemoWorkspace()} />
    </BrowserRouter>
  );
}

function seedSession(session: UserSession) {
  window.sessionStorage.setItem("fusion.auth.session", JSON.stringify(session));
  window.sessionStorage.setItem("fusion.auth.token", "fake-token");
}

function mockResponse(url: string): Response {
  if (url.endsWith("/api/v1/auth/users")) {
    return jsonResponse({
      users: [
        { id: "user-founder", email: "founder@dcraftlabs.com", label: "Founder Admin", scope: "tenant" },
        { id: "user-superadmin", email: "superadmin@dcraftlabs.com", label: "Platform Superadmin", scope: "platform" }
      ]
    });
  }
  if (url.endsWith("/.well-known/openid-configuration")) {
    return jsonResponse({
      issuer: "http://127.0.0.1:30173",
      authorization_endpoint: "http://127.0.0.1:30173/oidc/authorize",
      token_endpoint: "http://127.0.0.1:30173/oidc/token",
      jwks_uri: "http://127.0.0.1:30173/oidc/jwks"
    });
  }
  if (url.endsWith("/api/v1/bootstrap")) {
    return jsonResponse({
      connections: createDemoWorkspace().connections,
      datasets: createDemoWorkspace().datasets,
      runs: createDemoWorkspace().runs,
      policies: createDemoWorkspace().policies,
      aiRecommendations: createDemoWorkspace().aiRecommendations
    });
  }
  if (url.endsWith("/api/v1/platform/overview")) {
    return jsonResponse({
      identityProviders: [{ id: "idp-platform-oidc", name: "Platform OIDC", protocol: "oidc", issuer: "http://127.0.0.1:30173", audience: "dcraft-fusion", jwksUrl: "http://127.0.0.1:30173/oidc/jwks", ssoEnabled: true, samlEnabled: false, status: "configured" }],
      secretProviders: [{ id: "secret-vault", name: "HashiCorp Vault", kind: "vault", endpoint: "https://vault.example.com", scope: "platform", status: "configured" }],
      adapters: [{ id: "adapter-airflow", name: "Airflow", kind: "airflow", readOnly: true, status: "implemented", description: "Read-only DAG ingestion." }],
      backup: { target: "postgres", schedule: "0 2 * * *", status: "scheduled", runbook: "docs/operations/BACKUP_RESTORE_RUNBOOK.md" },
      observability: { tracesEnabled: true, metricsEndpoint: "/metrics", logFormat: "json", dashboards: [], alerts: [], openTelemetryReady: true },
      ha: { apiMinReplicas: 2, webMinReplicas: 2, autoscaling: true, podDisruption: true, loadTests: [], status: "configured" },
      jobs: [{ id: "job-1", type: "metadata.discovery", status: "succeeded", resourceId: "conn-postgres-brand-a", message: "Initial metadata discovery completed.", updatedAt: "2026-05-18T00:00:00Z" }],
      superAdmins: [{ id: "mem-superadmin", userId: "user-superadmin", email: "superadmin@dcraftlabs.com", organizationId: "platform", role: "owner" }],
      organizations: [
        {
          id: "org-b2b2c-1",
          name: "Commerce Platform Co.",
          type: "consumer_app",
          deploymentProfile: "customer_vpc",
          region: "us-east-1",
          dataPlaneLocation: "customer_managed",
          rawDataMovementAllowed: false,
          createdAt: "2026-05-18T00:00:00Z",
          tenants: [
            {
              id: "tenant-brand-a",
              name: "Brand A",
              model: "b2b2c",
              isolation: "customer_infra",
              region: "us-east-1",
              createdAt: "2026-05-18T00:00:00Z",
              projects: [{ id: "project-prod", name: "Production Data Platform", environment: "production", createdAt: "2026-05-18T00:00:00Z" }]
            }
          ]
        }
      ],
      recentAuditEvents: [{ id: "audit-1", action: "organization.created", actorId: "user-superadmin", correlationId: "corr-1", resourceType: "Organization", resourceId: "org-b2b2c-1", occurredAt: "2026-05-18T00:00:00Z" }]
    });
  }
  if (url.endsWith("/api/v1/connections") || url.includes("/api/v1/connections/")) {
    return jsonResponse({}, 200);
  }
  return jsonResponse({}, 200);
}

function jsonResponse(body: unknown, status = 200): Response {
  return new Response(JSON.stringify(body), {
    status,
    headers: { "Content-Type": "application/json" }
  });
}

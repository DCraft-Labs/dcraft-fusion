import { getAccessToken } from "./auth.js";
import type { Connection, FusionWorkspace, PlatformOverview } from "./workspace.js";

function headers(): HeadersInit {
  const token = getAccessToken();
  if (token === undefined) {
    throw new Error("Authentication is required");
  }
  return {
    "Content-Type": "application/json",
    Authorization: `Bearer ${token}`,
    "X-Fusion-Correlation-Id": `ui-${Date.now().toString()}`
  };
}

export async function createConnection(
  connection: Pick<Connection, "name" | "kind" | "readOnly"> & { readonly secretRef: string }
): Promise<Connection> {
  const response = await fetch("/api/v1/connections", {
    method: "POST",
    headers: headers(),
    body: JSON.stringify(connection)
  });
  if (!response.ok) {
    throw new Error(await response.text());
  }
  return (await response.json()) as Connection;
}

export async function testConnection(connectionID: string): Promise<void> {
  const response = await fetch(`/api/v1/connections/${connectionID}/test`, {
    method: "POST",
    headers: headers()
  });
  if (!response.ok) {
    throw new Error(await response.text());
  }
}

export async function loadBootstrap(): Promise<Partial<FusionWorkspace>> {
  const response = await fetch("/api/v1/bootstrap", {
    headers: headers()
  });
  if (!response.ok) {
    throw new Error(await response.text());
  }
  const body = (await response.json()) as {
    connections?: FusionWorkspace["connections"];
    datasets?: FusionWorkspace["datasets"];
    runs?: FusionWorkspace["runs"];
    policies?: FusionWorkspace["policies"];
    aiRecommendations?: FusionWorkspace["aiRecommendations"];
  };
  return {
    connections: body.connections ?? [],
    datasets: body.datasets ?? [],
    runs: body.runs ?? [],
    policies: body.policies ?? [],
    aiRecommendations: body.aiRecommendations ?? []
  };
}

export async function loadPlatformOverview(): Promise<PlatformOverview> {
  const response = await fetch("/api/v1/platform/overview", {
    headers: headers()
  });
  if (!response.ok) {
    throw new Error(await response.text());
  }
  return (await response.json()) as PlatformOverview;
}

export async function createPlatformOrganization(request: {
  readonly name: string;
  readonly type: string;
  readonly deploymentProfile: string;
  readonly region: string;
  readonly dataPlaneLocation: string;
  readonly rawDataMovementAllowed: boolean;
}): Promise<void> {
  const response = await fetch("/api/v1/platform/organizations", {
    method: "POST",
    headers: headers(),
    body: JSON.stringify(request)
  });
  if (!response.ok) {
    throw new Error(await response.text());
  }
}

export async function createPlatformTenant(request: {
  readonly organizationId: string;
  readonly name: string;
  readonly model: string;
  readonly isolation: string;
  readonly region: string;
}): Promise<void> {
  const response = await fetch("/api/v1/platform/tenants", {
    method: "POST",
    headers: headers(),
    body: JSON.stringify(request)
  });
  if (!response.ok) {
    throw new Error(await response.text());
  }
}

export async function createPlatformProject(request: {
  readonly organizationId: string;
  readonly tenantId: string;
  readonly name: string;
  readonly environment: string;
}): Promise<void> {
  const response = await fetch("/api/v1/platform/projects", {
    method: "POST",
    headers: headers(),
    body: JSON.stringify(request)
  });
  if (!response.ok) {
    throw new Error(await response.text());
  }
}

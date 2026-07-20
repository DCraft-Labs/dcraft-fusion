import type { AuthUser, OpenIDConfiguration, UserSession } from "./workspace.js";

const tokenKey = "fusion.auth.token";
const sessionKey = "fusion.auth.session";
const stateKey = "fusion.auth.state";

interface TokenClaims {
  readonly sub: string;
  readonly actor_id?: string;
  readonly email?: string;
  readonly name?: string;
  readonly scope?: "platform" | "tenant";
  readonly organization_id?: string;
  readonly tenant_id?: string;
  readonly project_id?: string;
}

interface LoginResponse {
  readonly ok?: boolean;
  readonly access_token?: string;
  readonly email?: string;
  readonly scope?: string;
}

export function getCurrentSession(): UserSession | undefined {
  const raw = window.sessionStorage.getItem(sessionKey);
  if (raw === null) return undefined;
  try {
    return JSON.parse(raw) as UserSession;
  } catch {
    clearSession();
    return undefined;
  }
}

export function getAccessToken(): string | undefined {
  return window.sessionStorage.getItem(tokenKey) ?? undefined;
}

export function clearSession(): void {
  window.sessionStorage.removeItem(tokenKey);
  window.sessionStorage.removeItem(sessionKey);
  window.sessionStorage.removeItem(stateKey);
}

export async function loadAvailableUsers(): Promise<readonly AuthUser[]> {
  const response = await fetch("/api/v1/auth/users");
  if (!response.ok) {
    throw new Error(await response.text());
  }
  const body = (await response.json()) as { readonly users?: readonly AuthUser[] };
  return body.users ?? [];
}

export async function loadOIDCConfiguration(): Promise<OpenIDConfiguration | undefined> {
  const response = await fetch("/.well-known/openid-configuration");
  if (!response.ok) {
    return undefined;
  }
  return (await response.json()) as OpenIDConfiguration;
}

/** Community login: password → JWT. Enterprise OIDC mode still returns cookie-only and callers may continue authorize. */
export async function loginWithPassword(email: string, password: string): Promise<UserSession | undefined> {
  const response = await fetch("/api/v1/auth/login", {
    method: "POST",
    credentials: "include",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password })
  });
  if (!response.ok) {
    throw new Error(await response.text());
  }
  const body = (await response.json()) as LoginResponse;
  if (body.access_token) {
    const session = sessionFromToken(body.access_token);
    window.sessionStorage.setItem(tokenKey, body.access_token);
    window.sessionStorage.setItem(sessionKey, JSON.stringify(session));
    return session;
  }
  return undefined;
}

export async function logoutFromProvider(): Promise<void> {
  await fetch("/api/v1/auth/logout", {
    method: "POST",
    credentials: "include"
  });
}

export function beginAuthorization(): void {
  const state = crypto.randomUUID();
  window.sessionStorage.setItem(stateKey, state);
  const redirectURI = new URL("/auth/callback", window.location.origin).toString();
  const authorizeURL = new URL("/oidc/authorize", window.location.origin);
  authorizeURL.searchParams.set("client_id", "fusion-web");
  authorizeURL.searchParams.set("redirect_uri", redirectURI);
  authorizeURL.searchParams.set("response_type", "code");
  authorizeURL.searchParams.set("scope", "openid profile email");
  authorizeURL.searchParams.set("state", state);
  window.location.assign(authorizeURL.toString());
}

export async function finishLogin(code: string, state: string): Promise<UserSession> {
  const expectedState = window.sessionStorage.getItem(stateKey);
  if (expectedState === null || expectedState !== state) {
    throw new Error("OIDC state validation failed");
  }
  const redirectURI = new URL("/auth/callback", window.location.origin).toString();
  const response = await fetch("/oidc/token", {
    method: "POST",
    headers: { "Content-Type": "application/x-www-form-urlencoded" },
    body: new URLSearchParams({
      grant_type: "authorization_code",
      client_id: "fusion-web",
      redirect_uri: redirectURI,
      code
    }).toString()
  });
  if (!response.ok) {
    throw new Error(await response.text());
  }
  const body = (await response.json()) as { readonly access_token: string };
  const session = sessionFromToken(body.access_token);
  window.sessionStorage.setItem(tokenKey, body.access_token);
  window.sessionStorage.setItem(sessionKey, JSON.stringify(session));
  window.sessionStorage.removeItem(stateKey);
  return session;
}

export function defaultRouteFor(session: UserSession): string {
  return session.scope === "platform" ? "/superadmin/overview" : "/workspace/command-center";
}

function sessionFromToken(token: string): UserSession {
  const claims = parseTokenClaims(token);
  return {
    actorId: claims.actor_id ?? claims.sub,
    email: claims.email ?? "unknown@local",
    label: claims.name ?? claims.actor_id ?? claims.sub,
    scope: claims.scope ?? "tenant",
    ...(claims.organization_id !== undefined ? { organizationId: claims.organization_id } : {}),
    ...(claims.tenant_id !== undefined ? { tenantId: claims.tenant_id } : {}),
    ...(claims.project_id !== undefined ? { projectId: claims.project_id } : {})
  };
}

function parseTokenClaims(token: string): TokenClaims {
  const parts = token.split(".");
  if (parts.length !== 3) {
    throw new Error("Invalid access token");
  }
  const claimsSegment = parts[1];
  if (claimsSegment === undefined) {
    throw new Error("Invalid access token");
  }
  const base64 = claimsSegment.replaceAll("-", "+").replaceAll("_", "/");
  const padding = base64.length % 4 === 0 ? "" : "=".repeat(4 - (base64.length % 4));
  const raw = window.atob(base64 + padding);
  return JSON.parse(raw) as TokenClaims;
}

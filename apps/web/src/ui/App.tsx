import { useEffect, useMemo, useRef, useState } from "react";
import { Link, NavLink, Navigate, Outlet, Route, Routes, useNavigate, useSearchParams } from "react-router-dom";
import {
  beginAuthorization,
  clearSession,
  defaultRouteFor,
  finishLogin,
  getCurrentSession,
  loadAvailableUsers,
  loadOIDCConfiguration,
  loginWithPassword,
  logoutFromProvider
} from "./auth.js";
import {
  createConnection as createConnectionAPI,
  createPlatformOrganization,
  createPlatformProject,
  createPlatformTenant,
  loadBootstrap,
  loadPlatformOverview,
  testConnection as testConnectionAPI
} from "./api.js";
import type {
  AuthUser,
  Connection,
  FusionWorkspace,
  IntegrationCard,
  OpenIDConfiguration,
  PlatformOrganization,
  PlatformOverview,
  UserSession
} from "./workspace.js";
import { validateWorkspaceContext } from "./workspace.js";

const workspaceNav = [
  { to: "/workspace/command-center", label: "Command Center" },
  { to: "/workspace/integrations", label: "Integration Hub" },
  { to: "/workspace/metadata", label: "Metadata Explorer" },
  { to: "/workspace/runs", label: "Workflow Runs" },
  { to: "/workspace/audit", label: "Audit Center" },
  { to: "/workspace/settings", label: "Settings" },
  { to: "/workspace/ai", label: "AI Operator" }
] as const;

const superadminNav = [
  { to: "/superadmin/overview", label: "Overview" },
  { to: "/superadmin/organizations", label: "Organizations" },
  { to: "/superadmin/identity", label: "Identity" },
  { to: "/superadmin/secrets", label: "Secrets" },
  { to: "/superadmin/adapters", label: "Adapters" },
  { to: "/superadmin/operations", label: "Operations" }
] as const;

const rolloutPhases = [
  {
    label: "Phase 1",
    title: "Private alpha control plane",
    body: "Read-only connectors, metadata discovery, audit, AI draft workflows, tenant and platform admin foundations."
  },
  {
    label: "Phase 2",
    title: "Identity and workflow hardening",
    body: "External IdP onboarding, adapter expansion, deeper policy controls, stronger platform operations and managed rollout flows."
  },
  {
    label: "Phase 3",
    title: "Enterprise scale posture",
    body: "Dedicated customer deployments, external secret orchestration, stronger observability and workload expansion beyond alpha."
  }
] as const;

const deploymentCards = [
  { title: "B2B", body: "Dedicated business tenants with shared, dedicated, or customer-managed infrastructure patterns." },
  { title: "B2C", body: "High-volume product operations with platform-level controls and project-by-project rollout." },
  { title: "B2B2C", body: "Brand and end-customer segmentation across shared control and isolated data-plane boundaries." },
  { title: "Dedicated", body: "Customer VPC, dedicated stack, or on-prem deployment topologies from day one." }
] as const;

type Theme = "dark" | "light";

export function App({ workspace }: { readonly workspace: FusionWorkspace }) {
  validateWorkspaceContext(workspace);

  const [theme, setTheme] = useState<Theme>(() => {
    const stored = window.localStorage.getItem("fusion.theme");
    if (stored === "dark" || stored === "light") return stored;
    return "light";
  });
  const [session, setSession] = useState<UserSession | undefined>(() => getCurrentSession());
  const [authUsers, setAuthUsers] = useState<readonly AuthUser[]>([]);
  const [authConfig, setAuthConfig] = useState<OpenIDConfiguration | undefined>(undefined);
  const [authError, setAuthError] = useState<string | undefined>(undefined);
  const [liveWorkspace, setLiveWorkspace] = useState<FusionWorkspace>(workspace);
  const [connections, setConnections] = useState<readonly Connection[]>(workspace.connections);
  const [platformOverview, setPlatformOverview] = useState<PlatformOverview | undefined>(undefined);
  const [platformError, setPlatformError] = useState<string | undefined>(undefined);
  const [platformLoading, setPlatformLoading] = useState(false);

  const workspaceContext = useMemo(() => liveWorkspace.context, [liveWorkspace.context]);

  useEffect(() => {
    document.documentElement.dataset["theme"] = theme;
    window.localStorage.setItem("fusion.theme", theme);
  }, [theme]);

  useEffect(() => {
    let cancelled = false;
    setAuthError(undefined);
    Promise.all([loadAvailableUsers(), loadOIDCConfiguration()])
      .then(([users, configuration]) => {
        if (cancelled) return;
        setAuthUsers(users);
        setAuthConfig(configuration);
      })
      .catch((error: unknown) => {
        if (!cancelled) {
          setAuthError(error instanceof Error ? error.message : "Auth bootstrap failed");
        }
      });
    return () => {
      cancelled = true;
    };
  }, []);

  useEffect(() => {
    if (session?.scope !== "tenant") return;
    let cancelled = false;
    loadBootstrap()
      .then((bootstrap) => {
        if (cancelled) return;
        setLiveWorkspace((current) => ({ ...current, ...bootstrap }));
        if (bootstrap.connections !== undefined) {
          setConnections(bootstrap.connections);
        }
      })
      .catch(() => {
        setLiveWorkspace(workspace);
        setConnections(workspace.connections);
      });
    return () => {
      cancelled = true;
    };
  }, [session, workspace]);

  useEffect(() => {
    if (session?.scope !== "platform") {
      setPlatformOverview(undefined);
      setPlatformError(undefined);
      return;
    }
    void refreshPlatformOverview();
  }, [session]);

  async function refreshPlatformOverview() {
    setPlatformLoading(true);
    setPlatformError(undefined);
    try {
      const overview = await loadPlatformOverview();
      setPlatformOverview(overview);
    } catch (error: unknown) {
      setPlatformOverview(undefined);
      setPlatformError(error instanceof Error ? error.message : "Platform overview is unavailable");
    } finally {
      setPlatformLoading(false);
    }
  }

  async function signOut() {
    clearSession();
    setSession(undefined);
    setLiveWorkspace(workspace);
    setConnections(workspace.connections);
    setPlatformOverview(undefined);
    setPlatformError(undefined);
    await logoutFromProvider();
  }

  return (
    <Routes>
      <Route
        path="/"
        element={<LandingPage authConfig={authConfig} session={session} theme={theme} setTheme={setTheme} />}
      />
      <Route
        path="/login"
        element={
          session ? (
            <Navigate replace to={defaultRouteFor(session)} />
          ) : (
            <LoginScreen
              authConfig={authConfig}
              authError={authError}
              onAuthenticated={setSession}
              theme={theme}
              setTheme={setTheme}
              users={authUsers}
            />
          )
        }
      />
      <Route path="/auth/callback" element={<AuthCallback onAuthenticated={setSession} theme={theme} setTheme={setTheme} />} />
      <Route element={<RequireSession session={session} scope="tenant" />}>
        <Route
          path="/workspace"
          element={<WorkspaceLayout context={workspaceContext} session={session} signOut={signOut} theme={theme} setTheme={setTheme} />}
        >
          <Route index element={<Navigate replace to="/workspace/command-center" />} />
          <Route path="command-center" element={<CommandCenter workspace={liveWorkspace} connections={connections} />} />
          <Route
            path="integrations"
            element={<IntegrationHub workspace={liveWorkspace} connections={connections} setConnections={setConnections} />}
          />
          <Route path="metadata" element={<MetadataExplorer workspace={liveWorkspace} />} />
          <Route path="runs" element={<RunCenter workspace={liveWorkspace} />} />
          <Route path="audit" element={<AuditCenter workspace={liveWorkspace} />} />
          <Route path="settings" element={<Settings workspace={liveWorkspace} />} />
          <Route path="ai" element={<AIOperator workspace={liveWorkspace} />} />
        </Route>
      </Route>
      <Route element={<RequireSession session={session} scope="platform" />}>
        <Route
          path="/superadmin"
          element={
            <SuperadminLayout
              error={platformError}
              loading={platformLoading}
              overview={platformOverview}
              session={session}
              signOut={signOut}
              theme={theme}
              setTheme={setTheme}
            />
          }
        >
          <Route index element={<Navigate replace to="/superadmin/overview" />} />
          <Route
            path="overview"
            element={<SuperadminOverviewScreen authConfig={authConfig} error={platformError} overview={platformOverview} />}
          />
          <Route path="organizations" element={<OrganizationsScreen overview={platformOverview} onRefresh={refreshPlatformOverview} />} />
          <Route path="identity" element={<IdentityScreen authConfig={authConfig} overview={platformOverview} users={authUsers} />} />
          <Route path="secrets" element={<SecretsScreen overview={platformOverview} />} />
          <Route path="adapters" element={<AdaptersScreen overview={platformOverview} />} />
          <Route path="operations" element={<OperationsScreen overview={platformOverview} />} />
        </Route>
      </Route>
      <Route path="*" element={<HomeRedirect session={session} />} />
    </Routes>
  );
}

function HomeRedirect({ session }: { readonly session: UserSession | undefined }) {
  if (session === undefined) {
    return <Navigate replace to="/" />;
  }
  return <Navigate replace to={defaultRouteFor(session)} />;
}

function RequireSession({
  session,
  scope
}: {
  readonly session: UserSession | undefined;
  readonly scope: UserSession["scope"];
}) {
  if (session === undefined) {
    return <Navigate replace to="/login" />;
  }
  if (session.scope !== scope) {
    return <Navigate replace to={defaultRouteFor(session)} />;
  }
  return <Outlet />;
}

function LandingPage({
  session,
  theme,
  setTheme
}: {
  readonly authConfig: OpenIDConfiguration | undefined;
  readonly session: UserSession | undefined;
  readonly theme: Theme;
  readonly setTheme: (theme: Theme) => void;
}) {
  return (
    <main className="landing-shell">
      <div className="landing-mesh" aria-hidden="true" />
      <header className="public-nav">
        <a className="brand-lockup" href="/">
          <img alt="" className="brand-mark" height={36} src="/logo.svg" width={36} />
          <span className="brand-word">DCraft Fusion</span>
        </a>
        <div className="public-actions">
          <a className="text-link" href="https://github.com/DCraft-Labs/dcraft-fusion" rel="noreferrer" target="_blank">
            GitHub
          </a>
          <a className="text-link" href="/docs/">
            Docs
          </a>
          <ThemeToggle theme={theme} setTheme={setTheme} />
          <Link className="btn btn-secondary" to={session ? defaultRouteFor(session) : "/login"}>
            {session ? "Open console" : "Sign in"}
          </Link>
        </div>
      </header>
      <section className="landing-hero">
        <div className="landing-copy reveal">
          <p className="brand-kicker">DCraft Fusion</p>
          <h1 className="hero-title">One brain. Many muscles.</h1>
          <p className="hero-copy">
            The open-source control plane for modern data platforms — govern connections, metadata, workflows, audit, and CDC from one operator surface.
          </p>
          <div className="cta-row">
            <Link className="btn btn-primary" to={session ? defaultRouteFor(session) : "/login"}>
              {session ? "Enter workspace" : "Get started"}
            </Link>
            <a className="btn btn-ghost" href="https://github.com/DCraft-Labs/dcraft-fusion" rel="noreferrer" target="_blank">
              Star on GitHub
            </a>
          </div>
        </div>
        <div className="hero-visual reveal reveal-delay" aria-hidden="true">
          <div className="hero-schematic">
            <div className="schematic-node schematic-core">
              <span>Control plane</span>
              <strong>Fusion Kernel</strong>
            </div>
            <div className="schematic-rails">
              <div className="schematic-node">
                <span>Engine</span>
                <strong>CDC</strong>
              </div>
              <div className="schematic-node">
                <span>Engine</span>
                <strong>dbt / Airflow</strong>
              </div>
              <div className="schematic-node">
                <span>Observe</span>
                <strong>Audit · Runs</strong>
              </div>
            </div>
          </div>
        </div>
      </section>
      <section className="landing-section" id="capabilities">
        <div className="section-intro">
          <h2>Built for operators who run real systems</h2>
          <p>Community edition includes the Fusion control plane and Fusion CDC Engine under Apache 2.0.</p>
        </div>
        <div className="capability-grid">
          {deploymentCards.map((item) => (
            <article className="capability" key={item.title}>
              <h3>{item.title}</h3>
              <p>{item.body}</p>
            </article>
          ))}
        </div>
      </section>
      <section className="landing-section" id="rollout">
        <div className="section-intro">
          <h2>Ship, observe, harden</h2>
          <p>A clear path from local alpha to enterprise posture — without burying the product under dashboard chrome.</p>
        </div>
        <div className="timeline-grid">
          {rolloutPhases.map((phase) => (
            <article className="timeline-item" key={phase.label}>
              <span className="timeline-label">{phase.label}</span>
              <h3>{phase.title}</h3>
              <p>{phase.body}</p>
            </article>
          ))}
        </div>
      </section>
      <footer className="landing-footer">
        <img alt="" height={28} src="/logo.svg" width={28} />
        <span>DCraft Labs · Apache 2.0 · One brain. Many muscles.</span>
      </footer>
    </main>
  );
}

function LoginScreen({
  authError,
  theme,
  setTheme,
  onAuthenticated
}: {
  readonly authConfig: OpenIDConfiguration | undefined;
  readonly authError: string | undefined;
  readonly theme: Theme;
  readonly setTheme: (theme: Theme) => void;
  readonly users: readonly AuthUser[];
  readonly onAuthenticated: (session: UserSession) => void;
}) {
  const navigate = useNavigate();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [loginError, setLoginError] = useState<string | undefined>(undefined);

  async function submit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setSubmitting(true);
    setLoginError(undefined);
    try {
      const session = await loginWithPassword(email, password);
      if (session) {
        onAuthenticated(session);
        navigate(defaultRouteFor(session), { replace: true });
        return;
      }
      // Enterprise OIDC mode: cookie session established — continue authorize redirect.
      beginAuthorization();
    } catch (error: unknown) {
      setLoginError(error instanceof Error ? error.message : "Login failed");
      setSubmitting(false);
    }
  }

  return (
    <main className="login-shell">
      <header className="public-nav compact-nav">
        <Link className="brand-lockup" to="/">
          <img alt="" className="brand-mark" height={32} src="/logo.svg" width={32} />
          <span className="brand-word">DCraft Fusion</span>
        </Link>
        <ThemeToggle theme={theme} setTheme={setTheme} />
      </header>
      <section className="login-card reveal">
        <p className="brand-kicker">Community access</p>
        <h1>Sign in</h1>
        <p className="hero-copy">Email and password. No OIDC redirect for the open-source default.</p>
        {authError ? <div className="message-banner error">{authError}</div> : null}
        {loginError ? <div className="message-banner error">{loginError}</div> : null}
        <form className="auth-form" onSubmit={(event) => void submit(event)}>
          <label>
            Email
            <input autoComplete="username" placeholder="you@company.com" value={email} onChange={(event) => setEmail(event.target.value)} />
          </label>
          <label>
            Password
            <input autoComplete="current-password" type="password" value={password} onChange={(event) => setPassword(event.target.value)} />
          </label>
          <button className="btn btn-primary btn-block" disabled={submitting} type="submit">
            {submitting ? "Signing in…" : "Sign in"}
          </button>
        </form>
        <p className="login-hint">
          Local seeds are configured via environment variables. See <code>infra/local-dev/.env.example</code>.
        </p>
      </section>
    </main>
  );
}

function AuthCallback({
  onAuthenticated,
  theme,
  setTheme
}: {
  readonly onAuthenticated: (session: UserSession) => void;
  readonly theme: Theme;
  readonly setTheme: (theme: Theme) => void;
}) {
  const [params] = useSearchParams();
  const navigate = useNavigate();
  const [error, setError] = useState<string | undefined>(undefined);
  const started = useRef(false);

  useEffect(() => {
    const code = params.get("code");
    const state = params.get("state");
    if (code === null || state === null) {
      setError("Missing OIDC callback parameters");
      return;
    }
    if (started.current) {
      return;
    }
    started.current = true;
    let cancelled = false;
    finishLogin(code, state)
      .then((session) => {
        if (cancelled) return;
        onAuthenticated(session);
        navigate(defaultRouteFor(session), { replace: true });
      })
      .catch((reason: unknown) => {
        if (!cancelled) {
          setError(reason instanceof Error ? reason.message : "OIDC callback failed");
        }
      });
    return () => {
      cancelled = true;
    };
  }, [navigate, onAuthenticated, params]);

  return (
    <main className="callback-shell">
      <header className="public-nav compact-nav">
        <Link className="secondary-link" to="/login">
          Back to login
        </Link>
        <ThemeToggle theme={theme} setTheme={setTheme} />
      </header>
      <section className="panel callback-panel">
        <p className="eyebrow">OIDC callback</p>
        <h2>Completing sign-in</h2>
        <p>{error ?? "Exchanging authorization code for signed access tokens."}</p>
      </section>
    </main>
  );
}

function WorkspaceLayout({
  context,
  session,
  signOut,
  theme,
  setTheme
}: {
  readonly context: FusionWorkspace["context"];
  readonly session: UserSession | undefined;
  readonly signOut: () => Promise<void>;
  readonly theme: Theme;
  readonly setTheme: (theme: Theme) => void;
}) {
  return (
    <div className="app-shell">
      <aside className="sidebar" aria-label="Workspace navigation">
        <BrandBlock />
        <SessionBlock session={session} signOut={signOut} />
        <nav className="nav-list">
          {workspaceNav.map((item) => (
            <NavLink className={({ isActive }) => navClassName(isActive)} key={item.to} to={item.to}>
              {item.label}
            </NavLink>
          ))}
        </nav>
      </aside>
      <main className="main">
        <header className="topbar">
          <div>
            <p className="eyebrow">Tenant control plane</p>
            <h2 className="topbar-title">{context.project.name}</h2>
          </div>
          <div className="topbar-actions">
            <div className="context-strip" aria-label="Active workspace context">
              <ContextPill label="Organization" value={context.organization.name} />
              <ContextPill label="Tenant" value={context.tenant.name} />
              <ContextPill label="Environment" value={context.project.environment.toUpperCase()} />
              <ContextPill label="Deployment" value={context.organization.deploymentProfile.replaceAll("_", " ")} />
              <ContextPill label="Region" value={context.organization.region} />
            </div>
            <ThemeToggle theme={theme} setTheme={setTheme} />
          </div>
        </header>
        <section className="page">
          <Outlet />
        </section>
      </main>
    </div>
  );
}

function SuperadminLayout({
  overview,
  error,
  loading,
  session,
  signOut,
  theme,
  setTheme
}: {
  readonly overview: PlatformOverview | undefined;
  readonly error: string | undefined;
  readonly loading: boolean;
  readonly session: UserSession | undefined;
  readonly signOut: () => Promise<void>;
  readonly theme: Theme;
  readonly setTheme: (theme: Theme) => void;
}) {
  return (
    <div className="app-shell">
      <aside className="sidebar" aria-label="Superadmin navigation">
        <BrandBlock />
        <SessionBlock session={session} signOut={signOut} />
        <nav className="nav-list">
          {superadminNav.map((item) => (
            <NavLink className={({ isActive }) => navClassName(isActive)} key={item.to} to={item.to}>
              {item.label}
            </NavLink>
          ))}
        </nav>
      </aside>
      <main className="main">
        <header className="topbar">
          <div>
            <p className="eyebrow">Platform superadmin</p>
            <h2 className="topbar-title">Control Plane</h2>
          </div>
          <div className="topbar-actions">
            <div className="context-strip">
              <ContextPill label="Organizations" value={(overview?.organizations.length ?? 0).toString()} />
              <ContextPill label="Identity" value={overview?.identityProviders[0]?.protocol ?? "loading"} />
              <ContextPill label="HA" value={overview?.ha.status ?? (loading ? "loading" : "unknown")} />
              <ContextPill label="Backup" value={overview?.backup.status ?? "unknown"} />
            </div>
            <ThemeToggle theme={theme} setTheme={setTheme} />
          </div>
        </header>
        <section className="page">
          {error ? <div className="message-banner error">{error}</div> : null}
          <Outlet />
        </section>
      </main>
    </div>
  );
}

function BrandBlock() {
  return (
    <div className="brand">
      <div className="brand-lockup">
        <img alt="" className="brand-mark" height={32} src="/logo.svg" width={32} />
        <div>
          <h1 className="brand-title">DCraft Fusion</h1>
          <p className="brand-subtitle">Control plane</p>
        </div>
      </div>
    </div>
  );
}

function SessionBlock({
  session,
  signOut
}: {
  readonly session: UserSession | undefined;
  readonly signOut: () => Promise<void>;
}) {
  return (
    <div className="signed-in">
      <div>
        <strong>{session?.label ?? "Signed out"}</strong>
        <span>{session?.email ?? ""}</span>
      </div>
      <button className="secondary-button" onClick={() => void signOut()} type="button">
        Sign out
      </button>
    </div>
  );
}

function ThemeToggle({ theme, setTheme }: { readonly theme: Theme; readonly setTheme: (theme: Theme) => void }) {
  return (
    <div className="theme-toggle" role="group" aria-label="Theme">
      <button className={theme === "dark" ? "theme-chip active" : "theme-chip"} onClick={() => setTheme("dark")} type="button">
        Dark
      </button>
      <button className={theme === "light" ? "theme-chip active" : "theme-chip"} onClick={() => setTheme("light")} type="button">
        Light
      </button>
    </div>
  );
}

function ContextPill({ label, value }: { readonly label: string; readonly value: string }) {
  return (
    <div className="context-pill">
      <span className="context-label">{label}</span>
      <span className="context-value">{value}</span>
    </div>
  );
}

function CommandCenter({
  workspace,
  connections
}: {
  readonly workspace: FusionWorkspace;
  readonly connections: readonly Connection[];
}) {
  const failedRuns = workspace.runs.filter((run) => run.status === "failed").length;
  const draftAI = workspace.aiRecommendations.filter((item) => item.state === "draft").length;

  return (
    <>
      <PageHeader
        title="Command Center"
        description="Operational state for the active tenant, environment, and production metadata workflows."
      />
      <div className="metric-grid">
        <MetricCard label="Healthy Connections" value={connections.filter((item) => item.status === "healthy").length.toString()} />
        <MetricCard label="Datasets" value={workspace.datasets.length.toString()} />
        <MetricCard label="Failed Runs" value={failedRuns.toString()} />
        <MetricCard label="AI Drafts" value={draftAI.toString()} />
      </div>
      <section className="panel">
        <div className="panel-heading">
          <h3>Recent Runs</h3>
          <p>Failure and freshness coverage stays visible from the first tenant screen.</p>
        </div>
        <DataTable
          headers={["Workflow", "Status", "Environment", "Artifact"]}
          rows={workspace.runs.map((run) => [run.workflowName, run.status, run.environment, run.artifactRef])}
        />
      </section>
    </>
  );
}

function IntegrationHub({
  workspace,
  connections,
  setConnections
}: {
  readonly workspace: FusionWorkspace;
  readonly connections: readonly Connection[];
  readonly setConnections: (connections: readonly Connection[]) => void;
}) {
  const [name, setName] = useState("Brand A Read Replica");
  const [secretRef, setSecretRef] = useState("vault/brand-a/postgres-readonly");
  const [kind, setKind] = useState<Connection["kind"]>("postgres");
  const [selectedConnectionId, setSelectedConnectionId] = useState(connections[0]?.id ?? "");
  const [message, setMessage] = useState<string | undefined>(undefined);
  const selectedConnection = connections.find((connection) => connection.id === selectedConnectionId) ?? connections[0];

  useEffect(() => {
    if (selectedConnectionId === "" && connections[0] !== undefined) {
      setSelectedConnectionId(connections[0].id);
    }
  }, [connections, selectedConnectionId]);

  async function createConnection() {
    const optimisticConnection: Connection = {
      id: `conn-ui-${connections.length + 1}`,
      name,
      kind,
      organizationId: workspace.context.organization.id,
      tenantId: workspace.context.tenant.id,
      projectId: workspace.context.project.id,
      secretRef: `secretref:${secretRef}`,
      readOnly: true,
      status: "untested"
    };
    setMessage(undefined);
    setConnections([...connections, optimisticConnection]);
    setSelectedConnectionId(optimisticConnection.id);
    try {
      const created = await createConnectionAPI({ name, kind, secretRef, readOnly: true });
      setConnections([...connections, created]);
      setSelectedConnectionId(created.id);
      setMessage("Read-only connection created through the API.");
    } catch (error: unknown) {
      setMessage(error instanceof Error ? error.message : "Connection creation failed");
    }
  }

  async function testConnection() {
    if (selectedConnection === undefined) return;
    setMessage(undefined);
    setConnections(
      connections.map((connection) =>
        connection.id === selectedConnection.id ? { ...connection, status: "healthy" as const } : connection
      )
    );
    try {
      await testConnectionAPI(selectedConnection.id);
      setMessage("Connection health validated through the API.");
    } catch (error: unknown) {
      setMessage(error instanceof Error ? error.message : "Connection test failed");
    }
  }

  return (
    <>
      <PageHeader title="Integration Hub" description="Read-only adapter entry point for Phase 1 connections and health checks." />
      <div className="integration-grid">
        {workspace.integrations.map((integration) => (
          <IntegrationCardView integration={integration} key={integration.id} />
        ))}
      </div>
      <div className="two-column">
        <section className="panel">
          <div className="panel-heading">
            <h3>Create Connection</h3>
            <p>Requests are now tied to a real login session and OIDC token exchange.</p>
          </div>
          <label>
            Name
            <input value={name} onChange={(event) => setName(event.target.value)} />
          </label>
          <label>
            Connector
            <select value={kind} onChange={(event) => setKind(event.target.value as Connection["kind"])}>
              <option value="postgres">Postgres</option>
              <option value="airflow">Airflow</option>
              <option value="github">GitHub</option>
              <option value="slack">Slack</option>
              <option value="fusion_cdc_engine">Fusion CDC Engine</option>
            </select>
          </label>
          <label>
            Secret Reference
            <input value={secretRef} onChange={(event) => setSecretRef(event.target.value)} />
          </label>
          <button className="primary-button" onClick={() => void createConnection()} type="button">
            Create Read-only Connection
          </button>
        </section>
        <section className="panel">
          <div className="panel-heading">
            <h3>Connection Health</h3>
            <p>Health checks run through the control plane and persist job + audit state.</p>
          </div>
          <select value={selectedConnection?.id ?? ""} onChange={(event) => setSelectedConnectionId(event.target.value)}>
            {connections.map((connection) => (
              <option key={connection.id} value={connection.id}>
                {connection.name}
              </option>
            ))}
          </select>
          {selectedConnection ? (
            <div className="detail-list">
              <p>
                <strong>Status</strong>
                <span className="status" data-tone={selectedConnection.status === "healthy" ? "ready" : "preview"}>
                  {selectedConnection.status}
                </span>
              </p>
              <p>
                <strong>Secret</strong>
                <span>{selectedConnection.secretRef}</span>
              </p>
            </div>
          ) : null}
          <button className="primary-button" onClick={() => void testConnection()} type="button">
            Test Connection
          </button>
          {message ? <div className="message-banner">{message}</div> : null}
        </section>
      </div>
    </>
  );
}

function IntegrationCardView({ integration }: { readonly integration: IntegrationCard }) {
  return (
    <article className="card">
      <div className="card-head">
        <h3>{integration.name}</h3>
        <span className="status" data-tone={integration.status}>
          {integration.status}
        </span>
      </div>
      <p>{integration.description}</p>
    </article>
  );
}

function MetadataExplorer({ workspace }: { readonly workspace: FusionWorkspace }) {
  const [query, setQuery] = useState("");
  const datasets = workspace.datasets.filter((dataset) => `${dataset.schema}.${dataset.name}`.includes(query.toLowerCase()));
  return (
    <>
      <PageHeader title="Metadata Explorer" description="Browse discovered datasets without exposing raw records." />
      <section className="panel">
        <input aria-label="Search datasets" placeholder="Search datasets" value={query} onChange={(event) => setQuery(event.target.value)} />
        <DataTable
          headers={["Dataset", "Source", "Owner", "Tags"]}
          rows={datasets.map((dataset) => [`${dataset.schema}.${dataset.name}`, dataset.sourceSystem, dataset.owner, dataset.tags.join(", ")])}
        />
      </section>
    </>
  );
}

function RunCenter({ workspace }: { readonly workspace: FusionWorkspace }) {
  return (
    <>
      <PageHeader title="Workflow Runs" description="Run state across orchestrated metadata and policy checks." />
      <section className="panel">
        <DataTable
          headers={["Workflow", "Status", "Correlation", "Artifact"]}
          rows={workspace.runs.map((run) => [run.workflowName, run.status, run.correlationId, run.artifactRef])}
        />
      </section>
    </>
  );
}

function AuditCenter({ workspace }: { readonly workspace: FusionWorkspace }) {
  return (
    <>
      <PageHeader title="Audit Center" description="Actor, resource, and correlation coverage for tenant operations." />
      <section className="panel">
        <DataTable
          headers={["Action", "Actor", "Resource", "Correlation"]}
          rows={workspace.auditEvents.map((event) => [event.action, event.actorId, `${event.resourceType}:${event.resourceId}`, event.correlationId])}
        />
      </section>
    </>
  );
}

function Settings({ workspace }: { readonly workspace: FusionWorkspace }) {
  return (
    <>
      <PageHeader title="Settings" description="Tenant deployment posture, isolation signals, and active guardrails." />
      <div className="metric-grid">
        <MetricCard label="Deployment Profile" value={workspace.context.organization.deploymentProfile.replaceAll("_", " ")} />
        <MetricCard label="Data Region" value={workspace.context.organization.region} />
        <MetricCard label="Tenants" value={workspace.tenants.length.toString()} />
        <MetricCard label="Projects" value={workspace.projects.length.toString()} />
      </div>
      <section className="panel">
        <DataTable headers={["Policy", "Type", "Mode"]} rows={workspace.policies.map((policy) => [policy.name, policy.type, policy.mode])} />
      </section>
    </>
  );
}

function AIOperator({ workspace }: { readonly workspace: FusionWorkspace }) {
  return (
    <>
      <PageHeader title="AI Operator" description="Draft-only recommendation surfaces with explicit human approval gates." />
      <div className="integration-grid">
        {workspace.aiRecommendations.map((recommendation) => (
          <article className="card" key={recommendation.id}>
            <div className="card-head">
              <h3>{recommendation.title}</h3>
              <span className="status" data-tone="preview">
                {recommendation.state}
              </span>
            </div>
            <p>{recommendation.summary}</p>
          </article>
        ))}
      </div>
    </>
  );
}

function SuperadminOverviewScreen({
  overview,
  authConfig,
  error
}: {
  readonly overview: PlatformOverview | undefined;
  readonly authConfig: OpenIDConfiguration | undefined;
  readonly error: string | undefined;
}) {
  if (overview === undefined) {
    return <AccessPanel error={error} title="Superadmin Overview" />;
  }
  return (
    <>
      <PageHeader title="Superadmin Overview" description="Platform-wide status across identity, org topology, HA, backup, and adapters." />
      <div className="metric-grid">
        <MetricCard label="Organizations" value={overview.organizations.length.toString()} />
        <MetricCard label="Superadmins" value={overview.superAdmins.length.toString()} />
        <MetricCard label="API Replicas" value={overview.ha.apiMinReplicas.toString()} />
        <MetricCard label="Backup Status" value={overview.backup.status} />
      </div>
      <div className="two-column">
        <section className="panel">
          <div className="panel-heading">
            <h3>Runtime Identity</h3>
            <p>Live OIDC configuration shared by the browser shell and the control plane.</p>
          </div>
          <div className="detail-list">
            <p>
              <strong>Issuer</strong>
              <span>{authConfig?.issuer ?? "Loading"}</span>
            </p>
            <p>
              <strong>Authorize</strong>
              <span>/oidc/authorize</span>
            </p>
            <p>
              <strong>Token</strong>
              <span>/oidc/token</span>
            </p>
            <p>
              <strong>JWKS</strong>
              <span>{authConfig?.jwks_uri ?? "Loading"}</span>
            </p>
          </div>
        </section>
        <section className="panel">
          <div className="panel-heading">
            <h3>Organization Coverage</h3>
            <p>Rollout spans B2B, B2C, B2B2C, shared infra, dedicated infra, and customer-managed profiles.</p>
          </div>
          <DataTable
            headers={["Organization", "Type", "Deployment", "Tenants"]}
            rows={overview.organizations.map((organization) => [
              organization.name,
              organization.type,
              organization.deploymentProfile.replaceAll("_", " "),
              organization.tenants.length.toString()
            ])}
          />
        </section>
      </div>
    </>
  );
}

function OrganizationsScreen({
  overview,
  onRefresh
}: {
  readonly overview: PlatformOverview | undefined;
  readonly onRefresh: () => Promise<void>;
}) {
  const [organizationName, setOrganizationName] = useState("New Enterprise Organization");
  const [organizationType, setOrganizationType] = useState("business");
  const [deploymentProfile, setDeploymentProfile] = useState("dedicated_stack");
  const [organizationRegion, setOrganizationRegion] = useState("us-east-1");
  const [dataPlaneLocation, setDataPlaneLocation] = useState("customer_managed");
  const [selectedOrganizationID, setSelectedOrganizationID] = useState("");
  const [tenantName, setTenantName] = useState("Primary Tenant");
  const [tenantModel, setTenantModel] = useState("b2b");
  const [tenantIsolation, setTenantIsolation] = useState("dedicated_infra");
  const [tenantRegion, setTenantRegion] = useState("us-east-1");
  const [selectedTenantID, setSelectedTenantID] = useState("");
  const [projectName, setProjectName] = useState("Production Workspace");
  const [projectEnvironment, setProjectEnvironment] = useState("production");
  const [message, setMessage] = useState<string | undefined>(undefined);

  useEffect(() => {
    if (overview?.organizations[0] !== undefined && selectedOrganizationID === "") {
      setSelectedOrganizationID(overview.organizations[0].id);
    }
  }, [overview, selectedOrganizationID]);

  useEffect(() => {
    const organization = overview?.organizations.find((item) => item.id === selectedOrganizationID);
    if (organization?.tenants[0] !== undefined && selectedTenantID === "") {
      setSelectedTenantID(organization.tenants[0].id);
    }
  }, [overview, selectedOrganizationID, selectedTenantID]);

  if (overview === undefined) {
    return <AccessPanel error="Platform overview is unavailable." title="Organizations" />;
  }

  const selectedOrganization = overview.organizations.find((item) => item.id === selectedOrganizationID) ?? overview.organizations[0];
  const selectedTenant = selectedOrganization?.tenants.find((item) => item.id === selectedTenantID) ?? selectedOrganization?.tenants[0];

  async function createOrganizationAction() {
    await createPlatformOrganization({
      name: organizationName,
      type: organizationType,
      deploymentProfile,
      region: organizationRegion,
      dataPlaneLocation,
      rawDataMovementAllowed: false
    });
    setMessage("Organization created.");
    await onRefresh();
  }

  async function createTenantAction() {
    if (selectedOrganization === undefined) return;
    await createPlatformTenant({
      organizationId: selectedOrganization.id,
      name: tenantName,
      model: tenantModel,
      isolation: tenantIsolation,
      region: tenantRegion
    });
    setMessage("Tenant created.");
    await onRefresh();
  }

  async function createProjectAction() {
    if (selectedOrganization === undefined || selectedTenant === undefined) return;
    await createPlatformProject({
      organizationId: selectedOrganization.id,
      tenantId: selectedTenant.id,
      name: projectName,
      environment: projectEnvironment
    });
    setMessage("Project created.");
    await onRefresh();
  }

  return (
    <>
      <PageHeader title="Organizations" description="Create and inspect organizations, tenants, and projects without touching the database." />
      {message ? <div className="message-banner">{message}</div> : null}
      <div className="three-column">
        <section className="panel">
          <div className="panel-heading">
            <h3>Create Organization</h3>
            <p>Foundation step for B2B, B2C, B2B2C, internal, and dedicated customer environments.</p>
          </div>
          <label>
            Name
            <input value={organizationName} onChange={(event) => setOrganizationName(event.target.value)} />
          </label>
          <label>
            Type
            <select value={organizationType} onChange={(event) => setOrganizationType(event.target.value)}>
              <option value="business">Business</option>
              <option value="consumer_app">Consumer App</option>
              <option value="internal_platform">Internal Platform</option>
              <option value="managed_service_provider">Managed Service Provider</option>
            </select>
          </label>
          <label>
            Deployment
            <select value={deploymentProfile} onChange={(event) => setDeploymentProfile(event.target.value)}>
              <option value="pooled_saas">Pooled SaaS</option>
              <option value="dedicated_database">Dedicated Database</option>
              <option value="dedicated_stack">Dedicated Stack</option>
              <option value="customer_vpc">Customer VPC</option>
              <option value="on_prem">On Prem</option>
            </select>
          </label>
          <label>
            Region
            <input value={organizationRegion} onChange={(event) => setOrganizationRegion(event.target.value)} />
          </label>
          <label>
            Data Plane
            <select value={dataPlaneLocation} onChange={(event) => setDataPlaneLocation(event.target.value)}>
              <option value="fusion_managed">Fusion Managed</option>
              <option value="customer_managed">Customer Managed</option>
              <option value="hybrid">Hybrid</option>
            </select>
          </label>
          <button className="primary-button" onClick={() => void createOrganizationAction()} type="button">
            Create Organization
          </button>
        </section>
        <section className="panel">
          <div className="panel-heading">
            <h3>Create Tenant</h3>
            <p>Assign tenant model and isolation posture from the platform layer.</p>
          </div>
          <label>
            Organization
            <select value={selectedOrganization?.id ?? ""} onChange={(event) => setSelectedOrganizationID(event.target.value)}>
              {overview.organizations.map((organization) => (
                <option key={organization.id} value={organization.id}>
                  {organization.name}
                </option>
              ))}
            </select>
          </label>
          <label>
            Name
            <input value={tenantName} onChange={(event) => setTenantName(event.target.value)} />
          </label>
          <label>
            Model
            <select value={tenantModel} onChange={(event) => setTenantModel(event.target.value)}>
              <option value="b2b">B2B</option>
              <option value="b2c">B2C</option>
              <option value="b2b2c">B2B2C</option>
              <option value="internal">Internal</option>
            </select>
          </label>
          <label>
            Isolation
            <select value={tenantIsolation} onChange={(event) => setTenantIsolation(event.target.value)}>
              <option value="shared_infra">Shared Infra</option>
              <option value="dedicated_infra">Dedicated Infra</option>
              <option value="customer_infra">Customer Infra</option>
            </select>
          </label>
          <label>
            Region
            <input value={tenantRegion} onChange={(event) => setTenantRegion(event.target.value)} />
          </label>
          <button className="primary-button" onClick={() => void createTenantAction()} type="button">
            Create Tenant
          </button>
        </section>
        <section className="panel">
          <div className="panel-heading">
            <h3>Create Project</h3>
            <p>Provision named workspaces for production, staging, or sandbox operations.</p>
          </div>
          <label>
            Tenant
            <select value={selectedTenant?.id ?? ""} onChange={(event) => setSelectedTenantID(event.target.value)}>
              {(selectedOrganization?.tenants ?? []).map((tenant) => (
                <option key={tenant.id} value={tenant.id}>
                  {tenant.name}
                </option>
              ))}
            </select>
          </label>
          <label>
            Name
            <input value={projectName} onChange={(event) => setProjectName(event.target.value)} />
          </label>
          <label>
            Environment
            <select value={projectEnvironment} onChange={(event) => setProjectEnvironment(event.target.value)}>
              <option value="production">Production</option>
              <option value="staging">Staging</option>
              <option value="development">Development</option>
            </select>
          </label>
          <button className="primary-button" onClick={() => void createProjectAction()} type="button">
            Create Project
          </button>
        </section>
      </div>
      <section className="panel">
        <div className="panel-heading">
          <h3>Organization Tree</h3>
          <p>One screen for org, tenant, and project topology.</p>
        </div>
        <div className="org-list">
          {overview.organizations.map((organization) => (
            <OrganizationCard key={organization.id} organization={organization} />
          ))}
        </div>
      </section>
    </>
  );
}

function IdentityScreen({
  overview,
  users,
  authConfig
}: {
  readonly overview: PlatformOverview | undefined;
  readonly users: readonly AuthUser[];
  readonly authConfig: OpenIDConfiguration | undefined;
}) {
  if (overview === undefined) {
    return <AccessPanel error="Platform overview is unavailable." title="Identity" />;
  }
  return (
    <>
      <PageHeader title="Identity" description="OIDC runtime, local login credentials, and future SSO provider posture." />
      <div className="two-column">
        <section className="panel">
          <div className="panel-heading">
            <h3>Local OIDC Provider</h3>
            <p>Login now uses typed credentials and a server-side auth session before the authorization code flow.</p>
          </div>
          <DataTable
            headers={["Field", "Value"]}
            rows={[
              ["Issuer", authConfig?.issuer ?? "Loading"],
              ["Authorize", "/oidc/authorize"],
              ["Token", "/oidc/token"],
              ["JWKS", authConfig?.jwks_uri ?? "Loading"]
            ]}
          />
        </section>
        <section className="panel">
          <div className="panel-heading">
            <h3>Seeded Login Credentials</h3>
            <p>Local runtime identities that can sign in through the real form.</p>
          </div>
          <DataTable
            headers={["User", "Email", "Scope"]}
            rows={users.map((user) => [user.label, user.email, user.scope])}
          />
        </section>
      </div>
      <section className="panel">
        <div className="panel-heading">
          <h3>Provider Inventory</h3>
          <p>Configured OIDC and planned SAML surfaces tracked in the platform dataset.</p>
        </div>
        <DataTable headers={["Provider", "Protocol", "Issuer", "Status"]} rows={overview.identityProviders.map((provider) => [provider.name, provider.protocol, provider.issuer, provider.status])} />
      </section>
    </>
  );
}

function SecretsScreen({ overview }: { readonly overview: PlatformOverview | undefined }) {
  if (overview === undefined) {
    return <AccessPanel error="Platform overview is unavailable." title="Secrets" />;
  }
  return (
    <>
      <PageHeader title="Secrets" description="Platform-wide secret manager coverage and deployment scope." />
      <div className="metric-grid">
        <MetricCard label="Secret Providers" value={overview.secretProviders.length.toString()} />
        <MetricCard label="Platform Scope" value={overview.secretProviders.filter((item) => item.scope === "platform").length.toString()} />
        <MetricCard label="Customer Scope" value={overview.secretProviders.filter((item) => item.scope !== "platform").length.toString()} />
      </div>
      <section className="panel">
        <DataTable
          headers={["Provider", "Kind", "Scope", "Endpoint", "Status"]}
          rows={overview.secretProviders.map((provider) => [provider.name, provider.kind, provider.scope, provider.endpoint, provider.status])}
        />
      </section>
    </>
  );
}

function AdaptersScreen({ overview }: { readonly overview: PlatformOverview | undefined }) {
  if (overview === undefined) {
    return <AccessPanel error="Platform overview is unavailable." title="Adapters" />;
  }
  return (
    <>
      <PageHeader title="Adapters" description="Read-only ingestion adapters and the async jobs that validate them." />
      <div className="two-column">
        <section className="panel">
          <DataTable headers={["Adapter", "Kind", "Mode", "Status"]} rows={overview.adapters.map((adapter) => [adapter.name, adapter.kind, adapter.readOnly ? "read-only" : "write", adapter.status])} />
        </section>
        <section className="panel">
          <DataTable headers={["Job", "Status", "Resource", "Message"]} rows={overview.jobs.map((job) => [job.type, job.status, job.resourceId, job.message])} />
        </section>
      </div>
    </>
  );
}

function OperationsScreen({ overview }: { readonly overview: PlatformOverview | undefined }) {
  if (overview === undefined) {
    return <AccessPanel error="Platform overview is unavailable." title="Operations" />;
  }
  return (
    <>
      <PageHeader title="Operations" description="Backup, observability, HA posture, and recent platform audit activity." />
      <div className="three-column">
        <section className="panel">
          <div className="panel-heading">
            <h3>Backup</h3>
          </div>
          <div className="detail-list">
            <p>
              <strong>Target</strong>
              <span>{overview.backup.target}</span>
            </p>
            <p>
              <strong>Schedule</strong>
              <span>{overview.backup.schedule}</span>
            </p>
            <p>
              <strong>Status</strong>
              <span>{overview.backup.status}</span>
            </p>
            <p>
              <strong>Runbook</strong>
              <span>{overview.backup.runbook}</span>
            </p>
          </div>
        </section>
        <section className="panel">
          <div className="panel-heading">
            <h3>Observability</h3>
          </div>
          <div className="detail-list">
            <p>
              <strong>Metrics</strong>
              <span>{overview.observability.metricsEndpoint}</span>
            </p>
            <p>
              <strong>Logs</strong>
              <span>{overview.observability.logFormat}</span>
            </p>
            <p>
              <strong>OTel Ready</strong>
              <span>{overview.observability.openTelemetryReady ? "yes" : "no"}</span>
            </p>
          </div>
        </section>
        <section className="panel">
          <div className="panel-heading">
            <h3>High Availability</h3>
          </div>
          <div className="detail-list">
            <p>
              <strong>API Replicas</strong>
              <span>{overview.ha.apiMinReplicas}</span>
            </p>
            <p>
              <strong>Web Replicas</strong>
              <span>{overview.ha.webMinReplicas}</span>
            </p>
            <p>
              <strong>Autoscaling</strong>
              <span>{overview.ha.autoscaling ? "enabled" : "disabled"}</span>
            </p>
            <p>
              <strong>Status</strong>
              <span>{overview.ha.status}</span>
            </p>
          </div>
        </section>
      </div>
      <section className="panel">
        <div className="panel-heading">
          <h3>Recent Audit Events</h3>
          <p>Newest platform-visible activity without opening the database.</p>
        </div>
        <DataTable
          headers={["Action", "Actor", "Resource", "Occurred"]}
          rows={overview.recentAuditEvents.map((event) => [event.action, event.actorId, `${event.resourceType}:${event.resourceId}`, new Date(event.occurredAt).toLocaleString()])}
        />
      </section>
    </>
  );
}

function AccessPanel({ title, error }: { readonly title: string; readonly error: string | undefined }) {
  return (
    <section className="panel">
      <h3>{title}</h3>
      <p>{error ?? "Access is unavailable."}</p>
    </section>
  );
}

function PageHeader({ title, description }: { readonly title: string; readonly description: string }) {
  return (
    <div className="page-header">
      <h2 className="page-title">{title}</h2>
      <p className="page-description">{description}</p>
    </div>
  );
}

function OrganizationCard({ organization }: { readonly organization: PlatformOrganization }) {
  return (
    <article className="org-card">
      <div className="card-head">
        <div>
          <h3>{organization.name}</h3>
          <p className="muted">
            {organization.type.replaceAll("_", " ")} · {organization.deploymentProfile.replaceAll("_", " ")} · {organization.region}
          </p>
        </div>
        <span className="status" data-tone="ready">
          {organization.tenants.length} tenants
        </span>
      </div>
      <div className="org-tenant-list">
        {organization.tenants.map((tenant) => (
          <div className="tenant-row" key={tenant.id}>
            <div>
              <strong>{tenant.name}</strong>
              <p className="muted">
                {tenant.model.toUpperCase()} · {tenant.isolation.replaceAll("_", " ")}
              </p>
            </div>
            <span>{tenant.projects.map((project) => project.name).join(", ") || "No projects"}</span>
          </div>
        ))}
      </div>
    </article>
  );
}

function MetricCard({ label, value }: { readonly label: string; readonly value: string }) {
  return (
    <article className="metric-card">
      <span className="metric-label">{label}</span>
      <p className="metric-value">{value}</p>
    </article>
  );
}

function MetricBadge({ label, value }: { readonly label: string; readonly value: string }) {
  return (
    <article className="metric-badge">
      <span className="metric-label">{label}</span>
      <strong>{value}</strong>
    </article>
  );
}

function InfoTile({ label, value }: { readonly label: string; readonly value: string }) {
  return (
    <article className="info-tile">
      <span className="metric-label">{label}</span>
      <p className="metric-value metric-value-small">{value}</p>
    </article>
  );
}

function DataTable({ headers, rows }: { readonly headers: readonly string[]; readonly rows: readonly (readonly string[])[] }) {
  return (
    <div className="table-wrap">
      <table>
        <thead>
          <tr>
            {headers.map((header) => (
              <th key={header}>{header}</th>
            ))}
          </tr>
        </thead>
        <tbody>
          {rows.map((row) => (
            <tr key={row.join(":")}>
              {row.map((cell) => (
                <td key={cell}>{cell}</td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

function navClassName(isActive: boolean) {
  return isActive ? "nav-item active" : "nav-item";
}

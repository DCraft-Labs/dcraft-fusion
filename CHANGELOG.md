# Changelog

All notable changes to DCraft Fusion (public repo) are documented here.
This project follows [Keep a Changelog](https://keepachangelog.com/) and
uses [Semantic Versioning](https://semver.org/).

## [1.2.10] — 2026-07-22

Coordinated release with the private `fusion-cdc-engine` v1.2.10. v1.2.9
verified the UI was live and solid, but the E2E CDC audit found the runtime
non-functional. v1.2.10 ships the local-dev infrastructure pods and the
re-tagged Helm charts that point at the v1.2.10 CDC images.

### Added
- **Local-dev infra pods (BLOCKER 7)** — `infra/local-dev/k8s/00-infra.yaml`
  now deploys `mysql-source` (MySQL 8.0 with binlog + self-seeding init SQL
  for customers/products/orders), `mongo-source` (MongoDB 7, no auth),
  `minio` (S3-compatible API :9000 + console :9001 + PVC + bucket-init Job
  that creates `iceberg-warehouse`), and `nessie` (Iceberg REST catalog
  :19120 + mgmt :19121). All with startup/readiness probes and small
  (128Mi–256Mi) resource requests. LOCAL-DEV ONLY — production must use
  managed equivalents (RDS, DocumentDB, S3, Nessie/Polaris, etc.).

### Changed
- Bumped `dcraft-fusion` and `fusion-cdc` Helm charts to `version: 1.2.10`
  / `appVersion: "1.2.10"` (`infra/helm/*/Chart.yaml`).
- Bumped all image tags from `1.2.9` → `1.2.10` in
  `infra/helm/dcraft-fusion/values.yaml`,
  `infra/helm/fusion-cdc/values.yaml`, the `examples/values-minimal.yaml`
  overlays, and `infra/local-dev/k8s/values-{cdc,fusion}-local.yaml`.
- Bumped `--version 1.2.10` in `infra/local-dev/k8s/deploy.ps1` and the
  `infra/helm/README.md` install snippets.

### Coordinated
- The CDC runtime fixes (worker assignment, MongoDB connections, Iceberg
  seed, pg discovery, UI form fields, port defaults) ship in the private
  `fusion-cdc-engine` repo's v1.2.10 release. This public repo re-tags the
  Fusion web + control-plane-kernel images to `1.2.10` and adds the
  local-dev infra pods so the seeded Mongo/MySQL/MinIO/Nessie hostnames
  resolve.

## [1.2.9] — 2026-07-22

UI polish pass on top of the verified-stable v1.2.8 backend. The v1.2.8 UX
audit (`E:\Dcraft\.tmp\v128-ux-audit\REPORT_v128_ux_audit.md`) found that
the Fusion workspace nav active state was invisible (no CSS for
`.nav-item.active`), the Audit Center table had no timestamp column, and
Recent Runs status had no semantic color. This release fixes those plus
the emilkowalski button `:active` scale finding. Coordinated with the
private `fusion-cdc-engine` v1.2.9 release which fixes the superadmin
role mapping (BLOCKER), the connectors "Used by" count, the localhost
GraphQL link, and adds keyboard-accessible dropdowns + button active scale
on the CDC frontend.

### Changed (Fusion public repo)
- Added `.nav-item` / `.nav-item.active` / `.nav-item:hover` styles in
  `apps/web/src/styles.css` — bolder font weight, brand-tinted background,
  and a 2px left accent border so the active workspace page is clearly
  distinguishable.
- Added `:active { transform: scale(0.97) }` to `.btn`, `.primary-button`,
  `.secondary-button`, `.primary-link`, `.secondary-link` in
  `apps/web/src/styles.css` (emilkowalski polish — buttons no longer feel
  dead on click).
- Added a "Timestamp" column to the Audit Center table
  (`apps/web/src/ui/App.tsx` `AuditCenter`) showing relative time
  ("2 min ago") with the absolute timestamp on hover via the `title`
  attribute. `AuditEvent.occurredAt` added to
  `apps/web/src/ui/workspace.ts` (optional) and populated in
  `apps/web/src/ui/demo-workspace.ts`. `DataTable` cell type relaxed from
  `string` to `ReactNode` so cells can carry hover titles.
- Added semantic color badges to the Recent Runs status column
  (`apps/web/src/ui/App.tsx` `RunCenter`): `succeeded` → green
  (`--success`/`--success-text`), `failed`/`canceled` → red
  (`--error-bg`/`--error-text`), `running` → brand blue/teal
  (`--brand-soft`/`--brand-strong`), `queued`/`unknown` → muted gray.
  New `.status[data-tone="success|failed|running|pending"]` rules in
  `apps/web/src/styles.css` use existing design tokens (no hardcoded hex).
- Bumped `dcraft-fusion` and `fusion-cdc` Helm charts to `version: 1.2.9`
  / `appVersion: "1.2.9"` (`infra/helm/*/Chart.yaml`).
- Bumped all image tags from `1.2.8` → `1.2.9` in
  `infra/helm/*/values.yaml`, `infra/helm/*/examples/values-minimal.yaml`,
  and `infra/local-dev/k8s/values-{cdc,fusion}-local.yaml`.
- Bumped `--version 1.2.9` in `infra/local-dev/k8s/deploy.ps1` and the
  helm install examples in `infra/helm/README.md`.

### Notes
- The CDC-side fixes (superadmin role mapping, connector usage_count,
  localhost GraphQL link, button active scale, keyboard-accessible
  dropdowns) ship in the private `fusion-cdc-engine` repo's v1.2.9
  release.
- No backend changes in this public repo — the v1.2.8 control-plane
  kernel is verified live and stable; v1.2.9 is UI + chart re-tag only.

## [1.2.8] — 2026-07-22

Coordinated release with the private `fusion-cdc-engine` v1.2.8. The
v1.2.7 live stability audit found that `/api/v1/alerts/suppressions`
still returned HTTP 500 (the v1.2.5 `f7a8b9c0d1e2` migration only added
`rule_ids`/`connection_ids`; the `is_recurring`, `recurrence_pattern`,
and `updated_by` columns the `AlertSuppression` model declares were
never created) and `/api/v1/data-quality/templates` returned HTTP 501
from a stub handler. The CDC-side fixes ship in the private repo's
v1.2.8 release (new migration `a8b9c0d1e2f3` + templates stub replaced
with an empty 200 response). This public release re-tags all images
and charts to `1.2.8` so the public `fusion-cdc` chart deploys the
fixed CDC images.

### Changed
- Bumped `dcraft-fusion` and `fusion-cdc` Helm charts to `version: 1.2.8`
  / `appVersion: "1.2.8"` (`infra/helm/*/Chart.yaml`).
- Bumped all image tags from `1.2.7` → `1.2.8` in
  `infra/helm/*/values.yaml`, `infra/helm/*/examples/values-minimal.yaml`,
  and `infra/local-dev/k8s/values-{cdc,fusion}-local.yaml`.
- Bumped `--version 1.2.8` in `infra/local-dev/k8s/deploy.ps1` and the
  helm install examples in `infra/helm/README.md`.

### Notes
- The CDC-side code fixes (alert_suppressions migration
  `a8b9c0d1e2f3`, data-quality/templates 501 → 200, FastAPI app version
  bump) ship in the private `fusion-cdc-engine` repo's v1.2.8 release.
- pg-source discovery returning 0 tables is operational and already
  documented in `docs/POST_DEPLOY_CHECKLIST.md` §2 — no code change.

## [1.2.7] — 2026-07-22

Follow-up to v1.2.6. The v1.2.6 `Publish images` and `Publish Helm charts`
workflows succeeded (the `dcraft-fusion-web` image was rebuilt for the
first time since v1.2.3), but the `CI` workflow failed at
`npm audit --audit-level=high` on newly-disclosed high-severity
dev-server-only vulnerabilities in transitive dev dependencies. This
release pins the patched versions so the CI gate is green.

### Fixed (CI — npm audit)
- **vite `server.fs.deny` bypass on Windows** (GHSA-fx2h-pf6j-xcff,
  CVE-2026-53571, high): the VitePress docs workspace pulled a nested
  `vite <=6.4.2` that had no auto-fix (VitePress 1.6.4 declares
  `vite: ^5.4.14`, and the advisory is only patched in `vite 6.4.3` /
  `7.3.5` / `8.0.16` — there is no 5.x fix). Added an npm `overrides`
  block in the root `package.json` forcing `vite` to `6.4.3` for both
  the web workspace and VitePress's nested dependency, and regenerated
  `package-lock.json`. VitePress 1.6.4 is compatible with Vite 6.4.3
  (docs build verified). This is a dev-server-only issue (Windows
  `vite dev --host` with sensitive files in `server.fs.allow`); it does
  not affect `vite build` or the shipped static bundle in the Docker
  image.
- `npm audit fix` also pulled in patched versions for the
  `@babel/core` (GHSA-4x5r-pxfx-6jf8), `form-data`
  (GHSA-hmw2-7cc7-3qxx), and `ws` (GHSA-96hv-2xvq-fx4p) advisories
  that had fixes available.

### Changed
- Bumped `dcraft-fusion` and `fusion-cdc` Helm charts to `version: 1.2.7`
  / `appVersion: "1.2.7"` (`infra/helm/*/Chart.yaml`).
- Bumped all image tags from `1.2.6` → `1.2.7` in
  `infra/helm/*/values.yaml`, `infra/helm/*/examples/values-minimal.yaml`,
  and `infra/local-dev/k8s/values-{cdc,fusion}-local.yaml`.
- Bumped `--version 1.2.7` in `infra/local-dev/k8s/deploy.ps1` and the
  helm install examples in `infra/helm/README.md`.

### Notes
- The v1.2.6 images remain valid in GHCR; v1.2.7 re-publishes identical
  source under a new tag so the CI gate is green for the release.
- The CDC-side version bump (`control-plane/app/main.py` and
  `helm/fusion-cdc/Chart.yaml`) ships in the private `fusion-cdc-engine`
  repo's v1.2.7 release.

## [1.2.6] — 2026-07-22

Follow-up to v1.2.5. The v1.2.4 and v1.2.5 `Publish images` and `CI`
workflows failed, so the `dcraft-fusion-web` image was never rebuilt and
the deployed frontend stayed on the v1.2.3 bundle. This release fixes the
two CI failures and re-publishes the images.

### Fixed (Fusion SPA — CI failures)
- **`auth.ts` UUID polyfill type** (`apps/web/src/ui/auth.ts:5-14`): the
  `crypto.randomUUID` polyfill returned `string` but `Crypto.randomUUID`
  is typed as the template-literal UUID type, so `tsc` failed with
  `TS2322` and the `Publish images` workflow could not build the web
  image. The polyfill now declares the return type as
  `` `${string}-${string}-${string}-${string}-${string}` `` and casts the
  generated string, satisfying the DOM lib type.
- **`SuperadminOverviewScreen` heading race** (`apps/web/src/ui/App.tsx:986-1050`):
  the screen rendered `<AccessPanel title="Superadmin Overview">` (an
  `<h3>`) while the platform overview was loading, then swapped to
  `<PageHeader>` (an `<h2>`) once `/api/v1/platform/overview` resolved.
  Under Node 24 / jsdom 25 the `App.test.tsx` superadmin test's
  `findByRole("heading", { name: "Superadmin Overview" })` resolved
  against the transient `<h3>`, which was then removed by the re-render,
  so `toBeInTheDocument()` failed. The screen now always renders the
  `<PageHeader>` heading and only swaps the body (loading message vs.
  metric grid / panels), so the heading element is stable across the
  loading → loaded transition. The v1.2.5 DEMO-banner and Test-Connection
  chip-flip fixes are unchanged.

### Changed
- Bumped `dcraft-fusion` and `fusion-cdc` Helm charts to `version: 1.2.6`
  / `appVersion: "1.2.6"` (`infra/helm/*/Chart.yaml`).
- Bumped all image tags from `1.2.5` → `1.2.6` in
  `infra/helm/*/values.yaml`, `infra/helm/*/examples/values-minimal.yaml`,
  and `infra/local-dev/k8s/values-{cdc,fusion}-local.yaml`.
- Bumped `--version 1.2.6` in `infra/local-dev/k8s/deploy.ps1` and the
  helm install examples in `infra/helm/README.md`.

### Notes
- The CDC-side version bump (`control-plane/app/main.py` and
  `helm/fusion-cdc/Chart.yaml`) ships in the private `fusion-cdc-engine`
  repo's v1.2.6 release.

## [1.2.5] — 2026-07-22

Follow-up to v1.2.4. The v1.2.4 live verification + codebase map audit
found that the v1.2.4 SPA fixes never took effect because the deployed
bundle was a stale v1.0.0 build (the image was never rebuilt from source),
and the CDC alert_rules/suppressions migrations were incomplete. This
release fixes the SPA build pipeline, hardens the kernel HTTP API
context scoping, and makes the local-dev Postgres hostPath cross-platform.

### Fixed (Fusion SPA — root cause of v1.2.4 not taking effect)
- **SPA bundle staleness** (`apps/web/Dockerfile`): the production bundle
  is now ALWAYS rebuilt from source inside the image. Added a defensive
  `rm -rf apps/web/dist` before `npm run build` so a cached/committed
  `dist/` can never leak into production, and a build-time check that
  fails if `dist/` is older than the newest `src/` file. Added a root
  `.dockerignore` that excludes `**/dist/`, `**/node_modules/`, and
  secrets from the build context. Fixes #1 (crypto polyfill), #2 (SPA-side
  context headers), #6 (Test Connection chip), #7 (DEMO banner) — all
  were already correct in source but missing from the deployed bundle.
- **CI freshness check** (`.github/workflows/ci.yml`): the `web` job now
  verifies `apps/web/dist` is fresh relative to `src/` and fails if any
  `dist/` artifacts are committed to git.
- **Makefile target** (`Makefile`): added `verify-fresh-bundle` so the
  same check can be run locally.

### Fixed (Fusion kernel — codebase map audit)
- **`listOrganizations` context scoping** (`services/control-plane-kernel/
  internal/httpapi/server.go`): now calls `fusioncontext.FromHeaders` and
  filters the result to the caller's organization when
  `X-Fusion-Organization-Id` is present. A platform superadmin (no org
  header) still sees all organizations. Previously this endpoint was
  effectively public after the auth middleware and bypassed tenant
  scoping.
- **`listAuditEvents` superadmin-only** (`services/control-plane-kernel/
  internal/httpapi/server.go`): now requires a valid request context AND
  platform superadmin access (via `phase1.PlatformOverview`). Audit
  events are platform-wide and may contain cross-tenant data.

### Fixed (local-dev infra)
- **HostPath PV cross-platform** (`infra/local-dev/k8s/00-infra.yaml`):
  the Postgres hostPath `path` is now parameterized via the
  `${DCRAFT_PGDATA_HOSTPATH}` placeholder. `deploy.ps1` substitutes it
  from the `DCRAFT_PGDATA_HOSTPATH` env var (default
  `/var/dcraft-local/postgres-data` for Docker Desktop Linux VM).
  Operators can override to a Windows host bind (e.g.
  `E:/Dcraft/.tmp/pgdata`) without editing the YAML.

### Changed
- Bumped `dcraft-fusion` and `fusion-cdc` Helm charts to `version: 1.2.5`
  / `appVersion: "1.2.5"` (`infra/helm/*/Chart.yaml`).
- Bumped all image tags from `1.2.4` → `1.2.5` in
  `infra/helm/*/values.yaml`, `infra/helm/*/examples/values-minimal.yaml`,
  and `infra/local-dev/k8s/values-{cdc,fusion}-local.yaml`.
- Bumped `--version 1.2.5` in `infra/local-dev/k8s/deploy.ps1` and the
  helm install examples in `infra/helm/README.md`.

### Notes
- The CDC-side fixes (alert_rules/suppressions migrations, system_alerts
  model cleanup, JWT/encryption/worker-secret fail-fast, seed health
  surface, dq_policies router dedup) live in the private
  `fusion-cdc-engine` repo and ship in its v1.2.5 release.

## [1.2.4] — 2026-07-22

Follow-up to v1.2.3. The v1.2.3 verification + Fusion UI audit (HTTP-only,
192.168.1.10:8088) found 7 issues. This release fixes the 4 Fusion-SPA-side
issues (1 HIGH, 1 HIGH, 2 LOW). The 3 CDC-side fixes live in the private
`fusion-cdc-engine` repo (alert_rules migration, route ordering, Kafka
widget).

### Fixed
- **Login crashed on HTTP with `crypto.randomUUID is not a function`
  (HIGH).** `apps/web/src/ui/auth.ts:92` called `crypto.randomUUID()` for
  the OIDC `state` nonce, which is `undefined` in non-secure HTTP contexts
  (browsers gate `crypto.randomUUID` behind `window.isSecureContext`). Every
  HTTP deployment — including `http://192.168.1.10:8088` — threw on the
  login redirect. Added a 5-line polyfill at the top of `auth.ts` that
  falls back to a `Math.random`-based RFC-4122 v4 generator when
  `crypto.randomUUID` is missing.
- **Bearer-token API calls returned 400 in `authMode: dev` (HIGH).** The
  local-dev gateway (`infra/local-dev/k8s/values-fusion-local.yaml:12`) runs
  in `authMode: dev`, which makes the middleware a pass-through that does
  NOT inject `X-Fusion-Actor-Id` / `X-Fusion-Organization-Id` /
  `X-Fusion-Tenant-Id` / `X-Fusion-Project-Id` from JWT claims. The SPA's
  `api.ts` only sent `Authorization` + `X-Fusion-Correlation-Id`, so
  tenant-scoped endpoints rejected requests. Added SPA-side defense-in-
  depth: `api.ts` now decodes the JWT payload (base64) and sends the four
  `X-Fusion-*` context headers from claims, so the SPA works regardless of
  gateway auth mode. `values-fusion-local.yaml` is left in `dev` mode (the
  JWKS URL config is fragile; SPA-side injection is the more robust fix).
- **"Test Connection" optimistic UI masked failures (LOW).**
  `apps/web/src/ui/App.tsx:760-767` flipped the connection status chip to
  `healthy` BEFORE calling `/api/v1/connections/{id}/test`, so a failing
  backend still showed a green chip with a small easy-to-miss banner. The
  handler now marks the chip `untested` while the call is in flight, flips
  to `healthy` only on HTTP 200, flips to `blocked` on error, and renders
  the error in a `message-banner error` (red) tone with the failure reason
  prefixed by `Connection test failed:`.
- **Demo data presented as real (LOW).** When `/api/v1/bootstrap` failed,
  `App.tsx` silently fell back to the bundled `createDemoWorkspace()` mock
  (199 lines of hardcoded data in `apps/web/src/ui/demo-workspace.ts`),
  making the UI look live when the backend was down. Added a
  `bootstrapFailed` state that flips on bootstrap error and renders a red
  `DEMO DATA — backend unreachable` banner at the top of the workspace
  shell explaining the situation and that changes will not persist.

### Changed
- Bumped `dcraft-fusion` and `fusion-cdc` Helm charts to `version: 1.2.4` /
  `appVersion: "1.2.4"` (`infra/helm/*/Chart.yaml`).
- Bumped all image tags from `1.2.3` → `1.2.4` in
  `infra/helm/dcraft-fusion/values.yaml`,
  `infra/helm/fusion-cdc/values.yaml`,
  `infra/helm/dcraft-fusion/examples/values-minimal.yaml`,
  `infra/helm/fusion-cdc/examples/values-minimal.yaml`,
  `infra/local-dev/k8s/values-fusion-local.yaml`,
  `infra/local-dev/k8s/values-cdc-local.yaml`.
- Bumped `--version 1.2.4` in `infra/local-dev/k8s/deploy.ps1` and the
  `infra/helm/README.md` install examples.

## [1.2.3] — 2026-07-21

Follow-up to v1.2.2. The v1.2.2 remote retest (192.168.1.10) confirmed the
self-healing seed worked (6 connectors, 1 source, 2 destinations, 1
connection all present), but found 4 remaining issues. This release fixes
the Fusion-kernel-side and chart-config issues; the CDC frontend / control-
plane fixes live in the private `fusion-cdc-engine` repo.

### Fixed
- **Fusion kernel JWKS misconfig — every protected kernel API route 401s
  with `Get "": unsupported protocol scheme ""` (HIGH).** The kernel runs in
  non-dev auth mode (`FUSION_AUTH_MODE=password`) but
  `FUSION_OIDC_JWKS_URL` was empty in the chart values, so the RS256 JWKS
  fetch failed for every bearer-token request. The Fusion SPA masked this
  by falling back to bundled seed data, so the UI looked alive while the
  kernel API was unusable by any real bearer client.
  - `infra/helm/dcraft-fusion/templates/configmap.yaml`: `FUSION_OIDC_JWKS_URL`
    now auto-defaults to the in-cluster kernel service URL
    (`http://<release>-control-plane-kernel:<port>/oidc/jwks`) when the
    value is empty, so RS256 verification works for any release name without
    operator intervention.
  - `infra/helm/dcraft-fusion/values.yaml`: documented the JWKS URL default
    and the `authMode` switch in `controlPlaneKernel.config`.
  - `infra/local-dev/k8s/values-fusion-local.yaml`: set
    `controlPlaneKernel.config.authMode: dev` for local-dev / community
    deployments, which bypasses the JWKS fetch entirely (per the kernel's
    auth middleware in
    `services/control-plane-kernel/internal/auth/middleware.go`). Prod
    values keep `authMode: password` and rely on the JWKS URL default.

### Changed
- **Helm chart + image tags bumped `1.2.2` → `1.2.3`** so `deploy.ps1`
  pulls the fixed images:
  `infra/helm/dcraft-fusion/Chart.yaml`,
  `infra/helm/fusion-cdc/Chart.yaml`,
  `infra/helm/dcraft-fusion/values.yaml`,
  `infra/helm/fusion-cdc/values.yaml`,
  `infra/helm/dcraft-fusion/examples/values-minimal.yaml`,
  `infra/helm/fusion-cdc/examples/values-minimal.yaml`,
  `infra/local-dev/k8s/values-fusion-local.yaml`,
  `infra/local-dev/k8s/values-cdc-local.yaml`,
  `infra/helm/README.md`, and the two `--version` flags in
  `infra/local-dev/k8s/deploy.ps1`.
- **`infra/helm/fusion-cdc/templates/control-plane.yaml`:** wire
  `KAFKA_BOOTSTRAP_SERVERS` into the control-plane deployment env (sourced
  from the same `fusion-cdc.kafkaBootstrapServers` helper used by the
  cdc-workers and spark-consumer) so the control-plane's
  `/api/v1/monitoring/health` endpoint can probe Kafka.

## [1.2.2] — 2026-07-21

Hotfix for the v1.2.1 regression: the CDC metadata DB stayed empty
(`connector-definitions`, `sources`, `destinations`, `connections` all
returned `total: 0`) even though `deploy.ps1` reported a successful seed
and the admin user existed. Root cause and self-healing fix live in the
private `fusion-cdc-engine` repo; this public release bumps chart/image
versions so `deploy.ps1` pulls the fixed control-plane image.

### Fixed
- **CDC seed still empty after v1.2.1 (BLOCKER).** The v1.2.1
  `deploy.ps1` fail-fast caught `psql` exit codes but the seed SQL itself
  was broken: `scripts/seed-admin.sql` INSERTed into `destinations` with
  columns (`host`, `port`, `database_name`, `schema_name`, `username`,
  `password_encrypted`, `ssl_enabled`, `ssl_config`, `config`) that do
  NOT exist on the `destinations` table — those fields live inside the
  `connection_config` JSONB column. The `connections` INSERT likewise
  listed `sync_enabled`, `replication_slot`, `publication`,
  `namespace_definition`, `namespace_format`, `stream_prefix`, `config`,
  which do not exist on `connections`. Because the whole seed runs as a
  single atomic `DO $$ ... $$;` block, the first non-existent-column
  error rolled back the entire transaction — including the
  `connector_definitions` INSERTs that ran earlier in the block —
  leaving every CDC table at `total: 0`. The admin user existed only
  because it had been registered manually via `/auth/register` (the
  control-plane has NO startup admin-creation hook). Fixed in the
  private repo by rewriting the `destinations`/`connections` INSERTs
  against the real schema and, more importantly, by adding a
  self-healing seed hook to the control-plane startup
  (`control-plane/app/seed/seed_admin.py`) that re-seeds
  `connector_definitions` whenever it finds an empty DB — so the
  deployment no longer depends on `kubectl cp` succeeding.
- **`infra/local-dev/k8s/deploy.ps1` — seed path is now a fallback.** The
  `kubectl cp + psql -f` seed step is retained as a manual re-seed
  fallback, but the primary seed mechanism is now the control-plane's
  startup hook (the seed SQL is idempotent, so running both is a no-op
  when the DB is already populated). Added a note in `deploy.ps1`
  explaining the new ordering.

### Changed
- **Helm chart + image tags bumped `1.2.1` → `1.2.2`** so `deploy.ps1`
  pulls the fixed control-plane image:
  `infra/helm/dcraft-fusion/Chart.yaml`,
  `infra/helm/fusion-cdc/Chart.yaml`,
  `infra/helm/dcraft-fusion/values.yaml`,
  `infra/helm/fusion-cdc/values.yaml`,
  `infra/helm/dcraft-fusion/examples/values-minimal.yaml`,
  `infra/helm/fusion-cdc/examples/values-minimal.yaml`,
  `infra/local-dev/k8s/values-fusion-local.yaml`,
  `infra/local-dev/k8s/values-cdc-local.yaml`,
  `infra/helm/README.md`, and the two `--version` flags in
  `infra/local-dev/k8s/deploy.ps1`.

## [1.2.1] — 2026-07-21

Hotfix release addressing the three hard blockers found when verifying v1.2.0
against the remote (192.168.1.10) deployment, plus the missing Kafka
dependency for the CDC pipeline.

### Fixed
- **`infra/local-dev/k8s/deploy.ps1` — silent seed failure (BLOCKER).** The
  `kubectl cp` and `kubectl exec ... psql -f` calls were piped to
  `| Out-Null` with no `$LASTEXITCODE` check, so a failed seed left the CDC
  metadata DB unpopulated and the UI showed empty connector lists with no
  error. Removed `| Out-Null`, added `-v ON_ERROR_STOP=1` to the psql call,
  and added `if ($LASTEXITCODE -ne 0) { throw "Seed failed — aborting deploy" }`.
  Audited every other `| Out-Null` / missing exit-code check in the script
  (infra apply, postgres/redis rollout, CREATE DATABASE, helm uninstall,
  kubectl patch, control-plane/frontend/web rollout) and made them all
  fail-fast with descriptive throws.
- **`infra/local-dev/k8s/deploy.ps1` — seed path fallback masked missing
  private repo.** The script silently fell back to the admin-only
  `seed-cdc-admin.sql` stub when the full
  `.tmp/fusion-cdc-engine-private/scripts/seed-admin.sql` was missing,
  producing a CDC DB with an admin user but no connectors. Changed the
  fallback to a hard error with a clear message telling the operator to
  git-fetch / clone the private fusion-cdc-engine repo.
- **`infra/local-dev/k8s/00-infra.yaml` — postgres data wiped on every pod
  restart (BLOCKER).** Postgres ran on an `emptyDir` because the Docker
  Desktop `storage-provisioner` is in CrashLoopBackOff and dynamic PVCs
  stay Pending. Added a `hostPath` PersistentVolume (`/var/dcraft-local/
  postgres-data`) plus a matching PersistentVolumeClaim (storageClassName
  `""`, reclaimPolicy `Retain`) and wired it into the postgres deployment
  with `PGDATA=/var/lib/postgresql/data/pgdata`. Data now survives pod
  restarts. **Local-dev only** — production must use a real storage class
  or external managed Postgres.
- **`infra/local-dev/k8s/seed-connectors-job.yaml` — dead code.** The
  ConfigMap body was a placeholder comment with no SQL and the Job was
  never wired into `deploy.ps1`. Deleted the file; the fixed `deploy.ps1`
  `kubectl cp + psql -f` path is now the single source of truth for
  seeding.

### Added
- **Kafka in the public `fusion-cdc` Helm chart (BLOCKER — user #1
  priority).** The CDC architecture is Debezium-style
  (cdc-workers → Kafka → spark-consumer) but Kafka previously lived only
  in the private repo's kustomize manifests (`kubernetes/base/kafka.yaml`)
  and never reached the published Helm chart. Ported the Apache Kafka 3.7.0
  KRaft StatefulSet (no ZooKeeper) into a new
  `infra/helm/fusion-cdc/templates/kafka.yaml`, gated by a new
  `kafka.enabled` value (default `true`). Added a
  `fusion-cdc.kafkaBootstrapServers` helper and wired
  `KAFKA_BOOTSTRAP_SERVERS` env into the `cdc-workers` and `spark-consumer`
  deployment templates (defaults to the in-cluster service
  `<release>-kafka:9092` when `kafka.enabled=true`, empty otherwise so
  operators can point at external MSK / Confluent Cloud via per-workload
  `env` overrides). Added a `kafka:` block to `values.yaml` (replicas,
  bootstrapServers, image, numPartitions, replicationFactor, retention,
  heap, externalNodePort, storage, resources, security contexts) and a
  matching entry to `values.schema.json`. Local-dev overrides in
  `values-cdc-local.yaml` pin a single small broker that fits the Docker
  Desktop VM.

### Changed
- Bumped chart versions: `dcraft-fusion` `1.1.2 → 1.2.1`,
  `fusion-cdc` `1.1.2 → 1.2.1` (both `version` and `appVersion`).
- Bumped all image tags in chart `values.yaml`, `examples/values-minimal.yaml`,
  and `infra/local-dev/k8s/values-*-local.yaml` from `1.1.1`/`1.1.2` to
  `1.2.1`.
- Bumped the `--version` flags in `deploy.ps1` and the install snippets in
  `infra/helm/README.md` from `1.1.2` to `1.2.1`.

### Notes
- The `dcraft-fusion-stack` umbrella chart (`infra/helm/dcraft-fusion-stack`)
  is a thin file-dependency guide not used by `deploy.ps1`; its dependency
  version pins (`1.0.1`) are pre-existing drift and not updated in this
  release.

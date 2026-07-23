# Deploy DCraft Fusion + Fusion CDC to local Docker Desktop Kubernetes.
# Prerequisites: Docker Desktop running with Kubernetes enabled.
$ErrorActionPreference = "Stop"
$helm = if (Test-Path "$env:TEMP\helm\windows-amd64\helm.exe") {
  "$env:TEMP\helm\windows-amd64\helm.exe"
} else {
  "helm"
}
$root = Split-Path (Split-Path $PSScriptRoot -Parent) -Parent
if (-not (Test-Path "$root\infra\helm\dcraft-fusion")) {
  $root = "E:\Dcraft"
}

Write-Host "==> checking cluster"
kubectl cluster-info | Select-Object -First 3
kubectl get nodes

Write-Host "==> applying infra (postgres/redis/secrets)"
# Render 00-infra.yaml with the DCRAFT_PGDATA_HOSTPATH placeholder substituted
# from the environment. Default to the Docker Desktop Linux VM path that works
# cross-platform; operators can override (e.g. to a Windows host bind path)
# by setting DCRAFT_PGDATA_HOSTPATH before running this script.
$pgDataHostPath = if ($env:DCRAFT_PGDATA_HOSTPATH) {
  $env:DCRAFT_PGDATA_HOSTPATH
} else {
  "/var/dcraft-local/postgres-data"
}
Write-Host "    hostPath: $pgDataHostPath (override via DCRAFT_PGDATA_HOSTPATH)"
$infraTemplate = Get-Content -Raw -Path "$PSScriptRoot\00-infra.yaml"
$infraRendered = $infraTemplate.Replace('${DCRAFT_PGDATA_HOSTPATH}', $pgDataHostPath)
$infraRenderedPath = Join-Path $env:TEMP "dcraft-00-infra-rendered.yaml"
Set-Content -Path $infraRenderedPath -Value $infraRendered -Encoding ascii
kubectl apply -f $infraRenderedPath
if ($LASTEXITCODE -ne 0) { throw "kubectl apply 00-infra.yaml failed (exit $LASTEXITCODE) — aborting deploy" }
kubectl -n dcraft-local rollout status deploy/postgres --timeout=180s
if ($LASTEXITCODE -ne 0) { throw "postgres rollout failed (exit $LASTEXITCODE) — aborting deploy" }
kubectl -n dcraft-local rollout status deploy/redis --timeout=120s
if ($LASTEXITCODE -ne 0) { throw "redis rollout failed (exit $LASTEXITCODE) — aborting deploy" }

# Ensure CDC metadata DB exists (init may have raced)
# PowerShell: keep kubectl args on one line (backslash is NOT a line continuation)
$exists = (kubectl -n dcraft-local exec deploy/postgres -- psql -U fusion -d fusion -Atc "SELECT 1 FROM pg_database WHERE datname='fusion_cdc_metadata'" 2>$null).Trim()
if ($exists -ne "1") {
  Write-Host "==> creating fusion_cdc_metadata database"
  kubectl -n dcraft-local exec deploy/postgres -- psql -U fusion -d fusion -c "CREATE DATABASE fusion_cdc_metadata;"
  if ($LASTEXITCODE -ne 0) { throw "CREATE DATABASE fusion_cdc_metadata failed (exit $LASTEXITCODE) — aborting deploy" }
}

# Clear stuck/failed Helm releases so upgrades are clean
$prevEap = $ErrorActionPreference
$ErrorActionPreference = "Continue"
foreach ($rel in @("fusion", "fusion-cdc")) {
  $raw = & $helm status $rel -n dcraft-local -o json 2>&1 | Out-String
  if ($LASTEXITCODE -ne 0 -or $raw -notmatch '"name"') { continue }
  $st = ($raw | ConvertFrom-Json)
  if ($st.info.status -in @("failed", "pending-install", "pending-upgrade", "pending-rollback", "uninstalling")) {
    Write-Host "==> clearing stuck release $rel ($($st.info.status))"
    & $helm uninstall $rel -n dcraft-local --wait --timeout 2m
    if ($LASTEXITCODE -ne 0) {
      Write-Host "WARN: helm uninstall $rel failed (exit $LASTEXITCODE) — manually clearing release secrets"
      kubectl -n dcraft-local delete secret -l name=$rel,owner=helm --ignore-not-found
      if ($LASTEXITCODE -ne 0) {
        throw "Failed to clear stuck Helm release $rel (uninstall + secret delete both failed, exit $LASTEXITCODE) — aborting deploy"
      }
      kubectl -n dcraft-local get secret -o name 2>$null |
        Select-String "sh.helm.release.v1.$rel" |
        ForEach-Object { kubectl -n dcraft-local delete $_ --ignore-not-found }
    }
  }
}
$ErrorActionPreference = $prevEap


Write-Host "==> helm install dcraft-fusion"
$prevEap = $ErrorActionPreference
$ErrorActionPreference = "Continue"
& $helm upgrade --install fusion oci://ghcr.io/dcraft-labs/charts/dcraft-fusion `
  --version 1.2.35 `
  --namespace dcraft-local `
  -f "$root\infra\helm\dcraft-fusion\examples\values-minimal.yaml" `
  -f "$PSScriptRoot\values-fusion-local.yaml" `
  --set global.secrets.existingSecret=fusion-secrets `
  --set externalRedis.addr=redis:6379 `
  --wait --timeout 5m
if ($LASTEXITCODE -ne 0) {
  Write-Host "WARN: fusion helm wait incomplete (exit $LASTEXITCODE) — continuing"
}

Write-Host "==> helm install fusion-cdc"
& $helm upgrade --install fusion-cdc oci://ghcr.io/dcraft-labs/charts/fusion-cdc `
  --version 1.2.35 `
  --namespace dcraft-local `
  -f "$root\infra\helm\fusion-cdc\examples\values-minimal.yaml" `
  -f "$PSScriptRoot\values-cdc-local.yaml" `
  --set global.secrets.existingSecret=fusion-cdc-secrets `
  --set externalRedis.url=redis://redis:6379/0 `
  --timeout 3m
# Image CMD uses 4 uvicorn workers and OOMs on small Docker Desktop VMs.
Write-Host "==> patching CDC control-plane to 1 uvicorn worker"
$patchFile = Join-Path $PSScriptRoot "patch-cdc-cp-workers.json"
if (-not (Test-Path $patchFile)) {
  @'
[
  {"op":"add","path":"/spec/template/spec/containers/0/command","value":["sh","-c","alembic upgrade head && uvicorn app.main:app --host 0.0.0.0 --port 8000 --workers 1 --loop uvloop"]}
]
'@ | Set-Content -Path $patchFile -Encoding ascii
}
kubectl -n dcraft-local patch deploy fusion-cdc-control-plane --type=json --patch-file $patchFile
if ($LASTEXITCODE -ne 0) {
  throw "kubectl patch deploy fusion-cdc-control-plane failed (exit $LASTEXITCODE) — aborting deploy"
}
kubectl -n dcraft-local rollout status deploy/fusion-cdc-control-plane --timeout=300s
if ($LASTEXITCODE -ne 0) { throw "fusion-cdc-control-plane rollout failed (exit $LASTEXITCODE) — aborting deploy" }
kubectl -n dcraft-local rollout status deploy/fusion-cdc-frontend --timeout=180s
if ($LASTEXITCODE -ne 0) { throw "fusion-cdc-frontend rollout failed (exit $LASTEXITCODE) — aborting deploy" }
kubectl -n dcraft-local rollout status deploy/fusion-dcraft-fusion-web --timeout=180s
if ($LASTEXITCODE -ne 0) { throw "fusion-dcraft-fusion-web rollout failed (exit $LASTEXITCODE) — aborting deploy" }

Write-Host "==> seeding CDC default admin + connectors (admin / Admin@123)"
# NOTE (v1.2.2): The PRIMARY seed mechanism is now the control-plane's
# self-healing startup hook (control-plane/app/seed/seed_admin.py), which runs
# AFTER Alembic migrations on every boot and re-seeds connector_definitions
# whenever it finds an empty DB. The kubectl cp + psql -f path below is kept
# as a FALLBACK for manual re-seeding and for environments where the
# control-plane image hasn't been rebuilt with the new seed module yet.
# If the control-plane already seeded the DB, the seed SQL is idempotent
# (ON CONFLICT DO NOTHING / WHERE NOT EXISTS) and will be a no-op.
$seedSql = "$root\.tmp\fusion-cdc-engine-private\scripts\seed-admin.sql"
if (-not (Test-Path $seedSql)) {
  # The full seed (admin user + connector definitions + sample source/destination)
  # lives in the private fusion-cdc-engine repo. Without it, the CDC DB is
  # missing connectors and the UI shows empty lists — a silent, hard-to-debug
  # failure mode. Fail fast instead of falling back to the admin-only stub.
  throw @"
Seed SQL not found: $seedSql

The full seed (admin + connectors) lives in the private fusion-cdc-engine
repo. Either:
  1. git-fetch / clone the private repo at:
     $root\.tmp\fusion-cdc-engine-private\
     (run: git clone <fusion-cdc-engine remote> .tmp\fusion-cdc-engine-private)
  2. Or restore the submodule if it was vendored as one.

Aborting deploy — refusing to silently fall back to the admin-only stub
(seed-cdc-admin.sql) which leaves the CDC UI without connectors.
"@
}
# NOTE: piping [System.IO.File]::ReadAllBytes(...) into `kubectl exec -i` via
# PowerShell's pipeline does NOT send a raw byte stream — PowerShell formats
# each byte object as text, corrupting the SQL (and mangling the → arrow,
# U+2192, into "???"). `kubectl cp` + `psql -f` (run entirely inside the pod)
# preserves the file's UTF-8 bytes exactly and is the only reliable way to do
# this from Windows PowerShell.
$pgPod = (kubectl -n dcraft-local get pods -l app=postgres -o jsonpath='{.items[0].metadata.name}').Trim()
if ([string]::IsNullOrWhiteSpace($pgPod)) {
  throw "Could not find a postgres pod (label app=postgres) to seed into — aborting."
}
Push-Location (Split-Path $seedSql -Parent)
$leaf = Split-Path $seedSql -Leaf
Write-Host "==> copying $leaf into pod $pgPod"
kubectl cp $leaf "dcraft-local/${pgPod}:/tmp/seed-admin.sql"
if ($LASTEXITCODE -ne 0) { Pop-Location; throw "kubectl cp seed-admin.sql failed (exit $LASTEXITCODE) — aborting deploy" }
Pop-Location
Write-Host "==> running psql -v ON_ERROR_STOP=1 -f /tmp/seed-admin.sql"
kubectl -n dcraft-local exec $pgPod -- psql -U fusion -d fusion_cdc_metadata -v ON_ERROR_STOP=1 -f /tmp/seed-admin.sql
if ($LASTEXITCODE -ne 0) {
  throw "Seed failed — psql exited $LASTEXITCODE. The CDC metadata DB was NOT seeded. Aborting deploy. Inspect the pod logs and the seed SQL at $seedSql."
}
Write-Host "==> seed applied successfully"
$ErrorActionPreference = $prevEap

Write-Host "==> pods"
kubectl -n dcraft-local get pods -o wide
Write-Host "==> services"
kubectl -n dcraft-local get svc

Write-Host @"

Port-forward (run in separate terminals):
  kubectl -n dcraft-local port-forward svc/fusion-dcraft-fusion-web 8088:80
  kubectl -n dcraft-local port-forward svc/fusion-dcraft-fusion-control-plane-kernel 18080:8080
  kubectl -n dcraft-local port-forward svc/fusion-cdc-frontend 8090:8080
  kubectl -n dcraft-local port-forward svc/fusion-cdc-control-plane 18000:8000

Then open:
  http://127.0.0.1:8088   (Fusion web)
  http://127.0.0.1:18080/healthz  (Fusion kernel)
  http://127.0.0.1:8090   (CDC frontend)
  http://127.0.0.1:18000/api/v1/monitoring/health  (CDC API)

CDC login (seeded):
  username: admin
  password: Admin@123

Fusion login (kernel seed — separate from CDC):
  founder@dcraftlabs.com / changeme-founder
  superadmin@dcraftlabs.com / changeme-superadmin

Register tip: use a real-looking email (not *.local), password with upper+lower+digit, min 8 chars.
"@

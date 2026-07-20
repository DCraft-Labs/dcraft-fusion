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
kubectl apply -f "$PSScriptRoot\00-infra.yaml"
kubectl -n dcraft-local rollout status deploy/postgres --timeout=180s
kubectl -n dcraft-local rollout status deploy/redis --timeout=120s

# Ensure CDC metadata DB exists (init may have raced)
# PowerShell: keep kubectl args on one line (backslash is NOT a line continuation)
$exists = (kubectl -n dcraft-local exec deploy/postgres -- psql -U fusion -d fusion -Atc "SELECT 1 FROM pg_database WHERE datname='fusion_cdc_metadata'" 2>$null).Trim()
if ($exists -ne "1") {
  Write-Host "==> creating fusion_cdc_metadata database"
  kubectl -n dcraft-local exec deploy/postgres -- psql -U fusion -d fusion -c "CREATE DATABASE fusion_cdc_metadata;"
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
    & $helm uninstall $rel -n dcraft-local --wait --timeout 2m 2>&1 | Out-Null
    if ($LASTEXITCODE -ne 0) {
      kubectl -n dcraft-local delete secret -l name=$rel,owner=helm --ignore-not-found 2>&1 | Out-Null
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
  --version 1.1.2 `
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
  --version 1.1.2 `
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
kubectl -n dcraft-local rollout status deploy/fusion-cdc-control-plane --timeout=300s
kubectl -n dcraft-local rollout status deploy/fusion-cdc-frontend --timeout=180s
kubectl -n dcraft-local rollout status deploy/fusion-dcraft-fusion-web --timeout=180s

Write-Host "==> seeding CDC default admin + connectors (admin / Admin@123)"
$seedSql = "$root\.tmp\fusion-cdc-engine-private\scripts\seed-admin.sql"
if (Test-Path $seedSql) {
  Get-Content -Raw $seedSql |
    kubectl -n dcraft-local exec -i deploy/postgres -- psql -U fusion -d fusion_cdc_metadata 2>&1 | Out-Null
} else {
  Write-Host "WARN: $seedSql not found — falling back to local seed-cdc-admin.sql"
  Get-Content -Raw "$PSScriptRoot\seed-cdc-admin.sql" |
    kubectl -n dcraft-local exec -i deploy/postgres -- psql -U fusion -d fusion_cdc_metadata 2>&1 | Out-Null
}
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

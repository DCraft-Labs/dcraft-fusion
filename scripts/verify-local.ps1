$ErrorActionPreference = "Stop"

Push-Location (Join-Path $PSScriptRoot "..")
try {
  $go = Get-Command go -ErrorAction SilentlyContinue
  if ($null -eq $go) {
    $fallbackGo = "C:\Program Files\Go\bin\go.exe"
    if (Test-Path $fallbackGo) {
      $goPath = $fallbackGo
    }
    else {
      $goPath = $null
    }
  }
  else {
    $goPath = $go.Source
  }

  Push-Location services\control-plane-kernel
  try {
    if ($null -eq $goPath) {
      docker run --rm -v "${PWD}:/src" -w /src/services/control-plane-kernel golang:1.26 sh -lc "/usr/local/go/bin/go test ./..."
    }
    else {
      & $goPath test ./...
    }
  }
  finally {
    Pop-Location
  }
  npm test
  npm run typecheck
  npm run build
  npm audit
  kubectl kustomize infra\kubernetes | Out-Null
  docker compose -f infra\local-dev\docker-compose.yml ps
  Invoke-WebRequest -UseBasicParsing http://127.0.0.1:8080/healthz | Out-Null
  Invoke-WebRequest -UseBasicParsing http://127.0.0.1:5174/healthz | Out-Null
}
finally {
  Pop-Location
}

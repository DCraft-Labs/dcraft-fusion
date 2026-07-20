$ErrorActionPreference = "Stop"

$BaseUrl = if ($env:FUSION_BASE_URL) { $env:FUSION_BASE_URL } else { "http://127.0.0.1:30173" }
$Headers = @{
  "X-Fusion-Actor-Id" = "load-test"
  "X-Fusion-Correlation-Id" = "load-test-$([DateTimeOffset]::UtcNow.ToUnixTimeMilliseconds())"
  "X-Fusion-Organization-Id" = "org-b2b2c-1"
  "X-Fusion-Tenant-Id" = "tenant-brand-a"
  "X-Fusion-Project-Id" = "project-prod"
}

1..100 | ForEach-Object {
  Invoke-RestMethod -Uri "$BaseUrl/api/v1/bootstrap" -Headers $Headers | Out-Null
}

Invoke-RestMethod -Uri "$BaseUrl/metrics" | Out-String

param(
  [ValidateSet("up", "down", "logs", "ps")]
  [string] $Command = "up"
)

$compose = Join-Path $PSScriptRoot "..\infra\local-dev\docker-compose.yml"

if ($Command -eq "up") {
  docker compose -f $compose up --build -d
  exit $LASTEXITCODE
}

if ($Command -eq "down") {
  docker compose -f $compose down
  exit $LASTEXITCODE
}

if ($Command -eq "logs") {
  docker compose -f $compose logs --tail=150
  exit $LASTEXITCODE
}

docker compose -f $compose ps

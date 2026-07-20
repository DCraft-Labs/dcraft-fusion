param(
  [ValidateSet("deploy", "delete")]
  [string] $Command = "deploy"
)

$root = Resolve-Path (Join-Path $PSScriptRoot "..")

if ($Command -eq "deploy") {
  docker build -t dcraft-fusion/control-plane-kernel:local -f "$root\services\control-plane-kernel\Dockerfile" "$root"
  docker build -t dcraft-fusion/web:local -f "$root\apps\web\Dockerfile" "$root"
  kubectl apply -k "$root\infra\kubernetes"
  exit $LASTEXITCODE
}

kubectl delete -k "$root\infra\kubernetes"

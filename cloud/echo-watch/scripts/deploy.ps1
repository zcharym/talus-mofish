# Production deploy helper for echo-watch Worker + PWA assets.
# Prerequisites: wrangler login, KV namespace, VAPID secrets.

$ErrorActionPreference = "Stop"
$workerDir = Join-Path $PSScriptRoot "..\worker"
Push-Location $workerDir
try {
  if (-not (Test-Path "node_modules")) {
    npm install
  }

  Write-Host "Deploying echo-watch Worker..."
  npx wrangler deploy

  Write-Host ""
  Write-Host "Deployed. Next steps:"
  Write-Host "  1. Open the Worker URL in Safari on iPhone"
  Write-Host "  2. Add to Home Screen, enable notifications with your pairing token"
  Write-Host "  3. curl -X POST https://<worker>/api/alert -H 'Authorization: Bearer <token>' -H 'Content-Type: application/json' -d '{\"rule_id\":\"test\",\"title\":\"Echo Watch\",\"body\":\"Hello\"}'"
}
finally {
  Pop-Location
}

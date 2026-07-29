# Simplified local verification — requires wrangler dev already running on :8787
$ErrorActionPreference = "Stop"

$health = Invoke-RestMethod -Uri "http://127.0.0.1:8787/health" -Method GET
if (-not $health.ok) { throw "health failed" }
Write-Host "OK /health"

$manifest = Invoke-WebRequest -Uri "http://127.0.0.1:8787/manifest.webmanifest" -UseBasicParsing
if ($manifest.StatusCode -ne 200) { throw "manifest failed" }
Write-Host "OK /manifest.webmanifest"

Write-Host "Local worker routes verified."

$ErrorActionPreference = "Stop"
$PSScriptRoot = Split-Path -Parent $MyInvocation.MyCommand.Path
$Root = Split-Path -Parent $PSScriptRoot
Set-Location $Root

Write-Host "=== Starting Restgeld Live Dev Environment (HMR) ===" -ForegroundColor Cyan
Write-Host "Frontend: http://localhost:5173" -ForegroundColor Green
Write-Host "Backend:  http://localhost:8080" -ForegroundColor Green
Write-Host "Database: localhost:5432" -ForegroundColor Green
docker compose -f docker-compose.dev.yml up --build

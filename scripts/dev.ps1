$ErrorActionPreference = "Stop"
$PSScriptRoot = Split-Path -Parent $MyInvocation.MyCommand.Path
$Root = Split-Path -Parent $PSScriptRoot
Set-Location $Root

Write-Host "=== Starting Restgeld Local Dev Environment ===" -ForegroundColor Cyan
docker compose up --build

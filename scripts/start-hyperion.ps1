$ErrorActionPreference = "Stop"
$exe = Join-Path $PSScriptRoot "..\bin\hyperion.exe"
if (-not (Test-Path $exe)) { throw "Hyperion executable not found. Run .\scripts\build.ps1 first." }
Start-Process -FilePath $exe -ArgumentList "serve" -WorkingDirectory (Split-Path -Parent $PSScriptRoot)
Start-Sleep -Seconds 2
Start-Process "http://127.0.0.1:8080"
Write-Host "Hyperion started. Keep its process running in the background." -ForegroundColor Green

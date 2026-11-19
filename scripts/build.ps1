$ErrorActionPreference = "Stop"
go mod tidy
New-Item -ItemType Directory -Force -Path .\bin | Out-Null
go build -trimpath -ldflags "-s -w" -o .\bin\hyperion.exe .\cmd\hyperion
Write-Host "Built .\bin\hyperion.exe" -ForegroundColor Green

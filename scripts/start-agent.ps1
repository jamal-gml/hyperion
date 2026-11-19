param(
  [Parameter(Mandatory=$true)][string]$Project,
  [Parameter(Mandatory=$true)][ValidateSet("dev","code-review","architect","qa","coordinator")][string]$Role,
  [Parameter(Mandatory=$true)][string]$WorkingDirectory,
  [string]$Provider = "claude",
  [string]$Instance = "",
  [string]$AnthropicaDirectory = "",
  [string]$HyperionExe = "$PSScriptRoot\..\bin\hyperion.exe"
)
$ErrorActionPreference = "Stop"
if (-not (Get-Command claude -ErrorAction SilentlyContinue)) { throw "Claude Code CLI was not found in PATH. Install/authenticate Claude Code first." }
if (-not (Test-Path $WorkingDirectory -PathType Container)) { throw "Working directory does not exist: $WorkingDirectory" }
if (-not (Test-Path $HyperionExe -PathType Leaf)) { throw "Hyperion executable not found: $HyperionExe. Build it first with go build -o .\bin\hyperion.exe .\cmd\hyperion" }
$agentId = if ($Instance) { "$Project-$Role-$Instance" } else { "$Project-$Role" }
$dataDir = Join-Path (Split-Path -Parent $PSScriptRoot) "data"
$promptFile = Join-Path (Split-Path -Parent $PSScriptRoot) "configs\roles\$Role.md"
$tempConfig = Join-Path $env:TEMP "hyperion-mcp-$agentId.json"
$config = @{
  mcpServers = @{
    hyperion = @{
      type = "stdio"
      command = (Resolve-Path $HyperionExe).Path
      args = @("mcp", "stdio")
      env = @{
        HYPERION_AGENT_ID = $agentId
        HYPERION_PROJECT = $Project
        HYPERION_ROLE = $Role
        HYPERION_PROVIDER = $Provider
        HYPERION_DATA_DIR = $dataDir
      }
    }
  }
} | ConvertTo-Json -Depth 8
[System.IO.File]::WriteAllText($tempConfig, $config, [System.Text.UTF8Encoding]::new($false))
Write-Host "Starting Hyperion agent: $agentId" -ForegroundColor Cyan
Write-Host "Project: $Project | Role: $Role | Working directory: $WorkingDirectory"
Write-Host "This is a separate Claude Code session. Shared state is exchanged only through Hyperion MCP tools."
Push-Location $WorkingDirectory
try {
  $claudeArgs = @("--mcp-config", $tempConfig, "--append-system-prompt-file", $promptFile)
  if ($AnthropicaDirectory) {
    if (-not (Test-Path $AnthropicaDirectory -PathType Container)) { throw "Anthropica directory not found: $AnthropicaDirectory" }
    $claudeArgs += @("--add-dir", (Resolve-Path $AnthropicaDirectory).Path)
  }
  & claude @claudeArgs
  if ($LASTEXITCODE -ne 0) { throw "Claude Code exited with code $LASTEXITCODE" }
}
finally {
  Pop-Location
  Remove-Item $tempConfig -ErrorAction SilentlyContinue
}

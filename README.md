# Hyperion

**Local-first MCP collaboration hub for multiple AI assistant instances.**

Hyperion lets separate AI client processes exchange tasks and messages through a shared local coordination service. It is designed for setups such as multiple Claude Code terminals working in tandem:

- `your-project-dev` — implementation
- `your-project-code-review` — independent review
- `your-project-architect` — design decisions

The first foundation includes a Go service, a local SQLite store, an MCP Streamable HTTP endpoint, an MCP STDIO mode for per-process identity, a small REST API, and Windows launch scripts.

> Important: MCP provides a shared tool interface. It does not automatically let Hyperion initiate work inside an arbitrary interactive AI session. The session must be running and connected, or a separately supported worker runtime must be used. The initial version uses user-started Claude Code processes and MCP tools; it does not use model-provider APIs.

## Requirements
- Go 1.23+
- Claude Code installed and authenticated for the Claude multi-instance launcher
- Windows PowerShell for the supplied scripts (the Go service itself is cross-platform)

## Quick start (Windows)

```powershell
# 1. Download dependencies and build
 go mod tidy
 go build -o .\bin\hyperion.exe .\cmd\hyperion

# 2. Start the local coordination service (terminal 1)
 .\bin\hyperion.exe serve

# 3. Open the dashboard
 start http://127.0.0.1:8080
```

The default data directory is `./data`; override with `HYPERION_DATA_DIR`. The service binds to `127.0.0.1` by default. Set `HYPERION_LISTEN_ADDR` only if you deliberately want another bind address.

## Start multiple Claude Code instances

From the repository root, open separate PowerShell terminals:

```powershell
.\scripts\start-agent.ps1 -Project your-project -Role dev -WorkingDirectory C:\work\your-project
.\scripts\start-agent.ps1 -Project your-project -Role code-review -WorkingDirectory C:\work\your-project
```

Each launcher creates a unique agent ID and a temporary MCP STDIO configuration. Both processes use the same local database and therefore see the same project task/message board. They retain separate Claude Code sessions and role instructions.

You can also start a third instance:

```powershell
.\scripts\start-agent.ps1 -Project your-project -Role architect -WorkingDirectory C:\work\your-project
```

Available roles: `dev`, `code-review`, `architect`, `qa`, `coordinator`. Use `-Instance 2` when you need multiple instances with the same role.

To make the separate Anthropica repository available as additional Claude Code context, pass its local path: `-AnthropicaDirectory C:\work\Anthropica`. The launcher adds it as an allowed working directory; each agent should consult only the relevant Anthropica instructions/skills. This does not automatically install Anthropica as a Claude Code plugin.

## MCP endpoint

- Streamable HTTP: `http://127.0.0.1:8080/mcp`
- REST health: `http://127.0.0.1:8080/healthz`
- Dashboard: `http://127.0.0.1:8080/`

The HTTP MCP endpoint is suitable for clients that support local Streamable HTTP. The STDIO launcher is the recommended route for multiple Claude Code instances because it assigns a distinct agent identity to each process.

## First collaboration test

1. Start Hyperion.
2. Start a `dev` Claude instance and a `code-review` Claude instance in separate terminals.
3. In the dev terminal, ask it to call `bridge_create_task` with a clear implementation task and assign it to `your-project-code-review` only after the implementation is ready.
4. In the review terminal, ask it to call `bridge_list_tasks` and `bridge_read_messages`, then review the requested change and call `bridge_submit_result`.
5. In the dev terminal, ask it to read the result and respond.

The bridge stores tasks and messages. It does not autonomously inject messages into a terminal that is not actively using the MCP tools. That event-driven worker/notification layer is a later milestone.

## Environment variables

| Variable | Default | Purpose |
|---|---|---|
| `HYPERION_DATA_DIR` | `./data` | SQLite database location |
| `HYPERION_LISTEN_ADDR` | `127.0.0.1:8080` | HTTP listen address |
| `HYPERION_AGENT_ID` | generated | Identity for STDIO MCP process |
| `HYPERION_PROJECT` | `default` | Project scope for STDIO agent |
| `HYPERION_ROLE` | `assistant` | Role for STDIO agent |
| `HYPERION_PROVIDER` | `unknown` | Provider label for STDIO agent |

## Anthropica companion repository

Hyperion treats Anthropica as a separate reusable instruction/skill project. Do not copy Anthropica's files into every target repository. Instead, configure the desired Anthropica instruction or skill paths in the consuming repository's Claude setup, and have Hyperion agents follow the target repo's `CLAUDE.md`/`AGENTS.md`. A future adapter can expose approved Anthropica packs as versioned project policies. Hyperion must not assume that the other repository exists at a fixed path.

## Current foundation limitations

- Local-only trust boundary; no multi-user/network deployment security yet.
- SQLite is used for the starter; PostgreSQL is a planned storage adapter.
- HTTP MCP endpoint is local and currently unauthenticated. Keep it bound to loopback.
- Agent identities supplied by STDIO environment are local process identity hints, not cryptographic authentication.
- No automatic model invocation, browser automation or provider API integration is included.
- Dashboard is intentionally minimal and reports real stored state only.

See `docs/ARCHITECTURE.md`, `docs/INTEGRATIONS.md`, `docs/SECURITY.md` and `docs/ROADMAP.md`.

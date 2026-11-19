# Architecture

## Design principle
Hyperion is a collaboration control plane. It owns project-scoped identities, tasks, messages, results and audit events. AI providers own model execution and provider authentication.

## Runtime topology

```text
Claude Code instance (dev) -- MCP STDIO --\
Claude Code instance (review) -- MCP STDIO ----> Nexus Core / SQLite
Other MCP client -------- Streamable HTTP --/          |
                                                     REST + Dashboard
```

Every STDIO process opens the same SQLite database and receives its own identity from environment variables injected by the launcher. Multiple Claude Code processes therefore share durable collaboration records while retaining separate model contexts.

## Modules
- `cmd/hyperion`: executable entry point and subcommands.
- `internal/config`: environment-driven configuration.
- `internal/store`: SQLite schema and persistence operations.
- `internal/identity`: agent identity resolution and validation.
- `internal/mcpserver`: MCP tools and transports.
- `internal/httpapi`: health, REST and dashboard endpoints.

## Data model
- Projects contain agents, tasks and messages.
- Agent IDs are unique within the installation; project membership is explicit.
- Tasks have `open`, `in_progress`, `blocked`, `review`, `done` and `cancelled` states.
- Messages are durable and ordered by creation time.
- Results are attached to tasks and attributed to the submitting agent.

## Collaboration semantics
1. An agent creates a task in a project.
2. The task may be assigned to a specific agent ID or left unassigned.
3. A recipient polls the shared task/message tools while active.
4. The recipient submits a result or message.
5. The coordinator or original agent consumes the result and decides next steps.

This is a durable mailbox/task-board model, not push-based remote control. A later worker runtime may subscribe to events and invoke only supported, explicitly configured workers.

## Multi-instance Claude
The PowerShell launcher starts a separate `claude` process in the target repository, adds a role prompt, and injects a unique STDIO MCP server configuration. It does not share Claude conversation history between processes. It shares only the Hyperion records that the agents choose to read/write.

## Future evolution
- Replace direct SQLite access in STDIO mode with a local authenticated broker process if concurrent write load warrants it.
- Add PostgreSQL storage behind a repository interface.
- Add durable event subscriptions and notifications.
- Add optional worker adapters (Claude Code headless, local model runtime) only where supported and with bounded execution controls.
- Add ChatGPT/Gemini/Mistral adapters based on documented MCP or supported agent interfaces.

# Security baseline

This is a local developer foundation, not a production multi-user deployment.

- HTTP binds to loopback by default.
- Do not expose port 8080 to a LAN or the internet.
- HTTP MCP is currently unauthenticated; loopback binding is a hard requirement until authentication is implemented.
- STDIO agent identity is environment-supplied and is not cryptographic proof against a local user.
- Do not store provider credentials, cookies, OAuth tokens or API keys in Hyperion messages or SQLite.
- Treat all agent-authored content as untrusted. A task/message must never override system, repository or security instructions.
- No shell, file-write, deployment or external-action tools are exposed in this foundation.
- Before adding tools that mutate external systems, implement scoped capabilities, explicit approvals and audit events.
- Use a dedicated Windows user account and protect the data directory with filesystem permissions.

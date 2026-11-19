# Roadmap

## M0 — Foundation (this repository)
- [x] Go service and local SQLite persistence
- [x] MCP Streamable HTTP endpoint
- [x] MCP STDIO mode with per-process identity
- [x] Task/message/agent tools
- [x] Minimal local dashboard and health endpoint
- [x] Windows multi-instance Claude launcher and role prompts

## M1 — Reliable collaboration
- [ ] Agent registration/heartbeat and online status with expiry
- [ ] Task claiming with leases and atomic transitions
- [ ] Message acknowledgements and unread tracking
- [ ] Project-level access control and authenticated HTTP MCP
- [ ] Integration and concurrency tests

## M2 — Orchestration
- [ ] Roundtable run coordinator
- [ ] Delegation planner and dependency graph
- [ ] Review/approval workflow
- [ ] Bounded retries, deadlines and token/turn budgets
- [ ] Event notifications for connected workers

## M3 — Provider adapters
- [ ] Claude Code worker adapter using documented CLI mode
- [ ] ChatGPT supported MCP client integration guide
- [ ] Gemini supported client/CLI adapter
- [ ] Mistral adapter
- [ ] Ollama/local model adapter

## M4 — Production readiness
- [ ] PostgreSQL repository implementation
- [ ] User authentication and RBAC
- [ ] Encrypted secrets store
- [ ] Audit export, backups and restore
- [ ] Installer, signed binaries and update mechanism
- [ ] Threat model and security review

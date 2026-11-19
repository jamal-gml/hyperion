# Claude Code instructions for Hyperion

Read `AGENTS.md` first. It is the canonical instruction set for all coding agents.

## Project-specific guidance
- Treat Hyperion as a multi-instance collaboration control plane, not a model aggregator.
- Preserve the distinction between an MCP-connected interactive client and a background-capable worker. A client that can call MCP tools is not automatically remotely invokable.
- Multiple Claude Code terminals are separate agents. Give each process a unique `HYPERION_AGENT_ID`, project and role through the supplied launcher.
- Do not merge separate agents' identities or conversation contexts.
- Do not introduce provider API calls as a hidden dependency. Any API-backed adapter must be optional and explicitly configured.
- Do not use `--dangerously-skip-permissions` in launcher scripts.
- Keep shell execution outside the Hyperion MCP tool surface in the initial release.

## Before implementation
1. Inspect the existing code and tests.
2. Explain the intended change briefly.
3. Implement the smallest complete change.
4. Run formatting, tests and static checks available in the environment.
5. Update documentation and report limitations.

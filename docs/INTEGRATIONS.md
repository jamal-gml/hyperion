# Integrations and API-less strategy

## Integration classes

### 1. MCP client (preferred)
A client connects to Hyperion and calls its tools. This works when the client supports MCP. It does not mean Hyperion can start a task inside that client; the client must be active and use the tools.

### 2. Local CLI worker (optional future adapter)
A local CLI may be launched as a distinct worker if its documented CLI supports the required mode and authentication. Existing user authentication may avoid a separate API key, but CLI capabilities, terms and account limits still apply. The adapter must not scrape browser sessions or reuse private tokens.

### 3. Provider API adapter (optional)
An explicit opt-in adapter may use a provider API key. This is not part of the API-less core and must disclose usage/cost.

## Current support
- Claude Code: supported through per-instance STDIO MCP launch configuration. The user must already have Claude Code installed and authenticated.
- ChatGPT: future/conditional. Connect through a supported custom MCP app/client surface when available to the user's account. Do not claim that this local service is automatically accessible from every ChatGPT plan or interface.
- Gemini: future/conditional. Use a documented MCP-capable Gemini client or supported CLI workflow.
- Mistral: future/conditional. Use a documented MCP-capable client or local runtime.
- Generic MCP: Streamable HTTP endpoint and STDIO server are provided.

## Claude Code launch model
The launcher generates an MCP config similar to:

```json
{
  "mcpServers": {
    "hyperion": {
      "type": "stdio",
      "command": "C:\path\to\hyperion.exe",
      "args": ["mcp", "stdio"],
      "env": {
        "HYPERION_AGENT_ID": "your-project-dev",
        "HYPERION_PROJECT": "your-project",
        "HYPERION_ROLE": "dev",
        "HYPERION_PROVIDER": "claude"
      }
    }
  }
}
```

The launcher uses `--mcp-config` and `--append-system-prompt-file`. These are Claude Code CLI options; verify compatibility with the installed Claude Code version. It deliberately does not use `--dangerously-skip-permissions`.

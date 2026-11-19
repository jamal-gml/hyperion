# your-project multi-instance workflow

Start Hyperion once, then open separate PowerShell terminals:

```powershell
# Terminal 1: implementation agent
.\scripts\start-agent.ps1 -Project your-project -Role dev -WorkingDirectory C:\work\your-project -AnthropicaDirectory C:\work\Anthropica

# Terminal 2: independent reviewer
.\scripts\start-agent.ps1 -Project your-project -Role code-review -WorkingDirectory C:\work\your-project -AnthropicaDirectory C:\work\Anthropica

# Terminal 3: another developer instance, with a distinct identity
.\scripts\start-agent.ps1 -Project your-project -Role dev -Instance 2 -WorkingDirectory C:\work\your-project -AnthropicaDirectory C:\work\Anthropica
```

## Suggested first prompt for the dev instance

"Register yourself with Hyperion. Create a task for implementing [feature]. Work only on the assigned feature. When the change is ready, send the reviewer a message with the task ID and summary, then set the task to review. Do not claim the reviewer has completed work until a result is submitted."

## Suggested first prompt for the reviewer

"Register yourself with Hyperion. Check tasks assigned to you and read messages. For the requested task, inspect the actual git diff and relevant tests independently. Do not edit the implementation. Submit findings with severity, file and line references. If no issues are found, state exactly what you reviewed."

## Coordination caveat

The agents exchange durable tasks and messages only when their active sessions call the MCP tools. The initial version does not push a message into an idle Claude terminal. The user can prompt each instance to check its inbox; event-driven worker coordination is planned for a later phase.

## Avoiding file collisions

Use separate Git worktrees for two agents that both need to edit code. A read-only reviewer may inspect the developer worktree, but independent implementation agents should work in separate worktrees/branches and submit changes for review.

# CLI commands used by this app

## Subcommands of the agent binary

The process reads `os.Args[1:]` and dispatches on the first argument (case-insensitive, trimmed). Default when omitted or empty is `start`.

| Invocation | Handler |
|------------|---------|
| *(no argument)*, `start`, or `""` | Start background agent (`startCommand`) |
| `agent` | Start HTTP API server (`startAPIServer`) |
| `monitor` | Monitoring loop (`runMonitor`) |
| `send-to-hub` | Hub send path (`runSendToHub` — referenced in router) |
| `stop` | Stop agent (`stopCommand`) |
| `restart` | Stop then start (`runRestart`) |
| `status` | Print running state (`runStatus`) |
| `version` | Print build version (`buildVersion`) |

## External shell processes

The application invokes exactly one external-style process via `os/exec`:

| What runs | Where | Purpose |
|-----------|--------|---------|
| `<executable> agent` | `cmd/agent-server/agent_commands.go` — `exec.Command(executable, "agent")` | `executable` is `os.Executable()` (this same binary); used to spawn the API server as a detached child when the user runs `start`. |

No other `exec.Command`, shell scripts, or subprocess calls appear in the repository’s Go sources.

## Notes

- Process control for `stop` uses `os.FindProcess` and `syscall.SIGTERM`, not a separate CLI tool.
- Dependencies such as MySQL or HTTP probing use Go libraries (`database/sql`, `net/http`), not CLI binaries, in the code under this repo.

---
name: tidewave-cli-usage
description: For Elixir/Phoenix/Ash development in worktree or non-MCP contexts. Provides CLI commands (via tidewave_cli) for live code evaluation, module navigation, SQL execution, schema introspection, and documentation access. Use when MCP tools are unavailable or port changes per worktree.
---

## Tidewave CLI — Elixir Runtime Tools

### When to Use

Use the `tidewave` CLI when Tidewave MCP tools are not available (e.g., worktree agents, non-MCP setups, or when the port changes). The CLI wraps the same 8 Tidewave tools as simple bash subcommands.

### Available Commands

| Command | Description | Usage |
|---------|-------------|-------|
| `eval` | Evaluate Elixir code | `tidewave eval '<code>' [--args '[1,2]'] [--timeout 5000]` |
| `logs` | Get application logs | `tidewave logs [--tail 20] [--grep pattern] [--level error]` |
| `source` | Find module source location | `tidewave source MyApp.Worker` |
| `docs` | Get module/function docs | `tidewave docs "Phoenix.Controller.render/3"` |
| `sql` | Execute SQL queries | `tidewave sql 'SELECT * FROM users' [--args '[1]'] [--repo MyApp.Repo]` |
| `schemas` | List Ecto schemas | `tidewave schemas` |
| `search` | Search Hex package docs | `tidewave search "genserver" [--packages phoenix,ecto]` |
| `ash` | List Ash resources | `tidewave ash` |

### Global Flags

All commands accept: `--port` (default: 4000 or `TIDEWAVE_PORT`), `--host` (default: localhost or `TIDEWAVE_HOST`), `--path` (default: /tidewave/mcp or `TIDEWAVE_PATH`)

### Tool Usage Hierarchy

#### 1. Code Evaluation — ALWAYS Use tidewave eval

**NEVER use Bash to run Elixir code!** Instead:

- WRONG: `bash: mix run -e "IO.inspect(MyModule.function())"`
- RIGHT: `bash: tidewave eval 'IO.inspect(MyModule.function())'`

Test function behavior, explore modules, access IEx helpers, capture IO output.

#### 2. Source Code Navigation — tidewave source First

Before using Grep/Glob/Read for Elixir code:

- `tidewave source MyApp.Worker` — find exact file path
- `tidewave docs MyModule.function/2` — get docs without reading files

- WRONG: `grep: pattern: "defmodule Worker"`
- RIGHT: `bash: tidewave source MyApp.Worker`

#### 3. Database Operations — Direct SQL

- `tidewave sql 'SELECT * FROM users LIMIT 10'` — run SQL directly
- `tidewave sql 'SELECT * FROM users WHERE id = $1' --args '[1]'` — parameterized
- `tidewave schemas` — list all Ecto schemas first

#### 4. Dependency Documentation

- `tidewave search "genserver" --packages phoenix,ecto` — search Hex docs

#### 5. Error Diagnosis

- `tidewave logs --tail 20 --grep "error"` — recent error logs

### Workflow Patterns

#### Understanding a Module
1. `tidewave docs MyModule`
2. `tidewave source MyModule`
3. `tidewave eval 'exports(MyModule)'`
4. Read file if needed

#### Testing Code Changes
1. `tidewave eval 'MyModule.new_function(:test) |> IO.inspect()'`
2. Verify behavior
3. Modify file

#### Database Work
1. `tidewave schemas`
2. `tidewave sql 'SELECT * FROM table LIMIT 10'`
3. `tidewave eval 'MyApp.Repo.all(MySchema)'`

#### Debugging Issues
1. `tidewave logs --tail 50 --grep "error"`
2. `tidewave source MyApp.ProblematicModule`
3. `tidewave eval 'reproduce_issue()'`

### Database Query Gotchas

- UUIDs return as 16-byte binaries — cast with `::text` (PostgreSQL)
- Results limited to 50 rows — use LIMIT/OFFSET
- Use parameterized queries: `tidewave sql 'SELECT * FROM users WHERE id = $1' --args '[123]'`

### Port Configuration

In worktree scenarios, the Phoenix server may run on a different port. Set it with:
- `--port 4001` on each command
- `export TIDEWAVE_PORT=4001` for the session

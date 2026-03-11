---
name: tidewave-cli-usage
description: CRITICAL for ALL Elixir/Phoenix/Ash development work. Invoke when working with Elixir code, Ecto schemas, Ash resources, Phoenix applications, or databases in Elixir projects. Provides CLI tools for live code evaluation (via IEx), instant module navigation, direct SQL execution, schema introspection, and documentation access.
---

## Tidewave CLI — Elixir Runtime Tools

### Policy

Use `tidewave_cli` as the single Tidewave interface in this repo.

### When to Use

Use `tidewave_cli` for Tidewave runtime introspection (eval, logs, docs, source, SQL, schemas, search, ash), especially when worktree-specific ports are required.

### Preflight

- `command -v tidewave_cli`
- `tidewave_cli --version`
- Optional connectivity check: `tidewave_cli schemas --host localhost --port 4001`

### Command Syntax

- Use `tidewave_cli <command> [command flags] [global flags]`.
- Put global flags after the subcommand.
- Invalid: `tidewave_cli --host localhost schemas`
- Valid: `tidewave_cli schemas --host localhost --port 4001`

### Available Commands

| Command | Description | Usage |
|---------|-------------|-------|
| `eval` | Evaluate Elixir code | `tidewave_cli eval '<code>' [--args '[1,2]'] [--timeout 5000]` |
| `logs` | Get application logs | `tidewave_cli logs [--tail 20] [--grep pattern] [--level error]` |
| `source` | Find module source location | `tidewave_cli source MyApp.Worker` |
| `docs` | Get module/function docs | `tidewave_cli docs "Phoenix.Controller.render/3"` |
| `sql` | Execute SQL queries | `tidewave_cli sql 'SELECT * FROM users' [--args '[1]'] [--repo MyApp.Repo]` |
| `schemas` | List Ecto schemas | `tidewave_cli schemas` |
| `search` | Search Hex package docs | `tidewave_cli search "genserver" [--packages phoenix,ecto]` |
| `ash` | List Ash resources | `tidewave_cli ash` |

### Global Flags

All commands accept: `--port` (default: 4000 or `TIDEWAVE_PORT`), `--host` (default: localhost or `TIDEWAVE_HOST`), `--path` (default: /tidewave/mcp or `TIDEWAVE_PATH`).

### Tool Usage Hierarchy

#### 1. Code Evaluation — ALWAYS Use tidewave_cli eval

**NEVER use Bash to run Elixir code!** Instead:

- WRONG: `bash: mix run -e "IO.inspect(MyModule.function())"`
- RIGHT: `bash: tidewave_cli eval 'IO.inspect(MyModule.function())'`

Test function behavior, explore modules, access IEx helpers, capture IO output, pass arguments with `--args`, set custom timeout for long-running operations.

#### 2. Source Code Navigation — tidewave_cli source First

Before using Grep/Glob/Read for Elixir code:

- `tidewave_cli source MyApp.Worker` — find exact file path
- `tidewave_cli source "dep:phoenix"` — find dependency source
- `tidewave_cli docs MyModule.function/2` — get docs without reading files
- `tidewave_cli docs "c:GenServer.init/1"` — get callback docs

- WRONG: `grep: pattern: "defmodule Worker"`
- RIGHT: `bash: tidewave_cli source MyApp.Worker`

#### 3. Database Operations — Direct SQL

- `tidewave_cli sql 'SELECT * FROM users LIMIT 10'` — run SQL directly
- `tidewave_cli sql 'SELECT * FROM users WHERE id = $1' --args '[1]'` — parameterized
- `tidewave_cli schemas` — list all Ecto schemas first
- Auto-detects available Ecto repositories; use `--repo` only to target a specific one

#### 4. Dependency Documentation

- `tidewave_cli search "genserver" --packages phoenix,ecto` — search Hex docs
- Searches project dependencies by default; use `--packages` to narrow scope

#### 5. Error Diagnosis

- `tidewave_cli logs --tail 20 --grep "error"` — recent error logs

### Workflow Patterns

#### Understanding a Module
1. `tidewave_cli docs MyModule`
2. `tidewave_cli source MyModule`
3. `tidewave_cli eval 'exports(MyModule)'`
4. Read file if needed

#### Testing Code Changes
1. `tidewave_cli eval 'MyModule.new_function(:test) |> IO.inspect()'`
2. Verify behavior
3. Modify file

#### Database Work
1. `tidewave_cli schemas`
2. `tidewave_cli sql 'SELECT * FROM table LIMIT 10'`
3. `tidewave_cli eval 'MyApp.Repo.all(MySchema)'`

#### Debugging Issues
1. `tidewave_cli logs --tail 50 --grep "error"`
2. `tidewave_cli source MyApp.ProblematicModule`
3. `tidewave_cli eval 'reproduce_issue()'`

### IEx Helpers Available in eval

- `h(Module)` — get help for a module
- `exports(Module)` — list all exported functions
- `i(value)` — inspect data structure info
- `t(Module)` — show types defined in module
- `b(Module)` — show behaviours module implements

### Common Mistakes to Avoid

- Don't use `bash` to run `mix` commands for code evaluation — use `tidewave_cli eval`
- Don't use `grep` to find module definitions when you know the module name — use `tidewave_cli source`
- Don't read entire files to find function documentation — use `tidewave_cli docs`
- Don't run `iex` in bash — use `tidewave_cli eval`
- Don't search the file system for Ecto schemas — use `tidewave_cli schemas` first

### Database Query Gotchas

- UUIDs return as 16-byte binaries — cast with `::text` (PostgreSQL)
- Results limited to 50 rows — use LIMIT/OFFSET
- Use parameterized queries: `tidewave_cli sql 'SELECT * FROM users WHERE id = $1' --args '[123]'`

### Port Configuration

In worktree scenarios, the Phoenix server may run on a different port. Set it with:
- `--port 4001` on each command, e.g. `tidewave_cli schemas --host localhost --port 4001`
- `export TIDEWAVE_PORT=4001` for the session
- Optional session defaults: `export TIDEWAVE_HOST=localhost` and `export TIDEWAVE_PATH=/tidewave/mcp`

# tidewave_cli

A Go CLI for [Tidewave](https://github.com/tidewave-elixir/tidewave) MCP tools. Wraps Tidewave's 8 Elixir/Phoenix runtime introspection tools as simple shell commands.

## Why?

Tidewave exposes powerful MCP tools for interacting with a running Elixir application — eval, logs, docs, source location, SQL, schemas, hex search, and Ash resources. These are normally accessed via MCP (Model Context Protocol) integration in editors and coding agents.

This CLI exists for situations where MCP isn't available:

- **Worktree agents** — the MCP port changes per worktree and can't be preconfigured
- **Non-MCP setups** — tools or agents that can run shell commands but don't support MCP
- **Quick terminal access** — faster than setting up an MCP client for a one-off query

Each tool is a subcommand that POSTs a JSON-RPC 2.0 request to the Tidewave server and prints the result.

## Install

Requires Go 1.24+.

```sh
go install tidewave_cli@latest
```

Or build from source:

```sh
git clone <repo-url>
cd tidewave_cli
go build -o tidewave_cli
```

## Codex Skill Install (symlink workflow)

To keep Codex skills in sync with repository updates, symlink the repo's `skills/` directory:

```sh
mkdir -p ~/.codex/skills
ln -s ~/.codex/tidewave_cli/skills ~/.codex/skills/tidewave
```

Then restart Codex so it refreshes skill discovery.

Detailed steps: [`.codex/INSTALL.md`](.codex/INSTALL.md).

## Usage

```
tidewave_cli <command> [flags] [args]
```

### Commands

| Command | Description | Example |
|---------|-------------|---------|
| `eval` | Evaluate Elixir code | `tidewave_cli eval 'Enum.map(1..5, &(&1*2))'` |
| `logs` | Get application logs | `tidewave_cli logs --tail 20 --grep error` |
| `source` | Find source file location | `tidewave_cli source MyApp.Worker` |
| `docs` | Get module/function docs | `tidewave_cli docs "Phoenix.Controller.render/3"` |
| `sql` | Execute SQL via Ecto | `tidewave_cli sql 'SELECT * FROM users LIMIT 10'` |
| `schemas` | List Ecto schemas | `tidewave_cli schemas` |
| `search` | Search Hex package docs | `tidewave_cli search "genserver" --packages phoenix` |
| `ash` | List Ash resources | `tidewave_cli ash` |

### Connection flags

All commands accept:

| Flag | Env var | Default | Description |
|------|---------|---------|-------------|
| `--host` | `TIDEWAVE_HOST` | `localhost` | Server host |
| `--port` | `TIDEWAVE_PORT` | `4000` | Server port |
| `--path` | `TIDEWAVE_PATH` | `/tidewave/mcp` | MCP endpoint path |

### Examples

```sh
# Evaluate Elixir code
tidewave_cli eval 'Enum.sum(1..100)'

# Evaluate with arguments
tidewave_cli eval 'Enum.map(arguments, &(&1 * 2))' --args '[1,2,3]'

# Get recent error logs
tidewave_cli logs --tail 50 --level error

# Find where a module is defined
tidewave_cli source MyApp.Accounts.User

# Look up documentation
tidewave_cli docs "Ecto.Changeset.cast/4"

# Run a parameterized SQL query
tidewave_cli sql 'SELECT * FROM users WHERE id = $1' --args '[1]'

# List all Ecto schemas
tidewave_cli schemas

# Search dependency docs
tidewave_cli search "plug" --packages phoenix

# Connect to a different port (e.g. worktree)
tidewave_cli --port 4001 eval '1 + 1'
# or
TIDEWAVE_PORT=4001 tidewave_cli eval '1 + 1'
```

## Skill Contents

This repository ships the skill at:

- `skills/tidewave-cli-usage/SKILL.md`

Optional UI metadata is in:

- `skills/tidewave-cli-usage/agents/openai.yaml`

## Requirements

- A running Elixir/Phoenix application with [Tidewave](https://github.com/tidewave-elixir/tidewave) installed
- Go 1.24+ (build only — the binary has zero runtime dependencies)

## License

MIT

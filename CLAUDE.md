# tidewave_cli

Go CLI wrapping Tidewave's MCP tools as shell subcommands via JSON-RPC 2.0.

## Build

```sh
go build -o tidewave_cli
```

Zero external dependencies — stdlib only.

## Testing

There are no unit tests. Validation is done against a live Tidewave server.

Run the full test suite after any change with `go test -v`. Requires a running Phoenix app with Tidewave at the default endpoint (`localhost:4000/tidewave/mcp`). See `docs/TESTING.md` for details.

## Architecture

- `client.go` — JSON-RPC 2.0 HTTP client (`CallTool`) and URL resolution (`resolveServerURL`)
- `commands.go` — all 8 subcommand implementations, each with its own `flag.FlagSet`
- `main.go` — entry point, subcommand dispatch, usage text

## MCP Tool Name Mapping

The CLI subcommand names differ from the MCP tool names:

| CLI command | MCP tool name |
|-------------|---------------|
| `eval` | `project_eval` |
| `logs` | `get_logs` |
| `source` | `get_source_location` |
| `docs` | `get_docs` |
| `sql` | `execute_sql_query` |
| `schemas` | `get_ecto_schemas` |
| `search` | `search_package_docs` |
| `ash` | `get_ash_resources` |

If Tidewave renames or changes tool parameter schemas, check with:

```sh
curl -s -X POST http://localhost:4000/tidewave/mcp \
  -H 'Content-Type: application/json' \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}' | python3 -m json.tool
```

## Commit Conventions

This repo uses [Conventional Commits](https://www.conventionalcommits.org/). Format:

```
<type>: <short summary>
```

Types: `feat`, `fix`, `docs`, `test`, `refactor`, `chore`

- Keep the subject line under 72 characters
- Use imperative mood ("add", "fix", not "added", "fixes")
- Body is optional — use it for context on *why*, not *what*

## Key Details

- Default path is `/tidewave/mcp` (not `/mcp`)
- `search_package_docs` uses `q` as the query parameter (not `query`)
- Response body is capped at 10 MB via `io.LimitReader`
- HTTP client timeout is 60s (server-side eval can take up to 30s)
- Optional params are only included in the JSON payload when explicitly set

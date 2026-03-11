# Testing Guide

Integration tests for tidewave_cli against a live Tidewave server.

## Prerequisites

- A running Phoenix application with Tidewave at `localhost:4000/tidewave/mcp`

## Running Tests

```sh
go test -v
```

To use a different port:

```sh
TIDEWAVE_PORT=4001 go test -v
```

If the Tidewave server is not reachable, the test runner will print a message and exit — no tests will be marked as failed or skipped.

## What's Covered

The test suite in `integration_test.go` covers:

| Test | What it verifies |
|------|------------------|
| `TestEval` | Basic arithmetic, arguments passing, timeout, error handling, string output |
| `TestLogs` | Default tail, grep filtering, level filtering |
| `TestSourceLocation` | App module resolution, stdlib rejection (expected error) |
| `TestDocs` | Module docs, function/arity docs |
| `TestSQL` | Basic query, parameterized queries |
| `TestSchemas` | Non-empty schema list |
| `TestSearch` | Basic hex search, packages filter |
| `TestAsh` | Ash resource listing (may be empty if app doesn't use Ash) |
| `TestConnectionOverride` | Wrong port produces connection error |

## Tidewave Compatibility Check

If tests start failing after a Tidewave update, inspect the current tool schemas:

```sh
curl -s -X POST http://localhost:4000/tidewave/mcp \
  -H 'Content-Type: application/json' \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}' \
  | python3 -c "
import sys, json
data = json.load(sys.stdin)
for t in data['result']['tools']:
    props = t.get('inputSchema', {}).get('properties', {})
    req = t.get('inputSchema', {}).get('required', [])
    print(f\"{t['name']}: params={list(props.keys())}, required={req}\")
"
```

Compare against the mapping in `CLAUDE.md`. Common breaking changes:
- Tool renamed (e.g. `project_eval` → `eval_code`)
- Parameter renamed (e.g. `q` → `query` or vice versa)
- Parameter type changed (e.g. `packages` from array to string)
- New required parameters added
- Endpoint path changed (e.g. `/tidewave/mcp` → `/mcp`)

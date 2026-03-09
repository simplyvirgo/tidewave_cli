package main

import (
	"fmt"
	"os"
)

var version = "0.1.0"

func main() {
	if len(os.Args) < 2 || os.Args[1] == "--help" || os.Args[1] == "-h" {
		printUsage()
		os.Exit(0)
	}

	if os.Args[1] == "--version" || os.Args[1] == "-v" {
		fmt.Println(version)
		os.Exit(0)
	}

	switch os.Args[1] {
	case "eval":
		cmdEval(os.Args[2:])
	case "logs":
		cmdLogs(os.Args[2:])
	case "source":
		cmdSource(os.Args[2:])
	case "docs":
		cmdDocs(os.Args[2:])
	case "sql":
		cmdSQL(os.Args[2:])
	case "schemas":
		cmdSchemas(os.Args[2:])
	case "search":
		cmdSearch(os.Args[2:])
	case "ash":
		cmdAsh(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Print(`tidewave - CLI for Tidewave MCP tools (Elixir/Phoenix runtime introspection)

Usage:
  tidewave <command> [options]

Commands:
  eval       Evaluate Elixir expression in the running application
  logs       Fetch recent application logs
  source     Find source file location for a module or function
  docs       Look up documentation for a module or function
  sql        Execute a SQL query via Ecto
  schemas    List Ecto schemas in the application
  search     Search Hex package documentation
  ash        List Ash resources

Global Flags:
  --host     Server host (env: TIDEWAVE_HOST, default: localhost)
  --port     Server port (env: TIDEWAVE_PORT, default: 4000)
  --path     MCP endpoint path (env: TIDEWAVE_PATH, default: /tidewave/mcp)

  -h, --help       Show this help message
  -v, --version    Show version
`)
}

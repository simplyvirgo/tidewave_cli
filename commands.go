package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

func cmdLogs(args []string) {
	fs := flag.NewFlagSet("logs", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: tidewave logs [flags]\n\nFetch logs from the running application.\n\nFlags:\n")
		fs.PrintDefaults()
	}

	host := fs.String("host", "", "server host")
	port := fs.String("port", "", "server port")
	path := fs.String("path", "", "server path")
	tail := fs.Int("tail", 20, "number of log lines to retrieve")
	grep := fs.String("grep", "", "filter logs by pattern")
	level := fs.String("level", "", "filter logs by level")

	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}

	arguments := map[string]any{
		"tail": *tail,
	}
	if *grep != "" {
		arguments["grep"] = *grep
	}
	if *level != "" {
		arguments["level"] = *level
	}

	serverURL := resolveServerURL(*host, *port, *path)
	result, isError, err := CallTool(serverURL, "get_logs", arguments)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	if isError {
		fmt.Fprintln(os.Stderr, result)
		os.Exit(1)
	}
	if !strings.HasSuffix(result, "\n") {
		result += "\n"
	}
	fmt.Print(result)
	os.Exit(0)
}

func cmdSource(args []string) {
	fs := flag.NewFlagSet("source", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: tidewave source [flags] <reference>\n\nGet source location for a module, function, or dependency.\n\nFlags:\n")
		fs.PrintDefaults()
	}

	host := fs.String("host", "", "server host")
	port := fs.String("port", "", "server port")
	path := fs.String("path", "", "server path")

	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}

	if len(fs.Args()) < 1 {
		fs.Usage()
		os.Exit(1)
	}
	reference := fs.Args()[0]

	arguments := map[string]any{
		"reference": reference,
	}

	serverURL := resolveServerURL(*host, *port, *path)
	result, isError, err := CallTool(serverURL, "get_source_location", arguments)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	if isError {
		fmt.Fprintln(os.Stderr, result)
		os.Exit(1)
	}
	if !strings.HasSuffix(result, "\n") {
		result += "\n"
	}
	fmt.Print(result)
	os.Exit(0)
}

func cmdDocs(args []string) {
	fs := flag.NewFlagSet("docs", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: tidewave docs [flags] <reference>\n\nGet documentation for a module, function, or callback.\n\nFlags:\n")
		fs.PrintDefaults()
	}

	host := fs.String("host", "", "server host")
	port := fs.String("port", "", "server port")
	path := fs.String("path", "", "server path")

	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}

	if len(fs.Args()) < 1 {
		fs.Usage()
		os.Exit(1)
	}
	reference := fs.Args()[0]

	arguments := map[string]any{
		"reference": reference,
	}

	serverURL := resolveServerURL(*host, *port, *path)
	result, isError, err := CallTool(serverURL, "get_docs", arguments)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	if isError {
		fmt.Fprintln(os.Stderr, result)
		os.Exit(1)
	}
	if !strings.HasSuffix(result, "\n") {
		result += "\n"
	}
	fmt.Print(result)
	os.Exit(0)
}

func cmdEval(args []string) {
	fs := flag.NewFlagSet("eval", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: tidewave eval [flags] '<code>'\n\nEvaluate Elixir code in the running application.\n\nFlags:\n")
		fs.PrintDefaults()
	}

	host := fs.String("host", "", "server host")
	port := fs.String("port", "", "server port")
	path := fs.String("path", "", "server path")
	argsJSON := fs.String("args", "", "arguments as a JSON array")
	timeout := fs.Int("timeout", 0, "evaluation timeout in milliseconds")

	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}

	if len(fs.Args()) < 1 {
		fs.Usage()
		os.Exit(1)
	}
	code := fs.Args()[0]

	arguments := map[string]any{
		"code": code,
	}
	if *argsJSON != "" {
		var parsed []any
		if err := json.Unmarshal([]byte(*argsJSON), &parsed); err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid JSON for --args: %v\n", err)
			os.Exit(1)
		}
		arguments["arguments"] = parsed
	}
	if *timeout != 0 {
		arguments["timeout"] = *timeout
	}

	serverURL := resolveServerURL(*host, *port, *path)
	result, isError, err := CallTool(serverURL, "project_eval", arguments)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	if isError {
		fmt.Fprintln(os.Stderr, result)
		os.Exit(1)
	}
	if !strings.HasSuffix(result, "\n") {
		result += "\n"
	}
	fmt.Print(result)
	os.Exit(0)
}

func cmdSQL(args []string) {
	fs := flag.NewFlagSet("sql", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: tidewave sql [flags] '<query>'\n\nExecute a SQL query against the application database.\n\nFlags:\n")
		fs.PrintDefaults()
	}

	host := fs.String("host", "", "server host")
	port := fs.String("port", "", "server port")
	path := fs.String("path", "", "server path")
	argsJSON := fs.String("args", "", "query arguments as a JSON array")
	repo := fs.String("repo", "", "Ecto repo module")

	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}

	if len(fs.Args()) < 1 {
		fs.Usage()
		os.Exit(1)
	}
	query := fs.Args()[0]

	arguments := map[string]any{
		"query": query,
	}
	if *argsJSON != "" {
		var parsed []any
		if err := json.Unmarshal([]byte(*argsJSON), &parsed); err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid JSON for --args: %v\n", err)
			os.Exit(1)
		}
		arguments["arguments"] = parsed
	}
	if *repo != "" {
		arguments["repo"] = *repo
	}

	serverURL := resolveServerURL(*host, *port, *path)
	result, isError, err := CallTool(serverURL, "execute_sql_query", arguments)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	if isError {
		fmt.Fprintln(os.Stderr, result)
		os.Exit(1)
	}
	if !strings.HasSuffix(result, "\n") {
		result += "\n"
	}
	fmt.Print(result)
	os.Exit(0)
}

func cmdSchemas(args []string) {
	fs := flag.NewFlagSet("schemas", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: tidewave schemas [flags]\n\nList all Ecto schemas in the application.\n\nFlags:\n")
		fs.PrintDefaults()
	}

	host := fs.String("host", "", "server host")
	port := fs.String("port", "", "server port")
	path := fs.String("path", "", "server path")

	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}

	arguments := map[string]any{}

	serverURL := resolveServerURL(*host, *port, *path)
	result, isError, err := CallTool(serverURL, "get_ecto_schemas", arguments)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	if isError {
		fmt.Fprintln(os.Stderr, result)
		os.Exit(1)
	}
	if !strings.HasSuffix(result, "\n") {
		result += "\n"
	}
	fmt.Print(result)
	os.Exit(0)
}

func cmdSearch(args []string) {
	fs := flag.NewFlagSet("search", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: tidewave search [flags] '<query>'\n\nSearch package documentation.\n\nFlags:\n")
		fs.PrintDefaults()
	}

	host := fs.String("host", "", "server host")
	port := fs.String("port", "", "server port")
	path := fs.String("path", "", "server path")
	packages := fs.String("packages", "", "comma-separated package names to search")

	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}

	if len(fs.Args()) < 1 {
		fs.Usage()
		os.Exit(1)
	}
	query := fs.Args()[0]

	arguments := map[string]any{
		"q": query,
	}
	if *packages != "" {
		arguments["packages"] = strings.Split(*packages, ",")
	}

	serverURL := resolveServerURL(*host, *port, *path)
	result, isError, err := CallTool(serverURL, "search_package_docs", arguments)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	if isError {
		fmt.Fprintln(os.Stderr, result)
		os.Exit(1)
	}
	if !strings.HasSuffix(result, "\n") {
		result += "\n"
	}
	fmt.Print(result)
	os.Exit(0)
}

func cmdAsh(args []string) {
	fs := flag.NewFlagSet("ash", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: tidewave ash [flags]\n\nList all Ash resources in the application.\n\nFlags:\n")
		fs.PrintDefaults()
	}

	host := fs.String("host", "", "server host")
	port := fs.String("port", "", "server port")
	path := fs.String("path", "", "server path")

	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}

	arguments := map[string]any{}

	serverURL := resolveServerURL(*host, *port, *path)
	result, isError, err := CallTool(serverURL, "get_ash_resources", arguments)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	if isError {
		fmt.Fprintln(os.Stderr, result)
		os.Exit(1)
	}
	if !strings.HasSuffix(result, "\n") {
		result += "\n"
	}
	fmt.Print(result)
	os.Exit(0)
}

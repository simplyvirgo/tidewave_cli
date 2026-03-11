package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

var testServerURL string

func TestMain(m *testing.M) {
	host := os.Getenv("TIDEWAVE_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("TIDEWAVE_PORT")
	if port == "" {
		port = "4000"
	}
	path := os.Getenv("TIDEWAVE_PATH")
	if path == "" {
		path = "/tidewave/mcp"
	}
	testServerURL = fmt.Sprintf("http://%s:%s%s", host, port, path)

	client := &http.Client{Timeout: 3 * time.Second}
	_, err := client.Post(testServerURL, "application/json",
		strings.NewReader(`{"jsonrpc":"2.0","id":0,"method":"tools/list","params":{}}`))
	if err != nil {
		fmt.Fprintf(os.Stderr, `
=============================================================
  Tidewave server is not reachable at %s

  Start a Phoenix application with Tidewave installed, then
  run the tests again:

    go test -v

  To use a different port:

    TIDEWAVE_PORT=4001 go test -v
=============================================================
`, testServerURL)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func TestEval(t *testing.T) {
	tests := []struct {
		name     string
		args     map[string]any
		contains string
		isError  bool
	}{
		{
			name:     "basic arithmetic",
			args:     map[string]any{"code": "1 + 1"},
			contains: "2",
		},
		{
			name:     "with arguments",
			args:     map[string]any{"code": "Enum.sum(arguments)", "arguments": []any{1, 2, 3}},
			contains: "6",
		},
		{
			name:     "with timeout",
			args:     map[string]any{"code": "Process.sleep(100); :ok", "timeout": 5000},
			contains: ":ok",
		},
		{
			name:     "raise returns stacktrace",
			args:     map[string]any{"code": `raise "boom"`},
			contains: "RuntimeError",
		},
		{
			name:     "string output",
			args:     map[string]any{"code": `"hello world"`},
			contains: "hello world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, isError, err := CallTool(testServerURL, "project_eval", tt.args)
			if err != nil {
				t.Fatalf("CallTool error: %v", err)
			}
			if isError != tt.isError {
				t.Errorf("isError = %v, want %v (result: %s)", isError, tt.isError, result)
			}
			if tt.contains != "" && !strings.Contains(result, tt.contains) {
				t.Errorf("result %q does not contain %q", result, tt.contains)
			}
		})
	}
}

func TestLogs(t *testing.T) {
	// Generate a log entry first
	CallTool(testServerURL, "project_eval", map[string]any{
		"code": `require Logger; Logger.info("tidewave_cli_test_marker")`,
	})

	tests := []struct {
		name string
		args map[string]any
	}{
		{
			name: "default tail",
			args: map[string]any{"tail": 5},
		},
		{
			name: "with grep",
			args: map[string]any{"tail": 50, "grep": "tidewave_cli_test_marker"},
		},
		{
			name: "with level",
			args: map[string]any{"tail": 10, "level": "info"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, isError, err := CallTool(testServerURL, "get_logs", tt.args)
			if err != nil {
				t.Fatalf("CallTool error: %v", err)
			}
			if isError {
				t.Errorf("unexpected tool error: %s", result)
			}
		})
	}
}

func TestSourceLocation(t *testing.T) {
	// Find an app module dynamically via schemas
	schemasResult, _, err := CallTool(testServerURL, "get_ecto_schemas", map[string]any{})
	if err != nil {
		t.Fatalf("failed to get schemas for test setup: %v", err)
	}

	// Extract first module name (lines like "* MyApp.Schema at lib/...}")
	var appModule string
	for _, line := range strings.Split(schemasResult, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "* ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				appModule = parts[1]
				break
			}
		}
	}

	t.Run("app module", func(t *testing.T) {
		if appModule == "" {
			t.Skip("no app module found from schemas")
		}
		result, isError, err := CallTool(testServerURL, "get_source_location", map[string]any{
			"reference": appModule,
		})
		if err != nil {
			t.Fatalf("CallTool error: %v", err)
		}
		if isError {
			t.Errorf("unexpected tool error: %s", result)
		}
		if !strings.Contains(result, ".ex") {
			t.Errorf("result %q does not contain .ex path", result)
		}
	})

	t.Run("stdlib rejected", func(t *testing.T) {
		_, isError, err := CallTool(testServerURL, "get_source_location", map[string]any{
			"reference": "Enum",
		})
		if err != nil {
			t.Fatalf("CallTool error: %v", err)
		}
		if !isError {
			t.Error("expected tool error for stdlib module")
		}
	})
}

func TestDocs(t *testing.T) {
	tests := []struct {
		name     string
		ref      string
		contains string
	}{
		{
			name:     "module docs",
			ref:      "Enum",
			contains: "enumerable",
		},
		{
			name:     "function with arity",
			ref:      "Enum.map/2",
			contains: "map",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, isError, err := CallTool(testServerURL, "get_docs", map[string]any{
				"reference": tt.ref,
			})
			if err != nil {
				t.Fatalf("CallTool error: %v", err)
			}
			if isError {
				t.Errorf("unexpected tool error: %s", result)
			}
			if !strings.Contains(strings.ToLower(result), tt.contains) {
				t.Errorf("result does not contain %q", tt.contains)
			}
		})
	}
}

func TestSQL(t *testing.T) {
	tests := []struct {
		name     string
		args     map[string]any
		contains string
	}{
		{
			name:     "basic query",
			args:     map[string]any{"query": "SELECT 1 AS test"},
			contains: "1",
		},
		{
			name:     "parameterized query",
			args:     map[string]any{"query": "SELECT $1::int AS val", "arguments": []any{42}},
			contains: "42",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, isError, err := CallTool(testServerURL, "execute_sql_query", tt.args)
			if err != nil {
				t.Fatalf("CallTool error: %v", err)
			}
			if isError {
				t.Errorf("unexpected tool error: %s", result)
			}
			if !strings.Contains(result, tt.contains) {
				t.Errorf("result %q does not contain %q", result, tt.contains)
			}
		})
	}
}

func TestSchemas(t *testing.T) {
	result, isError, err := CallTool(testServerURL, "get_ecto_schemas", map[string]any{})
	if err != nil {
		t.Fatalf("CallTool error: %v", err)
	}
	if isError {
		t.Errorf("unexpected tool error: %s", result)
	}
	if result == "" {
		t.Error("expected non-empty schema list")
	}
}

func TestSearch(t *testing.T) {
	tests := []struct {
		name     string
		args     map[string]any
		contains string
	}{
		{
			name:     "basic search",
			args:     map[string]any{"q": "genserver"},
			contains: "Results:",
		},
		{
			name:     "with packages filter",
			args:     map[string]any{"q": "plug", "packages": []string{"phoenix"}},
			contains: "Results:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, isError, err := CallTool(testServerURL, "search_package_docs", tt.args)
			if err != nil {
				t.Fatalf("CallTool error: %v", err)
			}
			if isError {
				t.Errorf("unexpected tool error: %s", result)
			}
			if !strings.Contains(result, tt.contains) {
				t.Errorf("result does not contain %q", tt.contains)
			}
		})
	}
}

func TestAsh(t *testing.T) {
	result, isError, err := CallTool(testServerURL, "get_ash_resources", map[string]any{})
	if err != nil {
		t.Fatalf("CallTool error: %v", err)
	}
	if isError {
		t.Errorf("unexpected tool error: %s", result)
	}
	// Result may be empty if the app doesn't use Ash, that's fine
}

func TestConnectionOverride(t *testing.T) {
	_, _, err := CallTool("http://localhost:9999/tidewave/mcp", "project_eval", map[string]any{
		"code": "1 + 1",
	})
	if err == nil {
		t.Error("expected connection error for port 9999, got nil")
	}
}

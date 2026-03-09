package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

func resolveServerURL(host, port, path string) string {
	if host == "" {
		host = os.Getenv("TIDEWAVE_HOST")
		if host == "" {
			host = "localhost"
		}
	}
	if port == "" {
		port = os.Getenv("TIDEWAVE_PORT")
		if port == "" {
			port = "4000"
		}
	}
	if path == "" {
		path = os.Getenv("TIDEWAVE_PATH")
		if path == "" {
			path = "/tidewave/mcp"
		}
	}
	return fmt.Sprintf("http://%s:%s%s", host, port, path)
}

func CallTool(serverURL, toolName string, arguments map[string]any) (string, bool, error) {
	reqBody := map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "tools/call",
		"params": map[string]any{
			"name":      toolName,
			"arguments": arguments,
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", false, fmt.Errorf("marshaling request: %w", err)
	}

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Post(serverURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return "", false, fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	var rpcResp struct {
		Error *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
		Result *struct {
			Content []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"content"`
			IsError bool `json:"isError"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
		return "", false, fmt.Errorf("decoding response: %w", err)
	}

	if rpcResp.Error != nil {
		return "", false, fmt.Errorf("rpc error %d: %s", rpcResp.Error.Code, rpcResp.Error.Message)
	}

	if rpcResp.Result == nil {
		return "", false, fmt.Errorf("response missing both result and error")
	}

	var texts []string
	for _, c := range rpcResp.Result.Content {
		if c.Text != "" {
			texts = append(texts, c.Text)
		}
	}

	return strings.Join(texts, "\n"), rpcResp.Result.IsError, nil
}

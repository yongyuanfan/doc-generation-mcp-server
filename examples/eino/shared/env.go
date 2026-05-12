package shared

import "os"

func MCPServerURL() string {
	if value := os.Getenv("MCP_SERVER_URL"); value != "" {
		return value
	}
	return "http://localhost:9103/mcp"
}

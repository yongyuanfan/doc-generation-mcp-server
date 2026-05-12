//go:build ignore

package main

import (
	"context"
	"log"
	"os"

	mcptool "github.com/cloudwego/eino-ext/components/tool/mcp"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	ctx := context.Background()
	transport := &mcp.StreamableClientTransport{
		Endpoint: envOrDefault("MCP_SERVER_URL", "http://localhost:9103/mcp"),
	}

	client := mcp.NewClient(&mcp.Implementation{Name: "eino-example-client", Version: "0.1.0"}, nil)
	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	tools, err := mcptool.GetTools(ctx, &mcptool.Config{Cli: &sessionClient{session: session}})
	if err != nil {
		log.Fatal(err)
	}

	for _, tool := range tools {
		log.Printf("discovered tool: %s", tool.Info(ctx).Name)
	}
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

type sessionClient struct {
	session *mcp.ClientSession
}

func (c *sessionClient) ListTools(ctx context.Context, params mcp.ListToolsParams) (*mcp.ListToolsResult, error) {
	return c.session.ListTools(ctx, &params)
}

func (c *sessionClient) CallTool(ctx context.Context, params mcp.CallToolParams) (*mcp.CallToolResult, error) {
	return c.session.CallTool(ctx, &params)
}

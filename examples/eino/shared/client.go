package shared

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Client struct {
	Session *mcp.ClientSession
}

func NewClient(ctx context.Context) (*Client, error) {
	transport := &mcp.StreamableClientTransport{Endpoint: MCPServerURL()}
	client := mcp.NewClient(&mcp.Implementation{Name: "eino-example-client", Version: "0.1.0"}, nil)
	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		return nil, err
	}
	return &Client{Session: session}, nil
}

func (c *Client) Close() error {
	if c == nil || c.Session == nil {
		return nil
	}
	return c.Session.Close()
}

func (c *Client) ListTools(ctx context.Context) ([]*mcp.Tool, error) {
	result, err := c.Session.ListTools(ctx, &mcp.ListToolsParams{})
	if err != nil {
		return nil, err
	}
	return result.Tools, nil
}

func (c *Client) CallToolJSON(ctx context.Context, name string, input any) (string, error) {
	result, err := c.Session.CallTool(ctx, &mcp.CallToolParams{Name: name, Arguments: input})
	if err != nil {
		return "", err
	}
	if result.IsError {
		return "", fmt.Errorf("%s", result.Content[0].(*mcp.TextContent).Text)
	}
	if result.StructuredContent != nil {
		data, err := json.MarshalIndent(result.StructuredContent, "", "  ")
		if err != nil {
			return "", err
		}
		return string(data), nil
	}
	if len(result.Content) == 0 {
		return "", nil
	}
	if text, ok := result.Content[0].(*mcp.TextContent); ok {
		return text.Text, nil
	}
	data, err := json.MarshalIndent(result.Content, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (c *Client) CallToolStructured(ctx context.Context, name string, input any) (map[string]any, error) {
	result, err := c.Session.CallTool(ctx, &mcp.CallToolParams{Name: name, Arguments: input})
	if err != nil {
		return nil, err
	}
	if result.IsError {
		if len(result.Content) > 0 {
			if text, ok := result.Content[0].(*mcp.TextContent); ok {
				return nil, fmt.Errorf("%s", text.Text)
			}
		}
		return nil, fmt.Errorf("tool call failed")
	}
	if result.StructuredContent == nil {
		return map[string]any{}, nil
	}
	if cast, ok := result.StructuredContent.(map[string]any); ok {
		return cast, nil
	}
	payload, err := json.Marshal(result.StructuredContent)
	if err != nil {
		return nil, err
	}
	var decoded map[string]any
	if err := json.Unmarshal(payload, &decoded); err != nil {
		return nil, err
	}
	return decoded, nil
}

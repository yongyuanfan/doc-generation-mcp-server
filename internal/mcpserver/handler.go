package mcpserver

import (
	"context"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/yong/doc-generation-mcp-server/internal/config"
	"github.com/yong/doc-generation-mcp-server/internal/model"
	docsvc "github.com/yong/doc-generation-mcp-server/internal/service/document"
)

func NewHandler(cfg config.Config, service *docsvc.Service) http.Handler {
	server := mcp.NewServer(&mcp.Implementation{Name: cfg.MCPServerName, Version: cfg.MCPServerVersion}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "generate_docx",
		Description: "Generate a DOCX document from structured content blocks.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input model.GenerateDocumentRequest) (*mcp.CallToolResult, model.DocumentResult, error) {
		result, err := service.Generate(ctx, input)
		return nil, result, err
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "render_docx_template",
		Description: "Render a DOCX file from a template in the templates directory.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input model.RenderTemplateRequest) (*mcp.CallToolResult, model.DocumentResult, error) {
		result, err := service.RenderTemplate(ctx, input)
		return nil, result, err
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_docx_templates",
		Description: "List available DOCX templates from the templates directory.",
	}, func(context.Context, *mcp.CallToolRequest, model.ListTemplatesRequest) (*mcp.CallToolResult, model.ListTemplatesResponse, error) {
		result, err := service.ListTemplates()
		return nil, result, err
	})

	return mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server { return server }, nil)
}

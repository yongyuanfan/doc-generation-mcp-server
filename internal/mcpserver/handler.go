package mcpserver

import (
	"context"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/yong/doc-generation-mcp-server/internal/config"
	"github.com/yong/doc-generation-mcp-server/internal/formaldoc"
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
		Name:        "generate_docx_from_draft",
		Description: "Generate a DOCX document from a FormalDocumentDraftV1 payload.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input formaldoc.Draft) (*mcp.CallToolResult, model.DraftDocumentResult, error) {
		result, err := service.GenerateFromDraft(ctx, input)
		return nil, result, err
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "validate_formal_document_draft",
		Description: "Validate a FormalDocumentDraftV1 payload without generating a DOCX file.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input formaldoc.Draft) (*mcp.CallToolResult, model.DraftValidationResult, error) {
		_ = ctx
		_ = req
		return nil, service.ValidateDraft(input), nil
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

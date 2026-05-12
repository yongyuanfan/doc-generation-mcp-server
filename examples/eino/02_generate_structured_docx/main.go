//go:build ignore

package main

import (
	"context"
	"log"

	"github.com/yong/doc-generation-mcp-server/examples/eino/shared"
	"github.com/yong/doc-generation-mcp-server/internal/model"
)

func main() {
	ctx := context.Background()
	client, err := shared.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	result, err := client.CallToolJSON(ctx, "generate_docx", model.GenerateDocumentRequest{
		FileName:         "eino-structured.docx",
		Title:            "Eino Structured Example",
		HeaderText:       "Eino Structured Example",
		FooterPageNumber: true,
		Content: []model.ContentBlock{
			{Type: "toc", Text: "Contents", Levels: "1-2"},
			{Type: "heading", Text: "Summary", Level: 1},
			{Type: "paragraph", Text: "This document is generated through the remote MCP tool call path."},
			{Type: "hyperlink", URL: "https://modelcontextprotocol.io", DisplayText: "MCP"},
			{Type: "table", Rows: [][]string{{"Task", "Status"}, {"Structured example", "Done"}}},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	shared.PrintJSON("generate_docx", result)
}

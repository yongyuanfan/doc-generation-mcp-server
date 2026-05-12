//go:build ignore

package main

import (
	"context"
	"log"

	"github.com/yong/doc-generation-mcp-server/examples/eino/shared"
)

func main() {
	ctx := context.Background()
	client, err := shared.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	draft := shared.WeeklyReportDraft()

	validation, err := client.CallToolJSON(ctx, "validate_formal_document_draft", draft)
	if err != nil {
		log.Fatal(err)
	}
	shared.PrintJSON("validate_formal_document_draft", validation)

	generated, err := client.CallToolJSON(ctx, "generate_docx_from_draft", draft)
	if err != nil {
		log.Fatal(err)
	}
	shared.PrintJSON("generate_docx_from_draft", generated)
}

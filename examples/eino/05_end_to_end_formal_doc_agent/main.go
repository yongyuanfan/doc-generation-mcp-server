//go:build ignore

package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/yong/doc-generation-mcp-server/examples/eino/shared"
	"github.com/yong/doc-generation-mcp-server/internal/formaldoc"
)

type UserInput struct {
	DocumentType string
	Goal         string
	Audience     string
	Facts        []string
}

func main() {
	ctx := context.Background()
	client, err := shared.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	input := UserInput{
		DocumentType: formaldoc.DocumentTypeWeeklyReport,
		Goal:         "提交管理层周报",
		Audience:     "management",
		Facts: []string{
			"已完成核心接口改造",
			"模板库仍需补充",
			"继续扩展标准文档类型",
		},
	}

	draft := buildDraftFromInput(input)
	shared.PrintJSON("mock_llm_draft", mustDraftJSON(draft))

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

func buildDraftFromInput(input UserInput) formaldoc.Draft {
	if input.DocumentType == formaldoc.DocumentTypeBusinessLetter {
		return shared.BusinessLetterDraft()
	}
	return shared.WeeklyReportDraft()
}

func mustDraftJSON(draft formaldoc.Draft) string {
	data, err := json.MarshalIndent(draft, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}

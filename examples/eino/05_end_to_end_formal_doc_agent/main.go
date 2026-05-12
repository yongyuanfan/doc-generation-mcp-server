//go:build ignore

package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/yong/doc-generation-mcp-server/examples/eino/shared"
	"github.com/yong/doc-generation-mcp-server/internal/formaldoc"
)

type UserInput struct {
	DocumentType string   `json:"document_type"`
	Goal         string   `json:"goal"`
	Audience     string   `json:"audience"`
	Facts        []string `json:"facts"`
}

func main() {
	documentType := flag.String("document-type", formaldoc.DocumentTypeWeeklyReport, "document type: weekly_report or business_letter")
	inputPath := flag.String("input", "", "optional path to a JSON input file")
	flag.Parse()

	ctx := context.Background()
	client, err := shared.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	input, err := loadUserInput(*documentType, *inputPath)
	if err != nil {
		log.Fatal(err)
	}

	llmConfig := shared.LoadLLMConfig()
	var draft formaldoc.Draft
	if llmConfig.Mode == "openai" {
		var raw string
		draft, raw, err = shared.BuildDraftWithLLM(ctx, llmConfig, input)
		if err != nil {
			log.Fatal(err)
		}
		shared.PrintJSON("openai_draft_json", raw)
	} else {
		draft = buildDraftFromInput(input)
		shared.PrintJSON("mock_llm_draft", mustDraftJSON(draft))
	}

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

func loadUserInput(documentType, inputPath string) (UserInput, error) {
	if inputPath != "" {
		data, err := os.ReadFile(inputPath)
		if err != nil {
			return UserInput{}, err
		}
		var input UserInput
		if err := json.Unmarshal(data, &input); err != nil {
			return UserInput{}, err
		}
		return input, nil
	}
	if documentType == formaldoc.DocumentTypeBusinessLetter {
		return UserInput{
			DocumentType: formaldoc.DocumentTypeBusinessLetter,
			Goal:         "向客户发送正式沟通说明函",
			Audience:     "customer",
			Facts: []string{
				"建议双方安排专项沟通会议",
				"会议议题包括需求确认、排期安排与责任分工",
				"请于收到本函后确认可行时间",
			},
		}, nil
	}
	return UserInput{
		DocumentType: formaldoc.DocumentTypeWeeklyReport,
		Goal:         "提交管理层周报",
		Audience:     "management",
		Facts: []string{
			"已完成核心接口改造",
			"模板库仍需补充",
			"继续扩展标准文档类型",
		},
	}, nil
}

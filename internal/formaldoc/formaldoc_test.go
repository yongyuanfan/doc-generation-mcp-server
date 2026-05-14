package formaldoc

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSchemaFileIsValidJSON(t *testing.T) {
	path := filepath.Join("..", "..", "schemas", "formal-document-draft-v1.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded["title"] != "FormalDocumentDraftV1" {
		t.Fatalf("unexpected schema title: %v", decoded["title"])
	}
}

func TestValidateDraftRejectsMissingProjectSections(t *testing.T) {
	issues := ValidateDraft(Draft{
		SchemaVersion: SchemaVersion,
		DocumentType:  DocumentTypeProjectProposal,
		Title:         "方案",
		Audience:      "management",
		Tone:          "formal",
		Language:      "zh-CN",
		IncludeTOC:    true,
		Sections: []Section{
			{Title: "一、项目背景", Level: 1, Blocks: []Block{{Type: "paragraph", Text: "内容"}}},
		},
	})
	if len(issues) == 0 {
		t.Fatal("expected validation issues")
	}
	joined := joinIssues(issues)
	if !strings.Contains(joined, "建设目标") {
		t.Fatalf("expected missing section issue, got %s", joined)
	}
	if !strings.Contains(joined, "header_text") {
		t.Fatalf("expected header_text issue, got %s", joined)
	}
}

func TestValidateDraftRejectsBusinessLetterTOC(t *testing.T) {
	issues := ValidateDraft(validBusinessLetterDraft(func(d *Draft) {
		d.IncludeTOC = true
	}))
	if len(issues) == 0 {
		t.Fatal("expected validation issue")
	}
	if !strings.Contains(joinIssues(issues), "include_toc") {
		t.Fatalf("expected include_toc issue, got %s", joinIssues(issues))
	}
}

func TestValidateDraftRejectsBannedPhrase(t *testing.T) {
	issues := ValidateDraft(validWeeklyReportDraft(func(d *Draft) {
		d.Sections[0].Blocks[0].Text = "我觉得本期推进情况良好。"
	}))
	if len(issues) == 0 {
		t.Fatal("expected style issue")
	}
	if !strings.Contains(joinIssues(issues), "我觉得") {
		t.Fatalf("expected banned phrase issue, got %s", joinIssues(issues))
	}
}

func TestToGenerateRequestBuildsExpectedBlocks(t *testing.T) {
	result, err := ToGenerateRequest(validProjectProposalDraft())
	if err != nil {
		t.Fatal(err)
	}
	if result.Request.HeaderText != "项目方案" {
		t.Fatalf("unexpected header text: %s", result.Request.HeaderText)
	}
	if len(result.Request.Content) < 5 {
		t.Fatalf("expected content blocks, got %d", len(result.Request.Content))
	}
	if result.Request.Content[0].Type != "toc" {
		t.Fatalf("expected first block to be toc, got %s", result.Request.Content[0].Type)
	}
	last := result.Request.Content[len(result.Request.Content)-1]
	if last.Type != "hyperlink" {
		t.Fatalf("expected reference mapped to hyperlink, got %s", last.Type)
	}
	if result.TemplateName != "project-proposal.docx" {
		t.Fatalf("unexpected template name: %s", result.TemplateName)
	}
}

func TestToTemplateRequestFillsCommonBusinessLetterPlaceholders(t *testing.T) {
	request, err := ToTemplateRequest(validBusinessLetterDraft(func(d *Draft) {
		d.TemplateName = "business-letter.docx"
		d.Organization = "某某科技有限公司"
		d.Summary = "关于工厂停产安排的通知"
	}))
	if err != nil {
		t.Fatal(err)
	}
	if got := request.Data["organization"]; got != "某某科技有限公司" {
		t.Fatalf("unexpected organization: %v", got)
	}
	if got := request.Data["sender_name"]; got != "某某科技有限公司" {
		t.Fatalf("unexpected sender_name: %v", got)
	}
	if got := request.Data["summary"]; got != "关于工厂停产安排的通知" {
		t.Fatalf("unexpected summary: %v", got)
	}
	if got := request.Data["title"]; got != "说明函" {
		t.Fatalf("unexpected title: %v", got)
	}
}

func validProjectProposalDraft() Draft {
	return Draft{
		SchemaVersion:    SchemaVersion,
		DocumentType:     DocumentTypeProjectProposal,
		Title:            "项目方案",
		Author:           "作者",
		Audience:         "management",
		Tone:             "formal",
		Language:         "zh-CN",
		HeaderText:       "项目方案",
		FooterPageNumber: true,
		IncludeTOC:       true,
		TemplateName:     "project-proposal.docx",
		Sections: []Section{
			{Title: "一、项目背景", Level: 1, Blocks: []Block{{Type: "paragraph", Text: "背景说明"}}},
			{Title: "二、建设目标", Level: 1, Blocks: []Block{{Type: "table", Rows: [][]string{{"项目", "说明"}, {"目标", "统一标准"}}}}},
			{Title: "三、建设内容", Level: 1, Blocks: []Block{{Type: "paragraph", Text: "建设内容说明"}}},
			{Title: "四、实施计划", Level: 1, Blocks: []Block{{Type: "paragraph", Text: "实施计划说明"}}},
			{Title: "五、资源需求", Level: 1, Blocks: []Block{{Type: "paragraph", Text: "资源需求待补充"}}},
			{Title: "六、风险与保障措施", Level: 1, Blocks: []Block{{Type: "paragraph", Text: "风险说明"}}},
			{Title: "七、预期成效", Level: 1, Blocks: []Block{{Type: "paragraph", Text: "成效说明"}}},
			{Title: "八、结论与建议", Level: 1, Blocks: []Block{{Type: "paragraph", Text: "结论说明"}}},
		},
		References:  []Reference{{Title: "参考资料", URL: "https://example.com/spec"}},
		ReviewNotes: []string{"预算待补充"},
	}
}

func validWeeklyReportDraft(mutator func(*Draft)) Draft {
	d := Draft{
		SchemaVersion:    SchemaVersion,
		DocumentType:     DocumentTypeWeeklyReport,
		Title:            "周报",
		Audience:         "management",
		Tone:             "formal",
		Language:         "zh-CN",
		FooterPageNumber: true,
		Sections: []Section{
			{Title: "一、本期工作概述", Level: 1, Blocks: []Block{{Type: "paragraph", Text: "概述"}}},
			{Title: "二、已完成事项", Level: 1, Blocks: []Block{{Type: "paragraph", Text: "完成事项"}}},
			{Title: "三、当前进展", Level: 1, Blocks: []Block{{Type: "paragraph", Text: "进展"}}},
			{Title: "四、存在问题", Level: 1, Blocks: []Block{{Type: "paragraph", Text: "问题"}}},
			{Title: "五、下阶段计划", Level: 1, Blocks: []Block{{Type: "paragraph", Text: "计划"}}},
			{Title: "六、需协调事项", Level: 1, Blocks: []Block{{Type: "paragraph", Text: "协调事项"}}},
		},
	}
	if mutator != nil {
		mutator(&d)
	}
	return d
}

func validBusinessLetterDraft(mutator func(*Draft)) Draft {
	d := Draft{
		SchemaVersion: SchemaVersion,
		DocumentType:  DocumentTypeBusinessLetter,
		Title:         "说明函",
		Audience:      "customer",
		Tone:          "formal",
		Language:      "zh-CN",
		Sections: []Section{
			{Title: "一、发函背景", Level: 1, Blocks: []Block{{Type: "paragraph", Text: "背景"}}},
			{Title: "二、发函事项", Level: 1, Blocks: []Block{{Type: "paragraph", Text: "事项"}}},
			{Title: "三、具体说明", Level: 1, Blocks: []Block{{Type: "paragraph", Text: "说明"}}},
			{Title: "四、后续安排", Level: 1, Blocks: []Block{{Type: "paragraph", Text: "安排"}}},
		},
	}
	if mutator != nil {
		mutator(&d)
	}
	return d
}

func joinIssues(issues []ValidationIssue) string {
	parts := make([]string, 0, len(issues))
	for _, issue := range issues {
		parts = append(parts, issue.Error())
	}
	return strings.Join(parts, " | ")
}

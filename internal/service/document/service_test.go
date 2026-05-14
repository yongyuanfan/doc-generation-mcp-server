package document

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yong/doc-generation-mcp-server/internal/config"
	"github.com/yong/doc-generation-mcp-server/internal/formaldoc"
	"github.com/yong/doc-generation-mcp-server/internal/model"
)

type stubProvider struct{}

func (stubProvider) Generate(_ context.Context, input model.GenerateDocumentRequest, outputName string) (model.DocumentResult, error) {
	return model.DocumentResult{FileName: outputName, Path: outputName}, nil
}

func (stubProvider) RenderTemplate(_ context.Context, templatePath string, _ map[string]any, outputName string) (model.DocumentResult, error) {
	return model.DocumentResult{FileName: outputName, Path: templatePath}, nil
}

func TestGenerateRejectsInvalidBlocks(t *testing.T) {
	tempDir := t.TempDir()
	service := NewService(testConfig(tempDir), stubProvider{}, nil)

	_, err := service.Generate(context.Background(), model.GenerateDocumentRequest{
		Content: []model.ContentBlock{{Type: "table", Rows: [][]string{{"a"}, {"a", "b"}}}},
	})
	if err == nil || !strings.Contains(err.Error(), "same column count") {
		t.Fatalf("expected table validation error, got %v", err)
	}
}

func TestGenerateRejectsHyperlinkWithoutURL(t *testing.T) {
	tempDir := t.TempDir()
	service := NewService(testConfig(tempDir), stubProvider{}, nil)

	_, err := service.Generate(context.Background(), model.GenerateDocumentRequest{
		Content: []model.ContentBlock{{Type: "hyperlink", DisplayText: "Open"}},
	})
	if err == nil || !strings.Contains(err.Error(), "url is required") {
		t.Fatalf("expected hyperlink validation error, got %v", err)
	}
}

func TestGenerateRejectsImageWithoutSource(t *testing.T) {
	tempDir := t.TempDir()
	service := NewService(testConfig(tempDir), stubProvider{}, nil)

	_, err := service.Generate(context.Background(), model.GenerateDocumentRequest{
		Content: []model.ContentBlock{{Type: "image"}},
	})
	if err == nil || !strings.Contains(err.Error(), "image_base64 or url is required") {
		t.Fatalf("expected image source validation error, got %v", err)
	}
}

func TestCapabilitiesIncludeExtendedBlocks(t *testing.T) {
	service := NewService(testConfig(t.TempDir()), stubProvider{}, nil)
	capabilities := service.Capabilities()
	if !capabilities.FooterPageNumber {
		t.Fatal("expected footer page number capability")
	}
	if !capabilities.HeaderText {
		t.Fatal("expected header text capability")
	}
	joined := strings.Join(capabilities.BlockTypes, ",")
	for _, blockType := range []string{"page_break", "hyperlink", "toc"} {
		if !strings.Contains(joined, blockType) {
			t.Fatalf("expected block type %s in %v", blockType, capabilities.BlockTypes)
		}
	}
}

func TestRenderTemplateRejectsMissingTemplate(t *testing.T) {
	tempDir := t.TempDir()
	service := NewService(testConfig(tempDir), stubProvider{}, nil)

	_, err := service.RenderTemplate(context.Background(), model.RenderTemplateRequest{
		TemplateName: "missing.docx",
	})
	if err == nil || !strings.Contains(err.Error(), "template not found") {
		t.Fatalf("expected missing template error, got %v", err)
	}
}

func TestListTemplatesFiltersNonDocx(t *testing.T) {
	tempDir := t.TempDir()
	cfg := testConfig(tempDir)
	service := NewService(cfg, stubProvider{}, nil)

	if err := os.MkdirAll(cfg.DocxTemplateDir, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"a.docx", "b.txt", "c.DOCX"} {
		if err := os.WriteFile(filepath.Join(cfg.DocxTemplateDir, name), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	result, err := service.ListTemplates()
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Templates) != 2 {
		t.Fatalf("expected 2 templates, got %d", len(result.Templates))
	}
}

func TestDownloadPathRejectsInvalidName(t *testing.T) {
	tempDir := t.TempDir()
	service := NewService(testConfig(tempDir), stubProvider{}, nil)

	_, err := service.DownloadPath("../secret.docx")
	if err == nil || !strings.Contains(err.Error(), "invalid file name") {
		t.Fatalf("expected invalid file name error, got %v", err)
	}
}

func TestDownloadPathFindsExistingFile(t *testing.T) {
	tempDir := t.TempDir()
	cfg := testConfig(tempDir)
	service := NewService(cfg, stubProvider{}, nil)

	if err := os.MkdirAll(cfg.DocxTempDir, 0o755); err != nil {
		t.Fatal(err)
	}
	fullPath := filepath.Join(cfg.DocxTempDir, "ready.docx")
	if err := os.WriteFile(fullPath, []byte("docx"), 0o644); err != nil {
		t.Fatal(err)
	}

	path, err := service.DownloadPath("ready.docx")
	if err != nil {
		t.Fatal(err)
	}
	if path != fullPath {
		t.Fatalf("expected %s, got %s", fullPath, path)
	}
}

func TestGenerateFromDraftReturnsReviewNotes(t *testing.T) {
	tempDir := t.TempDir()
	service := NewService(testConfig(tempDir), stubProvider{}, nil)

	result, err := service.GenerateFromDraft(context.Background(), formaldoc.Draft{
		SchemaVersion:    formaldoc.SchemaVersion,
		DocumentType:     formaldoc.DocumentTypeWeeklyReport,
		Title:            "周报",
		Audience:         "management",
		Tone:             "formal",
		Language:         "zh-CN",
		FooterPageNumber: true,
		Sections: []formaldoc.Section{
			{Title: "一、本期工作概述", Level: 1, Blocks: []formaldoc.Block{{Type: "paragraph", Text: "概述"}}},
			{Title: "二、已完成事项", Level: 1, Blocks: []formaldoc.Block{{Type: "paragraph", Text: "已完成"}}},
			{Title: "三、当前进展", Level: 1, Blocks: []formaldoc.Block{{Type: "paragraph", Text: "进展"}}},
			{Title: "四、存在问题", Level: 1, Blocks: []formaldoc.Block{{Type: "paragraph", Text: "问题"}}},
			{Title: "五、下阶段计划", Level: 1, Blocks: []formaldoc.Block{{Type: "paragraph", Text: "计划"}}},
			{Title: "六、需协调事项", Level: 1, Blocks: []formaldoc.Block{{Type: "paragraph", Text: "协调事项"}}},
		},
		ReviewNotes: []string{"资料待补充"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.ReviewNotes) != 1 || result.ReviewNotes[0] != "资料待补充" {
		t.Fatalf("unexpected review notes: %#v", result.ReviewNotes)
	}
	if result.FileName == "" {
		t.Fatal("expected generated file name")
	}
	if result.Route != "structured" {
		t.Fatalf("expected structured route, got %s", result.Route)
	}
}

func TestValidateDraftReturnsIssues(t *testing.T) {
	service := NewService(testConfig(t.TempDir()), stubProvider{}, nil)
	result := service.ValidateDraft(formaldoc.Draft{})
	if result.Valid {
		t.Fatal("expected invalid result")
	}
	if len(result.Issues) == 0 {
		t.Fatal("expected validation issues")
	}
}

func TestValidateDraftReturnsRecommendations(t *testing.T) {
	service := NewService(testConfig(t.TempDir()), stubProvider{}, nil)
	result := service.ValidateDraft(formaldoc.Draft{
		SchemaVersion: formaldoc.SchemaVersion,
		DocumentType:  formaldoc.DocumentTypeBusinessLetter,
		Title:         "说明函",
		Audience:      "customer",
		Tone:          "formal",
		Language:      "zh-CN",
		Sections: []formaldoc.Section{
			{Title: "一、发函背景", Level: 1, Blocks: []formaldoc.Block{{Type: "paragraph", Text: "背景"}}},
			{Title: "二、发函事项", Level: 1, Blocks: []formaldoc.Block{{Type: "paragraph", Text: "事项"}}},
			{Title: "三、具体说明", Level: 1, Blocks: []formaldoc.Block{{Type: "paragraph", Text: "说明"}}},
			{Title: "四、后续安排", Level: 1, Blocks: []formaldoc.Block{{Type: "paragraph", Text: "安排"}}},
		},
	})
	if result.RecommendedRoute != formaldoc.RouteTemplate {
		t.Fatalf("expected template recommendation, got %s", result.RecommendedRoute)
	}
	if result.RecommendedTemplate != "business-letter.docx" {
		t.Fatalf("unexpected template recommendation: %s", result.RecommendedTemplate)
	}
}

func TestValidateDraftUsesConfiguredTemplateMapping(t *testing.T) {
	cfg := testConfig(t.TempDir())
	cfg.DocumentTypeTemplateMap = map[string]string{"business_letter": "custom-letter.docx"}
	service := NewService(cfg, stubProvider{}, nil)
	result := service.ValidateDraft(formaldoc.Draft{
		SchemaVersion: formaldoc.SchemaVersion,
		DocumentType:  formaldoc.DocumentTypeBusinessLetter,
		Title:         "说明函",
		Audience:      "customer",
		Tone:          "formal",
		Language:      "zh-CN",
		Sections: []formaldoc.Section{
			{Title: "一、发函背景", Level: 1, Blocks: []formaldoc.Block{{Type: "paragraph", Text: "背景"}}},
			{Title: "二、发函事项", Level: 1, Blocks: []formaldoc.Block{{Type: "paragraph", Text: "事项"}}},
			{Title: "三、具体说明", Level: 1, Blocks: []formaldoc.Block{{Type: "paragraph", Text: "说明"}}},
			{Title: "四、后续安排", Level: 1, Blocks: []formaldoc.Block{{Type: "paragraph", Text: "安排"}}},
		},
	})
	if result.RecommendedTemplate != "custom-letter.docx" {
		t.Fatalf("unexpected template recommendation: %s", result.RecommendedTemplate)
	}
}

func TestGenerateFromDraftUsesTemplateRoute(t *testing.T) {
	tempDir := t.TempDir()
	cfg := testConfig(tempDir)
	service := NewService(cfg, stubProvider{}, nil)
	if err := os.MkdirAll(cfg.DocxTemplateDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cfg.DocxTemplateDir, "business-letter.docx"), []byte("template"), 0o644); err != nil {
		t.Fatal(err)
	}

	result, err := service.GenerateFromDraft(context.Background(), formaldoc.Draft{
		SchemaVersion: formaldoc.SchemaVersion,
		DocumentType:  formaldoc.DocumentTypeBusinessLetter,
		Title:         "说明函",
		Audience:      "customer",
		Tone:          "formal",
		Language:      "zh-CN",
		TemplateName:  "business-letter.docx",
		Sections: []formaldoc.Section{
			{Title: "一、发函背景", Level: 1, Blocks: []formaldoc.Block{{Type: "paragraph", Text: "背景"}}},
			{Title: "二、发函事项", Level: 1, Blocks: []formaldoc.Block{{Type: "paragraph", Text: "事项"}}},
			{Title: "三、具体说明", Level: 1, Blocks: []formaldoc.Block{{Type: "paragraph", Text: "说明"}}},
			{Title: "四、后续安排", Level: 1, Blocks: []formaldoc.Block{{Type: "paragraph", Text: "安排"}}},
		},
		Placeholders: map[string]string{"recipient_name": "张三"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Route != "template" {
		t.Fatalf("expected template route, got %s", result.Route)
	}
	if result.TemplateName != "business-letter.docx" {
		t.Fatalf("unexpected template name: %s", result.TemplateName)
	}
}

func testConfig(root string) config.Config {
	return config.Config{
		APIPrefix:           "/api/v1",
		DocxTempDir:         filepath.Join(root, "temp"),
		DocxTemplateDir:     filepath.Join(root, "templates"),
		DocxDefaultAuthor:   "tester",
		DocxDefaultFont:     "Calibri",
		DocxDefaultFontSize: 22,
		DocxMaxFileAge:      0,
		DocumentTypeTemplateMap: map[string]string{
			"business_letter": "business-letter.docx",
		},
	}
}

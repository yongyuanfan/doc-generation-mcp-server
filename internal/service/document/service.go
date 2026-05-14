package document

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/yong/doc-generation-mcp-server/internal/config"
	"github.com/yong/doc-generation-mcp-server/internal/formaldoc"
	"github.com/yong/doc-generation-mcp-server/internal/model"
	"github.com/yong/doc-generation-mcp-server/internal/storage"
)

type provider interface {
	Generate(context.Context, model.GenerateDocumentRequest, string) (model.DocumentResult, error)
	RenderTemplate(context.Context, string, map[string]any, string) (model.DocumentResult, error)
}

type Service struct {
	config   config.Config
	provider provider
	uploader storage.Uploader
}

func NewService(cfg config.Config, provider provider, uploader storage.Uploader) *Service {
	return &Service{config: cfg, provider: provider, uploader: uploader}
}

func (s *Service) Generate(ctx context.Context, input model.GenerateDocumentRequest) (model.DocumentResult, error) {
	if err := s.ensureRuntimeDirs(); err != nil {
		return model.DocumentResult{}, err
	}
	s.cleanupExpiredFiles()
	input = s.normalizeGenerateRequest(input)
	if err := validateGenerateRequest(input); err != nil {
		return model.DocumentResult{}, err
	}
	outputName := uniqueDocxName(input.FileName)
	result, err := s.provider.Generate(ctx, input, outputName)
	if err != nil {
		return model.DocumentResult{}, err
	}
	result.DownloadURL = s.config.APIPrefix + "/documents/files/" + result.FileName
	if err := s.uploadToMinIO(ctx, &result); err != nil {
		log.Printf("upload to minio failed for file=%s: %v", result.FileName, err)
	}
	return result, nil
}

func (s *Service) RenderTemplate(ctx context.Context, input model.RenderTemplateRequest) (model.DocumentResult, error) {
	if err := s.ensureRuntimeDirs(); err != nil {
		return model.DocumentResult{}, err
	}
	s.cleanupExpiredFiles()
	input = normalizeTemplateRequest(input)
	if err := validateTemplateRequest(input); err != nil {
		return model.DocumentResult{}, err
	}
	templatePath, err := s.templatePath(input.TemplateName)
	if err != nil {
		return model.DocumentResult{}, err
	}
	outputName := uniqueDocxName(input.FileName)
	result, err := s.provider.RenderTemplate(ctx, templatePath, input.Data, outputName)
	if err != nil {
		return model.DocumentResult{}, err
	}
	result.DownloadURL = s.config.APIPrefix + "/documents/files/" + result.FileName
	if err := s.uploadToMinIO(ctx, &result); err != nil {
		log.Printf("upload to minio failed for file=%s: %v", result.FileName, err)
	}
	return result, nil
}

func (s *Service) GenerateFromDraft(ctx context.Context, draft formaldoc.Draft) (model.DraftDocumentResult, error) {
	validation := s.ValidateDraft(draft)
	if !validation.Valid {
		return model.DraftDocumentResult{}, fmt.Errorf("%s", strings.Join(validation.Issues, "; "))
	}
	if strings.TrimSpace(draft.TemplateName) == "" {
		draft.TemplateName = formaldoc.RecommendedTemplate(draft, s.config.DocumentTypeTemplateMap)
	}
	if strings.TrimSpace(draft.TemplateName) != "" && formaldoc.RecommendedRoute(draft, s.config.DocumentTypeTemplateMap) == formaldoc.RouteTemplate {
		templateRequest, err := formaldoc.ToTemplateRequest(draft)
		if err != nil {
			return model.DraftDocumentResult{}, err
		}
		result, err := s.RenderTemplate(ctx, templateRequest)
		if err == nil {
			return model.DraftDocumentResult{
				FileName:     result.FileName,
				Path:         result.Path,
				DownloadURL:  result.DownloadURL,
				MIMEType:     result.MIMEType,
				SizeBytes:    result.SizeBytes,
				ReviewNotes:  validation.ReviewNotes,
				TemplateName: draft.TemplateName,
				Route:        formaldoc.RouteTemplate,
			}, nil
		}
		if !shouldFallbackToStructured(err) {
			return model.DraftDocumentResult{}, err
		}
	}
	converted, err := formaldoc.ToGenerateRequest(draft)
	if err != nil {
		return model.DraftDocumentResult{}, err
	}
	if strings.TrimSpace(draft.TemplateName) != "" && formaldoc.RecommendedRoute(draft, s.config.DocumentTypeTemplateMap) == formaldoc.RouteTemplate {
		converted.ReviewNotes = append(converted.ReviewNotes, fmt.Sprintf("模板渲染失败，已回退为结构化生成：%s", draft.TemplateName))
	}
	result, err := s.Generate(ctx, converted.Request)
	if err != nil {
		return model.DraftDocumentResult{}, err
	}
	return model.DraftDocumentResult{
		FileName:     result.FileName,
		Path:         result.Path,
		DownloadURL:  result.DownloadURL,
		MIMEType:     result.MIMEType,
		SizeBytes:    result.SizeBytes,
		ReviewNotes:  converted.ReviewNotes,
		TemplateName: converted.TemplateName,
		Route:        formaldoc.RouteStructured,
	}, nil
}

func (s *Service) ValidateDraft(draft formaldoc.Draft) model.DraftValidationResult {
	issues := formaldoc.ValidateDraft(draft)
	result := model.DraftValidationResult{
		Valid:               len(issues) == 0,
		ReviewNotes:         append([]string(nil), draft.ReviewNotes...),
		RecommendedRoute:    formaldoc.RecommendedRoute(draft, s.config.DocumentTypeTemplateMap),
		RecommendedTemplate: formaldoc.RecommendedTemplate(draft, s.config.DocumentTypeTemplateMap),
	}
	if len(issues) == 0 {
		return result
	}
	result.Issues = make([]string, 0, len(issues))
	for _, issue := range issues {
		result.Issues = append(result.Issues, issue.Error())
	}
	return result
}

func (s *Service) Capabilities() model.CapabilitiesResponse {
	return model.CapabilitiesResponse{
		Formats:          []string{"docx"},
		BlockTypes:       []string{"heading", "paragraph", "table", "image", "page_break", "hyperlink", "toc"},
		TemplateDir:      s.config.DocxTemplateDir,
		TempDir:          s.config.DocxTempDir,
		TemplateRender:   true,
		HeaderText:       true,
		FooterPageNumber: true,
	}
}

func (s *Service) ListTemplates() (model.ListTemplatesResponse, error) {
	if err := s.ensureRuntimeDirs(); err != nil {
		return model.ListTemplatesResponse{}, err
	}
	entries, err := os.ReadDir(s.config.DocxTemplateDir)
	if err != nil {
		return model.ListTemplatesResponse{}, err
	}
	templates := make([]model.TemplateInfo, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(strings.ToLower(entry.Name()), ".docx") {
			continue
		}
		templates = append(templates, model.TemplateInfo{
			Name: entry.Name(),
			Path: filepath.Join(s.config.DocxTemplateDir, entry.Name()),
		})
	}
	return model.ListTemplatesResponse{Templates: templates}, nil
}

func (s *Service) DownloadPath(name string) (string, error) {
	cleanName := sanitizeFileName(name)
	if cleanName == "" || cleanName != name {
		return "", fmt.Errorf("invalid file name")
	}
	path := filepath.Join(s.config.DocxTempDir, cleanName)
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("file not found")
		}
		return "", err
	}
	return path, nil
}

func (s *Service) ensureRuntimeDirs() error {
	if err := os.MkdirAll(s.config.DocxTempDir, 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll(s.config.DocxTemplateDir, 0o755); err != nil {
		return err
	}
	return nil
}

func (s *Service) templatePath(name string) (string, error) {
	cleanName := sanitizeFileName(name)
	if cleanName == "" || cleanName != name {
		return "", fmt.Errorf("invalid template name")
	}
	path := filepath.Join(s.config.DocxTemplateDir, cleanName)
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("template not found")
		}
		return "", err
	}
	return path, nil
}

func (s *Service) cleanupExpiredFiles() {
	entries, err := os.ReadDir(s.config.DocxTempDir)
	if err != nil {
		return
	}
	cutoff := time.Now().Add(-s.config.DocxMaxFileAge)
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".docx") {
			continue
		}
		info, err := entry.Info()
		if err != nil || info.ModTime().After(cutoff) {
			continue
		}
		_ = os.Remove(filepath.Join(s.config.DocxTempDir, entry.Name()))
	}
}

const docxMIMEType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"

func (s *Service) uploadToMinIO(ctx context.Context, result *model.DocumentResult) error {
	if s.uploader == nil {
		return nil
	}

	data, err := os.ReadFile(result.Path)
	if err != nil {
		return fmt.Errorf("read file for upload: %w", err)
	}

	objectName := buildDocxObjectName(s.config.MinIOObjectPrefix, result.FileName, time.Now())

	uploadedURL, err := s.uploader.Upload(ctx, objectName, docxMIMEType, data)
	if err != nil {
		return err
	}

	result.DownloadURL = uploadedURL
	log.Printf("uploaded file to minio: %s -> %s", result.FileName, uploadedURL)
	return nil
}

func buildDocxObjectName(prefix, fileName string, now time.Time) string {
	datePath := now.UTC().Format("2006/01/02")
	if prefix == "" {
		return path.Join(datePath, fileName)
	}
	return path.Join(prefix, datePath, fileName)
}

func (s *Service) normalizeGenerateRequest(input model.GenerateDocumentRequest) model.GenerateDocumentRequest {
	input.FileName = defaultFileName(input.FileName, "document.docx")
	input.Title = strings.TrimSpace(input.Title)
	input.Author = strings.TrimSpace(input.Author)
	if input.Author == "" {
		input.Author = s.config.DocxDefaultAuthor
	}
	input.HeaderText = strings.TrimSpace(input.HeaderText)
	input.Subject = strings.TrimSpace(input.Subject)
	for index := range input.Content {
		input.Content[index].Type = strings.ToLower(strings.TrimSpace(input.Content[index].Type))
		input.Content[index].Text = strings.TrimSpace(input.Content[index].Text)
		input.Content[index].Alignment = strings.ToLower(strings.TrimSpace(input.Content[index].Alignment))
		input.Content[index].URL = strings.TrimSpace(input.Content[index].URL)
		input.Content[index].DisplayText = strings.TrimSpace(input.Content[index].DisplayText)
		input.Content[index].Levels = strings.TrimSpace(input.Content[index].Levels)
	}
	// Filter out empty paragraph blocks: the auto-generated JSON Schema from
	// model.GenerateDocumentRequest marks text and runs as optional (omitempty),
	// so LLM-generated calls may include paragraphs without text or runs. We
	// silently skip these instead of rejecting with an error.
	input.Content = filterEmptyParagraphs(input.Content)
	return input
}

func normalizeTemplateRequest(input model.RenderTemplateRequest) model.RenderTemplateRequest {
	input.TemplateName = sanitizeFileName(input.TemplateName)
	input.FileName = defaultFileName(input.FileName, "rendered.docx")
	if input.Data == nil {
		input.Data = map[string]any{}
	}
	return input
}

func filterEmptyParagraphs(blocks []model.ContentBlock) []model.ContentBlock {
	if len(blocks) == 0 {
		return blocks
	}
	filtered := make([]model.ContentBlock, 0, len(blocks))
	for _, b := range blocks {
		if b.Type == "paragraph" && b.Text == "" && len(b.Runs) == 0 {
			continue
		}
		filtered = append(filtered, b)
	}
	return filtered
}

func validateGenerateRequest(input model.GenerateDocumentRequest) error {
	if len(input.Content) == 0 {
		return fmt.Errorf("content is required")
	}
	for _, block := range input.Content {
		switch block.Type {
		case "heading":
			if block.Text == "" {
				return fmt.Errorf("heading text is required")
			}
		case "paragraph":
			if block.Text == "" && len(block.Runs) == 0 {
				return fmt.Errorf("paragraph text or runs are required")
			}
		case "table":
			if len(block.Rows) == 0 {
				return fmt.Errorf("table rows are required")
			}
			width := len(block.Rows[0])
			if width == 0 {
				return fmt.Errorf("table must contain at least one column")
			}
			for _, row := range block.Rows {
				if len(row) != width {
					return fmt.Errorf("table rows must have the same column count")
				}
			}
		case "image":
			if strings.TrimSpace(block.ImageBase64) == "" && strings.TrimSpace(block.URL) == "" {
				return fmt.Errorf("image_base64 or url is required")
			}
		case "page_break":
		case "hyperlink":
			if strings.TrimSpace(block.URL) == "" {
				return fmt.Errorf("url is required")
			}
			if strings.TrimSpace(block.DisplayText) == "" && strings.TrimSpace(block.Text) == "" {
				return fmt.Errorf("display_text or text is required")
			}
		case "toc":
		default:
			return fmt.Errorf("unsupported block type: %s", block.Type)
		}
	}
	return nil
}

func validateTemplateRequest(input model.RenderTemplateRequest) error {
	if input.TemplateName == "" {
		return fmt.Errorf("template_name is required")
	}
	return nil
}

func shouldFallbackToStructured(err error) bool {
	if err == nil {
		return false
	}
	message := strings.ToLower(strings.TrimSpace(err.Error()))
	return strings.Contains(message, "template not found") || strings.Contains(message, "missing keys:")
}

func defaultFileName(value, fallback string) string {
	name := sanitizeFileName(value)
	if name == "" {
		name = fallback
	}
	if !strings.HasSuffix(strings.ToLower(name), ".docx") {
		name += ".docx"
	}
	return name
}

func uniqueDocxName(fileName string) string {
	ext := filepath.Ext(fileName)
	base := strings.TrimSuffix(fileName, ext)
	return fmt.Sprintf("%s-%d%s", base, time.Now().UnixNano(), ext)
}

func sanitizeFileName(value string) string {
	trimmed := strings.TrimSpace(filepath.Base(value))
	if trimmed == "." || trimmed == "/" {
		return ""
	}
	var builder strings.Builder
	for _, ch := range trimmed {
		switch {
		case unicode.IsLetter(ch), unicode.IsDigit(ch):
			builder.WriteRune(ch)
		case ch == '.', ch == '-', ch == '_':
			builder.WriteRune(ch)
		}
	}
	return builder.String()
}

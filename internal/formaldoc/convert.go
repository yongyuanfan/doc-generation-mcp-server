package formaldoc

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/yong/doc-generation-mcp-server/internal/model"
)

func ToGenerateRequest(d Draft) (ConversionResult, error) {
	if issues := ValidateDraft(d); len(issues) > 0 {
		return ConversionResult{}, issues[0]
	}

	request := model.GenerateDocumentRequest{
		FileName:         defaultFileNameFromTitle(d.Title),
		Title:            d.Title,
		Author:           d.Author,
		HeaderText:       d.HeaderText,
		FooterPageNumber: d.FooterPageNumber,
	}

	if d.IncludeTOC {
		request.Content = append(request.Content, model.ContentBlock{
			Type:   "toc",
			Text:   "目录",
			Levels: "1-3",
		})
	}

	for _, section := range d.Sections {
		request.Content = append(request.Content, model.ContentBlock{
			Type:  "heading",
			Text:  section.Title,
			Level: section.Level,
		})
		request.Content = append(request.Content, convertBlocks(section.Blocks)...)
	}

	if len(d.Appendices) > 0 {
		for _, appendix := range d.Appendices {
			request.Content = append(request.Content, model.ContentBlock{Type: "page_break"})
			request.Content = append(request.Content, model.ContentBlock{
				Type:  "heading",
				Text:  appendix.Title,
				Level: 1,
			})
			request.Content = append(request.Content, convertBlocks(appendix.Blocks)...)
		}
	}

	for _, reference := range d.References {
		request.Content = append(request.Content, model.ContentBlock{
			Type:        "hyperlink",
			URL:         reference.URL,
			DisplayText: reference.Title,
		})
	}

	return ConversionResult{
		Request:      request,
		ReviewNotes:  append([]string(nil), d.ReviewNotes...),
		TemplateName: d.TemplateName,
	}, nil
}

func convertBlocks(blocks []Block) []model.ContentBlock {
	converted := make([]model.ContentBlock, 0, len(blocks))
	for _, block := range blocks {
		converted = append(converted, model.ContentBlock{
			Type:        block.Type,
			Text:        block.Text,
			Level:       block.Level,
			Rows:        copyRows(block.Rows),
			ImageBase64: block.ImageBase64,
			Width:       block.Width,
			Height:      block.Height,
			URL:         block.URL,
			DisplayText: block.DisplayText,
			Levels:      block.Levels,
		})
	}
	return converted
}

func copyRows(rows [][]string) [][]string {
	if len(rows) == 0 {
		return nil
	}
	result := make([][]string, 0, len(rows))
	for _, row := range rows {
		result = append(result, append([]string(nil), row...))
	}
	return result
}

func defaultFileNameFromTitle(title string) string {
	trimmed := strings.TrimSpace(title)
	if trimmed == "" {
		return "document.docx"
	}
	var builder strings.Builder
	for _, ch := range trimmed {
		switch {
		case unicode.IsLetter(ch), unicode.IsDigit(ch):
			builder.WriteRune(ch)
		case ch == ' ', ch == '-', ch == '_':
			builder.WriteRune('-')
		}
	}
	result := strings.Trim(builder.String(), "-")
	if result == "" {
		return "document.docx"
	}
	return fmt.Sprintf("%s.docx", result)
}

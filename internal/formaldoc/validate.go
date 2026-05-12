package formaldoc

import (
	"fmt"
	"net/url"
	"strings"
)

var bannedPhrases = []string{
	"我觉得",
	"非常牛",
	"超级重要",
	"革命性",
	"接下来我将",
	"让我们来看看",
}

func ValidateDraft(d Draft) []ValidationIssue {
	var issues []ValidationIssue

	if strings.TrimSpace(d.SchemaVersion) != SchemaVersion {
		issues = append(issues, ValidationIssue{Field: "schema_version", Message: "must be 1.0"})
	}
	if !contains([]string{DocumentTypeProjectProposal, DocumentTypeWeeklyReport, DocumentTypeBusinessLetter}, d.DocumentType) {
		issues = append(issues, ValidationIssue{Field: "document_type", Message: "unsupported document type"})
	}
	if strings.TrimSpace(d.Title) == "" {
		issues = append(issues, ValidationIssue{Field: "title", Message: "is required"})
	}
	if strings.TrimSpace(d.Audience) == "" {
		issues = append(issues, ValidationIssue{Field: "audience", Message: "is required"})
	}
	if d.Tone != "formal" {
		issues = append(issues, ValidationIssue{Field: "tone", Message: "must be formal"})
	}
	if d.Language != "zh-CN" {
		issues = append(issues, ValidationIssue{Field: "language", Message: "must be zh-CN"})
	}
	if len(d.Sections) == 0 {
		issues = append(issues, ValidationIssue{Field: "sections", Message: "must not be empty"})
	}

	for i, section := range d.Sections {
		fieldPrefix := fmt.Sprintf("sections[%d]", i)
		if strings.TrimSpace(section.Title) == "" {
			issues = append(issues, ValidationIssue{Field: fieldPrefix + ".title", Message: "is required"})
		}
		if section.Level < 1 || section.Level > 3 {
			issues = append(issues, ValidationIssue{Field: fieldPrefix + ".level", Message: "must be between 1 and 3"})
		}
		issues = append(issues, validateBlocks(fieldPrefix+".blocks", section.Blocks)...)
	}

	for i, appendix := range d.Appendices {
		fieldPrefix := fmt.Sprintf("appendices[%d]", i)
		if strings.TrimSpace(appendix.Title) == "" {
			issues = append(issues, ValidationIssue{Field: fieldPrefix + ".title", Message: "is required"})
		}
		issues = append(issues, validateBlocks(fieldPrefix+".blocks", appendix.Blocks)...)
	}

	for i, reference := range d.References {
		fieldPrefix := fmt.Sprintf("references[%d]", i)
		if strings.TrimSpace(reference.Title) == "" {
			issues = append(issues, ValidationIssue{Field: fieldPrefix + ".title", Message: "is required"})
		}
		if !isValidURL(reference.URL) {
			issues = append(issues, ValidationIssue{Field: fieldPrefix + ".url", Message: "must be a valid URL"})
		}
	}

	issues = append(issues, validateDocumentTypeRules(d)...)
	issues = append(issues, validateFormalStyle(d)...)

	return issues
}

func validateBlocks(fieldPrefix string, blocks []Block) []ValidationIssue {
	var issues []ValidationIssue
	for i, block := range blocks {
		prefix := fmt.Sprintf("%s[%d]", fieldPrefix, i)
		switch block.Type {
		case "paragraph":
			if strings.TrimSpace(block.Text) == "" {
				issues = append(issues, ValidationIssue{Field: prefix + ".text", Message: "is required"})
			}
		case "heading":
			if strings.TrimSpace(block.Text) == "" {
				issues = append(issues, ValidationIssue{Field: prefix + ".text", Message: "is required"})
			}
			if block.Level < 1 || block.Level > 3 {
				issues = append(issues, ValidationIssue{Field: prefix + ".level", Message: "must be between 1 and 3"})
			}
		case "table":
			if len(block.Rows) == 0 {
				issues = append(issues, ValidationIssue{Field: prefix + ".rows", Message: "must not be empty"})
				continue
			}
			width := len(block.Rows[0])
			if width == 0 {
				issues = append(issues, ValidationIssue{Field: prefix + ".rows", Message: "must contain at least one column"})
			}
			for rowIndex, row := range block.Rows {
				if len(row) != width {
					issues = append(issues, ValidationIssue{Field: fmt.Sprintf("%s.rows[%d]", prefix, rowIndex), Message: "must match the first row column count"})
				}
			}
		case "image":
			if strings.TrimSpace(block.ImageBase64) == "" && !isValidURL(block.URL) {
				issues = append(issues, ValidationIssue{Field: prefix, Message: "image requires image_base64 or a valid url"})
			}
		case "hyperlink":
			if !isValidURL(block.URL) {
				issues = append(issues, ValidationIssue{Field: prefix + ".url", Message: "must be a valid URL"})
			}
			if strings.TrimSpace(block.DisplayText) == "" && strings.TrimSpace(block.Text) == "" {
				issues = append(issues, ValidationIssue{Field: prefix, Message: "display_text or text is required"})
			}
		case "toc":
		case "page_break":
		default:
			issues = append(issues, ValidationIssue{Field: prefix + ".type", Message: "unsupported block type"})
		}
	}
	return issues
}

func validateDocumentTypeRules(d Draft) []ValidationIssue {
	var issues []ValidationIssue
	sectionTitles := collectSectionTitles(d.Sections)
	tableCount := countBlocks(d, "table")

	switch d.DocumentType {
	case DocumentTypeProjectProposal:
		for _, title := range []string{"项目背景", "建设目标", "建设内容", "实施计划", "资源需求", "风险与保障措施", "预期成效", "结论与建议"} {
			if !containsSubstring(sectionTitles, title) {
				issues = append(issues, ValidationIssue{Field: "sections", Message: "missing required section: " + title})
			}
		}
		if !d.IncludeTOC {
			issues = append(issues, ValidationIssue{Field: "include_toc", Message: "should be enabled for project_proposal"})
		}
		if !d.FooterPageNumber {
			issues = append(issues, ValidationIssue{Field: "footer_page_number", Message: "should be enabled for project_proposal"})
		}
		if strings.TrimSpace(d.HeaderText) == "" {
			issues = append(issues, ValidationIssue{Field: "header_text", Message: "is required for project_proposal"})
		}
		if tableCount == 0 {
			issues = append(issues, ValidationIssue{Field: "sections", Message: "at least one table is required for project_proposal"})
		}
	case DocumentTypeWeeklyReport:
		for _, title := range []string{"本期工作概述", "已完成事项", "当前进展", "存在问题", "下阶段计划", "需协调事项"} {
			if !containsSubstring(sectionTitles, title) {
				issues = append(issues, ValidationIssue{Field: "sections", Message: "missing required section: " + title})
			}
		}
	case DocumentTypeBusinessLetter:
		for _, title := range []string{"发函背景", "发函事项", "具体说明", "后续安排"} {
			if !containsSubstring(sectionTitles, title) {
				issues = append(issues, ValidationIssue{Field: "sections", Message: "missing required section: " + title})
			}
		}
		if d.IncludeTOC {
			issues = append(issues, ValidationIssue{Field: "include_toc", Message: "should be disabled for business_letter"})
		}
	}

	return issues
}

func validateFormalStyle(d Draft) []ValidationIssue {
	var issues []ValidationIssue
	for _, text := range collectTexts(d) {
		for _, phrase := range bannedPhrases {
			if strings.Contains(text, phrase) {
				issues = append(issues, ValidationIssue{Field: "style", Message: "contains disallowed phrase: " + phrase})
			}
		}
	}
	return issues
}

func collectSectionTitles(sections []Section) []string {
	titles := make([]string, 0, len(sections))
	for _, section := range sections {
		titles = append(titles, section.Title)
	}
	return titles
}

func collectTexts(d Draft) []string {
	var texts []string
	texts = append(texts, d.Title, d.Subtitle, d.Summary)
	for _, section := range d.Sections {
		texts = append(texts, section.Title)
		for _, block := range section.Blocks {
			texts = append(texts, block.Text, block.DisplayText)
		}
	}
	for _, appendix := range d.Appendices {
		texts = append(texts, appendix.Title)
		for _, block := range appendix.Blocks {
			texts = append(texts, block.Text, block.DisplayText)
		}
	}
	for _, note := range d.ReviewNotes {
		texts = append(texts, note)
	}
	return texts
}

func countBlocks(d Draft, blockType string) int {
	count := 0
	for _, section := range d.Sections {
		for _, block := range section.Blocks {
			if block.Type == blockType {
				count++
			}
		}
	}
	for _, appendix := range d.Appendices {
		for _, block := range appendix.Blocks {
			if block.Type == blockType {
				count++
			}
		}
	}
	return count
}

func isValidURL(value string) bool {
	parsed, err := url.Parse(strings.TrimSpace(value))
	if err != nil {
		return false
	}
	return parsed.Scheme != "" && parsed.Host != ""
}

func contains(values []string, expected string) bool {
	for _, value := range values {
		if value == expected {
			return true
		}
	}
	return false
}

func containsSubstring(values []string, expected string) bool {
	for _, value := range values {
		if strings.Contains(value, expected) {
			return true
		}
	}
	return false
}

package formaldoc

import "github.com/yong/doc-generation-mcp-server/internal/model"

const SchemaVersion = "1.0"

const (
	DocumentTypeProjectProposal = "project_proposal"
	DocumentTypeWeeklyReport    = "weekly_report"
	DocumentTypeBusinessLetter  = "business_letter"
	RouteStructured             = "structured"
	RouteTemplate               = "template"
)

type Draft struct {
	SchemaVersion    string            `json:"schema_version,omitempty"`
	DocumentType     string            `json:"document_type,omitempty"`
	Title            string            `json:"title,omitempty"`
	Subtitle         string            `json:"subtitle,omitempty"`
	Author           string            `json:"author,omitempty"`
	Organization     string            `json:"organization,omitempty"`
	Audience         string            `json:"audience,omitempty"`
	Tone             string            `json:"tone,omitempty"`
	Language         string            `json:"language,omitempty"`
	HeaderText       string            `json:"header_text,omitempty"`
	FooterPageNumber bool              `json:"footer_page_number,omitempty"`
	IncludeTOC       bool              `json:"include_toc,omitempty"`
	TemplateName     string            `json:"template_name,omitempty"`
	Summary          string            `json:"summary,omitempty"`
	Sections         []Section         `json:"sections,omitempty"`
	Appendices       []Appendix        `json:"appendices,omitempty"`
	References       []Reference       `json:"references,omitempty"`
	Placeholders     map[string]string `json:"placeholders,omitempty"`
	ReviewNotes      []string          `json:"review_notes,omitempty"`
}

type Section struct {
	ID       string  `json:"id,omitempty"`
	Title    string  `json:"title,omitempty"`
	Level    int     `json:"level,omitempty"`
	Required bool    `json:"required,omitempty"`
	Blocks   []Block `json:"blocks,omitempty"`
}

type Appendix struct {
	Title  string  `json:"title,omitempty"`
	Blocks []Block `json:"blocks,omitempty"`
}

type Reference struct {
	Title string `json:"title,omitempty"`
	URL   string `json:"url,omitempty"`
}

type Block struct {
	Type        string     `json:"type,omitempty"`
	Text        string     `json:"text,omitempty"`
	Level       int        `json:"level,omitempty"`
	Rows        [][]string `json:"rows,omitempty"`
	URL         string     `json:"url,omitempty"`
	ImageBase64 string     `json:"image_base64,omitempty"`
	DisplayText string     `json:"display_text,omitempty"`
	Width       int        `json:"width,omitempty"`
	Height      int        `json:"height,omitempty"`
	Levels      string     `json:"levels,omitempty"`
}

type ValidationIssue struct {
	Field   string
	Message string
}

func (i ValidationIssue) Error() string {
	if i.Field == "" {
		return i.Message
	}
	return i.Field + ": " + i.Message
}

type ConversionResult struct {
	Request      model.GenerateDocumentRequest
	ReviewNotes  []string
	TemplateName string
}

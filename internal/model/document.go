package model

type GenerateDocumentRequest struct {
	FileName         string         `json:"file_name,omitempty"`
	Title            string         `json:"title,omitempty"`
	Author           string         `json:"author,omitempty"`
	Subject          string         `json:"subject,omitempty"`
	Keywords         []string       `json:"keywords,omitempty"`
	HeaderText       string         `json:"header_text,omitempty"`
	FooterPageNumber bool           `json:"footer_page_number,omitempty"`
	Content          []ContentBlock `json:"content"`
}

type RenderTemplateRequest struct {
	TemplateName string         `json:"template_name"`
	FileName     string         `json:"file_name,omitempty"`
	Data         map[string]any `json:"data"`
}

type ContentBlock struct {
	Type        string         `json:"type"`
	Text        string         `json:"text,omitempty"`
	Level       int            `json:"level,omitempty"`
	Alignment   string         `json:"alignment,omitempty"`
	Runs        []ParagraphRun `json:"runs,omitempty"`
	Rows        [][]string     `json:"rows,omitempty"`
	ImageBase64 string         `json:"image_base64,omitempty"`
	Width       int            `json:"width,omitempty"`
	Height      int            `json:"height,omitempty"`
	URL         string         `json:"url,omitempty"`
	DisplayText string         `json:"display_text,omitempty"`
	Levels      string         `json:"levels,omitempty"`
}

type ParagraphRun struct {
	Text       string `json:"text"`
	Bold       bool   `json:"bold,omitempty"`
	Italic     bool   `json:"italic,omitempty"`
	Underline  bool   `json:"underline,omitempty"`
	Color      string `json:"color,omitempty"`
	FontSize   int    `json:"font_size,omitempty"`
	FontFamily string `json:"font_family,omitempty"`
}

type DocumentResult struct {
	FileName    string `json:"file_name"`
	Path        string `json:"path"`
	DownloadURL string `json:"download_url"`
	MIMEType    string `json:"mime_type"`
	SizeBytes   int64  `json:"size_bytes"`
}

type ListTemplatesRequest struct{}

type TemplateInfo struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type ListTemplatesResponse struct {
	Templates []TemplateInfo `json:"templates"`
}

type CapabilitiesResponse struct {
	Formats          []string `json:"formats"`
	BlockTypes       []string `json:"block_types"`
	TemplateDir      string   `json:"template_dir"`
	TempDir          string   `json:"temp_dir"`
	TemplateRender   bool     `json:"template_render"`
	HeaderText       bool     `json:"header_text"`
	FooterPageNumber bool     `json:"footer_page_number"`
}

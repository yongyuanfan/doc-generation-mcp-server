package docx

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	docx "github.com/mmonterroca/docxgo/v2"
	"github.com/mmonterroca/docxgo/v2/domain"
	tpl "github.com/mmonterroca/docxgo/v2/pkg/template"
	"github.com/yong/doc-generation-mcp-server/internal/config"
	"github.com/yong/doc-generation-mcp-server/internal/model"
)

const mimeTypeDOCX = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"

type Client struct {
	config config.Config
}

func NewClient(cfg config.Config) *Client {
	return &Client{config: cfg}
}

func (c *Client) Generate(_ context.Context, input model.GenerateDocumentRequest, outputName string) (model.DocumentResult, error) {
	document := docx.NewDocument()

	if err := c.applyMetadata(document, input); err != nil {
		return model.DocumentResult{}, err
	}
	if strings.TrimSpace(input.HeaderText) != "" {
		if err := c.addHeaderText(document, input.HeaderText); err != nil {
			return model.DocumentResult{}, err
		}
	}
	if input.FooterPageNumber {
		if err := c.addFooterPageNumbers(document); err != nil {
			return model.DocumentResult{}, err
		}
	}

	for _, block := range input.Content {
		if err := c.addBlock(document, block); err != nil {
			return model.DocumentResult{}, err
		}
	}

	return c.saveDocument(document, outputName)
}

func (c *Client) RenderTemplate(_ context.Context, templatePath string, data map[string]any, outputName string) (model.DocumentResult, error) {
	document, err := docx.OpenDocument(templatePath)
	if err != nil {
		return model.DocumentResult{}, err
	}

	mergeData := make(tpl.MergeData, len(data))
	for key, value := range data {
		mergeData[key] = fmt.Sprint(value)
	}

	options := tpl.DefaultMergeOptions()
	options.StrictMode = true
	if err := tpl.MergeTemplate(document, mergeData, options); err != nil {
		return model.DocumentResult{}, err
	}

	return c.saveDocument(document, outputName)
}

func (c *Client) applyMetadata(document domain.Document, input model.GenerateDocumentRequest) error {
	meta := document.Metadata()
	if meta == nil {
		meta = &domain.Metadata{}
	}
	meta.Title = strings.TrimSpace(input.Title)
	meta.Subject = strings.TrimSpace(input.Subject)
	meta.Creator = strings.TrimSpace(input.Author)
	meta.Keywords = append([]string(nil), input.Keywords...)
	now := time.Now().UTC().Format(time.RFC3339)
	meta.Created = now
	meta.Modified = now
	return document.SetMetadata(meta)
}

func (c *Client) addBlock(document domain.Document, block model.ContentBlock) error {
	switch block.Type {
	case "heading":
		return c.addHeading(document, block)
	case "paragraph":
		return c.addParagraph(document, block)
	case "table":
		return c.addTable(document, block)
	case "image":
		return c.addImage(document, block)
	case "page_break":
		return document.AddPageBreak()
	case "hyperlink":
		return c.addHyperlink(document, block)
	case "toc":
		return c.addTOC(document, block)
	default:
		return fmt.Errorf("unsupported block type: %s", block.Type)
	}
}

func (c *Client) addHeading(document domain.Document, block model.ContentBlock) error {
	paragraph, err := document.AddParagraph()
	if err != nil {
		return err
	}
	if err := paragraph.SetStyle(headingStyle(block.Level)); err != nil {
		return err
	}
	run, err := paragraph.AddRun()
	if err != nil {
		return err
	}
	return run.SetText(block.Text)
}

func (c *Client) addParagraph(document domain.Document, block model.ContentBlock) error {
	paragraph, err := document.AddParagraph()
	if err != nil {
		return err
	}
	if err := paragraph.SetAlignment(parseAlignment(block.Alignment)); err != nil {
		return err
	}
	for _, item := range block.Runs {
		if err := addRun(paragraph, item, c.config); err != nil {
			return err
		}
	}
	if len(block.Runs) == 0 {
		run, err := paragraph.AddRun()
		if err != nil {
			return err
		}
		return run.SetText(block.Text)
	}
	return nil
}

func (c *Client) addTable(document domain.Document, block model.ContentBlock) error {
	rows := len(block.Rows)
	if rows == 0 {
		return fmt.Errorf("table rows are required")
	}
	cols := len(block.Rows[0])
	table, err := document.AddTable(rows, cols)
	if err != nil {
		return err
	}
	if err := table.SetStyle(domain.TableStyleGrid); err != nil {
		return err
	}
	for rowIndex, values := range block.Rows {
		row, err := table.Row(rowIndex)
		if err != nil {
			return err
		}
		for colIndex, value := range values {
			cell, err := row.Cell(colIndex)
			if err != nil {
				return err
			}
			paragraph, err := cell.AddParagraph()
			if err != nil {
				return err
			}
			run, err := paragraph.AddRun()
			if err != nil {
				return err
			}
			if err := run.SetText(value); err != nil {
				return err
			}
			if rowIndex == 0 {
				if err := run.SetBold(true); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (c *Client) addImage(document domain.Document, block model.ContentBlock) error {
	data, format, err := c.resolveImage(block)
	if err != nil {
		return err
	}
	paragraph, err := document.AddParagraph()
	if err != nil {
		return err
	}
	if block.Width > 0 || block.Height > 0 {
		_, err = paragraph.AddImageFromBytesWithSize(data, format, domain.NewImageSize(block.Width, block.Height))
		return err
	}
	_, err = paragraph.AddImageFromBytes(data, format)
	return err
}

func (c *Client) addHyperlink(document domain.Document, block model.ContentBlock) error {
	paragraph, err := document.AddParagraph()
	if err != nil {
		return err
	}
	if err := paragraph.SetAlignment(parseAlignment(block.Alignment)); err != nil {
		return err
	}
	displayText := strings.TrimSpace(block.DisplayText)
	if displayText == "" {
		displayText = strings.TrimSpace(block.Text)
	}
	_, err = paragraph.AddHyperlink(strings.TrimSpace(block.URL), displayText)
	return err
}

func (c *Client) addTOC(document domain.Document, block model.ContentBlock) error {
	if strings.TrimSpace(block.Text) != "" {
		heading, err := document.AddParagraph()
		if err != nil {
			return err
		}
		if err := heading.SetStyle("Heading1"); err != nil {
			return err
		}
		run, err := heading.AddRun()
		if err != nil {
			return err
		}
		if err := run.SetText(block.Text); err != nil {
			return err
		}
	}
	paragraph, err := document.AddParagraph()
	if err != nil {
		return err
	}
	run, err := paragraph.AddRun()
	if err != nil {
		return err
	}
	levels := strings.TrimSpace(block.Levels)
	if levels == "" {
		levels = "1-3"
	}
	return run.AddField(docx.NewTOCField(map[string]string{
		"levels":     levels,
		"hyperlinks": "true",
	}))
}

func (c *Client) addHeaderText(document domain.Document, value string) error {
	section, err := document.DefaultSection()
	if err != nil {
		return err
	}
	header, err := section.Header(domain.HeaderDefault)
	if err != nil {
		return err
	}
	paragraph, err := header.AddParagraph()
	if err != nil {
		return err
	}
	if err := paragraph.SetAlignment(domain.AlignmentCenter); err != nil {
		return err
	}
	run, err := paragraph.AddRun()
	if err != nil {
		return err
	}
	return run.SetText(strings.TrimSpace(value))
}

func (c *Client) addFooterPageNumbers(document domain.Document) error {
	section, err := document.DefaultSection()
	if err != nil {
		return err
	}
	footer, err := section.Footer(domain.FooterDefault)
	if err != nil {
		return err
	}
	paragraph, err := footer.AddParagraph()
	if err != nil {
		return err
	}
	if err := paragraph.SetAlignment(domain.AlignmentCenter); err != nil {
		return err
	}
	before, err := paragraph.AddRun()
	if err != nil {
		return err
	}
	if err := before.AddText("Page "); err != nil {
		return err
	}
	pageNumberRun, err := paragraph.AddRun()
	if err != nil {
		return err
	}
	if err := pageNumberRun.AddField(docx.NewPageNumberField()); err != nil {
		return err
	}
	middle, err := paragraph.AddRun()
	if err != nil {
		return err
	}
	if err := middle.AddText(" of "); err != nil {
		return err
	}
	pageCountRun, err := paragraph.AddRun()
	if err != nil {
		return err
	}
	return pageCountRun.AddField(docx.NewPageCountField())
}

func (c *Client) resolveImage(block model.ContentBlock) ([]byte, domain.ImageFormat, error) {
	if strings.TrimSpace(block.ImageBase64) != "" {
		return decodeImage(block.ImageBase64)
	}
	if strings.TrimSpace(block.URL) == "" {
		return nil, "", fmt.Errorf("image_base64 or url is required")
	}
	return c.fetchImage(strings.TrimSpace(block.URL))
}

func (c *Client) fetchImage(url string) ([]byte, domain.ImageFormat, error) {
	client := &http.Client{Timeout: c.config.RequestTimeout}
	resp, err := client.Get(url)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, "", fmt.Errorf("fetch image url: status %d", resp.StatusCode)
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, c.config.DocxMaxRequestBodyBytes))
	if err != nil {
		return nil, "", err
	}
	return data, imageFormatFromContentType(resp.Header.Get("Content-Type")), nil
}

func (c *Client) saveDocument(document domain.Document, outputName string) (model.DocumentResult, error) {
	if err := document.Validate(); err != nil {
		return model.DocumentResult{}, err
	}
	path := filepath.Join(c.config.DocxTempDir, outputName)
	if err := document.SaveAs(path); err != nil {
		return model.DocumentResult{}, err
	}
	stat, err := os.Stat(path)
	if err != nil {
		return model.DocumentResult{}, err
	}
	return model.DocumentResult{
		FileName:    outputName,
		Path:        path,
		DownloadURL: "/api/v1/documents/files/" + outputName,
		MIMEType:    mimeTypeDOCX,
		SizeBytes:   stat.Size(),
	}, nil
}

func addRun(paragraph domain.Paragraph, item model.ParagraphRun, cfg config.Config) error {
	run, err := paragraph.AddRun()
	if err != nil {
		return err
	}
	if err := run.SetText(item.Text); err != nil {
		return err
	}
	if item.Bold {
		if err := run.SetBold(true); err != nil {
			return err
		}
	}
	if item.Italic {
		if err := run.SetItalic(true); err != nil {
			return err
		}
	}
	if item.Underline {
		if err := run.SetUnderline(docx.UnderlineSingle); err != nil {
			return err
		}
	}
	if item.FontSize > 0 {
		if err := run.SetSize(item.FontSize); err != nil {
			return err
		}
	}
	fontName := strings.TrimSpace(item.FontFamily)
	if fontName == "" {
		fontName = cfg.DocxDefaultFont
	}
	if err := run.SetFont(domain.Font{Name: fontName}); err != nil {
		return err
	}
	if colorValue := strings.TrimSpace(item.Color); colorValue != "" {
		parsed, err := parseHexColor(colorValue)
		if err != nil {
			return err
		}
		if err := run.SetColor(parsed); err != nil {
			return err
		}
	}
	return nil
}

func parseAlignment(value string) domain.Alignment {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "center":
		return domain.AlignmentCenter
	case "right":
		return domain.AlignmentRight
	case "justify":
		return domain.AlignmentJustify
	default:
		return domain.AlignmentLeft
	}
}

func headingStyle(level int) string {
	switch level {
	case 2:
		return "Heading2"
	case 3:
		return "Heading3"
	case 4:
		return "Heading4"
	default:
		return "Heading1"
	}
}

func parseHexColor(value string) (domain.Color, error) {
	trimmed := strings.TrimPrefix(strings.TrimSpace(value), "#")
	if len(trimmed) != 6 {
		return domain.Color{}, fmt.Errorf("invalid color: %s", value)
	}
	var err error
	var rgb [3]uint8
	for i := 0; i < 3; i++ {
		var component uint64
		component, err = parseHexByte(trimmed[i*2 : i*2+2])
		if err != nil {
			return domain.Color{}, fmt.Errorf("invalid color: %s", value)
		}
		rgb[i] = uint8(component)
	}
	return domain.Color{R: rgb[0], G: rgb[1], B: rgb[2]}, nil
}

func parseHexByte(value string) (uint64, error) {
	var parsed uint64
	for _, ch := range value {
		parsed <<= 4
		switch {
		case ch >= '0' && ch <= '9':
			parsed += uint64(ch - '0')
		case ch >= 'a' && ch <= 'f':
			parsed += uint64(ch-'a') + 10
		case ch >= 'A' && ch <= 'F':
			parsed += uint64(ch-'A') + 10
		default:
			return 0, fmt.Errorf("invalid hex byte")
		}
	}
	return parsed, nil
}

func decodeImage(value string) ([]byte, domain.ImageFormat, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil, "", fmt.Errorf("image_base64 is required")
	}
	format := domain.ImageFormatPNG
	if strings.HasPrefix(trimmed, "data:") {
		prefix, payload, ok := strings.Cut(trimmed, ",")
		if !ok {
			return nil, "", fmt.Errorf("invalid image data url")
		}
		trimmed = payload
		format = imageFormatFromDataURL(prefix)
	}
	data, err := base64.StdEncoding.DecodeString(trimmed)
	if err != nil {
		data, err = base64.RawStdEncoding.DecodeString(trimmed)
	}
	if err != nil {
		return nil, "", fmt.Errorf("decode image_base64: %w", err)
	}
	return data, format, nil
}

func imageFormatFromDataURL(prefix string) domain.ImageFormat {
	switch {
	case strings.Contains(prefix, "image/jpeg"):
		return domain.ImageFormatJPEG
	case strings.Contains(prefix, "image/gif"):
		return domain.ImageFormatGIF
	case strings.Contains(prefix, "image/webp"):
		return domain.ImageFormatWEBP
	case strings.Contains(prefix, "image/svg+xml"):
		return domain.ImageFormatSVG
	default:
		return domain.ImageFormatPNG
	}
}

func imageFormatFromContentType(contentType string) domain.ImageFormat {
	switch {
	case strings.Contains(contentType, "image/jpeg"):
		return domain.ImageFormatJPEG
	case strings.Contains(contentType, "image/gif"):
		return domain.ImageFormatGIF
	case strings.Contains(contentType, "image/webp"):
		return domain.ImageFormatWEBP
	case strings.Contains(contentType, "image/svg+xml"):
		return domain.ImageFormatSVG
	case strings.Contains(contentType, "image/bmp"):
		return domain.ImageFormatBMP
	case strings.Contains(contentType, "image/tiff"):
		return domain.ImageFormatTIFF
	default:
		return domain.ImageFormatPNG
	}
}

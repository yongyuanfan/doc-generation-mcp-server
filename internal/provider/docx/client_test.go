package docx

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/yong/doc-generation-mcp-server/internal/config"
	"github.com/yong/doc-generation-mcp-server/internal/model"
)

var tinyPNG = []byte{
	0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a,
	0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52,
	0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
	0x08, 0x06, 0x00, 0x00, 0x00, 0x1f, 0x15, 0xc4,
	0x89, 0x00, 0x00, 0x00, 0x0d, 0x49, 0x44, 0x41,
	0x54, 0x78, 0x9c, 0x63, 0xf8, 0xcf, 0xc0, 0x00,
	0x00, 0x03, 0x01, 0x01, 0x00, 0xc9, 0xfe, 0x92,
	0xef, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4e,
	0x44, 0xae, 0x42, 0x60, 0x82,
}

func TestGenerateSupportsImageURLAndHeaderAndTOC(t *testing.T) {
	tempDir := t.TempDir()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		_, _ = w.Write(tinyPNG)
	}))
	defer server.Close()

	client := NewClient(config.Config{
		DocxTempDir:             tempDir,
		DocxDefaultFont:         "Calibri",
		DocxDefaultFontSize:     22,
		DocxMaxRequestBodyBytes: 2 << 20,
		RequestTimeout:          5 * time.Second,
	})

	result, err := client.Generate(context.Background(), model.GenerateDocumentRequest{
		Title:            "Advanced",
		Author:           "tester",
		HeaderText:       "Header",
		FooterPageNumber: true,
		Content: []model.ContentBlock{
			{Type: "toc", Text: "Contents", Levels: "1-2"},
			{Type: "heading", Text: "Section 1", Level: 1},
			{Type: "image", URL: server.URL},
			{Type: "paragraph", Text: "hello"},
		},
	}, "advanced.docx")
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(result.Path) != "advanced.docx" {
		t.Fatalf("unexpected output path: %s", result.Path)
	}
	if result.SizeBytes == 0 {
		t.Fatal("expected generated file size")
	}
}

func TestDecodeImageSupportsBase64DataURL(t *testing.T) {
	encoded := "data:image/png;base64," + base64.StdEncoding.EncodeToString(tinyPNG)
	data, format, err := decodeImage(encoded)
	if err != nil {
		t.Fatal(err)
	}
	if format != "png" {
		t.Fatalf("expected png format, got %s", format)
	}
	if len(data) == 0 {
		t.Fatal("expected decoded image data")
	}
}

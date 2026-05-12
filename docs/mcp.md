# MCP

This project exposes a Streamable HTTP MCP server at `/mcp`.

## Tools

### `generate_docx`

Generate a DOCX document from structured content blocks.

### `generate_docx_from_draft`

Generate a DOCX document from a validated `FormalDocumentDraftV1` payload.

If `template_name` is present, the service will automatically use template routing.

### `validate_formal_document_draft`

Validate a `FormalDocumentDraftV1` payload without generating a DOCX file.

### `render_docx_template`

Render a DOCX document from a template stored in `templates/`.

### `list_docx_templates`

List available `.docx` templates from `templates/`.

## Local Verification

Connect an MCP client to:

```text
http://localhost:9103/mcp
```

## JSON-RPC Example

```bash
curl -X POST http://localhost:9103/mcp \
  -H 'Content-Type: application/json' \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "tools/call",
    "params": {
      "name": "generate_docx",
      "arguments": {
        "file_name": "demo.docx",
        "header_text": "Demo Header",
        "footer_page_number": true,
        "content": [
          {"type": "toc", "text": "Contents", "levels": "1-2"},
          {"type": "heading", "text": "Demo", "level": 1},
          {"type": "paragraph", "text": "Hello from MCP"},
          {"type": "hyperlink", "url": "https://example.com", "display_text": "Example"},
          {"type": "page_break"}
        ]
      }
    }
  }'
```

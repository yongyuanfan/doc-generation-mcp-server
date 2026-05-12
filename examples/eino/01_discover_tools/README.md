# 01 Discover Tools

## Purpose

Connect to the remote MCP endpoint and list all available tools.

## Prerequisites

Start the server first:

```bash
cp .env.example .env
go run ./cmd/server
```

## Run

```bash
MCP_SERVER_URL=http://localhost:9103/mcp go run ./examples/eino/01_discover_tools
```

## Expected Output

The example prints the tool names and descriptions exposed by the service, such as:

- `generate_docx`
- `generate_docx_from_draft`
- `validate_formal_document_draft`
- `render_docx_template`
- `list_docx_templates`

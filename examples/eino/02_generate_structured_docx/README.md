# 02 Generate Structured DOCX

## Purpose

Call `generate_docx` directly with structured content blocks.

## Prerequisites

Start the server first:

```bash
cp .env.example .env
go run ./cmd/server
```

## Run

```bash
MCP_SERVER_URL=http://localhost:9103/mcp go run ./examples/eino/02_generate_structured_docx
```

## Expected Output

The example prints the generated file metadata, including:

- `file_name`
- `download_url`
- `mime_type`
- `size_bytes`

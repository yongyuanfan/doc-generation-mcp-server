# doc-generation-mcp-server

A remote MCP and REST service for generating `.docx` files with `github.com/mmonterroca/docxgo/v2`.

## Features

- Streamable HTTP MCP server
- REST API for structured DOCX generation
- REST API for FormalDocumentDraftV1 based generation
- DOCX template rendering from `templates/`
- Template listing via REST and MCP
- Structured blocks: heading, paragraph, table, image, page break, hyperlink, toc
- Header text support
- Optional footer page numbers
- Image embedding from base64 or URL
- Runtime temp-file storage with download endpoint
- Eino integration example
- Multiple Eino integration examples
- Docker deployment support

## Quick Start

```bash
cp .env.example .env
go run ./cmd/server
```

Default endpoints:

- MCP: `http://localhost:9101/mcp`
- API: `http://localhost:9101/api/v1`
- Health: `http://localhost:9101/healthz`

## Templates

Template files must be placed in `templates/`.

An example template is generated as `templates/example-letter.docx`.

## Smoke Test

Start the server, then run:

```bash
bash scripts/smoke_test.sh
```

## Documents

- `docs/api.md`
- `docs/mcp.md`
- `docs/deploy.md`
- `docs/eino.md`
- `docs/formal-document-schema.md`
- `docs/document-type-rules.md`
- `docs/agent-workflow.md`
- `docs/draft-to-docx-mapping.md`
- `docs/formal-document-examples.md`
- `docs/llm-integration-example.md`

## Internal Building Blocks

- `schemas/formal-document-draft-v1.json`
- `internal/formaldoc`

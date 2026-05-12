# doc-generation-mcp-server

A remote MCP and REST service for generating `.docx` files with `github.com/mmonterroca/docxgo/v2`.

## Features

- Streamable HTTP MCP server
- REST API for structured DOCX generation
- DOCX template rendering from `templates/`
- Template listing via REST and MCP
- Structured blocks: heading, paragraph, table, image, page break, hyperlink
- Structured blocks: heading, paragraph, table, image, page break, hyperlink, toc
- Header text support
- Optional footer page numbers
- Image embedding from base64 or URL
- Runtime temp-file storage with download endpoint
- Eino integration example
- Docker deployment support

## Quick Start

```bash
cp .env.example .env
go run ./cmd/server
```

Default endpoints:

- MCP: `http://localhost:9103/mcp`
- API: `http://localhost:9103/api/v1`
- Health: `http://localhost:9103/healthz`

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

# 03 Validate And Generate Draft

## Purpose

Build a `FormalDocumentDraftV1`, validate it first, then generate the final `.docx` file.

## Prerequisites

Start the server first:

```bash
cp .env.example .env
go run ./cmd/server
```

## Run

```bash
MCP_SERVER_URL=http://localhost:9103/mcp go run ./examples/eino/03_validate_and_generate_draft
```

## Expected Output

The example prints:

1. validation result
2. generated file metadata

This is the recommended baseline flow for draft-based document generation.

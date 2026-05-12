# Eino Integration

This repository provides a set of runnable Eino-oriented examples under `examples/eino/`.

These examples focus on the orchestration style used in Eino-based applications while calling the remote MCP server through the official MCP Go SDK.

Each example directory also includes its own local `README.md` for standalone usage.

## Start The Service

```bash
cp .env.example .env
go run ./cmd/server
```

Default MCP endpoint:

```text
http://localhost:9103/mcp
```

You can override it with:

```bash
MCP_SERVER_URL=http://localhost:9103/mcp
```

## Example List

### 1. Discover Tools

```bash
MCP_SERVER_URL=http://localhost:9103/mcp go run ./examples/eino/01_discover_tools
```

Shows how to connect to the remote MCP service and list the available tools.

### 2. Generate Structured DOCX

```bash
MCP_SERVER_URL=http://localhost:9103/mcp go run ./examples/eino/02_generate_structured_docx
```

Calls `generate_docx` directly with structured content blocks.

### 3. Validate And Generate Draft

```bash
MCP_SERVER_URL=http://localhost:9103/mcp go run ./examples/eino/03_validate_and_generate_draft
```

Builds a `FormalDocumentDraftV1`, validates it, then generates the `.docx` file.

### 4. Template Route Draft

```bash
MCP_SERVER_URL=http://localhost:9103/mcp go run ./examples/eino/04_template_route_draft
```

Builds a `business_letter` draft and demonstrates template-based routing.

### 5. End-To-End Formal Doc Agent

```bash
MCP_SERVER_URL=http://localhost:9103/mcp go run ./examples/eino/05_end_to_end_formal_doc_agent
```

Demonstrates a complete high-level flow:

1. mock LLM output draft
2. validate draft
3. generate document
4. print result metadata and download URL

This example also supports a real OpenAI-compatible model through:

- `LLM_MODE=openai`
- `OPENAI_API_KEY`
- `OPENAI_MODEL`
- `OPENAI_BASE_URL`

## Notes

- Each example `main.go` uses `//go:build ignore`
- The shared utilities under `examples/eino/shared` are used by all examples
- These examples are designed to be run individually with `go run`

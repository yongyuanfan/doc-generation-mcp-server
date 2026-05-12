# 05 End-To-End Formal Doc Agent

## Purpose

Demonstrate a complete high-level orchestration flow similar to an Eino agent:

1. mock LLM output draft
2. validate draft
3. generate document
4. print result metadata and download URL

## Prerequisites

Start the server first:

```bash
cp .env.example .env
go run ./cmd/server
```

## Run

```bash
MCP_SERVER_URL=http://localhost:9103/mcp go run ./examples/eino/05_end_to_end_formal_doc_agent
```

## Document Type Switch

Default document type is `weekly_report`.

To run the built-in `business_letter` scenario:

```bash
MCP_SERVER_URL=http://localhost:9103/mcp go run ./examples/eino/05_end_to_end_formal_doc_agent -- -document-type business_letter
```

## JSON Input File

You can provide a custom JSON input file:

```bash
MCP_SERVER_URL=http://localhost:9103/mcp go run ./examples/eino/05_end_to_end_formal_doc_agent -- -input ./examples/eino/05_end_to_end_formal_doc_agent/input.business_letter.json
```

Input schema:

```json
{
  "document_type": "weekly_report",
  "goal": "提交管理层周报",
  "audience": "management",
  "facts": ["事实1", "事实2"]
}
```

## Real LLM Mode

The example supports two modes:

1. `mock`
2. `openai`

Default mode is `mock`.

To run with a real OpenAI-compatible model:

```bash
LLM_MODE=openai \
OPENAI_API_KEY=your_api_key \
OPENAI_MODEL=gpt-4o-mini \
OPENAI_BASE_URL=https://api.openai.com/v1 \
MCP_SERVER_URL=http://localhost:9103/mcp \
go run ./examples/eino/05_end_to_end_formal_doc_agent -- -document-type business_letter
```

## Expected Output

The example prints:

1. the mock `FormalDocumentDraftV1` JSON
2. validation result
3. generated file metadata

When `LLM_MODE=openai` is enabled, the example prints the raw model-produced draft JSON before validation.

Use this example as the closest reference for integrating an upper-layer LLM workflow with the DOCX MCP service.

# 04 Template Route Draft

## Purpose

Show how a `business_letter` draft follows template-first routing.

## Prerequisites

Start the server first:

```bash
cp .env.example .env
go run ./cmd/server
```

## Run

```bash
MCP_SERVER_URL=http://localhost:9103/mcp go run ./examples/eino/04_template_route_draft
```

## Expected Output

The example prints:

1. validation result with template recommendation
2. generation result with `route: template`

This demonstrates automatic template-based rendering for formal letters.

# Eino Integration

Run the server:

```bash
cp .env.example .env
go run ./cmd/server
```

Run the example:

```bash
MCP_SERVER_URL=http://localhost:9103/mcp go run ./examples/eino
```

The example connects to the remote MCP endpoint and discovers:

1. `generate_docx`
2. `render_docx_template`
3. `list_docx_templates`

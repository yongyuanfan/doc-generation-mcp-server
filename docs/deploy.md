# Deployment

## Docker

```bash
cp .env.example .env
docker compose up --build
```

The service listens on port `9103` by default.

## Exposed Endpoints

- `GET /healthz`
- `GET /api/v1/capabilities`
- `GET /api/v1/templates`
- `POST /api/v1/documents/generate`
- `POST /api/v1/documents/render-template`
- `GET /api/v1/documents/files/{name}`
- `POST /mcp`

## Optional Environment Variables

- `HTTP_ADDR`
- `API_PREFIX`
- `MCP_PATH`
- `REQUEST_TIMEOUT_SECONDS`
- `DOCX_TEMP_DIR`
- `DOCX_TEMPLATE_DIR`
- `DOCX_DEFAULT_AUTHOR`
- `DOCX_DEFAULT_FONT`
- `DOCX_DEFAULT_FONT_SIZE`
- `DOCX_MAX_REQUEST_BODY_BYTES`
- `DOCX_MAX_FILE_AGE_MINUTES`

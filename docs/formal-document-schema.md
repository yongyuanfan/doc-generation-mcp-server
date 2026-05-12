# FormalDocumentDraftV1

`FormalDocumentDraftV1` is the recommended structured intermediate format for large-model generated formal documents.

The model should output this JSON draft first. The draft is then reviewed, validated, and converted into requests for this DOCX generation service.

## Design Goals

- Keep the schema simple enough for LLMs to follow reliably
- Keep the schema strict enough for deterministic validation
- Separate writing concerns from rendering concerns
- Make document type rules explicit instead of implicit in prompts

## Top-Level Object

Required fields:

- `schema_version`
- `document_type`
- `title`
- `audience`
- `tone`
- `language`
- `sections`

Optional fields:

- `subtitle`
- `author`
- `organization`
- `header_text`
- `footer_page_number`
- `include_toc`
- `template_name`
- `summary`
- `appendices`
- `references`
- `placeholders`
- `review_notes`

## Field Definitions

### `schema_version`

- Type: string
- Required: yes
- Allowed value: `1.0`

### `document_type`

- Type: string
- Required: yes
- Allowed values:
  - `project_proposal`
  - `weekly_report`
  - `business_letter`

### `title`

- Type: string
- Required: yes
- Minimum length: 1

### `subtitle`

- Type: string
- Required: no

### `author`

- Type: string
- Required: no

### `organization`

- Type: string
- Required: no

### `audience`

- Type: string
- Required: yes
- Allowed values:
  - `management`
  - `customer`
  - `internal_team`
  - `government`
  - `partner`

### `tone`

- Type: string
- Required: yes
- Allowed value: `formal`

### `language`

- Type: string
- Required: yes
- Allowed value: `zh-CN`

### `header_text`

- Type: string
- Required: no

### `footer_page_number`

- Type: boolean
- Required: no

### `include_toc`

- Type: boolean
- Required: no

### `template_name`

- Type: string
- Required: no
- Intended use: route to template-first rendering when needed

### `summary`

- Type: string
- Required: no

### `sections`

- Type: array
- Required: yes
- Minimum items: 1

Each section object contains:

- `id`: string, optional but recommended
- `title`: string, required
- `level`: integer, required, allowed values `1` to `3`
- `required`: boolean, optional
- `blocks`: array, required

### `appendices`

- Type: array
- Required: no

Each appendix object contains:

- `title`: string, required
- `blocks`: array, required

### `references`

- Type: array
- Required: no

Each reference object contains:

- `title`: string, required
- `url`: string, required

### `placeholders`

- Type: object
- Required: no
- Value type: string preferred

### `review_notes`

- Type: array
- Required: no
- Item type: string
- Intended use: missing facts, verification notes, follow-up reminders

## Block Types

Supported block types for `sections[].blocks[]` and `appendices[].blocks[]`:

- `paragraph`
- `heading`
- `table`
- `image`
- `hyperlink`
- `page_break`
- `toc`

### `paragraph`

```json
{
  "type": "paragraph",
  "text": "正文段落内容"
}
```

Rules:

- `text` is required
- `text` must not be empty

### `heading`

```json
{
  "type": "heading",
  "text": "二、建设目标",
  "level": 1
}
```

Rules:

- `text` is required
- `level` is required
- `level` must be `1`, `2`, or `3`

### `table`

```json
{
  "type": "table",
  "rows": [
    ["项目", "说明"],
    ["目标", "统一标准"]
  ]
}
```

Rules:

- `rows` is required
- `rows` must contain at least one row
- every row must have the same column count

### `image`

```json
{
  "type": "image",
  "url": "https://example.com/logo.png",
  "width": 160,
  "height": 80
}
```

or:

```json
{
  "type": "image",
  "image_base64": "data:image/png;base64,..."
}
```

Rules:

- either `url` or `image_base64` is required
- `width` and `height` are optional positive integers

### `hyperlink`

```json
{
  "type": "hyperlink",
  "url": "https://example.com",
  "display_text": "参考链接"
}
```

Rules:

- `url` is required
- at least one of `display_text` or `text` is required

### `page_break`

```json
{
  "type": "page_break"
}
```

Rules:

- no additional fields required

### `toc`

```json
{
  "type": "toc",
  "text": "目录",
  "levels": "1-3"
}
```

Rules:

- `text` is optional
- `levels` is optional, recommended format `1-3`

## Strict JSON Schema

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "https://example.com/schemas/formal-document-draft-v1.json",
  "title": "FormalDocumentDraftV1",
  "type": "object",
  "additionalProperties": false,
  "required": [
    "schema_version",
    "document_type",
    "title",
    "audience",
    "tone",
    "language",
    "sections"
  ],
  "properties": {
    "schema_version": {
      "type": "string",
      "const": "1.0"
    },
    "document_type": {
      "type": "string",
      "enum": ["project_proposal", "weekly_report", "business_letter"]
    },
    "title": {
      "type": "string",
      "minLength": 1
    },
    "subtitle": {
      "type": "string"
    },
    "author": {
      "type": "string"
    },
    "organization": {
      "type": "string"
    },
    "audience": {
      "type": "string",
      "enum": ["management", "customer", "internal_team", "government", "partner"]
    },
    "tone": {
      "type": "string",
      "const": "formal"
    },
    "language": {
      "type": "string",
      "const": "zh-CN"
    },
    "header_text": {
      "type": "string"
    },
    "footer_page_number": {
      "type": "boolean"
    },
    "include_toc": {
      "type": "boolean"
    },
    "template_name": {
      "type": "string"
    },
    "summary": {
      "type": "string"
    },
    "sections": {
      "type": "array",
      "minItems": 1,
      "items": {
        "$ref": "#/$defs/section"
      }
    },
    "appendices": {
      "type": "array",
      "items": {
        "$ref": "#/$defs/appendix"
      }
    },
    "references": {
      "type": "array",
      "items": {
        "$ref": "#/$defs/reference"
      }
    },
    "placeholders": {
      "type": "object",
      "additionalProperties": {
        "type": "string"
      }
    },
    "review_notes": {
      "type": "array",
      "items": {
        "type": "string"
      }
    }
  },
  "$defs": {
    "section": {
      "type": "object",
      "additionalProperties": false,
      "required": ["title", "level", "blocks"],
      "properties": {
        "id": {"type": "string"},
        "title": {"type": "string", "minLength": 1},
        "level": {"type": "integer", "minimum": 1, "maximum": 3},
        "required": {"type": "boolean"},
        "blocks": {
          "type": "array",
          "items": {"$ref": "#/$defs/block"}
        }
      }
    },
    "appendix": {
      "type": "object",
      "additionalProperties": false,
      "required": ["title", "blocks"],
      "properties": {
        "title": {"type": "string", "minLength": 1},
        "blocks": {
          "type": "array",
          "items": {"$ref": "#/$defs/block"}
        }
      }
    },
    "reference": {
      "type": "object",
      "additionalProperties": false,
      "required": ["title", "url"],
      "properties": {
        "title": {"type": "string", "minLength": 1},
        "url": {"type": "string", "format": "uri"}
      }
    },
    "block": {
      "type": "object",
      "required": ["type"],
      "properties": {
        "type": {
          "type": "string",
          "enum": ["paragraph", "heading", "table", "image", "hyperlink", "page_break", "toc"]
        },
        "text": {"type": "string"},
        "level": {"type": "integer", "minimum": 1, "maximum": 3},
        "rows": {
          "type": "array",
          "items": {
            "type": "array",
            "items": {"type": "string"}
          }
        },
        "url": {"type": "string", "format": "uri"},
        "image_base64": {"type": "string"},
        "display_text": {"type": "string"},
        "width": {"type": "integer", "minimum": 1},
        "height": {"type": "integer", "minimum": 1},
        "levels": {"type": "string"}
      }
    }
  }
}
```

## Notes

- Large models should not invent facts when source material is incomplete.
- Missing critical facts should be reflected in `review_notes` or explicit `待补充` text.
- Schema validation is necessary but not sufficient. Document-type rules and style-review rules should run after schema validation.

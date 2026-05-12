# Draft To DOCX Mapping

This document defines how a validated `FormalDocumentDraftV1` draft is converted into requests for this service.

## Mapping Goal

The draft format is designed for large-model output.

The DOCX service request is designed for deterministic rendering.

This mapping layer keeps those two concerns separate.

## Top-Level Field Mapping

`FormalDocumentDraftV1` -> `generate_docx`

- `title` -> `title`
- `author` -> `author`
- `header_text` -> `header_text`
- `footer_page_number` -> `footer_page_number`

Suggested output filename:

- use `title + .docx`
- or use a caller-provided file name policy

## TOC Mapping

If `include_toc = true`, prepend the following block to `content`:

```json
{
  "type": "toc",
  "text": "目录",
  "levels": "1-3"
}
```

## Section Mapping

Each section becomes:

1. a `heading` block for the section title
2. each child block appended in order

Example draft section:

```json
{
  "title": "一、项目背景",
  "level": 1,
  "blocks": [
    {
      "type": "paragraph",
      "text": "为提升正式文档生产效率，拟建设统一平台。"
    }
  ]
}
```

Mapped output:

```json
[
  {
    "type": "heading",
    "text": "一、项目背景",
    "level": 1
  },
  {
    "type": "paragraph",
    "text": "为提升正式文档生产效率，拟建设统一平台。"
  }
]
```

## Appendix Mapping

Recommended mapping:

1. insert a `page_break`
2. insert appendix heading
3. append appendix blocks

Example:

```json
[
  { "type": "page_break" },
  { "type": "heading", "text": "附录A 术语说明", "level": 1 },
  { "type": "table", "rows": [["术语", "说明"], ["MCP", "模型上下文协议"]] }
]
```

## Reference Mapping

References should be mapped to `hyperlink` blocks.

Example:

```json
{
  "title": "相关资料",
  "url": "https://example.com/spec"
}
```

becomes:

```json
{
  "type": "hyperlink",
  "url": "https://example.com/spec",
  "display_text": "相关资料"
}
```

## Template Mapping

If the render route is template-first:

- `template_name` -> `template_name`
- `placeholders` -> `data`

Example:

```json
{
  "template_name": "business-letter.docx",
  "placeholders": {
    "recipient_name": "张三",
    "subject": "项目沟通说明"
  }
}
```

becomes:

```json
{
  "template_name": "business-letter.docx",
  "file_name": "项目沟通说明.docx",
  "data": {
    "recipient_name": "张三",
    "subject": "项目沟通说明"
  }
}
```

## Full Example

Draft input:

```json
{
  "schema_version": "1.0",
  "document_type": "project_proposal",
  "title": "智能文档平台建设方案",
  "author": "产品与研发联合组",
  "audience": "management",
  "tone": "formal",
  "language": "zh-CN",
  "header_text": "智能文档平台建设方案",
  "footer_page_number": true,
  "include_toc": true,
  "sections": [
    {
      "title": "一、项目背景",
      "level": 1,
      "blocks": [
        {
          "type": "paragraph",
          "text": "为提升正式文档生产效率，拟建设统一平台。"
        }
      ]
    }
  ]
}
```

Mapped render request:

```json
{
  "file_name": "智能文档平台建设方案.docx",
  "title": "智能文档平台建设方案",
  "author": "产品与研发联合组",
  "header_text": "智能文档平台建设方案",
  "footer_page_number": true,
  "content": [
    {
      "type": "toc",
      "text": "目录",
      "levels": "1-3"
    },
    {
      "type": "heading",
      "text": "一、项目背景",
      "level": 1
    },
    {
      "type": "paragraph",
      "text": "为提升正式文档生产效率，拟建设统一平台。"
    }
  ]
}
```

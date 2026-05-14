# LLM Integration Example

This document shows a practical end-to-end flow for combining a large model with this service.

## Goal

Generate a formal weekly report through four explicit stages:

1. LLM produces `FormalDocumentDraftV1`
2. Client validates the draft
3. Client routes to document generation
4. Client returns the generated `.docx` download URL

## Step 1: Ask The LLM For Draft JSON

Prompt the model with a strict JSON-only instruction.

Example output:

```json
{
  "schema_version": "1.0",
  "document_type": "weekly_report",
  "title": "研发中心第20周工作周报",
  "audience": "management",
  "tone": "formal",
  "language": "zh-CN",
  "footer_page_number": true,
  "sections": [
    {
      "title": "一、本期工作概述",
      "level": 1,
      "blocks": [
        {
          "type": "paragraph",
          "text": "本期围绕智能文档平台推进。"
        }
      ]
    },
    {
      "title": "二、已完成事项",
      "level": 1,
      "blocks": [
        {
          "type": "paragraph",
          "text": "已完成核心接口改造。"
        }
      ]
    },
    {
      "title": "三、当前进展",
      "level": 1,
      "blocks": [
        {
          "type": "paragraph",
          "text": "正在补充规则校验器。"
        }
      ]
    },
    {
      "title": "四、存在问题",
      "level": 1,
      "blocks": [
        {
          "type": "paragraph",
          "text": "模板库仍需补充。"
        }
      ]
    },
    {
      "title": "五、下阶段计划",
      "level": 1,
      "blocks": [
        {
          "type": "paragraph",
          "text": "继续扩展标准文档类型。"
        }
      ]
    },
    {
      "title": "六、需协调事项",
      "level": 1,
      "blocks": [
        {
          "type": "paragraph",
          "text": "需协调模板规范。"
        }
      ]
    }
  ],
  "review_notes": [
    "模板清单待补充"
  ]
}
```

## Step 2: Validate Draft

Call the validation endpoint first.

```bash
curl -X POST http://localhost:9103/api/v1/documents/validate-draft \
  -H 'Content-Type: application/json' \
  -d @draft.json
```

Example response:

```json
{
  "valid": true,
  "review_notes": [
    "模板清单待补充"
  ],
  "recommended_route": "structured"
}
```

For business letters, the service may return:

```json
{
  "valid": true,
  "recommended_route": "template",
  "recommended_template": "business-letter.docx"
}
```

## Step 3: Generate DOCX From Draft

If validation succeeds, call the generation endpoint.

```bash
curl -X POST http://localhost:9103/api/v1/documents/generate-from-draft \
  -H 'Content-Type: application/json' \
  -d @draft.json
```

Example response:

```json
{
  "file_name": "研发中心第20周工作周报-1778592333617346000.docx",
  "path": "runtime/files/研发中心第20周工作周报-1778592333617346000.docx",
  "download_url": "/api/v1/documents/files/研发中心第20周工作周报-1778592333617346000.docx",
  "mime_type": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
  "size_bytes": 9470,
  "review_notes": [
    "模板清单待补充"
  ],
  "route": "structured"
}
```

## Step 4: Download The File

```bash
curl -O http://localhost:9103/api/v1/documents/files/研发中心第20周工作周报-1778592333617346000.docx
```

## MCP Flow

Equivalent MCP tools:

- `validate_formal_document_draft`
- `generate_docx_from_draft`

Recommended orchestration:

1. Ask the LLM for `FormalDocumentDraftV1`
2. Call `validate_formal_document_draft`
3. If valid, call `generate_docx_from_draft`
4. Return the `download_url` to the user

## Practical Recommendation

For formal production systems, always validate first and generate second.

This keeps the LLM output reviewable, the render path deterministic, and the final `.docx` more stable.

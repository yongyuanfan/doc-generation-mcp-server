# Formal Document Examples

This document contains example `FormalDocumentDraftV1` payloads for the first three supported document types.

## Project Proposal Example

```json
{
  "schema_version": "1.0",
  "document_type": "project_proposal",
  "title": "智能文档平台建设方案",
  "author": "产品与研发联合组",
  "organization": "某某科技有限公司",
  "audience": "management",
  "tone": "formal",
  "language": "zh-CN",
  "header_text": "智能文档平台建设方案",
  "footer_page_number": true,
  "include_toc": true,
  "summary": "本方案说明平台建设背景、目标、实施路径及预期成效。",
  "sections": [
    {
      "id": "background",
      "title": "一、项目背景",
      "level": 1,
      "required": true,
      "blocks": [
        {
          "type": "paragraph",
          "text": "为提升正式文档生产效率、统一文档格式标准，并降低人工排版成本，拟建设统一的智能文档生成平台。"
        }
      ]
    },
    {
      "id": "goal",
      "title": "二、建设目标",
      "level": 1,
      "required": true,
      "blocks": [
        {
          "type": "table",
          "rows": [
            ["目标项", "说明"],
            ["效率提升", "减少重复写作与人工排版工作量"],
            ["标准统一", "统一标题层级、页眉页脚与表格样式"]
          ]
        }
      ]
    }
  ],
  "review_notes": [
    "预算数据待补充"
  ]
}
```

## Weekly Report Example

```json
{
  "schema_version": "1.0",
  "document_type": "weekly_report",
  "title": "研发中心第20周工作周报",
  "author": "研发中心",
  "audience": "management",
  "tone": "formal",
  "language": "zh-CN",
  "header_text": "研发中心工作周报",
  "footer_page_number": true,
  "include_toc": false,
  "sections": [
    {
      "title": "一、本期工作概述",
      "level": 1,
      "blocks": [
        {
          "type": "paragraph",
          "text": "本期围绕智能文档平台核心生成功能、模板管理能力及服务联调工作推进。"
        }
      ]
    },
    {
      "title": "二、已完成事项",
      "level": 1,
      "blocks": [
        {
          "type": "table",
          "rows": [
            ["事项", "状态"],
            ["DOCX服务重构", "已完成"],
            ["模板渲染能力", "已完成"]
          ]
        }
      ]
    },
    {
      "title": "三、存在问题",
      "level": 1,
      "blocks": [
        {
          "type": "paragraph",
          "text": "正式文档的业务模板库仍需进一步补充。"
        }
      ]
    }
  ]
}
```

## Business Letter Example

```json
{
  "schema_version": "1.0",
  "document_type": "business_letter",
  "title": "关于项目沟通安排的说明函",
  "author": "某某科技有限公司",
  "audience": "customer",
  "tone": "formal",
  "language": "zh-CN",
  "template_name": "business-letter.docx",
  "sections": [
    {
      "title": "一、发函背景",
      "level": 1,
      "blocks": [
        {
          "type": "paragraph",
          "text": "为进一步推进双方项目沟通与实施协同，现就近期工作安排说明如下。"
        }
      ]
    },
    {
      "title": "二、发函事项",
      "level": 1,
      "blocks": [
        {
          "type": "paragraph",
          "text": "建议双方于下周组织专项沟通会议，围绕需求确认、排期安排及责任分工进行集中讨论。"
        }
      ]
    }
  ],
  "placeholders": {
    "recipient_name": "某某客户单位",
    "sender_name": "某某科技有限公司"
  }
}
```

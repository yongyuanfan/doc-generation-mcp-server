#!/usr/bin/env bash

set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:9103}"

printf 'Checking health endpoint...\n'
curl --fail --silent --show-error "${BASE_URL}/healthz"

printf '\n\nChecking capabilities endpoint...\n'
curl --fail --silent --show-error "${BASE_URL}/api/v1/capabilities"

printf '\n\nChecking templates endpoint...\n'
curl --fail --silent --show-error "${BASE_URL}/api/v1/templates"

printf '\n\nChecking generate validation...\n'
generate_status="$(curl --silent --show-error --output /tmp/doc_generation_validation.json --write-out '%{http_code}' \
  -X POST "${BASE_URL}/api/v1/documents/generate" \
  -H 'Content-Type: application/json' \
  -d '{}')"
if [[ "${generate_status}" != "400" ]]; then
	printf 'Unexpected status for generate validation: %s\n' "${generate_status}"
	exit 1
fi
cat /tmp/doc_generation_validation.json

printf '\n\nChecking template validation...\n'
template_status="$(curl --silent --show-error --output /tmp/doc_template_validation.json --write-out '%{http_code}' \
  -X POST "${BASE_URL}/api/v1/documents/render-template" \
  -H 'Content-Type: application/json' \
  -d '{}')"
if [[ "${template_status}" != "400" ]]; then
	printf 'Unexpected status for template validation: %s\n' "${template_status}"
	exit 1
fi
cat /tmp/doc_template_validation.json

printf '\n\nChecking advanced generation...\n'
curl --fail --silent --show-error -X POST "${BASE_URL}/api/v1/documents/generate" \
  -H 'Content-Type: application/json' \
  -d '{
    "file_name": "smoke-advanced.docx",
    "header_text": "Smoke Header",
    "footer_page_number": true,
    "content": [
      {"type": "toc", "text": "Contents", "levels": "1-2"},
      {"type": "heading", "text": "Smoke", "level": 1},
      {"type": "paragraph", "text": "Advanced generation works."},
      {"type": "hyperlink", "url": "https://example.com", "display_text": "Example"},
      {"type": "page_break"}
    ]
  }'

printf '\n\nChecking draft validation...\n'
curl --fail --silent --show-error -X POST "${BASE_URL}/api/v1/documents/validate-draft" \
  -H 'Content-Type: application/json' \
  -d '{
    "schema_version": "1.0",
    "document_type": "weekly_report",
    "title": "研发中心第20周工作周报",
    "audience": "management",
    "tone": "formal",
    "language": "zh-CN",
    "footer_page_number": true,
    "sections": [
      {"title": "一、本期工作概述", "level": 1, "blocks": [{"type": "paragraph", "text": "本期围绕智能文档平台推进。"}]},
      {"title": "二、已完成事项", "level": 1, "blocks": [{"type": "paragraph", "text": "已完成核心接口改造。"}]},
      {"title": "三、当前进展", "level": 1, "blocks": [{"type": "paragraph", "text": "正在补充规则校验器。"}]},
      {"title": "四、存在问题", "level": 1, "blocks": [{"type": "paragraph", "text": "模板库仍需补充。"}]},
      {"title": "五、下阶段计划", "level": 1, "blocks": [{"type": "paragraph", "text": "继续扩展标准文档类型。"}]},
      {"title": "六、需协调事项", "level": 1, "blocks": [{"type": "paragraph", "text": "需协调模板规范。"}]}
    ]
  }'

printf '\n\nSmoke test completed.\n'

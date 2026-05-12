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

printf '\n\nSmoke test completed.\n'

# Agent Workflow

This document defines the recommended orchestration flow for combining a large model with this DOCX generation service to produce standard, formal documents.

## Goal

Separate the responsibilities clearly:

- The model writes and structures content
- Validation rules enforce document standards
- This service renders the final DOCX output

## Workflow Stages

### 1. Input Collection

Collect the following inputs from the user or application:

- intended document type or business goal
- audience
- source material
- language and tone constraints
- whether template-first output is preferred
- whether TOC or page numbers are required

Suggested internal state:

- `input_collected`

### 2. Document Type Classification

If the user did not explicitly choose a type, run a classification step.

Expected outputs:

- `document_type`
- `prefer_template`
- `required_sections`

Suggested internal state:

- `type_classified`

### 3. Initial Draft Generation

Prompt the model to output a `FormalDocumentDraftV1` JSON object.

Rules:

- output must be JSON only
- markdown is not allowed
- explanatory text is not allowed
- facts must not be invented

Suggested internal state:

- `draft_generated`

### 4. Review and Formalization

Run a second review pass, either with the same model or another model.

The review step may:

- improve formality
- unify terminology
- normalize section titles
- mark missing facts as `待补充`
- add notes to `review_notes`

The review step must not:

- change schema shape
- fabricate facts
- convert output to markdown

Suggested internal state:

- `draft_reviewed`

### 5. Validation

Run validations in this order:

1. JSON Schema validation against `FormalDocumentDraftV1`
2. Document-type validation from `document-type-rules.md`
3. Formal-style validation

If validation fails:

- regenerate the draft
- or ask the user for missing facts
- or downgrade to a partial draft with `待补充`

Suggested internal state:

- `draft_validated`

### 6. Render Routing

Choose the render path.

#### Template-first route

Use when:

- the document structure is highly standardized
- the document type is `business_letter`
- a matching template already exists

MCP tools:

- `list_docx_templates`
- `render_docx_template`

#### Structured-generation route

Use when:

- the document is long and content-heavy
- heading hierarchy matters
- content varies significantly per request

MCP tool:

- `generate_docx`

#### Hybrid route

Optional later extension.

Use when:

- a fixed front matter or formal layout is needed
- but a large dynamic body must also be inserted

Suggested internal state:

- `render_route_selected`

### 7. Render Request Conversion

Transform the validated `FormalDocumentDraftV1` JSON into the service request.

Examples:

- `title` -> `title`
- `author` -> `author`
- `header_text` -> `header_text`
- `footer_page_number` -> `footer_page_number`
- `include_toc` -> prepend a `toc` block
- `sections` -> `heading` + child blocks
- `references` -> `hyperlink` blocks

Suggested internal state:

- `render_request_built`

### 8. MCP Invocation

Invoke the selected tool.

Possible tools:

- `generate_docx`
- `render_docx_template`

Successful outputs should include:

- `file_name`
- `download_url`
- `mime_type`
- `size_bytes`

Suggested internal state:

- `docx_rendered`

## Failure Handling

### Schema Failure

If the draft is not valid JSON or fails schema validation:

- retry with stronger JSON-only instructions
- optionally provide the schema excerpt to the model

### Type Rule Failure

If the draft misses required sections:

- trigger a targeted repair prompt
- or ask the user for missing information

### Template Failure

If a required template is missing:

- fall back to structured generation
- log the missing template name for operational follow-up

### Render Failure

If the service cannot render the result:

- surface the error
- preserve the validated draft JSON for debugging and retry

## Prompt Layers

Recommended prompts:

1. classification prompt
2. draft generation prompt
3. review prompt

Each prompt has one job. Do not mix classification, drafting, and review into one unstable call unless latency is more important than quality.

## Operational Notes

- Store the reviewed draft JSON for auditability
- Keep the final render request for reproducibility
- Return `review_notes` to upstream systems even when rendering succeeds
- Prefer templates for format-critical communications
- Prefer structured generation for long and analytical documents

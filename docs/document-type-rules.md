# Document Type Rules

This document defines the standard structure, style, and validation rules for the first three supported formal document types.

## Common Rules

Applicable to all document types:

- Language should be `zh-CN`
- Tone should be `formal`
- Facts must not be invented when source material is incomplete
- Missing critical facts should be marked as `待补充`
- Avoid colloquial expressions, exaggerated marketing language, and generic AI filler phrases

Disallowed examples:

- `我觉得`
- `非常牛`
- `超级重要`
- `革命性提升`
- `接下来我将`
- `让我们来看看`

## `project_proposal`

### Use Cases

- Project proposal
- Project initiation document
- Construction plan
- Implementation recommendation document

### Recommended Rendering Mode

- Structured generation first
- Template-assisted mode optional

### Required Sections

1. 项目背景
2. 建设目标
3. 建设内容
4. 实施计划
5. 资源需求
6. 风险与保障措施
7. 预期成效
8. 结论与建议

### Optional Sections

- 政策依据
- 技术架构
- 投资估算
- 附录

### Default Layout Policy

- `include_toc = true`
- `footer_page_number = true`
- `header_text` required

### Validation Rules

- At least 6 top-level sections should be present
- At least 1 table should appear in the document
- A document containing only broad conclusions without implementation detail is invalid
- Missing budget, schedule, or owner information should be marked as `待补充`

### Writing Style

- Objective
- Cautious
- Complete
- Decision-support oriented

## `weekly_report`

### Use Cases

- Weekly report
- Monthly report
- Periodic progress report

### Recommended Rendering Mode

- Structured generation first

### Required Sections

1. 本期工作概述
2. 已完成事项
3. 当前进展
4. 存在问题
5. 下阶段计划
6. 需协调事项

### Optional Sections

- 风险提示
- 关键指标
- 资源需求

### Default Layout Policy

- `include_toc = false`
- `footer_page_number = true`
- `header_text` recommended

### Validation Rules

- The draft must include both `存在问题` and `下阶段计划`
- A report containing only completed work and no risk or blocker content is incomplete
- At least 1 summary table is recommended

### Writing Style

- Concise
- Management-oriented
- Fact-based
- Action-focused

## `business_letter`

### Use Cases

- Business letter
- Formal notice letter
- External communication letter
- Official explanation letter

### Recommended Rendering Mode

- Template-first rendering

### Required Sections

1. 发函背景
2. 发函事项
3. 具体说明
4. 后续安排
5. 联系方式或落款

### Optional Sections

- 附件说明
- 时间节点
- 回复要求

### Default Layout Policy

- `include_toc = false`
- `footer_page_number = false`
- `header_text` optional
- Template mode preferred

### Validation Rules

- Deep nested headings should be avoided
- Each section should stay concise
- The document must include a clear request, explanation, or next action

### Writing Style

- Formal
- Polite
- Concise
- Non-promotional

## Missing Information Handling

When source material is incomplete:

- Use `待补充` for required facts that are missing
- Record unresolved issues in `review_notes`
- Do not fabricate schedules, amounts, organizational roles, names, or commitments

## Rule Application Order

1. Validate against `FormalDocumentDraftV1` schema
2. Validate against document-type rules in this document
3. Run formal-style review
4. Convert into the service render request

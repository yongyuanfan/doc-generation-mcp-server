package shared

import "github.com/yong/doc-generation-mcp-server/internal/formaldoc"

func WeeklyReportDraft() formaldoc.Draft {
	return formaldoc.Draft{
		SchemaVersion:    formaldoc.SchemaVersion,
		DocumentType:     formaldoc.DocumentTypeWeeklyReport,
		Title:            "研发中心第20周工作周报",
		Audience:         "management",
		Tone:             "formal",
		Language:         "zh-CN",
		FooterPageNumber: true,
		Sections: []formaldoc.Section{
			{Title: "一、本期工作概述", Level: 1, Blocks: []formaldoc.Block{{Type: "paragraph", Text: "本期围绕智能文档平台推进。"}}},
			{Title: "二、已完成事项", Level: 1, Blocks: []formaldoc.Block{{Type: "paragraph", Text: "已完成核心接口改造。"}}},
			{Title: "三、当前进展", Level: 1, Blocks: []formaldoc.Block{{Type: "paragraph", Text: "正在补充规则校验器。"}}},
			{Title: "四、存在问题", Level: 1, Blocks: []formaldoc.Block{{Type: "paragraph", Text: "模板库仍需补充。"}}},
			{Title: "五、下阶段计划", Level: 1, Blocks: []formaldoc.Block{{Type: "paragraph", Text: "继续扩展标准文档类型。"}}},
			{Title: "六、需协调事项", Level: 1, Blocks: []formaldoc.Block{{Type: "paragraph", Text: "需协调模板规范。"}}},
		},
		ReviewNotes: []string{"模板清单待补充"},
	}
}

func BusinessLetterDraft() formaldoc.Draft {
	return formaldoc.Draft{
		SchemaVersion: formaldoc.SchemaVersion,
		DocumentType:  formaldoc.DocumentTypeBusinessLetter,
		Title:         "项目沟通说明函",
		Audience:      "customer",
		Tone:          "formal",
		Language:      "zh-CN",
		Sections: []formaldoc.Section{
			{Title: "一、发函背景", Level: 1, Blocks: []formaldoc.Block{{Type: "paragraph", Text: "为进一步推进项目沟通，现说明如下。"}}},
			{Title: "二、发函事项", Level: 1, Blocks: []formaldoc.Block{{Type: "paragraph", Text: "建议双方安排专项沟通会议。"}}},
			{Title: "三、具体说明", Level: 1, Blocks: []formaldoc.Block{{Type: "paragraph", Text: "会议议题包括需求确认、排期安排与责任分工。"}}},
			{Title: "四、后续安排", Level: 1, Blocks: []formaldoc.Block{{Type: "paragraph", Text: "请于收到本函后确认可行时间。"}}},
		},
		Placeholders: map[string]string{
			"organization":   "某某科技有限公司",
			"recipient_name": "张三",
			"sender_name":    "某某科技有限公司",
			"summary":        "请确认近期项目沟通安排。",
			"date":           "2026-05-12",
		},
	}
}

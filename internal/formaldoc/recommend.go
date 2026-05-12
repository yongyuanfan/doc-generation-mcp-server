package formaldoc

func RecommendedRoute(d Draft, templateMap map[string]string) string {
	if d.TemplateName != "" {
		return RouteTemplate
	}
	switch d.DocumentType {
	case DocumentTypeBusinessLetter:
		if RecommendedTemplate(d, templateMap) != "" {
			return RouteTemplate
		}
		return RouteStructured
	default:
		return RouteStructured
	}
}

func RecommendedTemplate(d Draft, templateMap map[string]string) string {
	if d.TemplateName != "" {
		return d.TemplateName
	}
	if value, ok := templateMap[d.DocumentType]; ok && value != "" {
		return value
	}
	switch d.DocumentType {
	case DocumentTypeBusinessLetter:
		return "business-letter.docx"
	default:
		return ""
	}
}

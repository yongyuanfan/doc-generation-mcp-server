package formaldoc

func RecommendedRoute(d Draft) string {
	if d.TemplateName != "" {
		return RouteTemplate
	}
	switch d.DocumentType {
	case DocumentTypeBusinessLetter:
		return RouteTemplate
	default:
		return RouteStructured
	}
}

func RecommendedTemplate(d Draft) string {
	if d.TemplateName != "" {
		return d.TemplateName
	}
	switch d.DocumentType {
	case DocumentTypeBusinessLetter:
		return "business-letter.docx"
	default:
		return ""
	}
}

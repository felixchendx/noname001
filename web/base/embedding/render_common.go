package embedding

import (
	"strings"
	"html/template"
)

func (tmpl *BaseTemplating) RenderCommon(name string, dataMap map[string]any) (template.HTML, error) {
	dataMap["cacheBusterVar"] = tmpl.cacheBusterVar
	dataMap["appVersion"]     = tmpl.appVersion

	b := new(strings.Builder)

	err := tmpl.baseHTMLTemplate.ExecuteTemplate(b, name, dataMap)
	if err != nil {
		return template.HTML(""), err
	}

	return template.HTML(b.String()), nil
}

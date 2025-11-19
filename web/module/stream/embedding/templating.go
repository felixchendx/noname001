package embedding

import (
	"embed"
	"html/template"
	"maps"

	"noname001/logging"

	webBase "noname001/web/base"
	baseEmbedding "noname001/web/base/embedding"
)

var (
	//go:embed all:template
	htmlTemplateFS embed.FS
)

type ModuleTemplatingParams struct {
	Logger     *logging.WrappedLogger
	LogPrefix  string

	BaseBundle *webBase.BaseBundle
}

type ModuleTemplating struct {
	logger       *logging.WrappedLogger
	logPrefix    string

	baseBundle   *webBase.BaseBundle

	htmlTemplate *template.Template
}

func NewModuleTemplating(params *ModuleTemplatingParams) (*ModuleTemplating, error) {
	var err error

	var tmpl *ModuleTemplating = &ModuleTemplating{}
	tmpl.logger = params.Logger
	tmpl.logPrefix = params.LogPrefix + ".tmpl"
	tmpl.baseBundle = params.BaseBundle

	pattern := []string{
		"template/default/*.html.tmpl",
	}

	tmpl.htmlTemplate, err = template.
		New("module_template").
		Funcs(moduleTemplatingFunctions()).
		ParseFS(htmlTemplateFS, pattern...)
	if err != nil {
		tmpl.logger.Errorf("%s: new module template err %s", tmpl.logPrefix, err.Error())
		return nil, err
	}

	return tmpl, nil
}

func moduleTemplatingFunctions() template.FuncMap {
	moduleTemplatingFunctions := template.FuncMap{}

	maps.Insert(
		moduleTemplatingFunctions,
		maps.All(baseEmbedding.BaseTemplateFunctions()),
	)

	return moduleTemplatingFunctions
}

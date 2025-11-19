package embedding

import (
	"embed"
	"html/template"
	"maps"

	"noname001/logging"

	webBase "noname001/web/base"
	baseEmbedding "noname001/web/base/embedding"

	wallService "noname001/app/module/feature/wall/service"
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

	tmpl := &ModuleTemplating{}	
	tmpl.logger = params.Logger
	tmpl.logPrefix = params.LogPrefix + ".tmpl"

	tmpl.baseBundle = params.BaseBundle

	patterns := []string{
		"template/default/*.html.tmpl",
		// "template/template001/*.html.tmpl",
	}

	tmpl.htmlTemplate, err = template.
		New("module_template").
		Funcs(moduleTemplatingFunctions()).
		ParseFS(htmlTemplateFS, patterns...)

	if err != nil {
		tmpl.logger.Errorf("%s: new module template err %s", tmpl.logPrefix, err.Error())
		return nil, err
	}

	return tmpl, nil
}

func moduleTemplatingFunctions() template.FuncMap {
	moduleTemplatingFunctions := template.FuncMap{
		// credits: https://stackoverflow.com/questions/30065203/golang-html-templating-range-limit
		//          https://stackoverflow.com/a/30066033
		// TODO: FUCKIN EWW, move to base, but need reflections...
		"wtf_wallItemSlice": func(s []*wallService.WallItemDE, b, e int) []*wallService.WallItemDE { return s[b:e] },
	}

	maps.Insert(
		moduleTemplatingFunctions,
		maps.All(baseEmbedding.BaseTemplateFunctions()),
	)

	return moduleTemplatingFunctions
}

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

type UserTemplatingParams struct {
	Logger     *logging.WrappedLogger
	LogPrefix  string

	BaseBundle *webBase.BaseBundle
}
type UserTemplating struct {
	logger       *logging.WrappedLogger
	logPrefix    string

	baseBundle   *webBase.BaseBundle

	htmlTemplate *template.Template
}

func NewUserTemplating(params *UserTemplatingParams) (*UserTemplating, error) {
	var err error
	
	tmpl := &UserTemplating{}
	tmpl.logger = params.Logger
	tmpl.logPrefix = params.LogPrefix + ".tmpl"

	tmpl.baseBundle = params.BaseBundle

	patterns := []string{
		"template/*.html.tmpl",
		// "template/ehehe/*.html.tmpl",
	}

	tmpl.htmlTemplate, err = template.
		New("user_template").
		Funcs(userTemplatingFunctions()).
		ParseFS(htmlTemplateFS, patterns...)
	
	if err != nil {
		tmpl.logger.Errorf("%s: new user templating err %s", tmpl.logPrefix, err.Error())
		return nil, err
	}

	return tmpl, nil
}

func userTemplatingFunctions() template.FuncMap {
	userTemplatingFunctions := template.FuncMap{
	}

	maps.Insert(
		userTemplatingFunctions,
		maps.All(baseEmbedding.BaseTemplateFunctions()),
	)

	return userTemplatingFunctions
}

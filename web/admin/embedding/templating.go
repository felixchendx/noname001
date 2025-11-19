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

type AdminTemplatingParams struct {
	Logger     *logging.WrappedLogger
	LogPrefix  string

	BaseBundle *webBase.BaseBundle
}
type AdminTemplating struct {
	logger       *logging.WrappedLogger
	logPrefix    string

	baseBundle   *webBase.BaseBundle

	htmlTemplate *template.Template
	// textTemplate
}

func NewAdminTemplating(params *AdminTemplatingParams) (*AdminTemplating, error) {
	var err error
	
	tmpl := &AdminTemplating{}
	tmpl.logger = params.Logger
	tmpl.logPrefix = params.LogPrefix + ".tmpl"

	tmpl.baseBundle = params.BaseBundle

	patterns := []string{
		"template/*.html.tmpl",
		// "template/ehehe/*.html.tmpl",
	}

	tmpl.htmlTemplate, err = template.
		New("admin_template").
		Funcs(adminTemplatingFunctions()).
		ParseFS(htmlTemplateFS, patterns...)
	
	if err != nil {
		tmpl.logger.Errorf("%s: new admin templating err %s", tmpl.logPrefix, err.Error())
		return nil, err
	}

	return tmpl, nil
}

func adminTemplatingFunctions() template.FuncMap {
	adminTemplatingFunctions := template.FuncMap{
	}

	maps.Insert(
		adminTemplatingFunctions,
		maps.All(baseEmbedding.BaseTemplateFunctions()),
	)

	return adminTemplatingFunctions
}

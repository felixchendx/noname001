package embedding

import (
	"embed"
	"html/template"

	"noname001/logging"

	"noname001/web/base/auth"
	"noname001/web/base/navigation"
)

var (
	//go:embed all:template
	htmlTemplateFS embed.FS
)

type BaseTemplatingParams struct {
	Logger    *logging.WrappedLogger
	LogPrefix string

	AuthProvider *auth.AuthProvider
	Navi         *navigation.BaseNavi

	CacheBusterVar string
	AppVersion     string
}
type BaseTemplating struct {
	logger           *logging.WrappedLogger
	logPrefix        string

	authProvider *auth.AuthProvider
	navi         *navigation.BaseNavi

	baseHTMLTemplate *template.Template

	cacheBusterVar string
	appVersion     string
}

func NewBaseTemplating(params *BaseTemplatingParams) (*BaseTemplating, error) {
	var err error

	tmpl := &BaseTemplating{}
	tmpl.logger = params.Logger
	tmpl.logPrefix = params.LogPrefix + ".tmpl"

	tmpl.authProvider = params.AuthProvider
	tmpl.navi = params.Navi

	patterns := []string{
		"template/common/*.html.tmpl",
		"template/default/*.html.tmpl",
	}

	tmpl.baseHTMLTemplate, err = template.
		New("template_base").
		Funcs(BaseTemplateFunctions()).
		ParseFS(htmlTemplateFS, patterns...)

	if err != nil {
		tmpl.logger.Errorf("%s: new base templating err %s", tmpl.logPrefix, err.Error())
		return nil, err
	}

	tmpl.cacheBusterVar = params.CacheBusterVar
	tmpl.appVersion     = params.AppVersion

	return tmpl, nil
}

func BaseTemplateFunctions() (template.FuncMap) {
	return template.FuncMap{
		"btf_add": func(i1, i2 int) int { return i1 + i2 },
		"btf_isNil": func(v any) bool { return v == nil },
		"btf_notNil": func(v any) bool { return v != nil },
	}
}

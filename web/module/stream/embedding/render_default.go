package embedding

import (
	"html/template"
	"strings"

	"github.com/valyala/fasthttp"

	baseEmbedding "noname001/web/base/embedding"
)

type PageData_default = baseEmbedding.PageData_default

func (tmpl *ModuleTemplating) NewPageData_default() (*PageData_default) {
	return tmpl.baseBundle.Templating.NewPageData_default()
}

func (tmpl *ModuleTemplating) RenderContent_default(name string, pageData *PageData_default, ctx *fasthttp.RequestCtx) (template.HTML, error) {
	var dataMap map[string]any
	if pageData != nil {
		pageData.TopnavData = &baseEmbedding.TopnavData_default{
			WebUser: tmpl.baseBundle.Auth.CurrentWebUser(ctx),
		}
		pageData.SidenavData = &baseEmbedding.SidenavData_default{
			NavLinks: tmpl.baseBundle.Navi.MappedSiteNavLinks(ctx),
		}

		dataMap = pageData.AsDataMap()
	}

	contentSB := new(strings.Builder)
	err := tmpl.htmlTemplate.ExecuteTemplate(contentSB, name, dataMap["contentData"])
	if err != nil {
		tmpl.logger.Errorf("%s: RenderContent_default contentErr %s", tmpl.logPrefix, err.Error())
		return template.HTML(""), err
	}

	dataMap["contentHTML"] = template.HTML(contentSB.String())
	pageHTML, err := tmpl.baseBundle.Templating.RenderPage_default(dataMap)
	if err != nil {
		return template.HTML(""), err
	}

	return pageHTML, nil
}

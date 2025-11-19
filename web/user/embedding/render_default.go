package embedding

import (
	"html/template"
	"strings"

	"github.com/valyala/fasthttp"

	baseEmbedding "noname001/web/base/embedding"
)

type PageData_default = baseEmbedding.PageData_default

func (tmpl *UserTemplating) NewPageData() (*PageData_default) {
	return tmpl.baseBundle.Templating.NewPageData_default()
}

func (tmpl *UserTemplating) RenderContent_default(name string, pageData *PageData_default, ctx *fasthttp.RequestCtx) (template.HTML, error) {
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
	contentErr := tmpl.htmlTemplate.ExecuteTemplate(contentSB, name, dataMap["contentData"])
	if contentErr != nil {
		tmpl.logger.Errorf("%s: RenderContent_default contentErr %s", tmpl.logPrefix, contentErr.Error())
		return template.HTML(""), contentErr
	}

	dataMap["contentHTML"] = template.HTML(contentSB.String())
	pageHTML, pageErr := tmpl.baseBundle.Templating.RenderPage_default(dataMap)
	if pageErr != nil {
		return template.HTML(""), pageErr
	}

	return pageHTML, nil
}

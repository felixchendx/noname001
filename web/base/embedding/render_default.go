package embedding

import (
	"html/template"
	"strings"

	"github.com/valyala/fasthttp"

	"noname001/app/base/messaging"

	webConstant "noname001/web/constant"
	"noname001/web/base/auth"
)


// TODO: default template todo list
// - hide / show sidenav https://preview.tabler.io/offcanvas.html
// - breadcrumb

func (tmpl *BaseTemplating) RenderContent_default(name string, pageData *PageData_default, ctx *fasthttp.RequestCtx) (template.HTML, error) {
	var dataMap map[string]any
	if pageData != nil {
		pageData.TopnavData = &TopnavData_default{
			WebUser: tmpl.authProvider.CurrentWebUser(ctx),
		}
		pageData.SidenavData = &SidenavData_default{
			NavLinks: tmpl.navi.MappedSiteNavLinks(ctx),
		}

		dataMap = pageData.AsDataMap()
	}

	contentSB := new(strings.Builder)
	contentErr := tmpl.baseHTMLTemplate.ExecuteTemplate(contentSB, name, dataMap["contentData"])
	if contentErr != nil {
		tmpl.logger.Errorf("%s: RenderContent_default contentErr %s", tmpl.logPrefix, contentErr.Error())
		return template.HTML(""), contentErr
	}

	dataMap["contentHTML"] = template.HTML(contentSB.String())

	pageHTML, pageErr := tmpl.RenderPage_default(dataMap)
	if pageErr != nil {
		return template.HTML(""), pageErr
	}

	return pageHTML, nil
}

func (tmpl *BaseTemplating) RenderPage_default(dataMap map[string]any) (template.HTML, error) {
	dataMap["cacheBusterVar"] = tmpl.cacheBusterVar
	dataMap["appVersion"]     = tmpl.appVersion

	pageSB := new(strings.Builder)
	pageErr := tmpl.baseHTMLTemplate.ExecuteTemplate(pageSB, "d-page.html.tmpl", dataMap)
	if pageErr != nil {
		tmpl.logger.Errorf("%s: RenderPage_default pageErr %s", tmpl.logPrefix, pageErr.Error())
		return template.HTML(""), pageErr
	}

	return template.HTML(pageSB.String()), nil
}

func (tmpl *BaseTemplating) NewPageData_default() (*PageData_default) {
	return &PageData_default{
		// TODO: review other sane defaults

		Theme: webConstant.TBLR_THEME_LIGHT,
		RenderSidenav: true,
		RenderTopnav: true,
		RenderFooter: true,
	}
}

type PageData_default struct {
	Title         string
	Messages      *messaging.Messages
	ContentData   map[string]any

	TopnavData    *TopnavData_default
	SidenavData   *SidenavData_default

	Theme         string // temp ?
	RenderSidenav bool // temp ?
	RenderTopnav  bool // temp ?
	RenderFooter  bool // temp ?
	ExtraCssLinks []string
	ExtraJsLinks  []string
}
type TopnavData_default struct {
	WebUser *auth.WebUser
}
type SidenavData_default struct {
	NavLinks map[string]map[string]any
}

func (dataStruct *PageData_default) AsDataMap() (map[string]any) {
	var messages *messaging.Messages
	
	dataMap := map[string]any{
		"title": dataStruct.Title,
		"messages": messages,
		"contentData": map[string]any{},
		"topnavData": map[string]any{},
		"sidenavData": map[string]any{},
		"theme": dataStruct.Theme,
		"renderSidenav": dataStruct.RenderSidenav,
		"renderTopnav": dataStruct.RenderTopnav,
		"renderFooter": dataStruct.RenderFooter,
		"extraCssLinks": dataStruct.ExtraCssLinks,
		"extraJsLinks": dataStruct.ExtraJsLinks,
	}

	if _messages := dataStruct.Messages; _messages != nil && _messages.HasMessage() {
		dataMap["messages"] = _messages
	}
	if contentData := dataStruct.ContentData; contentData != nil {
		dataMap["contentData"] = contentData
	}
	if topnavData := dataStruct.TopnavData; topnavData != nil {
		dataMap["topnavData"] = topnavData.AsDataMap()
	}
	if sidenavData := dataStruct.SidenavData; sidenavData != nil {
		dataMap["sidenavData"] = sidenavData.AsDataMap()
	}

	return dataMap
}
func (dataStruct *TopnavData_default) AsDataMap() (map[string]any) {
	dataMap := map[string]any{
		"webuser": map[string]any{},
	}

	if webUser := dataStruct.WebUser; webUser != nil {
		dataMap["webuser"] = webUser.AsDataMap()
	}

	return dataMap
}
func (dataStruct *SidenavData_default) AsDataMap() (map[string]any) {
	dataMap := map[string]any{
		"nav_links": make(map[string]map[string]any),
	}

	if navLinks := dataStruct.NavLinks; navLinks != nil {
		dataMap["nav_links"] = navLinks
	}

	return dataMap
}

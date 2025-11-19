package auth

import (
	"github.com/valyala/fasthttp"

	"noname001/app/sys"
)

func (authProvider *AuthProvider) InjectWebSession(ctx *fasthttp.RequestCtx) {
	webSess := &WebSession{}
	webSess.webSessionID = authProvider.cookieStore.GetWebSessionID(ctx)
	webSess.internalSessionID = authProvider.cookieStore.GetInternalSessionID(ctx)

	webSess.sysSession = sys.Bundle.Service.FindSession(webSess.internalSessionID)

	ctx.SetUserValue("MuhWebSession", webSess)
}

func (authProvider *AuthProvider) RetrieveWebSession(ctx *fasthttp.RequestCtx) (*WebSession) {
	vany := ctx.UserValue("MuhWebSession")
	assertion, ok := vany.(*WebSession)
	_ = ok
	return assertion
}

func (authProvider *AuthProvider) TempCurrentUsername(ctx *fasthttp.RequestCtx) (string) {
	webSess := authProvider.RetrieveWebSession(ctx)
	if webSess.sysSession != nil && webSess.sysSession.User != nil {
		return webSess.sysSession.User.Username
	}

	return ""
}

type WebSession struct {
	webSessionID      string
	internalSessionID string

	sysSession *sys.SysSession
}

func (webSess *WebSession) IsSessionExpired() (bool) {
	return webSess.internalSessionID != "" && webSess.sysSession == nil
}

func (webSess *WebSession) IsLoggedIn()  (bool) { return webSess.sysSession != nil }
func (webSess *WebSession) IsLoggedOut() (bool) { return webSess.sysSession == nil }

func (webSess *WebSession) IsSuperadmin() (bool) {
	return webSess.sysSession != nil && webSess.sysSession.IsSuperadmin()
}

func (webSess *WebSession) IsAdmin() (bool) {
	return webSess.sysSession != nil && webSess.sysSession.IsAdmin()
}

func (webSess *WebSession) HasAdminAuthorization() (bool) {
	return webSess.sysSession != nil && webSess.sysSession.HasAdminAuthorization()
}
func (webSess *WebSession) DoesNotHaveAdminAuthorization() (bool) {
	return webSess.sysSession != nil && webSess.sysSession.DoesNotHaveAdminAuthorization()
}

func (webSess *WebSession) IsOperator() (bool) {
	return webSess.sysSession != nil && webSess.sysSession.IsOperator()
}

func (webSess *WebSession) IsViewer() (bool) {
	return webSess.sysSession != nil && webSess.sysSession.IsViewer()
}

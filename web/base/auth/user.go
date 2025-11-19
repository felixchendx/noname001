package auth

import (
	"github.com/valyala/fasthttp"
	
	"noname001/app/base/messaging"
	
	"noname001/app/sys"
)

func (authProvider *AuthProvider) WebLogin(ctx *fasthttp.RequestCtx, usernm, passwd string) (*WebUser, *messaging.Messages) {
	sysSession, messages := sys.Bundle.Service.WebLogin(usernm, passwd)
	if messages.HasError() {
		// TODO: temp
		// temp logging for cloud exposed setup, at least until further guard implementation such as reactive ip blacklist, or event triggers
		authProvider.logger.Errorf("%s: WebLogin FAILED usernm '%s'", authProvider.logPrefix, usernm)
		return nil, messages
	}

	authProvider.cookieStore.SetInternalSessionID(ctx, sysSession.ID())

	webUser := authProvider.CurrentWebUser(ctx)

	return webUser, messages
}

func (authProvider *AuthProvider) WebLogout(ctx *fasthttp.RequestCtx) {
	isid := authProvider.cookieStore.GetInternalSessionID(ctx)

	sys.Bundle.Service.WebLogout(isid)

	authProvider.cookieStore.DestroySessionCookie(ctx)
}

func (authProvider *AuthProvider) CurrentWebUser(ctx *fasthttp.RequestCtx) (*WebUser) {
	isid := authProvider.cookieStore.GetInternalSessionID(ctx)
	if isid == "" {
		return authProvider.anonymousWebUser()
	}

	sysSession := sys.Bundle.Service.FindSession(isid)
	if sysSession == nil {
		return authProvider.anonymousWebUser()
	}

	webUser := &WebUser{
		// ID userID necessary ?
		Name: sysSession.User.Username,
		Role: string(sysSession.User.SimpleRole),

		isid: isid,
	}

	return webUser
}

func (authProvider *AuthProvider) ChangePassword(ctx *fasthttp.RequestCtx, oldPass, newPass string) (*messaging.Messages) {
	webSess := authProvider.RetrieveWebSession(ctx)

	messages := sys.Bundle.Service.WebChangePassword(webSess.internalSessionID, oldPass, newPass)

	return messages
}

func (authProvider *AuthProvider) anonymousWebUser() (*WebUser) {
	return &WebUser{
		// ID: "",
		Name: "anon",
		Role: "anonymous",

		isid: "",
	}
}

type WebUser struct {
	// ID         string
	Name string
	Role string

	isid string
}
func (wu *WebUser) AsDataMap() (map[string]any) {
	return map[string]any{
		"name": wu.Name,
		"role": wu.Role,
	}
}

type Authn struct {

}


type Authz struct {

}

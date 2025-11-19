package auth

import (
	"github.com/valyala/fasthttp"
)

func (authProvider *AuthProvider) IsLoggedIn(ctx *fasthttp.RequestCtx) (bool) {
	return authProvider.RetrieveWebSession(ctx).IsLoggedIn()
}
func (authProvider *AuthProvider) IsLoggedOut(ctx *fasthttp.RequestCtx) (bool) {
	return authProvider.RetrieveWebSession(ctx).IsLoggedOut()
}

func (authProvider *AuthProvider) IsSuperadmin(ctx *fasthttp.RequestCtx) (bool) {
	return authProvider.RetrieveWebSession(ctx).IsSuperadmin()
}
func (authProvider *AuthProvider) IsAdmin(ctx *fasthttp.RequestCtx) (bool) {
	return authProvider.RetrieveWebSession(ctx).IsAdmin()
}
func (authProvider *AuthProvider) HasAdminAuthorization(ctx *fasthttp.RequestCtx) (bool) {
	return authProvider.RetrieveWebSession(ctx).HasAdminAuthorization()
}
func (authProvider *AuthProvider) DoesNotHaveAdminAuthorization(ctx *fasthttp.RequestCtx) (bool) {
	return authProvider.RetrieveWebSession(ctx).DoesNotHaveAdminAuthorization()
}

func (authProvider *AuthProvider) IsOperator(ctx *fasthttp.RequestCtx) (bool) {
	return authProvider.RetrieveWebSession(ctx).IsOperator()
}
func (authProvider *AuthProvider) IsViewer(ctx *fasthttp.RequestCtx) (bool) {
	return authProvider.RetrieveWebSession(ctx).IsViewer()
}

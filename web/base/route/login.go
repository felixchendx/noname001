package route

import (
	"github.com/valyala/fasthttp"

	webConstant "noname001/web/constant"
)

func (rh *BaseRouteHandler) RouteLogin(ctx *fasthttp.RequestCtx) {
	if rh.authProvider.IsLoggedIn(ctx) {
		rh.RedirectToIndex(ctx)
		return
	}

	switch string(ctx.Method()) {
	case fasthttp.MethodGet : rh.renderLogin(ctx)
	case fasthttp.MethodPost: rh.doLogin(ctx)
	default                 : rh.Route404(ctx)
	}
}

func (rh *BaseRouteHandler) renderLogin(ctx *fasthttp.RequestCtx) {
	var (
		webSess = rh.authProvider.RetrieveWebSession(ctx)

		flash        = rh.flashStore.GetFlashBundle(ctx)
		prevMessages = flash.Prev.Messages

		infoMsg, errorMsg = "", ""
	)


	if webSess.IsSessionExpired() {
		infoMsg = "Previous session has expired."
	}

	if prevMessages.HasError() {
		errorMsg = prevMessages.FirstErrorMessageDescription()
	}

	rendered, renderErr := rh.templating.RenderCommon("p-login.html.tmpl", map[string]any{
		"_form_login": flash.Prev.Data.Get("_form_login"),
		"_err_msg": errorMsg,
		"_info_msg": infoMsg,
	})
	if renderErr != nil {
		rh.Route500(ctx, renderErr)
		return
	}

	ctx.SetContentType(webConstant.CONTENT_TYPE_HTML)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte(rendered))
}

func (rh *BaseRouteHandler) doLogin(ctx *fasthttp.RequestCtx) {
	var (
		flash = rh.flashStore.GetFlashBundle(ctx)
	)

	usernm := ctx.FormValue("usernm")
	passwd := ctx.FormValue("passwd")

	_, messages := rh.authProvider.WebLogin(ctx, string(usernm), string(passwd))
	if messages.HasError() {
		flash.Next.Messages.Append(messages)
		flash.Next.Data.Set("_form_login", map[string]any{
			"usernm": string(usernm),
			"passwd": string(passwd),
		})
		rh.RedirectToLogin(ctx)

		// TODO: failure timeout per IP
		return
	}

	flash.Next.Messages.Append(messages)
	rh.RedirectToIndex(ctx)
}

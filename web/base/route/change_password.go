package route

import (
	"github.com/valyala/fasthttp"

	webConstant "noname001/web/constant"
)

func (rh *BaseRouteHandler) RouteChangePassword(ctx *fasthttp.RequestCtx) {
	if rh.authProvider.IsLoggedOut(ctx) {
		rh.RedirectToIndex(ctx)
		return
	}

	switch string(ctx.Method()) {
	case fasthttp.MethodGet : rh.renderChangePassword(ctx)
	case fasthttp.MethodPost: rh.doChangePassword(ctx)
	default                 : rh.Route404(ctx)
	}
}

type ChangePasswordForm struct {
	OldPassword string
	NewPassword string
}

func (rh *BaseRouteHandler) renderChangePassword(ctx *fasthttp.RequestCtx) {
	var (
		flash = rh.flashStore.GetFlashBundle(ctx)

		currInput    = ChangePasswordForm{}
	)

	prevError := flash.Prev.Data.Get("prev_error")
	for k, vany := range prevError {
		switch k {
		case "change_password":
			assertion, ok := vany.(ChangePasswordForm)
			if ok { currInput = assertion }
		}
	}

	pageData := rh.templating.NewPageData_default()
	pageData.Title = "Change Password"
	pageData.Messages = flash.Prev.Messages
	pageData.ContentData = map[string]any{
		"_title": "Change Password",
		"_curr_username": rh.authProvider.TempCurrentUsername(ctx),
		"_form" : currInput,
	}

	renderOut, renderErr := rh.templating.RenderContent_default("d-content--change-password.html.tmpl", pageData, ctx)
	if renderErr != nil {
		rh.Route500(ctx, renderErr)
		return
	}
	
	ctx.SetContentType(webConstant.CONTENT_TYPE_HTML)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte(renderOut))
}

func (rh *BaseRouteHandler) doChangePassword(ctx *fasthttp.RequestCtx) {
	var (
		flash = rh.flashStore.GetFlashBundle(ctx)

		formData = ChangePasswordForm{}
	)

	ctx.PostArgs().VisitAll(func(k, v []byte) {
		sk, sv := string(k), string(v)
		switch sk {
		case "old_password": formData.OldPassword = sv
		case "new_password": formData.NewPassword = sv
		}
	})

	messages := rh.authProvider.ChangePassword(ctx, formData.OldPassword, formData.NewPassword)
	if messages.HasError() {
		flash.Next.Data.Set("prev_error", map[string]any{
			"change_password": formData,
		})
	}
	flash.Next.Messages.Append(messages)

	ctx.Redirect("/change-password", fasthttp.StatusFound)
}

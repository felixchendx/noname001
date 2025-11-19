package route

import (
	"fmt"

	"github.com/valyala/fasthttp"

	webConstant "noname001/web/constant"

	"noname001/app/sys"
)

func (rh *AdminRouteHandler) renderUserDetail(ctx *fasthttp.RequestCtx) {
	var (
		flash           = rh.baseBundle.Flash.GetFlashBundle(ctx)
		showingMessages = flash.Prev.Messages

		dataID       string = string(ctx.QueryArgs().Peek("id"))
		isAddMode    bool   = (dataID == "")

		hasPrevError bool                = false
		prevInput    *sys.SysUserDE
		currInput    *sys.SysUserDE

		contentData map[string]any
	)

	prevError := flash.Prev.Data.Get("prev_error")
	for k, vany := range prevError {
		switch k {
		case "user":
			hasPrevError = true
			assertion, ok := vany.(*sys.SysUserDE)
			if ok { prevInput = assertion }
		}
	}

	switch {
	case isAddMode:
		if hasPrevError {
			currInput = prevInput
		} else {
			currInput = sys.Bundle.Service.SysUser__EmptyItem()
		}

		contentData = map[string]any{
			"_title": "New User",
			"_link": map[string]string{
				"back": "/admin/user/listing",
				"save": "/admin/user/detail/do/add",
				"delete": "",
			},
			"_is_add_mode" : isAddMode,
			"_is_edit_mode": !isAddMode,
			"_data": map[string]any{
				"user": currInput,
			},
		}

	case !isAddMode:
		sysUserDE, domainMessages := sys.Bundle.Service.SysUser__Get(dataID)
		if domainMessages.HasError() {
			flash.Next.Messages.Append(domainMessages)
			ctx.Redirect("/admin/user/listing", fasthttp.StatusFound)
			return
		}
		showingMessages.Append(domainMessages)

		if hasPrevError {
			currInput = prevInput
		} else {
			currInput = sysUserDE
		}

		contentData = map[string]any{
			"_title": fmt.Sprintf("User %s", sysUserDE.Username),
			"_link": map[string]string{
				"back": "/admin/user/listing",
				"save": "/admin/user/detail/do/edit",
				"delete": "/admin/user/detail/do/delete",
			},
			"_is_add_mode" : isAddMode,
			"_is_edit_mode": !isAddMode,
			"_data": map[string]any{
				"user": currInput,
			},
		}
	}

	pageData := rh.templating.NewPageData()
	pageData.Title = "Admin - User Detail"
	pageData.Messages = showingMessages
	pageData.ContentData = contentData
	pageData.ExtraJsLinks = []string{
		"/admin/assets/user-detail.js",
	}

	renderOut, renderErr := rh.templating.RenderContent_default("d-content--user-detail.html.tmpl", pageData, ctx)
	if renderErr != nil {
		rh.baseBundle.RouteHandler.Route500(ctx, renderErr)
		return
	}
	
	ctx.SetContentType(webConstant.CONTENT_TYPE_HTML)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte(renderOut))
}

func (rh *AdminRouteHandler) doAddUser(ctx *fasthttp.RequestCtx) {
	var (
		flash = rh.baseBundle.Flash.GetFlashBundle(ctx)

		sysUserDE = sys.Bundle.Service.SysUser__EmptyItem()

		redirectURI string
	)

	ctx.PostArgs().VisitAll(func(k, v []byte) {
		sk, sv := string(k), string(v)
		switch sk {
		case "username"   : sysUserDE.Username = sv
		case "password"   : sysUserDE.Password = sv
		case "role_simple": sysUserDE.RoleSimple = sv
		}
	})

	_sysUserDE, domainMessages := sys.Bundle.Service.SysUser__Add(sysUserDE)
	if domainMessages.HasError() {
		flash.Next.Data.Set("prev_error", map[string]any{
			"user": sysUserDE,
		})
		flash.Next.Messages.Append(domainMessages)
		redirectURI = "/admin/user/detail"
	} else {
		flash.Next.Messages.Append(domainMessages)
		redirectURI = fmt.Sprintf("/admin/user/detail?id=%s", _sysUserDE.ID)
	}

	ctx.Redirect(redirectURI, fasthttp.StatusFound)
}

func (rh *AdminRouteHandler) doEditUser(ctx *fasthttp.RequestCtx) {
	var (
		flash = rh.baseBundle.Flash.GetFlashBundle(ctx)

		dataID string = ""
		sysUserDE     = sys.Bundle.Service.SysUser__EmptyItem()

		redirectURI string
	)

	ctx.PostArgs().VisitAll(func(k, v []byte) {
		sk, sv := string(k), string(v)
		switch sk {
		case "id"         : dataID = sv
		case "username"   : sysUserDE.Username = sv
		case "password"   : sysUserDE.Password = sv
		case "role_simple": sysUserDE.RoleSimple = sv
		}
	})

	_, domainMessages := sys.Bundle.Service.SysUser__Edit(dataID, sysUserDE)
	if domainMessages.HasError() {
		flash.Next.Data.Set("prev_error", map[string]any{
			"user": sysUserDE,
		})
	}
	flash.Next.Messages.Append(domainMessages)
	redirectURI = fmt.Sprintf("/admin/user/detail?id=%s", dataID)

	ctx.Redirect(redirectURI, fasthttp.StatusFound)
}

func (rh *AdminRouteHandler) doDeleteUser(ctx *fasthttp.RequestCtx) {
	var (
		flash = rh.baseBundle.Flash.GetFlashBundle(ctx)

		dataID string = string(ctx.PostArgs().Peek("id"))

		redirectURI string
	)

	domainMessages := sys.Bundle.Service.SysUser__Delete(dataID)
	if domainMessages.HasError() {
		flash.Next.Messages.Append(domainMessages)
		redirectURI = fmt.Sprintf("/admin/user/detail?id=%s", dataID)
	} else {
		flash.Next.Messages.Append(domainMessages)
		redirectURI = "/admin/user/listing"
	}

	ctx.Redirect(redirectURI, fasthttp.StatusFound)
}

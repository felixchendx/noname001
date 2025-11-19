package route

import (
	"fmt"
	"strings"

	"github.com/valyala/fasthttp"

	webConstant "noname001/web/constant"

	wallService "noname001/app/module/feature/wall/service"
)

func (rh *ModuleRouteHandler) renderWallDetail(ctx *fasthttp.RequestCtx) {
	var (
		flash           = rh.baseBundle.Flash.GetFlashBundle(ctx)
		showingMessages = flash.Prev.Messages

		wallServiceInstance = wallService.Instance()

		host     []byte = ctx.Request.Header.Peek("Host")
		hostPart []string = strings.Split(string(host), ":")

		dataID       string = string(ctx.QueryArgs().Peek("id"))
		isAddMode    bool   = (dataID == "")

		hasPrevError bool                = false
		prevInput    *wallService.WallDE
		currInput    *wallService.WallDE

		wallLayoutDEList []*wallService.WallLayoutDE

		contentData map[string]any
	)

	prevError := flash.Prev.Data.Get("prev_error")
	for k, vany := range prevError {
		switch k {
		case "wall":
			hasPrevError = true
			assertion, ok := vany.(*wallService.WallDE)
			if ok { prevInput = assertion }
		}
	}

	sc := &wallService.WallLayout__SearchCriteria{
		State: []string{"active"},
	}
	wallLayoutSR, _ := wallServiceInstance.WallLayout__Search(sc)
	if wallLayoutSR != nil {
		wallLayoutDEList = wallLayoutSR.Data
	}

	switch {
	case isAddMode:
		if hasPrevError {
			currInput = prevInput
		} else {
			currInput = wallServiceInstance.Wall__Empty()
		}

		contentData = map[string]any{
			"_title": "New Wall",
			"_link": map[string]string{
				"wall_listing": "/wall/wall/listing",
				"wall_view": "#",
				"save": "/wall/wall/detail/do/add",
				"delete": "",
				"ws": "",
			},
			"_is_add_mode" : isAddMode,
			"_is_edit_mode": !isAddMode,
			"_data": map[string]any{
				"wall_layout": wallLayoutDEList,
				"wall": currInput,
			},
		}

	case !isAddMode:
		wallDE, domainMessages := wallServiceInstance.Wall__Find(dataID, true)
		if domainMessages.HasError() {
			flash.Next.Messages.Append(domainMessages)
			ctx.Redirect("/wall/wall/listing", fasthttp.StatusFound)
			return
		}
		showingMessages.Append(domainMessages)

		if hasPrevError {
			currInput = prevInput
		} else {
			for _, wallItem := range wallDE.Items {
				wallItem.TempRelayURL = wallServiceInstance.GetRelayedStreamViewURL(hostPart[0], wallItem.SourceNodeID, wallItem.StreamCode, "hls")
			}

			currInput = wallDE
		}

		contentData = map[string]any{
			"_title": fmt.Sprintf("Wall %s", wallDE.Code),
			"_link": map[string]string{
				"wall_listing": "/wall/wall/listing",
				"wall_view": fmt.Sprintf("/wall/wall/view?id=%s", wallDE.ID),
				"save": "/wall/wall/detail/do/edit",
				"delete": "/wall/wall/detail/do/delete",
				"ws": fmt.Sprintf("/wall/wall/detail/ws?id=%s", wallDE.ID),
			},
			"_is_add_mode" : isAddMode,
			"_is_edit_mode": !isAddMode,
			"_data": map[string]any{
				"wall_layout": wallLayoutDEList,
				"wall": currInput,
			},
		}
	}

	pageData := rh.templating.NewPageData_default()
	pageData.Title = "Wall - Wall Detail"
	pageData.Messages = showingMessages
	pageData.ContentData = contentData
	pageData.ExtraCssLinks = []string{
		"/assets/fu-component/live-stream/default/index.css",
		"/wall/assets/wall-detail.css",
	}
	pageData.ExtraJsLinks = []string{
		"/assets/hls/hls.js",
		"/assets/internal/ws.js",
		"/assets/fu-component/live-stream/default/index.js",
		"/wall/assets/wall-detail.js",
	}

	renderOut, renderErr := rh.templating.RenderContent_default("d-content--wall-detail.html.tmpl", pageData, ctx)
	if renderErr != nil {
		rh.baseBundle.RouteHandler.Route500(ctx, renderErr)
		return
	}
	
	ctx.SetContentType(webConstant.CONTENT_TYPE_HTML)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte(renderOut))
}

func (rh *ModuleRouteHandler) doAddWall(ctx *fasthttp.RequestCtx) {
	var (
		flash = rh.baseBundle.Flash.GetFlashBundle(ctx)

		wallDE = wallService.Instance().Wall__Empty()

		redirectURI string
	)

	ctx.PostArgs().VisitAll(func(k, v []byte) {
		sk, sv := string(k), string(v)
		switch sk {
		case "code"          : wallDE.Code = sv
		case "name"          : wallDE.Name = sv
		case "state"         : wallDE.State = sv
		case "note"          : wallDE.Note = sv
		case "wall_layout_id": wallDE.WallLayoutID = sv
		default: // ignored
		}
	})

	_wallDE, domainMessages := wallService.Instance().Wall__Add(wallDE)
	if domainMessages.HasError() {
		flash.Next.Data.Set("prev_error", map[string]any{
			"wall": wallDE,
		})
		flash.Next.Messages.Append(domainMessages)
		redirectURI = "/wall/wall/detail"
	} else {
		flash.Next.Messages.Append(domainMessages)
		redirectURI = fmt.Sprintf("/wall/wall/detail?id=%s", _wallDE.ID)
	}

	ctx.Redirect(redirectURI, fasthttp.StatusFound)
}

func (rh *ModuleRouteHandler) doEditWall(ctx *fasthttp.RequestCtx) {
	var (
		flash = rh.baseBundle.Flash.GetFlashBundle(ctx)

		dataID string = ""
		wallDE        = wallService.Instance().Wall__Empty()

		redirectURI string
	)

	ctx.PostArgs().VisitAll(func(k, v []byte) {
		sk, sv := string(k), string(v)
		switch sk {
		case "id"            : dataID = sv
		case "code"          : wallDE.Code = sv
		case "name"          : wallDE.Name = sv
		case "state"         : wallDE.State = sv
		case "note"          : wallDE.Note = sv
		case "wall_layout_id": wallDE.WallLayoutID = sv
		default: // ignored
		}
	})

	_, domainMessages := wallService.Instance().Wall__Edit(dataID, wallDE)
	if domainMessages.HasError() {
		flash.Next.Data.Set("prev_error", map[string]any{
			"wall": wallDE,
		})
	}
	flash.Next.Messages.Append(domainMessages)
	redirectURI = fmt.Sprintf("/wall/wall/detail?id=%s", dataID)

	ctx.Redirect(redirectURI, fasthttp.StatusFound)
}

func (rh *ModuleRouteHandler) doDeleteWall(ctx *fasthttp.RequestCtx) {
	var (
		flash = rh.baseBundle.Flash.GetFlashBundle(ctx)

		dataID string = string(ctx.PostArgs().Peek("id"))

		redirectURI string
	)

	domainMessages := wallService.Instance().Wall__Delete(dataID)
	if domainMessages.HasError() {
		flash.Next.Messages.Append(domainMessages)
		redirectURI = fmt.Sprintf("/wall/wall/detail?id=%s", dataID)
	} else {
		flash.Next.Messages.Append(domainMessages)
		redirectURI = "/wall/wall/listing"
	}

	ctx.Redirect(redirectURI, fasthttp.StatusFound)
}

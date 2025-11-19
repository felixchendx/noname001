package route

import (
	"fmt"

	"github.com/valyala/fasthttp"

	webConstant "noname001/web/constant"

	streamfs      "noname001/app/module/feature/stream/filesystem"
	streamService "noname001/app/module/feature/stream/service"
)

const (
	streamGroupDetail__pageTitle    = "Stream - Stream Group Detail"
	streamGroupDetail__prevErrorKey = "streamGroupDetail"
)

func (rh *ModuleRouteHandler) streamGropuDetail__render(ctx *fasthttp.RequestCtx) {
	var (
		flash           = rh.baseBundle.Flash.GetFlashBundle(ctx)
		showingMessages = flash.Prev.Messages

		dataID    string = string(ctx.QueryArgs().Peek("id"))
		isAddMode bool   = (dataID == "")

		hasPrevError bool = false
		prevInput    *streamService.StreamGroupDE
		currInput    *streamService.StreamGroupDE

		streamProfileDEList []*streamService.StreamProfileDE

		contentData map[string]any
	)

	prevError := flash.Prev.Data.Get("prev_error")
	for k, vany := range prevError {
		switch k {
		case streamGroupDetail__prevErrorKey:
			hasPrevError = true
			assertion, ok := vany.(*streamService.StreamGroupDE)
			if ok { prevInput = assertion }
		}
	}

	streamProfileSR, _ := streamService.Instance().SearchStreamProfile(
		&streamService.StreamProfile__SearchCriteria{
			State: []string{"active", "readonly"},
		},
	)
	if streamProfileSR != nil {
		streamProfileDEList = streamProfileSR.Data
	}

	switch {
	case isAddMode:
		if hasPrevError {
			currInput = prevInput
		} else {
			currInput = streamService.Instance().EmptyStreamGroup()
		}

		contentData = map[string]any {
			"_title": "New Stream Group",
			"_link" : map[string]string{
				"back"  : navLink__streamGroupListing,
				"save"  : actLink__streamGroupDoAdd,
				"delete": navLink__blank,
			},
			"_is_add_mode" : isAddMode,
			"_is_edit_mode": !isAddMode,
			"_data": map[string]any{
				"stream_profile": streamProfileDEList,
				"stream_group"  : currInput,

				"stream_local_dir_placeholder": streamfs.StreamLocalDirPlaceholder(),
			},
		}

	case !isAddMode:
		streamGroupDE, domainMessages := streamService.Instance().FindStreamGroup(dataID)
		if domainMessages.HasError() {
			flash.Next.Messages.Append(domainMessages)
			ctx.Redirect(navLink__streamGroupListing, fasthttp.StatusFound)
			return
		}
		showingMessages.Append(domainMessages)

		if hasPrevError {
			currInput = prevInput
		} else {
			currInput = streamGroupDE
		}

		contentData = map[string]any {
			"_title": fmt.Sprintf("Stream Group %s", streamGroupDE.Code),
			"_link": map[string]string{
				"back"  : navLink__streamGroupListing,
				"save"  : actLink__streamGroupDoEdit,
				"delete": actLink__streamGroupDoDelete,
			},
			"_is_add_mode" : isAddMode,
			"_is_edit_mode": !isAddMode,
			"_data": map[string]any{
				"stream_profile": streamProfileDEList,
				"stream_group"  : currInput,

				"stream_local_dir_placeholder": streamfs.StreamLocalDirPlaceholder(),
			},
		}
	}

	pageData := rh.templating.NewPageData_default()
	pageData.Title = streamGroupDetail__pageTitle
	pageData.Messages = showingMessages
	pageData.ContentData = contentData
	pageData.ExtraCssLinks = []string{
		"/assets/fu-component/live-stream/default/index.css",
		"/stream/assets/stream-group-detail-02.css",
	}
	pageData.ExtraJsLinks = []string{
		"/assets/hls/hls.js",
		"/assets/internal/ws.js",
		"/assets/fu-component/live-stream/default/index.js",
		"/stream/assets/stream-group-detail-02.js",
	}

	renderOut, renderErr := rh.templating.RenderContent_default("d-content--stream-group-detail-02.html.tmpl", pageData, ctx)
	if renderErr != nil {
		rh.baseBundle.RouteHandler.Route500(ctx, renderErr)
		return
	}

	ctx.SetContentType(webConstant.CONTENT_TYPE_HTML)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte(renderOut))
}

func (rh *ModuleRouteHandler) streamGroupDetail__doAdd(ctx *fasthttp.RequestCtx) {
	var (
		flash = rh.baseBundle.Flash.GetFlashBundle(ctx)

		currInput = streamService.Instance().EmptyStreamGroup()

		redirectURI string
	)

	ctx.PostArgs().VisitAll(func(k, v []byte) {
		sk, sv := string(k), string(v)
		switch sk {
		case "code"         : currInput.Code = sv
		case "name"         : currInput.Name = sv
		case "state"        : currInput.State = sv
		case "note"         : currInput.Note = sv
		case "streamProfile": currInput.StreamProfileID = sv
		}
	})

	_streamGroupDE, domainMessages := streamService.Instance().AddStreamGroup(currInput)
	if domainMessages.HasError() {
		flash.Next.Data.Set("prev_error", map[string]any{
			streamGroupDetail__prevErrorKey: currInput,
		})
		flash.Next.Messages.Append(domainMessages)
		redirectURI = navLink__streamGroupDetail
	} else {
		flash.Next.Messages.Append(domainMessages)
		redirectURI = navLink__streamGroupDetailWithID(_streamGroupDE.ID)
	}

	ctx.Redirect(redirectURI, fasthttp.StatusFound)
}

func (rh *ModuleRouteHandler) streamGroupDetail__doEdit(ctx *fasthttp.RequestCtx) {
	var (
		flash = rh.baseBundle.Flash.GetFlashBundle(ctx)

		dataID    string = ""
		currInput        = streamService.Instance().EmptyStreamGroup()

		redirectURI string
	)

	ctx.PostArgs().VisitAll(func(k, v []byte) {
		sk, sv := string(k), string(v)
		switch sk {
		case "id"           : dataID = sv
		case "code"         : currInput.Code = sv
		case "name"         : currInput.Name = sv
		case "state"        : currInput.State = sv
		case "note"         : currInput.Note = sv
		case "streamProfile": currInput.StreamProfileID = sv
		}
	})

	_, domainMessages := streamService.Instance().EditStreamGroup(dataID, currInput)
	if domainMessages.HasError() {
		flash.Next.Data.Set("prev_error", map[string]any{
			streamGroupDetail__prevErrorKey: currInput,
		})
	}
	flash.Next.Messages.Append(domainMessages)
	redirectURI = navLink__streamGroupDetailWithID(dataID)

	ctx.Redirect(redirectURI, fasthttp.StatusFound)
}

func (rh *ModuleRouteHandler) streamGroupDetail__doDelete(ctx *fasthttp.RequestCtx) {
	var (
		flash = rh.baseBundle.Flash.GetFlashBundle(ctx)

		dataID string = string(ctx.PostArgs().Peek("id"))

		redirectURI string
	)

	domainMessages := streamService.Instance().DeleteStreamGroup(dataID)
	if domainMessages.HasError() {
		flash.Next.Messages.Append(domainMessages)
		redirectURI = navLink__streamGroupDetailWithID(dataID)
	} else {
		flash.Next.Messages.Append(domainMessages)
		redirectURI = navLink__streamGroupListing
	}

	ctx.Redirect(redirectURI, fasthttp.StatusFound)
}

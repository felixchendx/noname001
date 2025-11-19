package route

import (
	"fmt"

	"github.com/valyala/fasthttp"

	webConstant "noname001/web/constant"

	streamService "noname001/app/module/feature/stream/service"
)

func (rh *ModuleRouteHandler) renderStreamGroupDetail(ctx *fasthttp.RequestCtx) {
	var (
		flash                       = rh.baseBundle.Flash.GetFlashBundle(ctx)
		showingMessages             = flash.Prev.Messages

		streamServiceInstance       = streamService.Instance()

		dataID               string = string(ctx.QueryArgs().Peek("id"))
		isAddMode            bool   = (dataID == "")

		hasPrevError         bool   = false
		prevInputStreamGroup, currInputStreamGroup *streamService.StreamGroupDE

		streamProfileDEList  []*streamService.StreamProfileDE

		contentData          map[string]any
	)

	prevError := flash.Prev.Data.Get("prev_error")
	for key, valueAny := range prevError {
		switch key {
		case "stream_group":
			hasPrevError = true
			assertion, ok := valueAny.(*streamService.StreamGroupDE)
			if ok { prevInputStreamGroup = assertion }
		}
	}

	sc := &streamService.StreamProfile__SearchCriteria{
		State: []string{"active", "readonly"},
	}
	streamProfileSR, _ := streamServiceInstance.SearchStreamProfile(sc)
	if streamProfileSR != nil && streamProfileSR.Data != nil {
		streamProfileDEList = streamProfileSR.Data
	}

	switch {
	case isAddMode:
		if hasPrevError {
			currInputStreamGroup = prevInputStreamGroup
		} else {
			currInputStreamGroup = streamServiceInstance.EmptyStreamGroup()
		}

		contentData = map[string]any{
			"_title": "New Stream Group",
			"_link": map[string]string{
				"back":   "/stream/stream-group/listing",
				"save":   "/stream/stream-group/detail/do/add",
				"delete": "",
			},
			"_is_add_mode":  isAddMode,
			"_is_edit_mode": !isAddMode,
			"_data": map[string]any{
				"stream_profile": streamProfileDEList,
				"stream_group":   currInputStreamGroup,
			},
		}

	case !isAddMode:
		streamGroupDE, domainMessages := streamServiceInstance.FindStreamGroup(dataID)
		if domainMessages.HasError() {
			flash.Next.Messages.Append(domainMessages)
			ctx.Redirect("/stream/stream-group/listing", fasthttp.StatusFound)
			return
		}
		showingMessages.Append(domainMessages)

		streamItemsDE, domainMessages := streamServiceInstance.FindStreamItemsByStreamGroupID(dataID)
		if domainMessages.HasError() {
			flash.Next.Messages.Append(domainMessages)
			ctx.Redirect("/stream/stream-group/listing", fasthttp.StatusFound)
			return
		}

		if hasPrevError {
			currInputStreamGroup = prevInputStreamGroup
		} else {
			currInputStreamGroup = streamGroupDE
		}

		contentData = map[string]any{
			"_title": fmt.Sprintf("Stream Group %s", streamGroupDE.Code),
			"_link": map[string]any{
				"back":   "/stream/stream-group/listing",
				"save":   "/stream/stream-group/detail/do/edit",
				"delete": "/stream/stream-group/detail/do/delete",
			},
			"_is_add_mode":  isAddMode,
			"_is_edit_mode": !isAddMode,
			"_data": map[string]any{
				"stream_profile": streamProfileDEList,
				"stream_group":   currInputStreamGroup,
			},
			"_datatable": map[string]any{
				"si_listing": map[string]any{
					"rows":   streamItemsDE,
				},
			},
		}
	}

	pageData := rh.templating.NewPageData_default()
	pageData.Title = "Stream - Stream Group Detail"
	pageData.Messages = showingMessages
	pageData.ContentData = contentData
	pageData.ExtraJsLinks = []string{
		"/assets/hls/hls.js",
		"/stream/assets/stream-group-detail.js",
	}

	renderOut, renderErr := rh.templating.RenderContent_default("d-content--stream-group-detail.html.tmpl", pageData, ctx)
	if renderErr != nil {
		rh.baseBundle.RouteHandler.Route500(ctx, renderErr)
		return
	}

	ctx.SetContentType(webConstant.CONTENT_TYPE_HTML)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte(renderOut))
}

func (rh *ModuleRouteHandler) doAddStreamGroup(ctx *fasthttp.RequestCtx) {
	var (
		flash         = rh.baseBundle.Flash.GetFlashBundle(ctx)

		streamGroupDE = streamService.Instance().EmptyStreamGroup()

		redirectURI   string
	)

	ctx.PostArgs().VisitAll(func(key, value []byte) {
		sKey, sValue := string(key), string(value)
		switch sKey {
		case "code"          : streamGroupDE.Code = sValue
		case "name"          : streamGroupDE.Name = sValue
		case "state"         : streamGroupDE.State = sValue
		case "note"          : streamGroupDE.Note = sValue
		case "streamProfile" : streamGroupDE.StreamProfileID = sValue
		}
	})

	_streamGroupDE, domainMessages := streamService.Instance().AddStreamGroup(streamGroupDE)
	if domainMessages.HasError() {
		flash.Next.Data.Set("prev_error", map[string]any{
			"stream_group": streamGroupDE,
		})
		flash.Next.Messages.Append(domainMessages)
		redirectURI = "/stream/stream-group/detail"
	} else {
		flash.Next.Messages.Append(domainMessages)
		redirectURI = fmt.Sprintf("/stream/stream-group/detail?id=%s", _streamGroupDE.ID)
	}

	ctx.Redirect(redirectURI, fasthttp.StatusFound)
}

func (rh *ModuleRouteHandler) doEditStreamGroup(ctx *fasthttp.RequestCtx) {
	var (
		flash                = rh.baseBundle.Flash.GetFlashBundle(ctx)

		dataID        string = ""
		streamGroupDE        = streamService.Instance().EmptyStreamGroup()

		redirectURI   string
	)

	ctx.PostArgs().VisitAll(func(key, value []byte) {
		sKey, sValue := string(key), string(value)
		switch sKey {
		case "id"            : dataID = sValue
		case "code"          : streamGroupDE.Code = sValue
		case "name"          : streamGroupDE.Name = sValue
		case "state"         : streamGroupDE.State = sValue
		case "note"          : streamGroupDE.Note = sValue
		case "streamProfile" : streamGroupDE.StreamProfileID = sValue
		}
	})

	_, domainMessages := streamService.Instance().EditStreamGroup(dataID, streamGroupDE)
	if domainMessages.HasError() {
		flash.Next.Data.Set("prev_error", map[string]any{
			"stream_group": streamGroupDE,
		})
	}
	flash.Next.Messages.Append(domainMessages)
	redirectURI = fmt.Sprintf("/stream/stream-group/detail?id=%s", dataID)

	ctx.Redirect(redirectURI, fasthttp.StatusFound)
}

func (rh *ModuleRouteHandler) doDeleteStreamGroup(ctx *fasthttp.RequestCtx) {
	var (
		flash              = rh.baseBundle.Flash.GetFlashBundle(ctx)

		dataID      string = string(ctx.PostArgs().Peek("id"))

		redirectURI string
	)

	domainMessages := streamService.Instance().DeleteStreamGroup(dataID)
	if domainMessages.HasError() {
		flash.Next.Messages.Append(domainMessages)
		redirectURI = fmt.Sprintf("/stream/stream-group/detail?id=%s", dataID)
	} else {
		flash.Next.Messages.Append(domainMessages)
		redirectURI = "/stream/stream-group/listing"
	}

	ctx.Redirect(redirectURI, fasthttp.StatusFound)
}

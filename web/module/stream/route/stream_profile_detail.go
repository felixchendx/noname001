package route

import (
	"fmt"
	"strconv"

	"github.com/valyala/fasthttp"

	webConstant "noname001/web/constant"

	streamService "noname001/app/module/feature/stream/service"
)

func (rh *ModuleRouteHandler) renderStreamProfileDetail(ctx *fasthttp.RequestCtx) {
	var (
		flash                       = rh.baseBundle.Flash.GetFlashBundle(ctx)
		showingMessages             = flash.Prev.Messages

		dataID               string = string(ctx.QueryArgs().Peek("id"))
		isAddMode            bool   = (dataID == "")

		contentData          map[string]any

		hasPrevError         bool   = false
		prevInput, currInput *streamService.StreamProfileDE
	)

	prevError := flash.Prev.Data.Get("prev_error")
	for key, value := range prevError {
		switch key {
		case "streamProfile":
			hasPrevError = true
			assertion, ok := value.(*streamService.StreamProfileDE)
			if ok {	prevInput = assertion }
		}
	}

	switch {
	case isAddMode:
		if hasPrevError {
			currInput = prevInput
		} else {
			currInput = streamService.Instance().EmptyStreamProfile()
			currInput.TargetVideoCompression = 80
			currInput.TargetVideoBitrate = 300000
			currInput.TargetAudioCompression = 80
			currInput.TargetAudioBitrate = 32000
		}
	
		contentData = map[string]any{
			"_title": "New Stream Profile",
			"_link" : map[string]string{
				"back"  : "/stream/stream-profile/listing",
				"save"  : "/stream/stream-profile/detail/do/add",
				"delete": "",
			},
			"_is_add_mode" : isAddMode,
			"_data": map[string]any{
				"sp": currInput,
			},
		}

	case !isAddMode:
		streamProfileDE, domainMessage := streamService.Instance().FindStreamProfile(dataID)
		if domainMessage.HasError() {
			flash.Next.Messages.Append(domainMessage)
			ctx.Redirect("/stream/stream-profile/listing", fasthttp.StatusFound)
			return
		}
		showingMessages.Append(domainMessage)
	
		if hasPrevError {
			currInput = prevInput
		} else {
			currInput = streamProfileDE
		}
	
		contentData = map[string]any{
			"_title": fmt.Sprintf("Stream Profile %s", streamProfileDE.Code),
			"_link" : map[string]string{
				"back"  : "/stream/stream-profile/listing",
				"save"  : "/stream/stream-profile/detail/do/edit",
				"delete": "/stream/stream-profile/detail/do/delete",
			},
			"_is_add_mode" : isAddMode,
			"_data": map[string]any{
				"sp": currInput,
			},
		}
	}

	pageData := rh.templating.NewPageData_default()
	pageData.Title = "Stream - Stream Profile Detail"
	pageData.Messages = showingMessages
	pageData.ContentData = contentData
	pageData.ExtraJsLinks = []string{
		"/stream/assets/stream-profile-detail.js",
	}

	renderOut, renderErr := rh.templating.RenderContent_default("d-content--stream-profile-detail.html.tmpl", pageData, ctx)
	if renderErr != nil {
		rh.baseBundle.RouteHandler.Route500(ctx, renderErr)
		return
	}

	ctx.SetContentType(webConstant.CONTENT_TYPE_HTML)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte(renderOut))
}

func (rh *ModuleRouteHandler) doAddStreamProfile(ctx *fasthttp.RequestCtx) {
	var (
		flash           = rh.baseBundle.Flash.GetFlashBundle(ctx)

		streamProfileDE = streamService.Instance().EmptyStreamProfile()

		redirectURI string
	)

	ctx.PostArgs().VisitAll(func(key, value []byte) {
		sKey, sValue := string(key), string(value)
		switch sKey {
		case "code"             : streamProfileDE.Code = sValue
		case "name"             : streamProfileDE.Name = sValue
		case "state"            : streamProfileDE.State = sValue
		case "note"             : streamProfileDE.Note = sValue
		case "targetVideoCodec" : streamProfileDE.TargetVideoCodec = sValue
		case "targetAudioCodec" : streamProfileDE.TargetAudioCodec = sValue
		
		case "targetVideoCompression" :
			vidCompression, err := strconv.Atoi(sValue)
			if err != nil { vidCompression = 0 }
			streamProfileDE.TargetVideoCompression = vidCompression

		case "targetVideoBitrate"     :
			vidBitrate, err := strconv.Atoi(sValue)
			if err != nil { vidBitrate = 0 }
			streamProfileDE.TargetVideoBitrate = vidBitrate

		case "targetAudioCompression" :
			audCompression, err := strconv.Atoi(sValue)
			if err != nil { audCompression = 0 }
			streamProfileDE.TargetAudioCompression = audCompression

		case "targetAudioBitrate"     :
			audBitrate, err := strconv.Atoi(sValue)
			if err != nil { audBitrate = 0 }
			streamProfileDE.TargetAudioBitrate = audBitrate
		}
	})

	_streamProfileDE, domainMessages := streamService.Instance().AddStreamProfile(streamProfileDE)
	if domainMessages.HasError() {
		flash.Next.Data.Set("prev_error", map[string]any{
			"streamProfile": streamProfileDE,
		})
		flash.Next.Messages.Append(domainMessages)
		redirectURI = "/stream/stream-profile/detail"
	} else {
		flash.Next.Messages.Append(domainMessages)
		redirectURI = fmt.Sprintf("/stream/stream-profile/detail?id=%s", _streamProfileDE.ID)
	}

	ctx.Redirect(redirectURI, fasthttp.StatusFound)
}

func (rh *ModuleRouteHandler) doEditStreamProfile(ctx *fasthttp.RequestCtx) {
	var (
		flash                  = rh.baseBundle.Flash.GetFlashBundle(ctx)

		dataID          string = ""
		streamProfileDE        = streamService.Instance().EmptyStreamProfile()

		redirectURI     string = ""
	)

	ctx.PostArgs().VisitAll(func(key, value []byte) {
		sKey, sValue := string(key), string(value)
		switch sKey {
		case "id"               : dataID = sValue
		case "code"             : streamProfileDE.Code = sValue
		case "name"             : streamProfileDE.Name = sValue
		case "state"            : streamProfileDE.State = sValue
		case "note"             : streamProfileDE.Note = sValue
		case "targetVideoCodec" : streamProfileDE.TargetVideoCodec = sValue
		case "targetAudioCodec" : streamProfileDE.TargetAudioCodec = sValue

		case "targetVideoCompression" :
			vidCompression, err := strconv.Atoi(sValue)
			if err != nil { vidCompression = 0 }
			streamProfileDE.TargetVideoCompression = vidCompression

		case "targetVideoBitrate"     :
			vidBitrate, err := strconv.Atoi(sValue)
			if err != nil { vidBitrate = 0 }
			streamProfileDE.TargetVideoBitrate = vidBitrate

		case "targetAudioCompression" :
			audCompression, err := strconv.Atoi(sValue)
			if err != nil { audCompression = 0 }
			streamProfileDE.TargetAudioCompression = audCompression

		case "targetAudioBitrate"     :
			audBitrate, err := strconv.Atoi(sValue)
			if err != nil { audBitrate = 0 }
			streamProfileDE.TargetAudioBitrate = audBitrate
		}
	})

	_, domainMessages := streamService.Instance().EditStreamProfile(dataID, streamProfileDE)
	if domainMessages.HasError() {
		flash.Next.Data.Set("prev_error", map[string]any{
			"streamProfile": streamProfileDE,
		})
	}
	flash.Next.Messages.Append(domainMessages)
	redirectURI = fmt.Sprintf("/stream/stream-profile/detail?id=%s", dataID)

	ctx.Redirect(redirectURI, fasthttp.StatusFound)
}

func (rh *ModuleRouteHandler) doDeleteStreamProfile(ctx *fasthttp.RequestCtx) {
	var (
		flash = rh.baseBundle.Flash.GetFlashBundle(ctx)

		dataID = string(ctx.PostArgs().Peek("id"))

		redirectURI string
	)

	domainMessage := streamService.Instance().DeleteStreamProfile(dataID)
	if domainMessage.HasError() {
		flash.Next.Messages.Append(domainMessage)
		redirectURI = fmt.Sprintf("/stream/stream-profile/detail?id=%s", dataID)
	} else {
		flash.Next.Messages.Append(domainMessage)
		redirectURI = "/stream/stream-profile/listing"
	}

	ctx.Redirect(redirectURI, fasthttp.StatusFound)
}

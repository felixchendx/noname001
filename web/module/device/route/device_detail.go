package route

import (
	"fmt"

	"github.com/valyala/fasthttp"

	webConstant "noname001/web/constant"

	deviceService "noname001/app/module/common/device/service"
)

func (rh *ModuleRouteHandler) renderDeviceDetail(ctx *fasthttp.RequestCtx) {
	var (
		flash           = rh.baseBundle.Flash.GetFlashBundle(ctx)
		showingMessages = flash.Prev.Messages

		dataID               string = string(ctx.QueryArgs().Peek("id"))
		isAddMode            bool   = (dataID == "")

		contentData          map[string]any

		hasPrevError         bool   = false
		prevInput, currInput *deviceService.DeviceDE
	)

	prevError := flash.Prev.Data.Get("prev_error")
	for key, valueAny := range prevError {
		switch key {
		case "device":
			hasPrevError = true
			assertion, ok := valueAny.(*deviceService.DeviceDE)
			if ok { prevInput = assertion }
		}
	}

	switch {
	case isAddMode:
		if hasPrevError {
			currInput = prevInput
		} else {
			currInput = deviceService.Instance().EmptyDevice()
		}

		contentData = map[string]any{
			"_title": "New Device",
			"_link": map[string]string{
				"back":   "/device/device/listing",
				"save":   "/device/device/detail/do/add",
				"delete": "",
			},
			"_is_add_mode" : isAddMode,
			"_is_edit_mode": !isAddMode,
			"_data": map[string]any{
				"device": currInput,
			},
		}

	case !isAddMode:
		deviceDE, domainMessages := deviceService.Instance().FindDevice(dataID)
		if domainMessages.HasError() {
			flash.Next.Messages.Append(domainMessages)
			ctx.Redirect("/device/device/listing", fasthttp.StatusFound)
			return
		}
		showingMessages.Append(domainMessages)

		if hasPrevError {
			currInput = prevInput
		} else {
			currInput = deviceDE
		}

		contentData = map[string]any{
			"_title": fmt.Sprintf("Device %s", deviceDE.Code),
			"_link": map[string]string{
				"back":   "/device/device/listing",
				"save":   "/device/device/detail/do/edit",
				"delete": "/device/device/detail/do/delete",
			},
			"_is_add_mode" : isAddMode,
			"_is_edit_mode": !isAddMode,
			"_data": map[string]any{
				"device": currInput,
			},
		}
	}

	pageData := rh.templating.NewPageData_default()
	pageData.Title = "Device - Device Detail"
	pageData.Messages = showingMessages
	pageData.ContentData = contentData
	pageData.ExtraCssLinks = []string{
		"/device/assets/device-detail.css",
	}
	pageData.ExtraJsLinks = []string{
		"/assets/internal/ws.js",
		"/device/assets/device-detail.js",
	}

	renderOut, renderErr := rh.templating.RenderContent_default("d-content--device-detail.html.tmpl", pageData, ctx)
	if renderErr != nil {
		rh.baseBundle.RouteHandler.Route500(ctx, renderErr)
		return
	}

	ctx.SetContentType(webConstant.CONTENT_TYPE_HTML)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte(renderOut))
}

func (rh *ModuleRouteHandler) doAddDevice(ctx *fasthttp.RequestCtx) {
	var (
		flash              = rh.baseBundle.Flash.GetFlashBundle(ctx)

		deviceDE           = deviceService.Instance().EmptyDevice()
		redirectURI string = ""
	)

	ctx.PostArgs().VisitAll(func(key, value []byte) {
		sKey, sValue := string(key), string(value)
		switch sKey {
		case "code"    : deviceDE.Code = sValue
		case "name"    : deviceDE.Name = sValue
		case "state"   : deviceDE.State = sValue
		case "note"    : deviceDE.Note = sValue
		case "protocol": deviceDE.Protocol = sValue
		case "hostname": deviceDE.Hostname = sValue
		case "port"    : deviceDE.Port = sValue
		case "username": deviceDE.Username = sValue
		case "password": deviceDE.Password = sValue
		case "brand"   : deviceDE.Brand = sValue

		case "fallback_rtsp_port": deviceDE.FallbackRTSPPort = sValue
		}
	})

	_deviceDE, domainMessages := deviceService.Instance().AddDevice(deviceDE)
	if domainMessages.HasError() {
		flash.Next.Data.Set("prev_error", map[string]any{
			"device": deviceDE,
		})
		flash.Next.Messages.Append(domainMessages)
		redirectURI = "/device/device/detail"
	} else {
		flash.Next.Messages.Append(domainMessages)
		redirectURI = fmt.Sprintf("/device/device/detail?id=%s", _deviceDE.ID)
	}

	ctx.Redirect(redirectURI, fasthttp.StatusFound)
}

func (rh *ModuleRouteHandler) doEditDevice(ctx *fasthttp.RequestCtx) {
	var (
		flash    = rh.baseBundle.Flash.GetFlashBundle(ctx)

		dataID   = ""
		deviceDE = deviceService.Instance().EmptyDevice()
	)

	ctx.PostArgs().VisitAll(func(key, value []byte) {
		sKey, sValue := string(key), string(value)
		switch sKey {
		case "id"      : dataID = sValue
		case "code"    : deviceDE.Code = sValue
		case "name"    : deviceDE.Name = sValue
		case "state"   : deviceDE.State = sValue
		case "note"    : deviceDE.Note = sValue
		case "protocol": deviceDE.Protocol = sValue
		case "hostname": deviceDE.Hostname = sValue
		case "port"    : deviceDE.Port = sValue
		case "username": deviceDE.Username = sValue
		case "password": deviceDE.Password = sValue
		case "brand"   : deviceDE.Brand = sValue

		case "fallback_rtsp_port": deviceDE.FallbackRTSPPort = sValue
		}
	})

	_, domainMessages := deviceService.Instance().EditDevice(dataID, deviceDE)
	if domainMessages.HasError() {
		flash.Next.Data.Set("prev_error", map[string]any{
			"device": deviceDE,
		})
	}
	flash.Next.Messages.Append(domainMessages)
	redirectURI := fmt.Sprintf("/device/device/detail?id=%s", dataID)

	ctx.Redirect(redirectURI, fasthttp.StatusFound)
}

func (rh *ModuleRouteHandler) doDeleteDevice(ctx *fasthttp.RequestCtx) {
	var (
		flash  = rh.baseBundle.Flash.GetFlashBundle(ctx)

		dataID = string(ctx.PostArgs().Peek("id"))

		redirectURI string
	)

	domainMessages := deviceService.Instance().DeleteDevice(dataID)
	if domainMessages.HasError() {
		flash.Next.Messages.Append(domainMessages)
		redirectURI = fmt.Sprintf("/device/device/detail?id=%s", dataID)
	} else {
		flash.Next.Messages.Append(domainMessages)
		redirectURI = "/device/device/listing"
	}

	ctx.Redirect(redirectURI, fasthttp.StatusFound)
}

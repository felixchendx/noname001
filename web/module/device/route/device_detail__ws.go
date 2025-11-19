package route

import (
	"github.com/valyala/fasthttp"

	wsUtil "noname001/web/base/comm/ws/util"

	deviceService "noname001/app/module/common/device/service"
	deviceIntface "noname001/app/module/common/device/intface"
)

// TODO: err handling via ws

func (rh *ModuleRouteHandler) deviceDetail__ws(ctx *fasthttp.RequestCtx) {
	var handlerName = "deviceDetail__ws"
	
	wsClient, err := rh.baseBundle.WSHub.UpgradeToWebsocketConnection(ctx)
	if err != nil {
		rh.logger.Errorf("%s: %s, upgrade err: %s", rh.logPrefix, handlerName, err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	deviceID := string(ctx.QueryArgs().Peek("id"))

	go func() {
		var (
			deviceServiceInstance = deviceService.Instance()

			liveDeviceEvSub = deviceIntface.EventProvider().SubscribeToLiveDeviceEvent()
		)

		defer func() {
			deviceIntface.EventProvider().UnsubscribeFromLiveDeviceEvent(liveDeviceEvSub)
		}()

		wsLoop:
		for {
			select {
			case recvBytes, open := <- wsClient.ReceiverChannel():
				if !open { break wsLoop }

				basicReqRep := wsUtil.ExtractBasicReqRep(recvBytes)

				switch basicReqRep.RequestCode {
				case dd__wsRequestCode__deviceSnapshot:
					deviceSnapshot, _ := deviceServiceInstance.GetDeviceSnapshot(deviceID)
					repStruct := dd__wsReply__deviceSnapshot{basicReqRep, deviceSnapshot}
					_ = wsClient.Send(wsUtil.NewReqRepJson(repStruct))

				case dd__wsRequestCode__tempErrorDetails:
					tempErrorDetails := deviceServiceInstance.GetTempErrorDetails(deviceID)
					repStruct := dd___wsReply__tempErrorDetails{basicReqRep, tempErrorDetails}
					_ = wsClient.Send(wsUtil.NewReqRepJson(repStruct))

				case dd__wsRequestCode__deviceReload:
					_ = deviceServiceInstance.ReloadDevice(deviceID)

				default:
					rh.logger.Warnf("%s: %s, no handler for %s", rh.logPrefix, handlerName, basicReqRep)
					break
				}

			case ev := <- liveDeviceEvSub.Channel:
				if ev.DeviceID != deviceID { break } // TODO: move filtering cap to inner stuffs

				// switch ev.EventCode {}
				evStruct := dd__wsEvent{string(ev.EventCode)}
				_ = wsClient.Send(wsUtil.NewEventJson(evStruct))
			}
		}
	}()
}

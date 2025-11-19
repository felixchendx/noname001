package route

import (
	"strings"
	
	"github.com/valyala/fasthttp"

	wsUtil "noname001/web/base/comm/ws/util"

	streamService "noname001/app/module/feature/stream/service"
	deviceService "noname001/app/module/common/device/service"
)

// TODO: err handling via ws

func (rh *ModuleRouteHandler) streamGroupDetail__ws(ctx *fasthttp.RequestCtx) {
	var handlerName = "streamGroupDetail__ws"

	wsClient, err := rh.baseBundle.WSHub.UpgradeToWebsocketConnection(ctx)
	if err != nil {
		rh.logger.Errorf("%s: %s, upgrade err: %s", rh.logPrefix, handlerName, err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	host := ctx.Request.Header.Peek("Host")
	hostPart := strings.Split(string(host), ":")
	streamGroupID := string(ctx.QueryArgs().Peek("id"))

	go func() {
		var (
			streamServiceInstance = streamService.Instance()
			deviceServiceInstance = deviceService.Instance()
		)

		defer func() {

		}()

		wsLoop:
		for {
			select {
			case recvBytes, open := <- wsClient.ReceiverChannel():
				if !open { break wsLoop }

				basicReqRep := wsUtil.ExtractBasicReqRep(recvBytes)

				switch basicReqRep.RequestCode {
				case sgd__wsRequestCode__streamItemListing:
					streamItems, domainMessages := streamServiceInstance.FindStreamItemsByStreamGroupID(streamGroupID)
					if domainMessages.HasError() {
						rh.logger.Errorf("%s: %s, FindStreamItemsByStreamGroupID err: %s", rh.logPrefix, handlerName, domainMessages.Dump())
						break
					}

					sgdStreamItems := make([]sgd__streamItem, 0)
					for _, _streamItem := range streamItems {
						sgdStreamItems = append(sgdStreamItems, sgd__streamItem{
							ID   : _streamItem.ID,
							Code : _streamItem.Code,
							Name : _streamItem.Name,
							State: _streamItem.State,
							Note : _streamItem.Note,

							SourceType      : _streamItem.SourceType,
							DeviceCode      : _streamItem.DeviceCode,
							DeviceChannelID : _streamItem.DeviceChannelID,
							DeviceStreamType: _streamItem.DeviceStreamType,
							ExternalURL     : _streamItem.ExternalURL,
							Filepath        : _streamItem.Filepath,

							StreamURL: streamServiceInstance.GetStreamViewURL(hostPart[0], _streamItem.Code, "hls"),
						})
					}

					repStruct := sgd__wsReply__streamItemListing{basicReqRep, sgdStreamItems}
					_ = wsClient.Send(wsUtil.NewReqRepJson(repStruct))

				case sgd__wsRequestCode__deviceSnapshotListing:
					deviceSnapshots := deviceServiceInstance.GetDeviceSnapshots()

					repStruct := sgd__wsReply__deviceSnapshotListing{basicReqRep, deviceSnapshots}
					_ = wsClient.Send(wsUtil.NewReqRepJson(repStruct))

				default:
					rh.logger.Warnf("%s: %s, no handler for %s", rh.logPrefix, handlerName, basicReqRep)
					break
				}
			}
		}
	}()
}

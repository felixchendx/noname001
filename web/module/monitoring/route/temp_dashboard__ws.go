package route

import (
	"github.com/valyala/fasthttp"

	"noname001/web/base/comm/ws"
	wsUtil "noname001/web/base/comm/ws/util"

	cacheIntface "noname001/app/module/common/cache/intface"
)

func (rh *ModuleRouteHandler) tempDashboard__wsHandler(ctx *fasthttp.RequestCtx) {
	wsClient, err := rh.baseBundle.WSHub.UpgradeToWebsocketConnection(ctx)
	if err != nil {
		rh.logger.Errorf("%s: tempDashboard__ws, upgrade err: %s", rh.logPrefix, err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	go func() {
		var (
			cachedNodeStatusFeedSub = cacheIntface.DataProvider().SubscribeToCachedNodeStatusFeed()
			cachedNodeResourceFeedSub = cacheIntface.DataProvider().SubscribeToCachedNodeResourceFeed()
			cachedDeviceStatusFeedSub = cacheIntface.DataProvider().SubscribeToCachedDeviceStatusFeed()
			cachedStreamStatusFeedSub = cacheIntface.DataProvider().SubscribeToCachedStreamStatusFeed()
		)

		defer func() {
			cacheIntface.DataProvider().UnsubscribeFromCachedNodeStatusFeed(cachedNodeStatusFeedSub)
			cacheIntface.DataProvider().UnsubscribeFromCachedNodeResourceFeed(cachedNodeResourceFeedSub)
			cacheIntface.DataProvider().UnsubscribeFromCachedDeviceStatusFeed(cachedDeviceStatusFeedSub)
			cacheIntface.DataProvider().UnsubscribeFromCachedStreamStatusFeed(cachedStreamStatusFeedSub)

			rh.baseBundle.WSHub.CloseWebsocketConnection(wsClient)
		}()

		wsLoop:
		for {
			select {
			case recvBytes, open := <- wsClient.ReceiverChannel():
				if !open { break wsLoop }

				var reqHeader = wsUtil.ExtractBasicReqHeader(recvBytes)

				switch reqHeader.ReqCode {
				case td__wsReqCode__nodeInfoListing:
					rh.td__wsHandleReqNodeInfoListing(wsClient, reqHeader)

				case td__wsReqCode__deviceInfoListing:
					rh.td__wsHandleReqDeviceInfoListing(wsClient, reqHeader)

				case td__wsReqCode__streamInfoListing:
					rh.td__wsHandleReqStreamInfoListing(wsClient, reqHeader)

				case td__wsReqCode__streamInfoItem:
					rh.td__wsHandleReqStreamInfoItem(wsClient, recvBytes)

				default:
					rh.logger.Warnf("%s: tempDashboard__wsHandler, no handler for %s", rh.logPrefix, reqHeader.ReqCode)
					break
				}

			case dat := <- cachedNodeStatusFeedSub.Channel:
				var nodeStatusInfo = &td__nodeStatusInfo{}
				nodeStatusInfo.mapFromCachedNodeStatus(dat)

				_ = wsClient.Send(wsUtil.EncodeBasicDatafeedMessage(
					wsUtil.BasicDataFeedHeader{To: "n", Topic: "stat"},
					&td__wsDatafeed__nodeStatusInfo{nodeStatusInfo},
				))

			case dat := <- cachedNodeResourceFeedSub.Channel:
				var nodeResourceInfo = &td__nodeResourceInfo{}
				nodeResourceInfo.mapFromCachedNodeResource(dat)

				_ = wsClient.Send(wsUtil.EncodeBasicDatafeedMessage(
					wsUtil.BasicDataFeedHeader{To: "n", Topic: "res"},
					nodeResourceInfo,
				))

			case dat := <- cachedDeviceStatusFeedSub.Channel:
				var deviceStatusInfo = &td__deviceStatusInfo{}
				deviceStatusInfo.mapFromCachedDeviceStatus(dat)

				_ = wsClient.Send(wsUtil.EncodeBasicDatafeedMessage(
					wsUtil.BasicDataFeedHeader{To: "d", Topic: "stat"},
					&td__wsDatafeed__deviceStatusInfo{deviceStatusInfo},
				))

			case dat := <- cachedStreamStatusFeedSub.Channel:
				var streamStatusInfo = &td__streamStatusInfo{}
				streamStatusInfo.mapFromCachedStreamStatus(dat)

				_ = wsClient.Send(wsUtil.EncodeBasicDatafeedMessage(
					wsUtil.BasicDataFeedHeader{To: "s", Topic: "stat"},
					&td__wsDatafeed__streamStatusInfo{streamStatusInfo},
				))
			}
		}
	}()
}

func (rh *ModuleRouteHandler) td__wsHandleReqNodeInfoListing(wsClient *ws.WSClient, reqHeader *wsUtil.BasicReqHeader) {
	var cachedNodes = cacheIntface.DataProvider().CachedNodes()
	var nodeInfoListing = make([]*td__nodeInfo, 0, len(cachedNodes))

	for _, _cachedNode := range cachedNodes {
		var nodeInfo = &td__nodeInfo{}
		nodeInfo.mapFromCachedNode(_cachedNode)

		nodeInfoListing = append(nodeInfoListing, nodeInfo)
	}

	_ = wsClient.Send(wsUtil.EncodeBasicRepMessage(
		reqHeader,
		td__wsRepPayload__nodeInfoListing{nodeInfoListing},
	))
}

func (rh *ModuleRouteHandler) td__wsHandleReqDeviceInfoListing(wsClient *ws.WSClient, reqHeader *wsUtil.BasicReqHeader) {
	var cachedDevices = cacheIntface.DataProvider().TempCachedDevicesAll()
	var deviceInfoListing = make([]*td__deviceInfo, 0, len(cachedDevices))

	for _, _cachedDevice := range cachedDevices {
		var deviceInfo = &td__deviceInfo{}
		deviceInfo.mapFromCachedDevice(_cachedDevice)

		deviceInfoListing = append(deviceInfoListing, deviceInfo)
	}

	_ = wsClient.Send(wsUtil.EncodeBasicRepMessage(
		reqHeader,
		td__wsRepPayload__deviceInfoListing{deviceInfoListing},
	))
}

func (rh *ModuleRouteHandler) td__wsHandleReqStreamInfoListing(wsClient *ws.WSClient, reqHeader *wsUtil.BasicReqHeader) {
	var cachedStreams = cacheIntface.DataProvider().TempCachedStreamsAll()
	var streamInfoListing = make([]*td__streamInfo, 0, len(cachedStreams))

	for _, _cachedStream := range cachedStreams {
		var streamInfo = &td__streamInfo{}
		streamInfo.mapFromCachedStream(_cachedStream)

		streamInfoListing = append(streamInfoListing, streamInfo)
	}

	_ = wsClient.Send(wsUtil.EncodeBasicRepMessage(
		reqHeader,
		td__wsRepPayload__streamInfoListing{streamInfoListing},
	))
}

func (rh *ModuleRouteHandler) td__wsHandleReqStreamInfoItem(wsClient *ws.WSClient, recvBytes []byte) {
	var reqMessage = &td__wsReqMessage__streamInfoItem{}
	wsUtil.DecodeTypedReqMessage(recvBytes, reqMessage)

	var cachedStream = cacheIntface.DataProvider().CachedStream(reqMessage.Payload.NodeID, reqMessage.Payload.StreamCode)

	var streamInfo = &td__streamInfo{}
	streamInfo.mapFromCachedStream(cachedStream)

	_ = wsClient.Send(wsUtil.EncodeBasicRepMessage(
		reqMessage.Header,
		td__wsRepPayload__streamInfoItem{streamInfo},
	))
}

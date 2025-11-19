package route

import (
	"encoding/json"
	"strings"

	"github.com/valyala/fasthttp"

	wsUtil "noname001/web/base/comm/ws/util"

	wallService "noname001/app/module/feature/wall/service"
)

// TODO: ALL err handling via ws

func (rh *ModuleRouteHandler) wallDetail__ws(ctx *fasthttp.RequestCtx) {
	var handlerName = "wallDetail__ws"
	var (
		host     = ctx.Request.Header.Peek("Host")
		hostPart = strings.Split(string(host), ":")

		wallID   = string(ctx.QueryArgs().Peek("id"))
	)

	wsClient, err := rh.baseBundle.WSHub.UpgradeToWebsocketConnection(ctx)
	if err != nil {
		rh.logger.Errorf("%s: %s, upgrade err: %s", rh.logPrefix, handlerName, err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	go func() {
		var (
			wallServiceInstance = wallService.Instance()

			fnHandleReqWallInfo func([]byte) = nil
			fnHandleReqWallItemInfo func([]byte) = nil

			fnHandleReqNodeInfoListing func([]byte) = nil
			fnHandleReqStreamInfoListing func([]byte) = nil
		)

		defer func() {
		}()

		fnHandleReqWallInfo = func(recvBytes []byte) {
			basicReqRep := wsUtil.ExtractBasicReqRep(recvBytes)

			wall, wallMessages := wallServiceInstance.Wall__Find(wallID, true)
			if wallMessages.HasError() {
				rh.logger.Errorf("%s: %s, Wall__Find err: %s", rh.logPrefix, handlerName, wallMessages.Dump())
				return
			}

			wallLayout, wallLayoutMessages := wallServiceInstance.WallLayout__Find(wall.WallLayoutID)
			if wallLayoutMessages.HasError() {
				rh.logger.Errorf("%s: %s, WallLayout__Find err: %s", rh.logPrefix, handlerName, wallLayoutMessages.Dump())
				return
			}

			wallItemInfoList := make([]*wd__wallItemInfo, 0, len(wall.Items))
			for _, _wallItem := range wall.Items {
				var (
					nodeInfo   *wd__nodeInfo
					streamInfo *wd__streamInfo

					cachedNode   = wallServiceInstance.GetCachedNode(_wallItem.SourceNodeID)
					cachedStream = wallServiceInstance.GetCachedStream(_wallItem.SourceNodeID, _wallItem.StreamCode)
				)

				if cachedNode != nil {
					nodeInfo = &wd__nodeInfo{
						ID   : cachedNode.NodeSnapshot.ID,
						State: string(cachedNode.NodeSnapshot.State),
	
						StreamCount   : len(cachedNode.CachedStreams),
	
						LastActivityAt: cachedNode.LastActivityAt,
					}
				}

				if cachedStream != nil {
					streamInfo = &wd__streamInfo{
						ID: cachedStream.ID,
						Code: cachedStream.Code,

						SourceType: cachedStream.StreamSnapshot.Persistence.SourceType,

						StreamerState: string(cachedStream.StreamSnapshot.Live.StreamerState),
						EstimatedVideoBitrate: cachedStream.StreamSnapshot.Live.EstimatedOutputVideoBitrate,

						LastActivityAt: cachedStream.LastActivityAt,

						// PreviewURL: wallServiceInstance.GetRelayedStreamViewURL(hostPart[0], _wallItem.SourceNodeID, _wallItem.StreamCode, "hls"),
					}
				}

				wallItemInfoList = append(wallItemInfoList, &wd__wallItemInfo{
					ID   : _wallItem.ID,
					Index: _wallItem.Index,

					SourceNode: _wallItem.SourceNodeID,
					SourceStream: _wallItem.StreamCode,

					NodeInfo: nodeInfo,
					StreamInfo: streamInfo,
				})
			}

			wallInfo := &wd__wallInfo{
				LayoutCode: wallLayout.Code,
				Items: wallItemInfoList,
			}

			repStruct := wd__wsReply__wallInfo{basicReqRep, wallInfo}
			_ = wsClient.Send(wsUtil.NewReqRepJson(repStruct))
		}

		fnHandleReqWallItemInfo = func(recvBytes []byte) {
			var reqStruct *wd__wsRequest__wallItemInfo

			unmarshalErr := json.Unmarshal(recvBytes, &reqStruct)
			if unmarshalErr != nil {
				rh.logger.Errorf("%s: %s, unmarshalErr on %s", rh.logPrefix, handlerName, wd__wsRequestCode__streamInfoListing)
				return
			}

			basicReqRep := wsUtil.BasicReqRep{reqStruct.RequestID, reqStruct.RequestCode}

			_wallItem, wallItemMessages := wallServiceInstance.WallItem__Find(reqStruct.WallItemID)
			if wallItemMessages.HasError() {
				rh.logger.Errorf("%s: %s, WallItem__Find err: %s", rh.logPrefix, handlerName, wallItemMessages.Dump())
				return
			}

			var (
				nodeInfo   *wd__nodeInfo
				streamInfo *wd__streamInfo

				cachedNode   = wallServiceInstance.GetCachedNode(_wallItem.SourceNodeID)
				cachedStream = wallServiceInstance.GetCachedStream(_wallItem.SourceNodeID, _wallItem.StreamCode)
			)

			if cachedNode != nil {
				nodeInfo = &wd__nodeInfo{
					ID   : cachedNode.NodeSnapshot.ID,
					State: string(cachedNode.NodeSnapshot.State),

					StreamCount   : len(cachedNode.CachedStreams),

					LastActivityAt: cachedNode.LastActivityAt,
				}
			}

			if cachedStream != nil {
				streamInfo = &wd__streamInfo{
					ID: cachedStream.ID,
					Code: cachedStream.Code,

					SourceType: cachedStream.StreamSnapshot.Persistence.SourceType,

					StreamerState: string(cachedStream.StreamSnapshot.Live.StreamerState),
					EstimatedVideoBitrate: cachedStream.StreamSnapshot.Live.EstimatedOutputVideoBitrate,

					LastActivityAt: cachedStream.LastActivityAt,

					// PreviewURL: wallServiceInstance.GetRelayedStreamViewURL(hostPart[0], _wallItem.SourceNodeID, _wallItem.StreamCode, "hls"),
				}
			}

			wallItemInfo := &wd__wallItemInfo{
				ID   : _wallItem.ID,
				Index: _wallItem.Index,

				SourceNode: _wallItem.SourceNodeID,
				SourceStream: _wallItem.StreamCode,

				NodeInfo: nodeInfo,
				StreamInfo: streamInfo,
			}

			repStruct := wd__wsReply__wallItemInfo{basicReqRep, wallItemInfo}
			_ = wsClient.Send(wsUtil.NewReqRepJson(repStruct))
		}

		fnHandleReqNodeInfoListing = func(recvBytes []byte) {
			basicReqRep := wsUtil.ExtractBasicReqRep(recvBytes)

			cachedNodes := wallServiceInstance.GetCachedNodes()
			nodeInfoListing := make([]*wd__nodeInfo, 0, len(cachedNodes))
			for _, _cachedNode := range cachedNodes {
				nodeInfoListing = append(nodeInfoListing, &wd__nodeInfo{
					ID   : _cachedNode.NodeSnapshot.ID,
					State: string(_cachedNode.NodeSnapshot.State),

					StreamCount   : len(_cachedNode.CachedStreams),

					LastActivityAt: _cachedNode.LastActivityAt,
				})
			}

			repStruct := wd__wsReply__nodeInfoListing{basicReqRep, nodeInfoListing}
			_ = wsClient.Send(wsUtil.NewReqRepJson(repStruct))
		}

		fnHandleReqStreamInfoListing = func(recvBytes []byte) {
			var reqStruct *wd__wsRequest__streamInfoListing

			unmarshalErr := json.Unmarshal(recvBytes, &reqStruct)
			if unmarshalErr != nil {
				rh.logger.Errorf("%s: %s, unmarshalErr on %s", rh.logPrefix, handlerName, wd__wsRequestCode__streamInfoListing)
				return
			}

			basicReqRep := wsUtil.BasicReqRep{reqStruct.RequestID, reqStruct.RequestCode}

			cachedStreams := wallServiceInstance.GetCachedStreams(reqStruct.NodeID)
			streamInfoListing := make([]*wd__streamInfo, 0, len(cachedStreams))
			for _, _cachedStream := range cachedStreams {
				streamInfoListing = append(streamInfoListing, &wd__streamInfo{
					ID: _cachedStream.ID,
					Code: _cachedStream.Code,

					SourceType: _cachedStream.StreamSnapshot.Persistence.SourceType,

					StreamerState: string(_cachedStream.StreamSnapshot.Live.StreamerState),
					EstimatedVideoBitrate: _cachedStream.StreamSnapshot.Live.EstimatedOutputVideoBitrate,

					LastActivityAt: _cachedStream.LastActivityAt,

					PreviewURL: wallServiceInstance.GetRelayedStreamViewURL(hostPart[0], reqStruct.NodeID, _cachedStream.Code, "hls"),
				})
			}

			repStruct := wd__wsReply__streamInfoListing{basicReqRep, streamInfoListing}
			_ = wsClient.Send(wsUtil.NewReqRepJson(repStruct))
		}

		wsLoop:
		for {
			select {
			case recvBytes, open := <- wsClient.ReceiverChannel():
				if !open { break wsLoop }

				basicReqRep := wsUtil.ExtractBasicReqRep(recvBytes)

				switch basicReqRep.RequestCode {
				case wd__wsRequestCode__wallInfo:     fnHandleReqWallInfo(recvBytes)
				case wd__wsRequestCode__wallItemInfo: fnHandleReqWallItemInfo(recvBytes)

				case wd__wsRequestCode__nodeInfoListing:   fnHandleReqNodeInfoListing(recvBytes)
				case wd__wsRequestCode__streamInfoListing: fnHandleReqStreamInfoListing(recvBytes)

				default:
					rh.logger.Warnf("%s: %s, no handler for %s", rh.logPrefix, handlerName, basicReqRep)
					break
				}
			}
		}
	}()
}

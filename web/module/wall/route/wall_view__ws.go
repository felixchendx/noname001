package route

import (
	"encoding/json"
	"time"

	"github.com/valyala/fasthttp"

	wsUtil "noname001/web/base/comm/ws/util"

	cacheIntface "noname001/app/module/common/cache/intface"

	wallService "noname001/app/module/feature/wall/service"
)

// TODO: err handling via ws

// TODO: retype stuffs that's used in js, as constant dependency reminders
//       and use either ws-onopen or localapi to generate dynamic spec for initial js setup
const (
	// case "/spec":
	wallView__wsRequestCode__wallInfo     string = "/wall/info"
	wallView__wsRequestCode__wallItemInfo string = "/wall/item/info"
)

type wallView__wsEvent struct {
	Timestamp time.Time `json:"ts"`
	EventCode string    `json:"ev_code"`

	ItemIndex   int      `json:"item_index"`
	NodeID      string   `json:"node_id"`
	StreamCode  string   `json:"stream_code"`
	StreamState string   `json:"stream_state"`
}

type wallView__wsReply__wallInfo struct {
	wsUtil.BasicReqRep

	WallInfo *wallView__wallInfo `json:"wall_info"`
}

type wallView__wsRequest__wallItemInfo struct {
	wsUtil.BasicReqRep

	ItemIndex int `json:"item_index"`
}
type wallView__wsReply__wallItemInfo struct {
	wsUtil.BasicReqRep

	WallItemInfo *wallView__wallItemInfo `json:"wall_item_info"`
}

type wallView__wallInfo struct {
	Code string `json:"code"`
	Name string `json:"name"`

	Items []*wallView__wallItemInfo `json:"items"`
}
type wallView__wallItemInfo struct {
	ItemIndex int `json:"item_index"`

	SourceNode string `json:"source_node"`
	StreamCode string `json:"stream_code"`

	StreamState string `json:"stream_state"`
}

// this whole ws handler needs to be simplify, do it after cache + ev is kinda established
// which means, after minimal monitoring dashboard is done...
func (rh *ModuleRouteHandler) wallView__ws(ctx *fasthttp.RequestCtx) {
	wsClient, err := rh.baseBundle.WSHub.UpgradeToWebsocketConnection(ctx)
	if err != nil {
		rh.logger.Errorf("%s: wallView__ws, upgrade err: %s", rh.logPrefix, err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	wallID := string(ctx.QueryArgs().Peek("id"))

	go func() {
		var (
			wallServiceInstance = wallService.Instance()

			currWall          *wallService.WallDE                  = nil
			relItemsByNodes   map[string][]*wallService.WallItemDE = make(map[string][]*wallService.WallItemDE)
			relItemsByStreams map[string][]*wallService.WallItemDE = make(map[string][]*wallService.WallItemDE)

			wallInfo *wallView__wallInfo = nil
			// cachedWall *wallView__cachedWall = nil

			cachedStreamEvSub = cacheIntface.EventProvider().SubscribeToCachedStreamEvent()
		)

		defer func() {
			cacheIntface.EventProvider().UnsubscribeFromCachedStreamEvent(cachedStreamEvSub)
		}()

		// TODO: repopulate these on wall info changes T_T
		{ // pre-process relevant infos
			_wall, _wallMessages := wallServiceInstance.Wall__Find(wallID, true)
			if _wallMessages.HasError() {
				rh.logger.Errorf("%s: handleWallViewWSConn, Wall__Find err: %s", rh.logPrefix, _wallMessages.Dump())
				return
			}
			currWall = _wall

			wallItemInfoList := make([]*wallView__wallItemInfo, 0, len(currWall.Items))
			for _index, _wallItem := range currWall.Items {
				{
					relItems, inMap := relItemsByNodes[_wallItem.SourceNodeID]
					if inMap {
						relItems = append(relItems, _wallItem)
						relItemsByNodes[_wallItem.SourceNodeID] = relItems
					} else {
						newList := make([]*wallService.WallItemDE, 0, 2)
						newList = append(newList, _wallItem)
						relItemsByNodes[_wallItem.SourceNodeID] = newList
					}
				}

				{
					relItems, inMap := relItemsByStreams[_wallItem.StreamCode]
					if inMap {
						relItems = append(relItems, _wallItem)
						relItemsByStreams[_wallItem.StreamCode] = relItems
					} else {
						newList := make([]*wallService.WallItemDE, 0, 4)
						newList = append(newList, _wallItem)
						relItemsByStreams[_wallItem.StreamCode] = newList
					}
				}

				{
					var streamState string = ""

					if _wallItem.SourceNodeID != "" && _wallItem.StreamCode != "" {
						cachedStream := wallServiceInstance.GetCachedStream(_wallItem.SourceNodeID, _wallItem.StreamCode)
						if cachedStream != nil {
							streamState = string(cachedStream.StreamSnapshot.Live.State)
						}
					}

					wallItemInfoList = append(wallItemInfoList, &wallView__wallItemInfo{
						ItemIndex: _index,
	
						SourceNode: _wallItem.SourceNodeID,
						StreamCode: _wallItem.StreamCode,
	
						StreamState: streamState,
					})
				}
			}

			wallInfo = &wallView__wallInfo{
				Code: currWall.Code,
				Name: currWall.Name,

				Items: wallItemInfoList,
			}
		}

		loopityLoop:
		for {
			select {
			case recvBytes, open := <- wsClient.ReceiverChannel():
				if !open { break loopityLoop }

				basicReqRep := wsUtil.ExtractBasicReqRep(recvBytes)

				switch basicReqRep.RequestCode {
				case wallView__wsRequestCode__wallInfo:
					repStruct := wallView__wsReply__wallInfo{basicReqRep, wallInfo}
					_ = wsClient.Send(wsUtil.NewReqRepJson(repStruct))

				case wallView__wsRequestCode__wallItemInfo:
					var itemInfoReq *wallView__wsRequest__wallItemInfo

					err := json.Unmarshal(recvBytes, &itemInfoReq)
					if err != nil {
						rh.logger.Errorf("%s: wallView__ws, itemInfoReq unmarshal err %s", rh.logPrefix, err.Error())
						break
					}

					if itemInfoReq.ItemIndex < 0 || itemInfoReq.ItemIndex >= len(wallInfo.Items) {
						rh.logger.Warnf(
							"%s: wallView__ws, itemInfoReq index out of range. given: %v, len: %v",
							rh.logPrefix, itemInfoReq.ItemIndex, len(wallInfo.Items),
						)
						break
					}

					repStruct := wallView__wsReply__wallItemInfo{
						basicReqRep,
						wallInfo.Items[itemInfoReq.ItemIndex],
					}
					_ = wsClient.Send(wsUtil.NewReqRepJson(repStruct))

				default:
					rh.logger.Warnf("%s: wallView__ws, no handler for %s", rh.logPrefix, basicReqRep)
					break
				}

			case ev := <- cachedStreamEvSub.Channel:
				_, nodeRelated := relItemsByNodes[ev.OriginalStreamEvent.NodeID]
				if !nodeRelated { break }

				_, streamRelated := relItemsByStreams[ev.OriginalStreamEvent.StreamCode]
				if !streamRelated { break }

				// wonky
				for _index, _wallItemInfo := range wallInfo.Items {
					if _wallItemInfo.SourceNode == ev.OriginalStreamEvent.NodeID && _wallItemInfo.StreamCode == ev.OriginalStreamEvent.StreamCode {
						cachedStream := wallServiceInstance.GetCachedStream(_wallItemInfo.SourceNode, _wallItemInfo.StreamCode)
						if cachedStream != nil {
							_wallItemInfo.StreamState = string(cachedStream.StreamSnapshot.Live.State)

							evStruct := wallView__wsEvent{
								Timestamp: ev.OriginalStreamEvent.Timestamp,
								EventCode: string(ev.OriginalStreamEvent.EventCode),

								ItemIndex  : _index,
								NodeID     : ev.OriginalStreamEvent.NodeID,
								StreamCode : ev.OriginalStreamEvent.StreamCode,
								StreamState: _wallItemInfo.StreamState,
							}

							_ = wsClient.Send(wsUtil.NewEventJson(evStruct))
						}

						break // for
					}
				}
			}
		}
	}()
}

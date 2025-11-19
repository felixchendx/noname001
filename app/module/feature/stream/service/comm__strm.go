package service

import (
	"noname001/app/base/messaging"

	streamTyping "noname001/app/base/typing/stream"

	streamMdp "noname001/app/node/comm/stream/mdpservice"
)

// ============ VVV conform to streamMdp.DataHandlerIntface VVV ============= //
func (svc *Service) ProvideServiceInfo() (*streamMdp.ServiceInfo, *messaging.Messages) {
	messages := messaging.NewMessages()
	serviceInfo := &streamMdp.ServiceInfo{
		Placeholder: "s-nwoo",
	}

	return serviceInfo, messages
}

func (svc *Service) ProvideStreamSnapshotList() ([]*streamTyping.StreamSnapshot, *messaging.Messages) {
	messages := messaging.NewMessages()
	streamSnapshotList := make([]*streamTyping.StreamSnapshot, 0)

	liveStreams := svc.coordinator.GetLiveStreams()
	for _, liveStream := range liveStreams {
		streamSnapshotList = append(streamSnapshotList, liveStream.Snapshot())
	}

	return streamSnapshotList, messages
}

func (svc *Service) ProvideStreamSnapshot(streamCode string) (*streamTyping.StreamSnapshot, *messaging.Messages) {
	messages := messaging.NewMessages()

	liveStream, found := svc.coordinator.GetLiveStreamByCode(streamCode)
	if !found {
		messages.AddError(COMMSTRM_ERR_90003.NewMessage(streamCode))
		return nil, messages
	}

	return liveStream.Snapshot(), messages
}
// ============ ^^^ conform to streamMdp.DataHandlerIntface ^^^ ============= //

func (svc *Service) _getStreamItemWithIdentifier(ider *StreamItemIdentifier) (*StreamItemDE, *messaging.Messages) {
	messages := messaging.NewMessages()

	if ider == nil || (ider.ID == "" && ider.Code == "") {
		messages.AddError(COMMSTRM_ERR_90001.NewMessage())
		return nil, messages
	}

	var de *StreamItemDE

	if ider.ID != "" {
		pe, dbev := svc.store.DB.StreamItem__Get(ider.ID)
		if dbev.IsError() {
			messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
			return nil, messages
		}
		if pe == nil {
			messages.AddError(COMMSTRM_ERR_90002.NewMessage(ider.ID))
			return nil, messages
		}

		de = (&StreamItemDE{}).fromPE(pe)

	} else if ider.Code != "" {
		pe, dbev := svc.store.DB.StreamItem__GetByCode(ider.Code)
		if dbev.IsError() {
			messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
			return nil, messages
		}
		if pe == nil {
			messages.AddError(COMMSTRM_ERR_90003.NewMessage(ider.Code))
			return nil, messages
		}

		de = (&StreamItemDE{}).fromPE(pe)
	}

	return de, messages
}

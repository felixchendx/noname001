package service

import (
	"github.com/google/uuid"

	"noname001/app/base/messaging"
	"noname001/app/base/sanitation"
	"noname001/app/constant"

	streamEv "noname001/app/module/feature/stream/event"
)

func (svc *Service) EmptyStreamItem() (*StreamItemDE) {
	return &StreamItemDE{}
}
func (svc *Service) EmptyStreamItemList() ([]*StreamItemDE) {
	return make([]*StreamItemDE, 0)
}

func (svc *Service) AddStreamItem(de *StreamItemDE) (*StreamItemDE, *messaging.Messages) {
	messages := messaging.NewMessages()

	// sanitizing
	de.sanitize()

	// standard validating
	if de.Code == "" {
		messages.AddError(STRMITMBSC_ERR_11501.NewMessage())
	} else {
		isIllegal, illegalChar := sanitation.Code_ContainsIllegalChar(de.Code)
		if isIllegal {
			messages.AddError(STRMITMBSC_ERR_11511.NewMessage(illegalChar, sanitation.CODE__LEGAL_CHARS))
		}
	}
	if de.SourceType == "" {
		messages.AddError(STRMITMBSC_ERR_11502.NewMessage())
	} else {
		switch de.SourceType {
		case string(constant.STREAM__SOURCE_TYPE__MOD_DEVICE):
			if de.DeviceCode == ""      { messages.AddError(STRMITMBSC_ERR_11503.NewMessage()) }
			if de.DeviceChannelID == "" { messages.AddError(STRMITMBSC_ERR_11504.NewMessage()) }
		case string(constant.STREAM__SOURCE_TYPE__EXTERNAL):
			if de.ExternalURL == "" { messages.AddError(STRMITMBSC_ERR_11505.NewMessage()) }
		case string(constant.STREAM__SOURCE_TYPE__FILE):
			if de.Filepath == "" { messages.AddError(STRMITMBSC_ERR_11506.NewMessage()) }
		// case string(constant.STREAM__SOURCE_TYPE__EMBEDDED_FILE):
		}
	}
	if messages.HasError() { return nil, messages }

	// TODO: others
	// contextual validating
	// TODO: validate stream group ref
	// TODO: validate file and external stream path
	// TODO: validate stream type for each brand...
	existingSI, dbev1 := svc.store.DB.StreamItem__GetByCode(de.Code)
	switch {
	case dbev1.IsError()  : messages.AddError(STRDB_ERR_00001.NewMessage(dbev1.EventID()))
	case existingSI != nil: messages.AddError(STRMITMBSC_ERR_11551.NewMessage(de.Code))
	}
	if messages.HasError() { return nil, messages }

	// defaulting
	de.ID = uuid.New().String()
	if de.Name == ""  { de.Name = de.Code }
	if de.State == "" { de.State = constant.ENTITY__STATE_INACTIVE }

	pe, dbev := svc.store.DB.StreamItem__AtomicInsert(de.toPE())
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return nil, messages
	}

	svc.evHub.PublishStreamItemEvent(
		streamEv.STREAM_ITEM_EVENT_CODE__CREATE,
		pe.ID, pe.Code,
	)

	messages.AddNotice(STRMITMBSC_NTC_11101.NewMessage(pe.Code))
	return de.fromPE(pe), messages
}

func (svc *Service) FindStreamItem(id string) (*StreamItemDE, *messaging.Messages) {
	messages := messaging.NewMessages()

	pe, dbev := svc.store.DB.StreamItem__Get(id)
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return nil, messages
	}
	if pe == nil {
		messages.AddError(STRMITMBSC_ERR_03001.NewMessage(id))
		return nil, messages
	}

	return (&StreamItemDE{}).fromPE(pe), messages
}

func (svc *Service) FindStreamItemByCode(code string) (*StreamItemDE, *messaging.Messages) {
	messages := messaging.NewMessages()

	pe, dbev := svc.store.DB.StreamItem__GetByCode(code)
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return nil, messages
	}
	if pe == nil {
		messages.AddError(STRMITMBSC_ERR_04001.NewMessage(code))
		return nil, messages
	}

	return (&StreamItemDE{}).fromPE(pe), messages
}

func (svc *Service) FindStreamItemsByStreamGroupID(groupID string) ([]*StreamItemDE, *messaging.Messages) {
	messages := messaging.NewMessages()

	listDE := make([]*StreamItemDE, 0)
	listPE, dbev := svc.store.DB.StreamItem__GetByStreamGroupID(groupID)
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return nil, messages
	}

	for _, pe := range listPE {
		de := (&StreamItemDE{}).fromPE(pe)
		listDE = append(listDE, de)
	}

	return listDE, messages
}

// func (svc *Service) EditStreamItemPrecheck() (*messaging.Messages) {
// 	messages := messaging.NewMessages()
// 	return messages
// }

func (svc *Service) EditStreamItem(id string, de *StreamItemDE) (*StreamItemDE, *messaging.Messages) {
	messages := messaging.NewMessages()

	pe, dbev := svc.store.DB.StreamItem__Get(id)
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return nil, messages
	}
	if pe == nil {
		messages.AddError(STRMITMBSC_ERR_14501.NewMessage(id))
		return nil, messages
	}

	// sanitizing
	de.sanitize()

	// standard validating
	if de.SourceType == "" {
		messages.AddError(STRMITMBSC_ERR_14502.NewMessage())
	} else {
		switch de.SourceType {
		case string(constant.STREAM__SOURCE_TYPE__MOD_DEVICE):
			if de.DeviceCode == ""       { messages.AddError(STRMITMBSC_ERR_14503.NewMessage()) }
			if de.DeviceChannelID == ""  { messages.AddError(STRMITMBSC_ERR_14504.NewMessage()) }
		case string(constant.STREAM__SOURCE_TYPE__EXTERNAL):
			if de.ExternalURL == "" { messages.AddError(STRMITMBSC_ERR_14505.NewMessage()) }
		case string(constant.STREAM__SOURCE_TYPE__FILE):
			if de.Filepath == "" { messages.AddError(STRMITMBSC_ERR_14506.NewMessage()) }
		// case string(constant.STREAM__SOURCE_TYPE__EMBEDDED_FILE):
		}
	}
	if messages.HasError() { return nil, messages }

	// TODO: validate file and external stream path
	// TODO: on brand change, also convert stream type
	// TODO: validate stream type for each brand...

	// TODO: changes that affect downstreams

	// merge editable fields
	pe.Name  = de.Name
	pe.State = de.State
	pe.Note  = de.Note

	pe.SourceType       = de.SourceType
	pe.DeviceCode       = de.DeviceCode
	pe.DeviceChannelID  = de.DeviceChannelID
	pe.DeviceStreamType = de.DeviceStreamType
	pe.ExternalURL      = de.ExternalURL
	pe.Filepath         = de.Filepath
	pe.EmbeddedFilepath = de.EmbeddedFilepath

	// defaulting
	if pe.Name == ""  { pe.Name = pe.Code }
	if pe.State == "" { pe.State = constant.ENTITY__STATE_INACTIVE }

	pe, dbev = svc.store.DB.StreamItem__AtomicUpdate(pe)
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return nil, messages
	}

	svc.evHub.PublishStreamItemEvent(
		streamEv.STREAM_ITEM_EVENT_CODE__UPDATE,
		pe.ID, pe.Code,
	)

	messages.AddNotice(STRMITMBSC_NTC_14101.NewMessage(pe.Code))
	return de.fromPE(pe), messages
}

func (svc *Service) DeleteStreamItemPrecheck() (*messaging.Messages) {
	messages := messaging.NewMessages()
	return messages
}

func (svc *Service) DeleteStreamItem(id string) (*messaging.Messages) {
	messages := messaging.NewMessages()

	pe, dbev := svc.store.DB.StreamItem__Get(id)
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return messages
	}
	if pe == nil {
		messages.AddError(STRMITMBSC_ERR_06001.NewMessage(id))
		return messages
	}

	// TODO: validation

	if messages.HasError() { return messages }

	// TODO: changes that affect downstreams

	dbev = svc.store.DB.StreamItem__AtomicDelete(id)
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return messages
	}

	svc.evHub.PublishStreamItemEvent(
		streamEv.STREAM_ITEM_EVENT_CODE__DELETE,
		pe.ID, pe.Code,
	)

	messages.AddNotice(STRMITMBSC_NTC_06002.NewMessage(pe.Code))
	return messages
}

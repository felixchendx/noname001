package service

import (
	"github.com/google/uuid"

	"noname001/app/base/messaging"
	"noname001/app/base/sanitation"
	"noname001/app/constant"

	streamEv "noname001/app/module/feature/stream/event"
)

func (svc *Service) EmptyStreamGroup() (*StreamGroupDE) {
	return &StreamGroupDE{}
}
func (svc *Service) EmptyStreamGroupList() ([]*StreamGroupDE) {
	return make([]*StreamGroupDE, 0)
}

func (svc *Service) AddStreamGroup(de *StreamGroupDE) (*StreamGroupDE, *messaging.Messages) {
	messages := messaging.NewMessages()

	// sanitizing
	de.sanitize()

	// standard validating
	if de.Code == "" {
		messages.AddError(STRMGRPBSC_ERR_11501.NewMessage())
	} else {
		isIllegal, illegalChar := sanitation.Code_ContainsIllegalChar(de.Code)
		if isIllegal {
			messages.AddError(STRMGRPBSC_ERR_11511.NewMessage(illegalChar, sanitation.CODE__LEGAL_CHARS))
		}
	}
	if messages.HasError() { return nil, messages }

	// TODO: other contextual stuffs, ie reference validating
	// contextual validating
	existingSG, dbev1 := svc.store.DB.StreamGroup__GetByCode(de.Code)
	switch {
	case dbev1.IsError()  : messages.AddError(STRDB_ERR_00001.NewMessage(dbev1.EventID()))
	case existingSG != nil: messages.AddError(STRMGRPBSC_ERR_11551.NewMessage(de.Code))
	}
	if messages.HasError() { return nil, messages }

	// defaulting
	de.ID = uuid.New().String()
	if de.Name == "" { de.Name = de.Code }
	if de.State == "" { de.State = constant.ENTITY__STATE_INACTIVE }
	// TODO: defaulting de.StreamProfileID

	pe, dbev := svc.store.DB.StreamGroup__AtomicInsert(de.toPE())
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return nil, messages
	}

	svc.evHub.PublishStreamGroupEvent(
		streamEv.STREAM_GROUP_EVENT_CODE__CREATE,
		pe.ID, pe.Code,
	)

	messages.AddNotice(STRMGRPBSC_NTC_11101.NewMessage(pe.Code))
	return de.fromPE(pe), messages
}

func (svc *Service) FindStreamGroup(id string) (*StreamGroupDE, *messaging.Messages) {
	messages := messaging.NewMessages()

	pe, dbev := svc.store.DB.StreamGroup__Get(id)
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return nil, messages
	}
	if pe == nil {
		messages.AddError(STRMGRPBSC_ERR_03001.NewMessage(id))
		return nil, messages
	}

	return (&StreamGroupDE{}).fromPE(pe), messages
}

func (svc *Service) FindStreamGroupByCode(code string) (*StreamGroupDE, *messaging.Messages) {
	messages := messaging.NewMessages()

	pe, dbev := svc.store.DB.StreamGroup__GetByCode(code)
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return nil, messages
	}
	if pe == nil {
		messages.AddError(STRMGRPBSC_ERR_04001.NewMessage(code))
		return nil, messages
	}

	return (&StreamGroupDE{}).fromPE(pe), messages
}

// func (svc *Service) EditStreamGroupPrecheck() (*messaging.Messages) {
// 	messages := messaging.NewMessages()
// 	return messages
// }

func (svc *Service) EditStreamGroup(id string, de *StreamGroupDE) (*StreamGroupDE, *messaging.Messages) {
	messages := messaging.NewMessages()

	pe, dbev := svc.store.DB.StreamGroup__Get(id)
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return nil, messages
	}
	if pe == nil {
		messages.AddError(STRMGRPBSC_ERR_14501.NewMessage(id))
		return nil, messages
	}

	// sanitizing
	de.sanitize()

	// standard validating: none

	// TODO: contextual validating

	// TODO: changes that affect downstreams

	// merge editable fields
	pe.Name  = de.Name
	pe.State = de.State
	pe.Note  = de.Note
	pe.StreamProfileID = de.StreamProfileID

	// defaulting
	if pe.Name == ""  { pe.Name = pe.Code }
	if pe.State == "" { pe.State = constant.ENTITY__STATE_INACTIVE }
	// TODO: defaulting de.StreamProfileID

	pe, dbev = svc.store.DB.StreamGroup__AtomicUpdate(pe)
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return nil, messages
	}

	svc.evHub.PublishStreamGroupEvent(
		streamEv.STREAM_GROUP_EVENT_CODE__UPDATE,
		pe.ID, pe.Code,
	)

	messages.AddNotice(STRMGRPBSC_NTC_14101.NewMessage(pe.Code))
	return de.fromPE(pe), messages
}

// func (svc *Service) DeleteStreamGroupPrecheck() (*messaging.Messages) {
// 	messages := messaging.NewMessages()
// 	return messages
// }

func (svc *Service) DeleteStreamGroup(id string) (*messaging.Messages) {
	messages := messaging.NewMessages()

	pe, dbev := svc.store.DB.StreamGroup__Get(id)
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return messages
	}
	if pe == nil {
		messages.AddError(STRMGRPBSC_ERR_06001.NewMessage(id))
		return messages
	}

	// TODO: validation

	if messages.HasError() { return messages }

	// TODO: changes that affect downstreams

	dbev = svc.store.DB.StreamGroup__CascadingDelete(id)
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return messages
	}

	svc.evHub.PublishStreamGroupEvent(
		streamEv.STREAM_GROUP_EVENT_CODE__DELETE,
		pe.ID, pe.Code,
	)

	messages.AddNotice(STRMGRPBSC_NTC_06002.NewMessage(pe.Code))
	return messages
}

package service

import (
	"github.com/google/uuid"

	"noname001/app/base/messaging"
	"noname001/app/base/sanitation"
	"noname001/app/constant"

	streamEv "noname001/app/module/feature/stream/event"
)

func (svc *Service) EmptyStreamProfile() (*StreamProfileDE) {
	return &StreamProfileDE{}
}
func (svc *Service) EmptyStreamProfileList() ([]*StreamProfileDE) {
	return make([]*StreamProfileDE, 0)
}

func (svc *Service) AddStreamProfile(de *StreamProfileDE) (*StreamProfileDE, *messaging.Messages) {
	messages := messaging.NewMessages()

	// sanitizing
	de.sanitize()

	// standard validating
	if de.Code == "" {
		messages.AddError(STRMPRFLBSC_ERR_11501.NewMessage())
	} else {
		isIllegal, illegalChar := sanitation.Code_ContainsIllegalChar(de.Code)
		if isIllegal {
			messages.AddError(STRMPRFLBSC_ERR_11511.NewMessage(illegalChar, sanitation.CODE__LEGAL_CHARS))
		}
	}
	if messages.HasError() { return nil, messages }

	// TODO: other contextual stuffs
	// contextual validating
	existingSP, dbev1 := svc.store.DB.StreamProfile__GetByCode(de.Code)
	switch {
	case dbev1.IsError()  : messages.AddError(STRDB_ERR_00001.NewMessage(dbev1.EventID()))
	case existingSP != nil: messages.AddError(STRMPRFLBSC_ERR_11551.NewMessage(de.Code))
	}
	if messages.HasError() { return nil, messages }

	// defaulting
	de.ID = uuid.New().String()
	if de.Name == ""  { de.Name = de.Code }
	if de.State == "" { de.State = constant.ENTITY__STATE_INACTIVE }
	if de.TargetVideoCodec == "" { de.TargetVideoCodec = "h264" }
	if de.TargetAudioCodec == "" { de.TargetAudioCodec = "aac" }

	pe, dbev := svc.store.DB.StreamProfile__AtomicInsert(de.toPE())
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return nil, messages
	}

	svc.evHub.PublishStreamProfileEvent(
		streamEv.STREAM_PROFILE_EVENT_CODE__CREATE,
		pe.ID, pe.Code,
	)

	messages.AddNotice(STRMPRFLBSC_NTC_11101.NewMessage(pe.Code))
	return de.fromPE(pe), messages
}

func (svc *Service) FindStreamProfile(id string) (*StreamProfileDE, *messaging.Messages) {
	messages := messaging.NewMessages()

	pe, dbev := svc.store.DB.StreamProfile__Get(id)
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return nil, messages
	}
	if pe == nil {
		messages.AddError(STRMPRFLBSC_ERR_03001.NewMessage(id))
		return nil, messages
	}

	return (&StreamProfileDE{}).fromPE(pe), messages
}

func (svc *Service) FindStreamProfileByCode(code string) (*StreamProfileDE, *messaging.Messages) {
	messages := messaging.NewMessages()

	pe, dbev := svc.store.DB.StreamProfile__GetByCode(code)
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return nil, messages
	}
	if pe == nil {
		messages.AddError(STRMPRFLBSC_ERR_04001.NewMessage(code))
		return nil, messages
	}

	return (&StreamProfileDE{}).fromPE(pe), messages
}

// func (svc *Service) EditStreamProfilePrecheck() (*messaging.Messages) {
// 	messages := messaging.NewMessages()
// 	return messages
// }

func (svc *Service) EditStreamProfile(id string, de *StreamProfileDE) (*StreamProfileDE, *messaging.Messages) {
	messages := messaging.NewMessages()

	pe, dbev := svc.store.DB.StreamProfile__Get(id)
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return nil, messages
	}
	if pe == nil {
		messages.AddError(STRMPRFLBSC_ERR_14501.NewMessage(id))
		return nil, messages
	}

	// sanitizing
	de.sanitize()

	// standard validating: none ?

	// TODO: contextual validating

	// TODO: changes that affect downstreams

	// merge editable fields
	pe.Name  = de.Name
	pe.State = de.State
	pe.Note  = de.Note

	pe.TargetVideoCodec       = de.TargetVideoCodec
	pe.TargetVideoCompression = de.TargetVideoCompression
	pe.TargetVideoBitrate     = de.TargetVideoBitrate

	pe.TargetAudioCodec       = de.TargetAudioCodec
	pe.TargetAudioCompression = de.TargetAudioCompression
	pe.TargetAudioBitrate     = de.TargetAudioBitrate

	// defaulting
	if pe.Name == ""  { pe.Name = pe.Code }
	if pe.State == "" { pe.State = constant.ENTITY__STATE_INACTIVE }
	if pe.TargetVideoCodec == "" { pe.TargetVideoCodec = "h264" }
	if pe.TargetAudioCodec == "" { pe.TargetAudioCodec = "aac" }

	pe, dbev = svc.store.DB.StreamProfile__AtomicUpdate(pe)
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return nil, messages
	}

	svc.evHub.PublishStreamProfileEvent(
		streamEv.STREAM_PROFILE_EVENT_CODE__UPDATE,
		pe.ID, pe.Code,
	)

	messages.AddNotice(STRMPRFLBSC_NTC_14101.NewMessage(pe.Code))
	return de.fromPE(pe), messages
}

func (svc *Service) DeleteStreamProfilePrecheck(id string) (*messaging.Messages) {
	messages := messaging.NewMessages()

	streamGroupList, dbev := svc.store.DB.StreamGroup__GetByStreamProfileRef(id)
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return messages
	}

	if len(streamGroupList) > 0 {
		sgCodes := ""
		for _, item := range streamGroupList {
			if sgCodes != "" { sgCodes = sgCodes + ", " }
			sgCodes = sgCodes + item.Code
		}

		messages.AddError(STRMPRFLBSC_NTC_06010.NewMessage(sgCodes))
		return messages
	}

	return messages
}

func (svc *Service) DeleteStreamProfile(id string) (*messaging.Messages) {
	messages := messaging.NewMessages()

	pe, dbev := svc.store.DB.StreamProfile__Get(id)
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return messages
	}
	if pe == nil {
		messages.AddError(STRMPRFLBSC_ERR_06001.NewMessage(id))
		return messages
	}

	if _messages := svc.DeleteStreamProfilePrecheck(id); _messages.HasError() {
		messages.Append(_messages)
		return messages
	}

	if messages.HasError() { return messages }

	// TODO: changes that affect downstreams

	dbev = svc.store.DB.StreamProfile__AtomicDelete(id)
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return messages
	}

	svc.evHub.PublishStreamProfileEvent(
		streamEv.STREAM_PROFILE_EVENT_CODE__DELETE,
		pe.ID, pe.Code,
	)

	messages.AddNotice(STRMPRFLBSC_NTC_06002.NewMessage(pe.Code))
	return messages
}

func (svc *Service) SearchStreamProfile(sc *StreamProfile__SearchCriteria) (*StreamProfile__SearchResult, *messaging.Messages) {
	messages := messaging.NewMessages()

	sr, dbev := svc.store.DB.StreamProfile__Search(sc)
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return nil, messages
	}

	return (&StreamProfile__SearchResult{}).fromStore(sr), messages
}

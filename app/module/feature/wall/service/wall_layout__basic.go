package service

import (
	"noname001/app/base/messaging"
)

func (svc *Service) WallLayout__Search(sc *WallLayout__SearchCriteria) (*WallLayout__SearchResult, *messaging.Messages) {
	messages := messaging.NewMessages()

	sr, dbev := svc.store.DB.WallLayout__Search(sc)
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return nil, messages
	}

	return (&WallLayout__SearchResult{}).fromStore(sr), messages
}

func (svc *Service) WallLayout__Find(id string) (*WallLayoutDE, *messaging.Messages) {
	messages := messaging.NewMessages()

	pe, dbev := svc.store.DB.WallLayout__Get(id)
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return nil, messages
	}
	if pe == nil {
		messages.AddError(WLLLYTBSC_ERR_13001.NewMessage(id))
		return nil, messages
	}

	return (&WallLayoutDE{}).fromPE(pe), messages
}

func (svc *Service) WallLayout__FindByCode(code string) (*WallLayoutDE, *messaging.Messages) {
	messages := messaging.NewMessages()

	pe, dbev := svc.store.DB.WallLayout__GetByCode(code)
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return nil, messages
	}
	if pe == nil {
		messages.AddError(WLLLYTBSC_ERR_14001.NewMessage(code))
		return nil, messages
	}

	return (&WallLayoutDE{}).fromPE(pe), messages
}

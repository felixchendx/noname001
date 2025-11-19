package service

import (
	"noname001/app/base/messaging"

	"noname001/app/module/common/device/store/sqlite"
)

type Web__DeviceListingSearchCriteria = sqlite.Web__DeviceListingSearchCriteria
type Web__DeviceListingSearchResult = sqlite.Web__DeviceListingSearchResult

func (svc *Service) Web__DeviceListing(sc *Web__DeviceListingSearchCriteria) (*Web__DeviceListingSearchResult, *messaging.Messages) {
	messages := messaging.NewMessages()

	searchResult, dbEv := svc.store.DB.Web__DeviceListing(sc)
	if dbEv.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbEv.EventID()))
		return nil, messages
	}

	return searchResult, messages
}

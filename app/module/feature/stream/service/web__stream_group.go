package service

import (
	"noname001/app/base/messaging"
	"noname001/app/module/feature/stream/store/sqlite"
)

// === search params ===
type Web__StreamGroupListingSearchCriteria = sqlite.Web__StreamGroupListingSearchCriteria
type Web__StreamGroupListingSearchResult   = sqlite.Web__StreamGroupListingSearchResult

func (svc *Service) Web__StreamGroupListing(sc *Web__StreamGroupListingSearchCriteria) (*Web__StreamGroupListingSearchResult, *messaging.Messages) {
	messages := messaging.NewMessages()

	searchResult, dbev := svc.store.DB.Web__StreamGroupListing(sc)
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return nil, messages
	}

	return searchResult, messages
}

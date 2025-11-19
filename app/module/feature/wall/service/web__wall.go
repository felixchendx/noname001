package service

import (
	"noname001/app/base/messaging"

	"noname001/app/module/feature/wall/store/sqlite"
)

type Web__WallListingSearchCriteria = sqlite.Web__WallListingSearchCriteria
type Web__WallListingSearchResult   = sqlite.Web__WallListingSearchResult

func (svc *Service) Web__WallListing(sc *Web__WallListingSearchCriteria) (*Web__WallListingSearchResult, *messaging.Messages) {
	messages := messaging.NewMessages()

	searchResult, dbev := svc.store.DB.Web__WallListing(sc)
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return nil, messages
	}

	return searchResult, messages
}

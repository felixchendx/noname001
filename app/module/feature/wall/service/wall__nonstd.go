package service

import (
	"github.com/google/uuid"

	"noname001/app/base/messaging"
	"noname001/app/base/sanitation"

	"noname001/app/module/feature/wall/store/sqlite"
)

/**
	non standard behavior (db transactional 1 header + many items)
	- creating header
	  ~ will generate items with empty streams
	  ~ if partial items are provided, will pass information accordingly, see foot note

	- updating header
	  ~ will GROW   and RETAIN original items (i.e. 4  items -> 16 items)
	  ~ will SHRINK and TRIM   original items (i.e. 16 items -> 4  items)
	  ~ if partial items are provided, will pass information accordingly, see foot note

	- deleting header
	  ~ will delete items

	foot note:
	- partial items pass information to the denoted index (not the array index)

	  = layout items =    = provided partial items on add / edit =
	  [0] idx 1           [0] idx 4
	  [1] idx 2           [1] idx 3
	  [2] idx 3                   ^ is denoting which index to pass info to
	  [3] idx 4
	          ^ is the target index

	  the scenario above will pass information
	  - partial [0] --> layout [3]
	  - partial [1] --> layout [2]
**/

func (svc *Service) Wall__Empty() (*WallDE) {
	return &WallDE{
		Items: make([]*WallItemDE, 0),
	}
}
func (svc *Service) WallItem__Empty() (*WallItemDE) {
	return &WallItemDE{}
}
func (svc *Service) WallItem__EmptyList() ([]*WallItemDE) {
	return make([]*WallItemDE, 0)
}

func (svc *Service) Wall__Add(de *WallDE) (*WallDE, *messaging.Messages) {
	var (
		messages      *messaging.Messages = messaging.NewMessages()
		dbev          *sqlite.DBEvent

		headerPE      *sqlite.WallPE
		// refLayoutPE   *sqlite.WallLayoutPE
		itemPEList    []*sqlite.WallItemPE
		newItemPEList []*sqlite.WallItemPE
	)

	de.sanitize()
	
	// TODO: validation and stuffs
	// validate item index > 0 && < layout item count
	if de.Code == "" {
		messages.AddError(WLLNONSTD_ERR_12501.NewMessage())
	} else {
		isIllegal, illegalChar := sanitation.Code_ContainsIllegalChar(de.Code)
		if isIllegal {
			messages.AddError(WLLNONSTD_ERR_12511.NewMessage(illegalChar, sanitation.CODE__LEGAL_CHARS))
		}
	}
	if messages.HasError() {
		return nil, messages
	}

	// TODO: contextual validations
	// TODO: GET directly from DB for better control and message
	wallLayoutDE, wallLayoutMessages := svc.WallLayout__Find(de.WallLayoutID)
	if wallLayoutMessages.HasError() {
		messages.Append(wallLayoutMessages)
		return nil, messages
	}

	// defaulting
	de.ID = uuid.New().String()
	if de.Name == "" { de.Name = de.Code }

	headerPE, itemPEList = de.toPE()

	newItemPEList = make([]*sqlite.WallItemPE, 0, wallLayoutDE.LayoutItemCount)
	for i := 0; i < wallLayoutDE.LayoutItemCount; i++ {
		newItemPEList = append(newItemPEList, &sqlite.WallItemPE{
			ID: uuid.New().String(),
			WallID: headerPE.ID,
			Index: i + 1,
		})
	}

	for _, itemPE := range itemPEList {
		idx := itemPE.Index - 1

		newItemPE := newItemPEList[idx]
		newItemPE.SourceNodeID = itemPE.SourceNodeID
		newItemPE.StreamCode   = itemPE.StreamCode
	}

	dbev = svc.store.DB.Wall__Insert(headerPE, newItemPEList)
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return nil, messages
	}

	headerPE, itemPEList, dbev = svc.store.DB.Wall__Get(headerPE.ID, true)
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return nil, messages
	}

	messages.AddNotice(WLLNONSTD_NTC_12001.NewMessage(headerPE.Code))
	return (&WallDE{}).fromPE(headerPE, itemPEList), messages
}

func (svc *Service) Wall__Find(id string, withItems bool) (*WallDE, *messaging.Messages) {
	messages := messaging.NewMessages()

	headerPE, itemPEList, dbev := svc.store.DB.Wall__Get(id, withItems)
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return nil, messages
	}
	if headerPE == nil {
		messages.AddError(WLLNONSTD_ERR_13501.NewMessage(id))
		return nil, messages
	}

	return (&WallDE{}).fromPE(headerPE, itemPEList), messages
}

func (svc *Service) Wall__FindByCode(code string, withItems bool) (*WallDE, *messaging.Messages) {
	messages := messaging.NewMessages()

	headerPE, itemPEList, dbev := svc.store.DB.Wall__GetByCode(code, withItems)
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return nil, messages
	}
	if headerPE == nil {
		messages.AddError(WLLNONSTD_ERR_14501.NewMessage(code))
		return nil, messages
	}

	return (&WallDE{}).fromPE(headerPE, itemPEList), messages
}

func (svc *Service) Wall__Edit(id string, de *WallDE) (*WallDE, *messaging.Messages) {
	var (
		messages      *messaging.Messages = messaging.NewMessages()
		dbev          *sqlite.DBEvent

		headerPE      *sqlite.WallPE
		// refLayoutPE   *sqlite.WallLayoutPE
		itemPEList    []*sqlite.WallItemPE
		newItemPEList []*sqlite.WallItemPE
	)

	headerPE, itemPEList, dbev = svc.store.DB.Wall__Get(id, true)
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return nil, messages
	}
	if headerPE == nil {
		messages.AddError(WLLNONSTD_ERR_15501.NewMessage(id))
		return nil, messages
	}

	// TODO: use direct query
	wallLayoutDE, wallLayoutMessages := svc.WallLayout__Find(de.WallLayoutID)
	if wallLayoutMessages.HasError() {
		messages.Append(wallLayoutMessages)
		return nil, messages
	}

	// TODO validation and stuffs

	if messages.HasError() {
		return nil, messages
	}

	headerPE.Name  = de.Name
	headerPE.State = de.State
	headerPE.Note  = de.Note
	headerPE.WallLayoutID = wallLayoutDE.ID

	newItemPEList = make([]*sqlite.WallItemPE, 0, wallLayoutDE.LayoutItemCount)
	for i := 0; i < wallLayoutDE.LayoutItemCount; i++ {
		var itemPE *sqlite.WallItemPE

		if i  < len(itemPEList) {
			itemPE = itemPEList[i]
		} else {
			itemPE = &sqlite.WallItemPE{
				ID: uuid.New().String(),
				WallID: headerPE.ID,
				Index: i + 1,
			}
		}

		newItemPEList = append(newItemPEList, itemPE)
	}

	for _, itemDE := range de.Items {
		idx := itemDE.Index - 1

		newItemPE := newItemPEList[idx]
		newItemPE.SourceNodeID = itemDE.SourceNodeID
		newItemPE.StreamCode   = itemDE.StreamCode
	}

	dbev = svc.store.DB.Wall__Update(headerPE, newItemPEList)
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return nil, messages
	}

	headerPE, itemPEList, dbev = svc.store.DB.Wall__Get(headerPE.ID, true)
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return nil, messages
	}

	messages.AddNotice(WLLNONSTD_NTC_15001.NewMessage(headerPE.Code))
	return (&WallDE{}).fromPE(headerPE, itemPEList), messages
}

func (svc *Service) Wall__Delete(id string) (*messaging.Messages) {
	messages := messaging.NewMessages()

	headerPE, _, dbev := svc.store.DB.Wall__Get(id, false)
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return messages
	}
	if headerPE == nil {
		messages.AddError(WLLNONSTD_ERR_16501.NewMessage(id))
		return messages
	}

	dbev = svc.store.DB.Wall__Delete(headerPE.ID)
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return messages
	}

	messages.AddNotice(WLLNONSTD_NTC_16001.NewMessage(headerPE.Code))
	return messages
}


func (svc *Service) WallItem__Find(id string) (*WallItemDE, *messaging.Messages) {
	messages := messaging.NewMessages()

	pe, dbev := svc.store.DB.WallItem__Get(id)
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return nil, messages
	}
	if pe == nil {
		messages.AddError(WLLNONSTD_ERR_18501.NewMessage(id))
		return nil, messages
	}

	return (&WallItemDE{}).fromPE(pe), messages
}

func (svc *Service) WallItem__Edit(id string, de *WallItemDE) (*WallItemDE, *messaging.Messages) {
	messages := messaging.NewMessages()

	pe, dbev := svc.store.DB.WallItem__Get(id)
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return nil, messages
	}
	if pe == nil {
		messages.AddError(WLLNONSTD_ERR_17501.NewMessage(id))
		return nil, messages
	}

	// TOOD: sanit valid
	// TODO: cross check reference

	if messages.HasError() {
		return nil, messages
	}


	pe.SourceNodeID = de.SourceNodeID
	pe.StreamCode = de.StreamCode

	pe, dbev = svc.store.DB.WallItem__Update(id, pe)
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return nil, messages
	}

	messages.AddNotice(WLLNONSTD_NTC_17001.NewMessage(pe.Index))
	return de.fromPE(pe), messages
}

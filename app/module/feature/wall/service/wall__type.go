package service

import (
	"strings"

	"noname001/app/module/feature/wall/store/sqlite"
)

type WallIdentifier struct {
	ID   string
	Code string
}

// === search params ===

// ============================= VVV header VVV ============================= //
type WallDE struct {
	ID           string
	Code         string
	Name         string
	State        string
	Note         string

	WallLayoutID string

	// === VVV de only VVV ===
	Items        []*WallItemDE
}

func (de *WallDE) sanitize() {
	de.Code = strings.TrimSpace(de.Code)

	de.Name  = strings.TrimSpace(de.Name)
	de.State = strings.TrimSpace(de.State)
	de.Note  = strings.TrimSpace(de.Note)
}

func (de *WallDE) toPE() (*sqlite.WallPE, []*sqlite.WallItemPE) {
	pe := &sqlite.WallPE{}
	pe.ID    = de.ID
	pe.Code  = de.Code
	pe.Name  = de.Name
	pe.State = de.State
	pe.Note  = de.Note

	pe.WallLayoutID = de.WallLayoutID

	itemPEList := make([]*sqlite.WallItemPE, 0) // capped slices better performance ?
	if de.Items != nil {
		for _, itemDE := range de.Items {
			itemPEList = append(itemPEList, itemDE.toPE())
		}
	}

	return pe, itemPEList
}

func (de *WallDE) fromPE(pe *sqlite.WallPE, itemPEList []*sqlite.WallItemPE) (*WallDE) {
	if pe == nil { return nil }
	
	de = &WallDE{}
	de.ID    = pe.ID
	de.Code  = pe.Code
	de.Name  = pe.Name
	de.State = pe.State
	de.Note  = pe.Note

	de.WallLayoutID = pe.WallLayoutID

	de.Items = make([]*WallItemDE, 0)
	if itemPEList != nil {
		for _, itemPE := range itemPEList {
			de.Items = append(de.Items, (&WallItemDE{}).fromPE(itemPE))
		}
	}

	return de
}
// ============================= ^^^ header ^^^ ============================= //

// ============================== VVV item VVV ============================== //
type WallItemDE struct {
	ID           string
	WallID       string
	Index        int

	SourceNodeID string
	StreamCode   string

	TempRelayURL string
}

func (de *WallItemDE) toPE() (*sqlite.WallItemPE) {
	pe := &sqlite.WallItemPE{}
	pe.ID     = de.ID
	pe.WallID = de.WallID
	pe.Index  = de.Index

	pe.SourceNodeID = de.SourceNodeID
	pe.StreamCode   = de.StreamCode

	return pe
}

func (de *WallItemDE) fromPE(pe *sqlite.WallItemPE) (*WallItemDE) {
	if pe == nil { return nil }
	
	de = &WallItemDE{}
	de.ID     = pe.ID
	de.WallID = pe.WallID
	de.Index  = pe.Index

	de.SourceNodeID = pe.SourceNodeID
	de.StreamCode   = pe.StreamCode

	return de
}
// ============================== ^^^ item ^^^ ============================== //

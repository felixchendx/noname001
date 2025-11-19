package service

import (
	"noname001/app/module/feature/wall/store/sqlite"
)

type WallLayoutIdentifier struct {
	ID   string
	Code string
}

// === search ===
type WallLayout__SearchCriteria = sqlite.WallLayout__SearchCriteria
type WallLayout__SearchResult struct {
	Data       []*WallLayoutDE

	Pagination *SearchPagination
}

func (sr *WallLayout__SearchResult) fromStore(_sr *sqlite.WallLayout__SearchResult) (*WallLayout__SearchResult) {
	sr.Data       = make([]*WallLayoutDE, 0, len(_sr.Data))
	sr.Pagination = _sr.Pagination

	for _, pe := range _sr.Data {
		sr.Data = append(sr.Data, (&WallLayoutDE{}).fromPE(pe))
	}

	return sr
}

// === domain entities ===
type WallLayoutDE struct {
	ID              string
	Code            string
	Name            string
	State           string
	Note            string

	LayoutFormation string
	LayoutItemCount int
}

func (de *WallLayoutDE) toPE() (*sqlite.WallLayoutPE) {
	pe := &sqlite.WallLayoutPE{}
	pe.ID    = de.ID
	pe.Code  = de.Code
	pe.Name  = de.Name
	pe.State = de.State
	pe.Note  = de.Note

	pe.LayoutFormation = de.LayoutFormation
	pe.LayoutItemCount = de.LayoutItemCount

	return pe
}

func (de *WallLayoutDE) fromPE(pe *sqlite.WallLayoutPE) (*WallLayoutDE) {
	if pe == nil { return nil }
	
	de = &WallLayoutDE{}
	de.ID    = pe.ID
	de.Code  = pe.Code
	de.Name  = pe.Name
	de.State = pe.State
	de.Note  = pe.Note

	de.LayoutFormation = pe.LayoutFormation
	de.LayoutItemCount = pe.LayoutItemCount

	return de
}

package service

import (
	"strings"

	"noname001/app/module/feature/stream/store/sqlite"
)

// === domain entities ===
type StreamGroupDE struct {
	ID              string
	Code            string
	Name            string
	State           string
	Note            string

	StreamProfileID string
}

func (de *StreamGroupDE) sanitize() {
	de.Code = strings.TrimSpace(de.Code)

	de.Name  = strings.TrimSpace(de.Name)
	de.State = strings.TrimSpace(de.State)
	de.Note  = strings.TrimSpace(de.Note)

	de.StreamProfileID = strings.TrimSpace(de.StreamProfileID)
}

func (de *StreamGroupDE) toPE() (*sqlite.StreamGroupPE) {
	pe := &sqlite.StreamGroupPE{}
	pe.ID = de.ID
	pe.Code = de.Code
	pe.Name = de.Name
	pe.State = de.State
	pe.Note = de.Note

	pe.StreamProfileID = de.StreamProfileID

	return pe
}

func (de *StreamGroupDE) fromPE(pe *sqlite.StreamGroupPE) (*StreamGroupDE) {
	if pe == nil { return nil }
	
	de = &StreamGroupDE{}
	de.ID = pe.ID
	de.Code = pe.Code
	de.Name = pe.Name
	de.State = pe.State
	de.Note = pe.Note

	de.StreamProfileID = pe.StreamProfileID

	return de
}

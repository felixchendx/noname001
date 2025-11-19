package service

import (
	"strings"

	"noname001/app/module/feature/stream/store/sqlite"
)

// === search params ===
type StreamProfile__SearchCriteria = sqlite.StreamProfile__SearchCriteria
type StreamProfile__SearchResult   struct {
	Data       []*StreamProfileDE

	Pagination *SearchPagination
}

func (sr *StreamProfile__SearchResult) fromStore(_sr *sqlite.StreamProfile__SearchResult) (*StreamProfile__SearchResult) {
	sr.Data       = make([]*StreamProfileDE, 0, len(_sr.Data))
	sr.Pagination = _sr.Pagination

	for _, pe := range _sr.Data {
		sr.Data = append(sr.Data, (&StreamProfileDE{}).fromPE(pe))
	}

	return sr
}

// === domain entities ===
type StreamProfileDE struct {
	ID                     string
	Code                   string
	Name                   string
	State                  string
	Note                   string

	TargetVideoCodec       string
	TargetVideoCompression int
	TargetVideoBitrate     int // Custom convenience type here ?

	TargetAudioCodec       string
	TargetAudioCompression int
	TargetAudioBitrate     int
}

func (de *StreamProfileDE) sanitize() {
	de.Code = strings.TrimSpace(de.Code)
	de.Name  = strings.TrimSpace(de.Name)
	de.State = strings.TrimSpace(de.State)
	de.Note  = strings.TrimSpace(de.Note)

	de.TargetVideoCodec = strings.TrimSpace(de.TargetVideoCodec)
	if de.TargetVideoCompression < 0 { de.TargetVideoCompression = 0 }
	if de.TargetVideoCompression > 95 { de.TargetVideoCompression = 95 }

	if de.TargetVideoBitrate < 0 { de.TargetVideoBitrate = 0 }

	de.TargetAudioCodec = strings.TrimSpace(de.TargetAudioCodec)
	if de.TargetAudioCompression < 0 { de.TargetAudioCompression = 0 }
	if de.TargetAudioCompression > 95 { de.TargetAudioCompression = 95 }

	if de.TargetAudioBitrate < 0 { de.TargetAudioBitrate = 0 }
}

func (de *StreamProfileDE) toPE() (*sqlite.StreamProfilePE) {
	pe := &sqlite.StreamProfilePE{}
	pe.ID = de.ID
	pe.Code = de.Code
	pe.Name = de.Name
	pe.State = de.State
	pe.Note = de.Note

	pe.TargetVideoCodec = de.TargetVideoCodec
	pe.TargetVideoCompression = de.TargetVideoCompression
	pe.TargetVideoBitrate = de.TargetVideoBitrate

	pe.TargetAudioCodec = de.TargetAudioCodec
	pe.TargetAudioCompression = de.TargetAudioCompression
	pe.TargetAudioBitrate = de.TargetAudioBitrate

	return pe
}

func (de *StreamProfileDE) fromPE(pe *sqlite.StreamProfilePE) (*StreamProfileDE) {
	if pe == nil { return nil }
	
	de = &StreamProfileDE{}
	de.ID = pe.ID
	de.Code = pe.Code
	de.Name = pe.Name
	de.State = pe.State
	de.Note = pe.Note

	de.TargetVideoCodec = pe.TargetVideoCodec
	de.TargetVideoCompression = pe.TargetVideoCompression
	de.TargetVideoBitrate = pe.TargetVideoBitrate

	de.TargetAudioCodec = pe.TargetAudioCodec
	de.TargetAudioCompression = pe.TargetAudioCompression
	de.TargetAudioBitrate = pe.TargetAudioBitrate

	return de
}

package service

import (
	"strings"

	"noname001/app/module/feature/stream/store/sqlite"
)

type StreamItemIdentifier struct {
	ID   string
	Code string
}

// === search params ===
type StreamItemSearchCriteria struct {

}

// === domain entities ===
type StreamItemDE struct {
	ID               string
	StreamGroupID    string
	Code             string
	Name             string
	State            string
	Note             string

	SourceType       string
	DeviceCode       string
	DeviceChannelID  string
	DeviceStreamType string
	ExternalURL      string
	Filepath         string
	EmbeddedFilepath string
}

func (de *StreamItemDE) sanitize() {
	de.StreamGroupID = strings.TrimSpace(de.StreamGroupID)

	de.Code = strings.TrimSpace(de.Code)
	de.Name  = strings.TrimSpace(de.Name)
	de.State = strings.TrimSpace(de.State)
	de.Note  = strings.TrimSpace(de.Note)

	de.SourceType       = strings.TrimSpace(de.SourceType)
	de.DeviceCode       = strings.TrimSpace(de.DeviceCode)
	de.DeviceChannelID  = strings.TrimSpace(de.DeviceChannelID)
	de.DeviceStreamType = strings.TrimSpace(de.DeviceStreamType)
	de.ExternalURL      = strings.TrimSpace(de.ExternalURL)
	de.Filepath         = strings.TrimSpace(de.Filepath)
	de.EmbeddedFilepath = strings.TrimSpace(de.EmbeddedFilepath)
}

func (de *StreamItemDE) toPE() (*sqlite.StreamItemPE) {
	pe := &sqlite.StreamItemPE{}
	pe.ID = de.ID
	pe.StreamGroupID = de.StreamGroupID
	pe.Code = de.Code
	pe.Name = de.Name
	pe.State = de.State
	pe.Note = de.Note

	pe.SourceType = de.SourceType
	pe.DeviceCode = de.DeviceCode
	pe.DeviceChannelID = de.DeviceChannelID
	pe.DeviceStreamType = de.DeviceStreamType
	pe.ExternalURL = de.ExternalURL
	pe.Filepath = de.Filepath
	pe.EmbeddedFilepath = de.EmbeddedFilepath

	return pe
}

func (de *StreamItemDE) fromPE(pe *sqlite.StreamItemPE) (*StreamItemDE) {
	if pe == nil { return nil }
	
	de = &StreamItemDE{}
	de.ID = pe.ID
	de.StreamGroupID = pe.StreamGroupID
	de.Code = pe.Code
	de.Name = pe.Name
	de.State = pe.State
	de.Note = pe.Note

	de.SourceType = pe.SourceType
	de.DeviceCode = pe.DeviceCode
	de.DeviceChannelID = pe.DeviceChannelID
	de.DeviceStreamType = pe.DeviceStreamType
	de.ExternalURL = pe.ExternalURL
	de.Filepath = pe.Filepath
	de.EmbeddedFilepath = pe.EmbeddedFilepath

	return de
}

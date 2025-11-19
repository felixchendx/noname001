package service

import (
	"strings"

	"noname001/app/module/common/device/store/sqlite"
)

type DeviceIdentifier struct {
	ID   string
	Code string
}

// === domain entities ===
type DeviceDE struct {
	ID    string
	Code  string
	Name  string
	State string
	Note  string

	Protocol string
	Hostname string
	Port     string
	Username string
	Password string
	Brand    string

	FallbackRTSPPort string
}

func (de *DeviceDE) sanitize() {
	de.Code = strings.TrimSpace(de.Code)

	de.Name  = strings.TrimSpace(de.Name)
	de.State = strings.TrimSpace(de.State)
	de.Note  = strings.TrimSpace(de.Note)

	de.Protocol = strings.TrimSpace(de.Protocol)
	de.Hostname = strings.TrimSpace(de.Hostname)
	de.Port     = strings.TrimSpace(de.Port)
	de.Username = strings.TrimSpace(de.Username)
	de.Password = strings.TrimSpace(de.Password)
	de.Brand    = strings.TrimSpace(de.Brand)

	de.FallbackRTSPPort = strings.TrimSpace(de.FallbackRTSPPort)
}

func (de *DeviceDE) toPE() (*sqlite.DevicePE) {
	pe := &sqlite.DevicePE{}
	pe.ID   = de.ID
	pe.Code = de.Code

	pe.Name  = de.Name
	pe.State = de.State
	pe.Note  = de.Note

	pe.Protocol = de.Protocol
	pe.Hostname = de.Hostname
	pe.Port     = de.Port
	pe.Username = de.Username
	pe.Password = de.Password
	pe.Brand    = de.Brand

	pe.FallbackRTSPPort = de.FallbackRTSPPort

	return pe
}

func (de *DeviceDE) fromPE(pe *sqlite.DevicePE) (*DeviceDE) {
	if pe == nil { return nil }
	
	de = &DeviceDE{}
	de.ID   = pe.ID
	de.Code = pe.Code

	de.Name  = pe.Name
	de.State = pe.State
	de.Note  = pe.Note

	de.Protocol = pe.Protocol
	de.Hostname = pe.Hostname
	de.Port     = pe.Port
	de.Username = pe.Username
	de.Password = pe.Password
	de.Brand    = pe.Brand

	de.FallbackRTSPPort = pe.FallbackRTSPPort

	return de
}

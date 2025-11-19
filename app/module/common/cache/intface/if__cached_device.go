package intface

import (
	"time"

	baseTyping "noname001/app/base/typing"
	deviceTyping "noname001/app/base/typing/device"
)

type CachedDeviceFilter struct {}
type CachedDevice struct {
	LastUpdated time.Time

	DeviceSnapshot baseTyping.BaseDeviceSnapshot
	DeviceStatus   *CachedDeviceStatus
}

type CachedDeviceStatus struct {
	DeviceID   string

	NodeID     string
	DeviceCode string
	Timestamp  time.Time

	TextualIndicator  string
	VisualIndicator   string
	AuditoryIndicator string
}

type CachedDeviceEvent struct {
	OriginalDeviceEvent *deviceTyping.LiveDeviceEvent
}


type CachedDeviceStatusFeedSubscription struct {
	ID string

	Channel chan *CachedDeviceStatus
}

type CachedDeviceEventSubscription struct {
	ID string

	Channel chan CachedDeviceEvent
}

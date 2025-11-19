package device

import (
	"time"
)

type LiveDeviceEventCode string

const (
	LIVE_DEVICE_EVENT_CODE__DEACTIVATED  LiveDeviceEventCode = "ld:deactivated"
	LIVE_DEVICE_EVENT_CODE__INIT_BEGIN   LiveDeviceEventCode = "ld:init:begin"
	LIVE_DEVICE_EVENT_CODE__INIT_FAIL    LiveDeviceEventCode = "ld:init:fail"
	LIVE_DEVICE_EVENT_CODE__INIT_OK      LiveDeviceEventCode = "ld:init:ok"
	LIVE_DEVICE_EVENT_CODE__DISCONNECTED LiveDeviceEventCode = "ld:disconnected"
	LIVE_DEVICE_EVENT_CODE__RELOAD_BEGIN LiveDeviceEventCode = "ld:reload:begin"
	LIVE_DEVICE_EVENT_CODE__RELOAD_FAIL  LiveDeviceEventCode = "ld:reload:fail"
	LIVE_DEVICE_EVENT_CODE__RELOAD_OK    LiveDeviceEventCode = "ld:reload:ok"
	LIVE_DEVICE_EVENT_CODE__DESTROYED    LiveDeviceEventCode = "ld:destroyed"
)

type LiveDeviceEvent struct {
	Timestamp time.Time
	// EventID   string
	EventCode LiveDeviceEventCode

	NodeID     string
	DeviceID   string
	DeviceCode string
}

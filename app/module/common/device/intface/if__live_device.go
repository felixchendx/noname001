package intface

import (
	deviceTyping "noname001/app/base/typing/device"
)

type LiveDeviceEventSubscription struct {
	ID string

	Channel chan deviceTyping.LiveDeviceEvent
}

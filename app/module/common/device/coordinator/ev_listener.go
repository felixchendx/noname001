package coordinator

import (
	deviceEv "noname001/app/module/common/device/event"
)

func (coord *Coordinator) eventListeners() {
	deviceEvSub := coord.evHub.NewDeviceEventSubscription()

	go func() {
		evListenerLoop:
		for {
			select {
			case <- coord.context.Done():
				break evListenerLoop

			case ev := <- deviceEvSub.Channel:
				switch ev.EventCode {
				case deviceEv.DEVICE_EVENT_CODE__CREATE: coord.initLiveDevice(ev.DeviceID)
				case deviceEv.DEVICE_EVENT_CODE__UPDATE: coord.reloadLiveDevice(ev.DeviceID)
				case deviceEv.DEVICE_EVENT_CODE__DELETE: coord.destroyLiveDevice(ev.DeviceID)
				}
			}
		}
	}()
}

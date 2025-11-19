package live

import (
	"time"

	baseTyping   "noname001/app/base/typing"
	deviceTyping "noname001/app/base/typing/device"
)

func (lc *LiveCache) interpretDeviceEventToDeviceStatus(device *t_device, deviceEv *deviceTyping.LiveDeviceEvent) {
	var deviceStatus = &t_deviceStatusInterpretation{
		timestamp: time.Now(),
	}

	switch deviceEv.EventCode {
	case deviceTyping.LIVE_DEVICE_EVENT_CODE__DEACTIVATED:
		deviceStatus.textualIndicator = "deactivated"
		deviceStatus.visualIndicator  = visual_indicator__off

	case deviceTyping.LIVE_DEVICE_EVENT_CODE__INIT_BEGIN:
		deviceStatus.textualIndicator = "initializing"
		deviceStatus.visualIndicator  = visual_indicator__green_blink

	case deviceTyping.LIVE_DEVICE_EVENT_CODE__INIT_FAIL:
		deviceStatus.textualIndicator = "init_fail"
		deviceStatus.visualIndicator  = visual_indicator__red_blink

	case deviceTyping.LIVE_DEVICE_EVENT_CODE__INIT_OK:
		deviceStatus.textualIndicator = "init_ok"
		deviceStatus.visualIndicator  = visual_indicator__green_steady
		
	case deviceTyping.LIVE_DEVICE_EVENT_CODE__DISCONNECTED:
		deviceStatus.textualIndicator = "disconnect"
		deviceStatus.visualIndicator  = visual_indicator__yellow_blink

	case deviceTyping.LIVE_DEVICE_EVENT_CODE__RELOAD_BEGIN:
		deviceStatus.textualIndicator = "reloading"
		deviceStatus.visualIndicator  = visual_indicator__green_blink

	case deviceTyping.LIVE_DEVICE_EVENT_CODE__RELOAD_FAIL:
		deviceStatus.textualIndicator = "reload_fail"
		deviceStatus.visualIndicator  = visual_indicator__red_blink

	case deviceTyping.LIVE_DEVICE_EVENT_CODE__RELOAD_OK:
		deviceStatus.textualIndicator = "reload_ok"
		deviceStatus.visualIndicator  = visual_indicator__green_steady

	case deviceTyping.LIVE_DEVICE_EVENT_CODE__DESTROYED:
		deviceStatus.textualIndicator = "deleted"
		deviceStatus.visualIndicator  = visual_indicator__off

	default:
		// noop
	}

	device.deviceStatus = deviceStatus
}

func (lc *LiveCache) interpretDeviceStateToDeviceStatus(device *t_device) {
	var deviceStatus = &t_deviceStatusInterpretation{
		timestamp: time.Now(),
	}

	switch device.deviceSnapshot.Live.State {
	case baseTyping.DEVICE_LIVE_STATE__NEW:
		// noop

	case baseTyping.DEVICE_LIVE_STATE__INACTIVE:
		deviceStatus.textualIndicator = "deactivated"
		deviceStatus.visualIndicator  = visual_indicator__off

	case baseTyping.DEVICE_LIVE_STATE__INIT_BEGIN:
		deviceStatus.textualIndicator = "initializing"
		deviceStatus.visualIndicator  = visual_indicator__green_blink

	case baseTyping.DEVICE_LIVE_STATE__INIT_FAIL:
		deviceStatus.textualIndicator = "init_fail"
		deviceStatus.visualIndicator  = visual_indicator__red_blink

	case baseTyping.DEVICE_LIVE_STATE__INIT_OK:
		deviceStatus.textualIndicator = "init_ok"
		deviceStatus.visualIndicator  = visual_indicator__green_steady

	case baseTyping.DEVICE_LIVE_STATE__DISCONNECTED:
		deviceStatus.textualIndicator = "disconnect"
		deviceStatus.visualIndicator  = visual_indicator__yellow_blink

	case baseTyping.DEVICE_LIVE_STATE__RELOAD_BEGIN:
		deviceStatus.textualIndicator = "reloading"
		deviceStatus.visualIndicator  = visual_indicator__green_blink

	case baseTyping.DEVICE_LIVE_STATE__RELOAD_FAIL:
		deviceStatus.textualIndicator = "reload_fail"
		deviceStatus.visualIndicator  = visual_indicator__red_blink

	case baseTyping.DEVICE_LIVE_STATE__RELOAD_OK:
		deviceStatus.textualIndicator = "reload_ok"
		deviceStatus.visualIndicator  = visual_indicator__green_steady

	case baseTyping.DEVICE_LIVE_STATE__DESTROY:
		deviceStatus.textualIndicator = "deleted"
		deviceStatus.visualIndicator  = visual_indicator__off

	default:
		// noop
	}

	device.deviceStatus = deviceStatus
}

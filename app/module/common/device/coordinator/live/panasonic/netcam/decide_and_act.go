package netcam

import (
	"time"

	"noname001/node/commconf"

	baseTyping "noname001/app/base/typing"
	deviceTyping "noname001/app/base/typing/device"
)

// TODO: periodic check to detect replaced hardware

func (dev *PanasonicNetworkCamera) decideOnWhatToDo() {
	switch dev.cache.lDat.State { // if current state is...
	case baseTyping.DEVICE_LIVE_STATE__NEW         : // do nothing, awaiting init
	case baseTyping.DEVICE_LIVE_STATE__INACTIVE    : dev.doDeactivate()

	case baseTyping.DEVICE_LIVE_STATE__INIT_BEGIN  : dev.doInit()
	case baseTyping.DEVICE_LIVE_STATE__INIT_FAIL   : dev.doRetryInit()
	case baseTyping.DEVICE_LIVE_STATE__INIT_OK     :
		switch dev.cache.lDat.ConnState {
		case baseTyping.DEVICE_CONN_STATE_ALIVE: dev.doKeepAlive()
		default                                : dev.doRetryInit()
		}

	case baseTyping.DEVICE_LIVE_STATE__DISCONNECTED: dev.doReload()

	case baseTyping.DEVICE_LIVE_STATE__RELOAD_BEGIN: dev.doReload()
	case baseTyping.DEVICE_LIVE_STATE__RELOAD_FAIL : dev.doRetryReload()
	case baseTyping.DEVICE_LIVE_STATE__RELOAD_OK   :
		switch dev.cache.lDat.ConnState {
		case baseTyping.DEVICE_CONN_STATE_ALIVE: dev.doKeepAlive()
		default                                : dev.doRetryReload()
		}

	case baseTyping.DEVICE_LIVE_STATE__DESTROY     : dev.doDestroy()
	}
}

func (dev *PanasonicNetworkCamera) doInit() {
	go dev._init()
}
func (dev *PanasonicNetworkCamera) doRetryInit() {
	go dev._init()
}

func (dev *PanasonicNetworkCamera) doReload() {
	go dev._reload()
}
func (dev *PanasonicNetworkCamera) doRetryReload() {
	go dev._reload()
}

func (dev *PanasonicNetworkCamera) doKeepAlive() {
	go dev._keepAlive()
}

func (dev *PanasonicNetworkCamera) doDeactivate() {
	go dev._deactivate()
}

func (dev *PanasonicNetworkCamera) doDestroy() {
	go dev._destroy()
}

func (dev *PanasonicNetworkCamera) _init() {
	dev.cache.lDat.State     = baseTyping.DEVICE_LIVE_STATE__INIT_BEGIN
	dev.cache.lDat.ConnState = baseTyping.DEVICE_CONN_STATE_NEVER
	dev._announce(deviceTyping.LIVE_DEVICE_EVENT_CODE__INIT_BEGIN)

	dev.gatherOperationalCapabilities()
	dev.determineOperationalCapabilities()

	switch dev.cache.lDat.ConnState {
	case baseTyping.DEVICE_CONN_STATE_ALIVE:
		dev.cache.lDat.State = baseTyping.DEVICE_LIVE_STATE__INIT_OK
		dev._announce(deviceTyping.LIVE_DEVICE_EVENT_CODE__INIT_OK)

	default:
		dev.cache.lDat.State = baseTyping.DEVICE_LIVE_STATE__INIT_FAIL
		dev._announce(deviceTyping.LIVE_DEVICE_EVENT_CODE__INIT_FAIL)
	}
}

func (dev *PanasonicNetworkCamera) _reload() {
	dev.cache.lDat.State = baseTyping.DEVICE_LIVE_STATE__RELOAD_BEGIN
	dev._announce(deviceTyping.LIVE_DEVICE_EVENT_CODE__RELOAD_BEGIN)

	dev.gatherOperationalCapabilities()
	dev.determineOperationalCapabilities()

	switch dev.cache.lDat.ConnState {
	case baseTyping.DEVICE_CONN_STATE_ALIVE:
		dev.cache.lDat.State = baseTyping.DEVICE_LIVE_STATE__RELOAD_OK
		dev._announce(deviceTyping.LIVE_DEVICE_EVENT_CODE__RELOAD_OK)

	default:
		dev.cache.lDat.State = baseTyping.DEVICE_LIVE_STATE__RELOAD_FAIL
		dev._announce(deviceTyping.LIVE_DEVICE_EVENT_CODE__RELOAD_FAIL)
	}
}

func (dev *PanasonicNetworkCamera) _keepAlive() {
	connOK, _ := dev.api.TestConnection()
	if connOK {
		dev.cache.lDat.ConnState = baseTyping.DEVICE_CONN_STATE_ALIVE
		dev.cache.lDat.LastSeen  = time.Now()

	} else {
		dev.cache.lDat.State     = baseTyping.DEVICE_LIVE_STATE__DISCONNECTED
		dev.cache.lDat.ConnState = baseTyping.DEVICE_CONN_STATE_LOST
		dev._announce(deviceTyping.LIVE_DEVICE_EVENT_CODE__DISCONNECTED)
	}
}

func (dev *PanasonicNetworkCamera) _deactivate() {
	dev.cache.lDat.State = baseTyping.DEVICE_LIVE_STATE__INACTIVE
	dev._announce(deviceTyping.LIVE_DEVICE_EVENT_CODE__DEACTIVATED)
}

func (dev *PanasonicNetworkCamera) _destroy() {
	dev.cache.lDat.State = baseTyping.DEVICE_LIVE_STATE__DESTROY
	dev._announce(deviceTyping.LIVE_DEVICE_EVENT_CODE__DESTROYED)
	dev.cancel()
}

func (dev *PanasonicNetworkCamera) _announce(evCode deviceTyping.LiveDeviceEventCode) {
	dev.evHub.PublishLiveDeviceEvent(evCode, commconf.ID(), dev.id, dev.code)

	_ = dev.commBundle.DevicePublisher.PublishDeviceEvent(evCode, commconf.ID(), dev.id, dev.code, 60)
}

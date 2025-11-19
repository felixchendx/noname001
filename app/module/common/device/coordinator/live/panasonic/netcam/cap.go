package netcam

import (
	// "time"

	baseTyping "noname001/app/base/typing"
)

// this calls blocking calls
func (dev *PanasonicNetworkCamera) gatherOperationalCapabilities() {
	// fnDelayInBetweenAPICall := func() { time.Sleep(333 * time.Millisecond) }

	dev.cache.opCap.State = baseTyping.DEVICE_CAP_STATE_UPDATING

	{
		// TODO: checks for rtsp stream

		dev.cache.opCap.CanReadRTSPStream = true
		// dev.cache.lDat.ConnState = baseTyping.DEVICE_CONN_STATE_ALIVE
		// dev.cache.lDat.LastSeen  = time.Now()
	}

	// fnDelayInBetweenAPICall()
	dev.capTest__fetchDeviceInfo()

	if !dev.cache.opCap.CanReadDeviceInfo {
		return
	}

	// no need to delay, already available in device info cache
	dev.capTest__fetchDigitalChannelsInfo()
	dev.capTest__fetchStreamInfo()
}

func (dev *PanasonicNetworkCamera) determineOperationalCapabilities() {
	dev.cache.opCap.State = baseTyping.DEVICE_CAP_STATE_NONE

	if dev.cache.opCap.CanReadRTSPStream {
		dev.cache.opCap.State = baseTyping.DEVICE_CAP_STATE_ONLY_STREAM

		if dev.cache.opCap.CanReadDeviceInfo {
			dev.cache.opCap.State = baseTyping.DEVICE_CAP_STATE_PARTIAL

			if dev.cache.opCap.CanReadStreamInfo {
				if dev.cache.opCap.CanReadAnalogInputChannels || dev.cache.opCap.CanReadDigitalInputChannels {
					dev.cache.opCap.State = baseTyping.DEVICE_CAP_STATE_FULL
				}
			}
		}
	}
}

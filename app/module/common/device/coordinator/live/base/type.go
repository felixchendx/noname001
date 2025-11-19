package base

import (
	"noname001/corebiz/integration/base/apicall"
	"noname001/corebiz/integration/base/response"
)

type DeviceResponseIntface = apicall.APICallEventIntface

type DeviceInfo = response.DeviceInfo
type AnalogInputChannel  = response.AnalogInputChannel
type DigitalInputChannel = response.DigitalInputChannel

type StreamInfo struct {
	*response.StreamInfo

	StreamURL string
}

package route

import (
	baseTyping "noname001/app/base/typing"

	wsUtil "noname001/web/base/comm/ws/util"
)

const (
	dd__wsRequestCode__deviceSnapshot string = "/device/snapshot"
	dd__wsRequestCode__tempErrorDetails string = "/device/temp-error-details"

	dd__wsRequestCode__deviceReload string = "/device/reload"
)

type dd__wsEvent struct {
	EventCode string `json:"ev_code"`
}

type dd__wsReply__deviceSnapshot struct {
	wsUtil.BasicReqRep

	DeviceSnapshot *baseTyping.BaseDeviceSnapshot `json:"device_snapshot"`
}

type dd___wsReply__tempErrorDetails struct {
	wsUtil.BasicReqRep

	TempErrorDetails map[string]string `json:"temp_error_details"`
}

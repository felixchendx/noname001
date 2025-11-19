package typing

import (
	"time"
)

// TODO: move these to app/base/typing/device/

type DeviceLiveState string
type DeviceConnState string
type DeviceCapState  string

const (
	DEVICE_LIVE_STATE__NEW               DeviceLiveState = "dls:new"
	DEVICE_LIVE_STATE__INACTIVE          DeviceLiveState = "dls:inactive"
	DEVICE_LIVE_STATE__INIT_BEGIN        DeviceLiveState = "dls:init:begin"
	// DEVICE_LIVE_STATE__INIT_RTN001_BEGIN DeviceLiveState = "dls:init:rtn001:begin"
	// DEVICE_LIVE_STATE__INIT_RTN001_FAIL  DeviceLiveState = "dls:init:rtn001:fail"
	// DEVICE_LIVE_STATE__INIT_RTN001_OK    DeviceLiveState = "dls:init:rtn001:ok"
	// DEVICE_LIVE_STATE__INIT_RTN002_BEGIN DeviceLiveState = "dls:init:rtn002:begin"
	// DEVICE_LIVE_STATE__INIT_RTN002_END   DeviceLiveState = "dls:init:rtn002:end"
	DEVICE_LIVE_STATE__INIT_FAIL         DeviceLiveState = "dls:init:fail"
	DEVICE_LIVE_STATE__INIT_OK           DeviceLiveState = "dls:init:ok"
	DEVICE_LIVE_STATE__DISCONNECTED      DeviceLiveState = "dls:disconnected"
	DEVICE_LIVE_STATE__RELOAD_BEGIN      DeviceLiveState = "dls:reload:begin"
	DEVICE_LIVE_STATE__RELOAD_FAIL       DeviceLiveState = "dls:reload:fail"
	DEVICE_LIVE_STATE__RELOAD_OK         DeviceLiveState = "dls:reload:ok"
	DEVICE_LIVE_STATE__DESTROY           DeviceLiveState = "dls:destroy"

	DEVICE_CONN_STATE_NEVER DeviceConnState = "never"
	DEVICE_CONN_STATE_ALIVE DeviceConnState = "alive"
	DEVICE_CONN_STATE_LOST  DeviceConnState = "lost"

	DEVICE_CAP_STATE_UPDATING    DeviceCapState = "updating"
	DEVICE_CAP_STATE_FULL        DeviceCapState = "full"
	DEVICE_CAP_STATE_PARTIAL     DeviceCapState = "partial"
	DEVICE_CAP_STATE_ONLY_STREAM DeviceCapState = "only_stream"
	DEVICE_CAP_STATE_NONE        DeviceCapState = "none"
)

type BaseDeviceSnapshot struct {
	Persistence BaseDevicePersistenceData `json:"persistence"`
	Live        BaseDeviceLiveData        `json:"live"`

	OpCap       BaseDeviceOpCap           `json:"op_cap"`
	Hardware    BaseDeviceHardwareData    `json:"hardware"`
}

// data that's persisted from user input
type BaseDevicePersistenceData struct {
	ID    string `json:"id"`
	Code  string `json:"code"`
	Name  string `json:"name"`
	State string `json:"state"`
	Brand string `json:"brand"`
}

// data that's describing "liveness" stuffs
type BaseDeviceLiveData struct {
	State DeviceLiveState `json:"state"`

	ConnState        DeviceConnState `json:"conn_state"`
	ConnStateMessage string          `json:"conn_state_msg"`
	LastSeen         time.Time       `json:"last_seen"`
}

// data that's indicating operational capabilities and
// what hardware data are available
// TODO: find better structure for cap stuffs...
type BaseDeviceOpCap struct {
	State DeviceCapState   `json:"state"`

	CanReadDeviceInfo bool `json:"can_read_device_info"`

	CanReadAnalogInputChannels  bool `json:"can_read_analog_input_channels"`
	CanReadDigitalInputChannels bool `json:"can_read_digital_input_channels"`

	CanReadStreamInfo bool `json:"can_read_stream_info"`
	CanReadRTSPStream bool `json:"can_read_rtsp_stream"`
}

// data that's originating from hardware
type BaseDeviceHardwareData struct {
	LastUpdated time.Time `json:"last_updated"`

	DeviceID   string `json:"device_id"`
	DeviceName string `json:"device_name"`
	Model      string `json:"model"`
	DeviceType string `json:"device_type"`

	AnalogChannels  []BaseDeviceAnalogChannel  `json:"analog_channels"`
	DigitalChannels []BaseDeviceDigitalChannel `json:"digital_channels"`
}

type BaseDeviceAnalogChannel struct {
	ChannelID   string `json:"channel_id"`
	ChannelName string `json:"channel_name"`
	Enabled     bool   `json:"enabled"`
}

type BaseDeviceDigitalChannel struct {
	ChannelID   string `json:"channel_id"`
	ChannelName string `json:"channel_name"`
	Enabled     bool   `json:"enabled"`
}

type BaseDeviceStreamInfo struct {
	ChannelID   string
	ChannelName string
	Enabled     bool

	StreamURL string

	VideoEnabled          bool
	VideoResolutionWidth  int
	VideoResolutionHeight int
	VideoCodec            string
	VideoFPS              float32
	VideoBitrate          int     // bit/s

	AudioEnabled          bool
	AudioCodec            string
	AudioBitrate          int
}

package intface

type DeviceServiceProviderIntface interface {
	// hmm... cumbersome
}

type DeviceEventProviderIntface interface {
	SubscribeToLiveDeviceEvent() (*LiveDeviceEventSubscription)
	UnsubscribeFromLiveDeviceEvent(*LiveDeviceEventSubscription)
}

func EventProvider() (DeviceEventProviderIntface) {
	return _deviceEventProvider
}

// === ^^^ for those that uses     ^^^ ===
// =======================================
// === VVV for those that provides VVV ===

var (
	_deviceServiceProvider DeviceServiceProviderIntface
	_deviceEventProvider   DeviceEventProviderIntface
)

func AssignDeviceEventProvider(_something DeviceEventProviderIntface) {
	_deviceEventProvider = _something
}

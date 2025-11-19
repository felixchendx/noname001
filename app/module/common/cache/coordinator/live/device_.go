package live

import (
	"sync"
	"time"

	baseTyping "noname001/app/base/typing"
)

const (
	device__external_activity_logs_limit = 20
	device__internal_activity_logs_limit = 10

	device__last_external_activity_threshold int64 = 60 * 3 // 3 mins

	device__ttl int64 = 60 * 60 * 1 // 1 hour
)

func (lc *LiveCache) newDeviceService() (*t_deviceService) {
	return &t_deviceService{
		devices: make(map[string]*t_device),
		sortedDevices: make([]*t_device, 0),

		opFlag  : false,
		opStatus: "",
	}
}

// -----------------------------------------------------------------------------
type t_deviceService struct {
	// === devicesMutex ===
	devicesMutex  sync.Mutex
	devices       map[string]*t_device
	sortedDevices []*t_device
	// === devicesMutex ===

	opFlag   bool
	opStatus string
}

func (deviceService *t_deviceService) addNewDevice(deviceSnapshot *baseTyping.BaseDeviceSnapshot) (*t_device) {
	device := &t_device{}
	device.id = deviceSnapshot.Persistence.ID
	device.code = deviceSnapshot.Persistence.Code

	device.deviceSnapshot = deviceSnapshot
	device.deviceStatus = nil

	device.lastInternalActivityAt = time.Now()
	device.internalActivityLogs = make([]*t_internalActivityLog, 0, device__internal_activity_logs_limit + 1)

	device.lastExternalActivityAt = time.Now()
	device.externalActivityLogs = make([]*t_externalActivityLog, 0, device__external_activity_logs_limit + 1)

	device.staleSince = zero_time
	device.ttl = device__ttl

	deviceService.registerDevice(device)

	return device
}

func (deviceService *t_deviceService) updateDeviceData(device *t_device, deviceSnapshot *baseTyping.BaseDeviceSnapshot) {
	device.deviceSnapshot = deviceSnapshot

	device.defunctAt = zero_time
	device.defunctReason = ""

	device.logExternalActivity("data_update", nil)
}

func (deviceService *t_deviceService) markDeviceAsDefunct(device *t_device, reason string) {
	device.defunctAt = time.Now()
	device.defunctReason = reason

	device.logExternalActivity("defunct", []string{reason})
}

func (deviceService *t_deviceService) checkDeviceStaleness(device *t_device, checkTime time.Time) {
	if (checkTime.Unix() - device.lastExternalActivityAt.Unix()) > device__last_external_activity_threshold {
		if device.staleSince.IsZero() {
			device.staleSince = checkTime
		}
		device.logInternalActivity("stale_check", "stale")

	} else {
		device.staleSince = zero_time
		device.logInternalActivity("stale_check", "ok")
	}
}

func (deviceService *t_deviceService) isDeviceExpired(device *t_device, checkTime time.Time) (bool) {
	var expired bool = false

	if device.staleSince.IsZero() {
		expired = false

	} else {
		expired = (checkTime.Unix() - device.lastExternalActivityAt.Unix()) > device.ttl
	}

	return expired
}

func (deviceService *t_deviceService) removeDevice(device *t_device) {
	deviceService.deregisterDevice(device)
}

// -----------------------------------------------------------------------------
type t_device struct {
	id   string
	code string

	// === data holder ===
	deviceSnapshot *baseTyping.BaseDeviceSnapshot
	deviceStatus   *t_deviceStatusInterpretation
	// === data holder ===

	// === internal ===
	lastInternalActivityAt time.Time
	internalActivityLogs []*t_internalActivityLog

	lastExternalActivityAt time.Time
	externalActivityLogs []*t_externalActivityLog

	defunctAt     time.Time
	defunctReason string
	
	staleSince time.Time
	ttl int64
}
type t_deviceStatusInterpretation struct {
	timestamp time.Time

	textualIndicator  string
	visualIndicator   string
	auditoryIndicator string
}

func (device *t_device) logInternalActivity(activity, result string) {
	var naw = time.Now()

	device.lastInternalActivityAt = naw
	device.internalActivityLogs = append(device.internalActivityLogs, &t_internalActivityLog{
		ts: naw,
		activity: activity,
		result: result,
	})

	if len(device.internalActivityLogs) > device__internal_activity_logs_limit {
		device.internalActivityLogs = device.internalActivityLogs[1:len(device.internalActivityLogs)]
	}
}

func (device *t_device) logExternalActivity(activity string, extra []string) {
	var naw = time.Now()

	device.lastExternalActivityAt = naw
	device.externalActivityLogs = append(device.externalActivityLogs, &t_externalActivityLog{
		ts: naw,
		activity: activity,
		extra: extra,
	})

	if len(device.externalActivityLogs) > device__external_activity_logs_limit {
		device.externalActivityLogs = device.externalActivityLogs[1:len(device.externalActivityLogs)]
	}
}

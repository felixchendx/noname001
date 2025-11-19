package live

import (
	"slices"
)

func (deviceService *t_deviceService) registerDevice(device *t_device) {
	deviceService.devicesMutex.Lock()

	deviceService.devices[device.code] = device
	deviceService._generateSortedDevices()

	deviceService.devicesMutex.Unlock()
}

func (deviceService *t_deviceService) deregisterDevice(device *t_device) {
	deviceService.devicesMutex.Lock()

	delete(deviceService.devices, device.code)
	deviceService._generateSortedDevices()

	deviceService.devicesMutex.Unlock()
}

func (deviceService *t_deviceService) _generateSortedDevices() {
	deviceCodes := make([]string, len(deviceService.devices))
	i := 0
	for _deviceCode, _ := range deviceService.devices {
		deviceCodes[i] = _deviceCode
		i++
	}

	slices.SortFunc(deviceCodes, _caseInsensitiveSort)

	sortedDevices := make([]*t_device, len(deviceCodes))
	for _idx, _deviceCode := range deviceCodes {
		sortedDevices[_idx] = deviceService.devices[_deviceCode]
	}

	deviceService.sortedDevices = sortedDevices
}

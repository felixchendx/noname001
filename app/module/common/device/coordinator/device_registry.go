package coordinator

import (
	"slices"
	"strings"

	liveBase "noname001/app/module/common/device/coordinator/live/base"
)

func (coord *Coordinator) addLiveDeviceToRegistry(deviceID string, liveDevice liveBase.LiveDeviceIntface) {
	deviceCode := liveDevice.PersistenceData().Code

	coord.liveDevicesMutex.Lock()
	coord.liveDevices[deviceID] = liveDevice
	coord.liveDevicesCode[deviceCode] = deviceID
	coord._generateOrderedLiveDevices()
	coord.liveDevicesMutex.Unlock()
}

func (coord *Coordinator) removeLiveDeviceFromRegistry(deviceID string) {
	existingLiveDevice, _ := coord.liveDevices[deviceID]
	deviceCode := existingLiveDevice.PersistenceData().Code

	coord.liveDevicesMutex.Lock()
	delete(coord.liveDevices, deviceID)
	delete(coord.liveDevicesCode, deviceCode)
	coord._generateOrderedLiveDevices()
	coord.liveDevicesMutex.Unlock()
}

// alphabetically ordered, case insensitive
func (coord *Coordinator) _generateOrderedLiveDevices() {
	deviceCodes := make([]string, len(coord.liveDevicesCode))
	i := 0
	for _ldCode, _ := range coord.liveDevicesCode {
		deviceCodes[i] = _ldCode
		i++
	}

	slices.SortFunc(deviceCodes, coord._caseInsensitiveSort)

	coord.liveDevicesOrdered = make([]liveBase.LiveDeviceIntface, len(coord.liveDevicesCode))
	j := 0
	for _, k := range deviceCodes {
		deviceID := coord.liveDevicesCode[k]
		coord.liveDevicesOrdered[j] = coord.liveDevices[deviceID]
		j++
	}
}

// https://pkg.go.dev/slices#SortFunc
func (coord *Coordinator) _caseInsensitiveSort(a, b string) (int) {
	return strings.Compare(strings.ToLower(a), strings.ToLower(b))
}


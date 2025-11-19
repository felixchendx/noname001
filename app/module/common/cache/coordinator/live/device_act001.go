package live

import (
	baseTyping "noname001/app/base/typing"
)

func (lc *LiveCache) refreshDeviceService(node *t_node, cascadingRefresh bool) {
	deviceServiceInfo, err := lc.fetchDeviceServiceInfo(node.id)
	if err != nil {
		// activity log
		return
	}

	_ = deviceServiceInfo

	node.deviceService.opFlag = true
	node.deviceService.opStatus = "ok"

	if cascadingRefresh {
		lc.refreshDevices(node)
	}
}

func (lc *LiveCache) refreshDevices(node *t_node) {
	deviceSnapshotList, err := lc.fetchDeviceSnapshotList(node.id)
	if err != nil {
		// activity log
		return
	}

	for _, deviceSnapshot := range deviceSnapshotList {
		lc._deviceRefreshRoutine(node, deviceSnapshot)
	}
}

func (lc *LiveCache) refreshDevice(node *t_node, deviceCode string) {
	deviceSnapshot, err := lc.fetchDeviceSnapshot(node.id, deviceCode)
	if err != nil {
		// activity log
		return
	}

	lc._deviceRefreshRoutine(node, deviceSnapshot)
}


func (lc *LiveCache) defunctDeviceService(node *t_node, reason string) {
	node.deviceService.opFlag = false
	node.deviceService.opStatus = reason

	lc.defunctDevices(node, reason)
}

func (lc *LiveCache) defunctDevices(node *t_node, reason string) {
	for _, device := range node.deviceService.devices {
		lc.defunctDevice(node, device, reason)
	}
}

func (lc *LiveCache) defunctDevice(node *t_node, device *t_device, defunctReason string) {
	node.deviceService.markDeviceAsDefunct(device, defunctReason)
}


func (lc *LiveCache) _deviceRefreshRoutine(node *t_node, deviceSnapshot *baseTyping.BaseDeviceSnapshot) {
	seenDevice, alreadySeen := node.deviceService.devices[deviceSnapshot.Persistence.Code]

	if alreadySeen {
		node.deviceService.updateDeviceData(seenDevice, deviceSnapshot)

	} else {
		seenDevice = node.deviceService.addNewDevice(deviceSnapshot)
	}

	lc.interpretDeviceStateToDeviceStatus(seenDevice)

	_ = seenDevice
}

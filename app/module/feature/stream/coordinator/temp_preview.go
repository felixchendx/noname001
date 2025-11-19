package coordinator

import (
	"fmt"
	"time"

	"noname001/app/base/messaging"

	appConstant "noname001/app/constant"

	deviceService "noname001/app/module/common/device/service"

	"noname001/app/module/feature/stream/coordinator/preview"

	mediasrvIntface "noname001/app/module/common/mediasrv/intface"
)

// TEMP
// TODO: move this functionality to mod_device

func (coord *Coordinator) DeviceChannelPreview(requesterHostname, deviceCode, channelID, streamType, streamProtocol string) (string, *messaging.Messages) {
	messages := messaging.NewMessages()
	
	streamInfoModDev, messages01 := deviceService.Instance().GetStreamInfo(
		&deviceService.DeviceIdentifier{Code: deviceCode},
		channelID,
		appConstant.BrandStreamType(streamType),
	)

	if messages01.HasError() {
		messages.Append(messages01)
		return "", messages
	}


	mapKey := fmt.Sprintf("%s_%s_%s", deviceCode, channelID, streamType)
	dcp, ok := coord.dyingPreviews[mapKey]
	if ok {
		dcp.UpdateLastRequested()
	} else {
		dcp = preview.NewDeviceChannelPreview(&preview.DeviceChannelPreviewParams{
			ParentContext: coord.context,
			ParentLogger: coord.logger,
			ParentLogPrefix: coord.logPrefix,
	
			DeviceCode: deviceCode,
			DeviceChannelID: channelID,
			DeviceStreamType: streamType,
	
			StreamSource: streamInfoModDev.StreamURL,
			StreamDestination: mediasrvIntface.Provider().PublishDeviceChannelPreviewURL(deviceCode, channelID, streamType),
		})
		dcp.Start()

		coord.dyingPreviews[mapKey] = dcp
	}


	previewURL := mediasrvIntface.Provider().DeviceChannelPreviewURL(requesterHostname, deviceCode, channelID, streamType, streamProtocol)

	return previewURL, messages
}

func (coord *Coordinator) dcpCleanupWorker() {
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer func() {
			ticker.Stop()
		}()

		WorkerLoop:
		for {
			select {
			case <- coord.context.Done():
				break WorkerLoop
			case <- ticker.C:
				coord.dcpCleanup()
			}
		}
	}()
}

func (coord *Coordinator) dcpCleanup() {
	naw := time.Now().Unix()
	var ttl int64 = 30
	cleanupList := make([]string, 0)

	for k, dcp := range coord.dyingPreviews {
		if naw - dcp.LastRequested() > ttl {
			cleanupList = append(cleanupList, k)
		}
	}

	// mutex aaaaa
	for _, k := range cleanupList {
		dcp, ok := coord.dyingPreviews[k]
		if ok {
			dcp.Stop()
			delete(coord.dyingPreviews, k)
		}
	}
}

package service

import (
	"noname001/app/base/messaging"
)

func (svc *Service) DeviceChannelPreview(requesterHostname, deviceCode, channelID, streamType, streamProtocol string) (string, *messaging.Messages) {
	messages := messaging.NewMessages()

	previewURL, messages01 := svc.coordinator.DeviceChannelPreview(requesterHostname, deviceCode, channelID, streamType, "hls")
	if messages01.HasError() {
		messages.Append(messages01)
		return "", messages
	}

	return previewURL, messages
}

func (svc *Service) GetStreamViewURL(requesterHostname, streamItemCode, streamProtocol string) (string) {
	return svc.coordinator.GetStreamViewURL(requesterHostname, streamItemCode, streamProtocol)
}

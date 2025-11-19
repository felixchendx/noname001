package intface

var (
	mediasrvProvider MediasrvProviderIntface
)

func Provider() (MediasrvProviderIntface) {
	return mediasrvProvider
}

func AssignMediasrvProvider(_instance MediasrvProviderIntface) {
	mediasrvProvider = _instance
}

type MediasrvProviderIntface interface {
	// mod_cache
	SourceStreamURL(mediasrvAuthnPair, mediasrvIP, rtspPort, streamCode string) (rtspURL string)
	RelayedStreamURI(nodeID, streamCode string) (streamURI string)
	RelayedStreamViewURL(requesterHostname, nodeID, streamCode, streamProtocol string) (streamURL string)

	// mod_device
	PublishDeviceChannelPreviewURL(deviceCode, channelID, streamType string) (string)
	DeviceChannelPreviewURL(requesterHostname, deviceCode, channelID, streamType, streamProtocol string) (string)

	// mod_stream
	PublishStreamViewURL(string) (string)
	StreamViewURL(string, string, string) (string)

	// tempe ?
	StreamingPorts() (map[string]string)
	RelayAuthnPair() (string)


	AddPathConfiguration(pathName, source string, onDemand bool) (error)
	ReplacePathConfiguration(pathName, source string, onDemand bool) (error)
	DeletePathConfiguration(pathName string) (error)
}

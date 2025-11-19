package v1

import (
	"time"
)

// API docs: https://bluenviron.github.io/mediamtx/

// === VVV copied from mediamtx's source code, file: internal/defs/api.go VVV ===
// APIError is a generic error.
type APIError struct {
	Error string `json:"error"`
}

// APIPathSourceOrReader is a source or a reader.
type APIPathSourceOrReader struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// APIPath is a path.
type APIPath struct {
	Name          string                  `json:"name"`
	ConfName      string                  `json:"confName"`
	Source        *APIPathSourceOrReader  `json:"source"`
	Ready         bool                    `json:"ready"`
	ReadyTime     *time.Time              `json:"readyTime"`
	Tracks        []string                `json:"tracks"`
	BytesReceived uint64                  `json:"bytesReceived"`
	BytesSent     uint64                  `json:"bytesSent"`
	Readers       []APIPathSourceOrReader `json:"readers"`
}

// APIPathList is a list of paths.
type APIPathList struct {
	ItemCount int        `json:"itemCount"`
	PageCount int        `json:"pageCount"`
	Items     []*APIPath `json:"items"`
}
// === ^^^ copied from mediamtx's source code, file: internal/defs/api.go ^^^ ===

// === VVV copied from mediamtx's source code, file: internal/conf/conf.go VVV ===
type Conf struct {
	// // General
	// LogLevel            LogLevel        `json:"logLevel"`
	// LogDestinations     LogDestinations `json:"logDestinations"`
	// LogFile             string          `json:"logFile"`
	// ReadTimeout         StringDuration  `json:"readTimeout"`
	// WriteTimeout        StringDuration  `json:"writeTimeout"`
	// ReadBufferCount     *int            `json:"readBufferCount,omitempty"` // deprecated
	// WriteQueueSize      int             `json:"writeQueueSize"`
	// UDPMaxPayloadSize   int             `json:"udpMaxPayloadSize"`
	// RunOnConnect        string          `json:"runOnConnect"`
	// RunOnConnectRestart bool            `json:"runOnConnectRestart"`
	// RunOnDisconnect     string          `json:"runOnDisconnect"`

	// // Authentication
	// AuthMethod                AuthMethod                  `json:"authMethod"`
	// AuthInternalUsers         AuthInternalUsers           `json:"authInternalUsers"`
	// AuthHTTPAddress           string                      `json:"authHTTPAddress"`
	// ExternalAuthenticationURL *string                     `json:"externalAuthenticationURL,omitempty"` // deprecated
	// AuthHTTPExclude           AuthInternalUserPermissions `json:"authHTTPExclude"`
	// AuthJWTJWKS               string                      `json:"authJWTJWKS"`

	// // Control API
	// API               bool       `json:"api"`
	// APIAddress        string     `json:"apiAddress"`
	// APIEncryption     bool       `json:"apiEncryption"`
	// APIServerKey      string     `json:"apiServerKey"`
	// APIServerCert     string     `json:"apiServerCert"`
	// APIAllowOrigin    string     `json:"apiAllowOrigin"`
	// APITrustedProxies IPNetworks `json:"apiTrustedProxies"`

	// // Metrics
	// Metrics               bool       `json:"metrics"`
	// MetricsAddress        string     `json:"metricsAddress"`
	// MetricsEncryption     bool       `json:"metricsEncryption"`
	// MetricsServerKey      string     `json:"metricsServerKey"`
	// MetricsServerCert     string     `json:"metricsServerCert"`
	// MetricsAllowOrigin    string     `json:"metricsAllowOrigin"`
	// MetricsTrustedProxies IPNetworks `json:"metricsTrustedProxies"`

	// // PPROF
	// PPROF               bool       `json:"pprof"`
	// PPROFAddress        string     `json:"pprofAddress"`
	// PPROFEncryption     bool       `json:"pprofEncryption"`
	// PPROFServerKey      string     `json:"pprofServerKey"`
	// PPROFServerCert     string     `json:"pprofServerCert"`
	// PPROFAllowOrigin    string     `json:"pprofAllowOrigin"`
	// PPROFTrustedProxies IPNetworks `json:"pprofTrustedProxies"`

	// // Playback
	// Playback               bool       `json:"playback"`
	// PlaybackAddress        string     `json:"playbackAddress"`
	// PlaybackEncryption     bool       `json:"playbackEncryption"`
	// PlaybackServerKey      string     `json:"playbackServerKey"`
	// PlaybackServerCert     string     `json:"playbackServerCert"`
	// PlaybackAllowOrigin    string     `json:"playbackAllowOrigin"`
	// PlaybackTrustedProxies IPNetworks `json:"playbackTrustedProxies"`

	// RTSP server
	RTSP              bool             `json:"rtsp"`
	RTSPDisable       *bool            `json:"rtspDisable,omitempty"` // deprecated
	// Protocols         Protocols        `json:"protocols"` <<< custom type with custom json marshal
	Protocols         []string         `json:"protocols"`
	// Encryption        Encryption       `json:"encryption"` <<< custom type with custom json marshal
	Encryption        string           `json:"encryption"`
	RTSPAddress       string           `json:"rtspAddress"`
	RTSPSAddress      string           `json:"rtspsAddress"`
	RTPAddress        string           `json:"rtpAddress"`
	RTCPAddress       string           `json:"rtcpAddress"`
	MulticastIPRange  string           `json:"multicastIPRange"`
	MulticastRTPPort  int              `json:"multicastRTPPort"`
	MulticastRTCPPort int              `json:"multicastRTCPPort"`
	ServerKey         string           `json:"serverKey"`
	ServerCert        string           `json:"serverCert"`
	// AuthMethods       *RTSPAuthMethods `json:"authMethods,omitempty"` // deprecated
	// RTSPAuthMethods   RTSPAuthMethods  `json:"rtspAuthMethods"` <<< custom type with custom json marshal
	RTSPAuthMethods   []string         `json:"rtspAuthMethods"`

	// // RTMP server
	// RTMP           bool       `json:"rtmp"`
	// RTMPDisable    *bool      `json:"rtmpDisable,omitempty"` // deprecated
	// RTMPAddress    string     `json:"rtmpAddress"`
	// RTMPEncryption Encryption `json:"rtmpEncryption"`
	// RTMPSAddress   string     `json:"rtmpsAddress"`
	// RTMPServerKey  string     `json:"rtmpServerKey"`
	// RTMPServerCert string     `json:"rtmpServerCert"`

	// HLS server
	HLS                bool           `json:"hls"`
	HLSDisable         *bool          `json:"hlsDisable,omitempty"` // deprecated
	HLSAddress         string         `json:"hlsAddress"`
	HLSEncryption      bool           `json:"hlsEncryption"`
	HLSServerKey       string         `json:"hlsServerKey"`
	HLSServerCert      string         `json:"hlsServerCert"`
	HLSAllowOrigin     string         `json:"hlsAllowOrigin"`
	// HLSTrustedProxies  IPNetworks     `json:"hlsTrustedProxies"`
	HLSTrustedProxies  []string       `json:"hlsTrustedProxies"`
	HLSAlwaysRemux     bool           `json:"hlsAlwaysRemux"`
	// HLSVariant         HLSVariant     `json:"hlsVariant"`
	HLSVariant         string         `json:"hlsVariant"`
	HLSSegmentCount    int            `json:"hlsSegmentCount"`
	// HLSSegmentDuration StringDuration `json:"hlsSegmentDuration"`
	HLSSegmentDuration string         `json:"hlsSegmentDuration"`
	// HLSPartDuration    StringDuration `json:"hlsPartDuration"`
	HLSPartDuration    string         `json:"hlsPartDuration"`
	// HLSSegmentMaxSize  StringSize     `json:"hlsSegmentMaxSize"`
	HLSSegmentMaxSize  string         `json:"hlsSegmentMaxSize"`
	HLSDirectory       string         `json:"hlsDirectory"`
	// HLSMuxerCloseAfter StringDuration `json:"hlsMuxerCloseAfter"`
	HLSMuxerCloseAfter string         `json:"hlsMuxerCloseAfter"`

	// // WebRTC server
	// WebRTC                      bool             `json:"webrtc"`
	// WebRTCDisable               *bool            `json:"webrtcDisable,omitempty"` // deprecated
	// WebRTCAddress               string           `json:"webrtcAddress"`
	// WebRTCEncryption            bool             `json:"webrtcEncryption"`
	// WebRTCServerKey             string           `json:"webrtcServerKey"`
	// WebRTCServerCert            string           `json:"webrtcServerCert"`
	// WebRTCAllowOrigin           string           `json:"webrtcAllowOrigin"`
	// WebRTCTrustedProxies        IPNetworks       `json:"webrtcTrustedProxies"`
	// WebRTCLocalUDPAddress       string           `json:"webrtcLocalUDPAddress"`
	// WebRTCLocalTCPAddress       string           `json:"webrtcLocalTCPAddress"`
	// WebRTCIPsFromInterfaces     bool             `json:"webrtcIPsFromInterfaces"`
	// WebRTCIPsFromInterfacesList []string         `json:"webrtcIPsFromInterfacesList"`
	// WebRTCAdditionalHosts       []string         `json:"webrtcAdditionalHosts"`
	// WebRTCICEServers2           WebRTCICEServers `json:"webrtcICEServers2"`
	// WebRTCHandshakeTimeout      StringDuration   `json:"webrtcHandshakeTimeout"`
	// WebRTCTrackGatherTimeout    StringDuration   `json:"webrtcTrackGatherTimeout"`
	// WebRTCICEUDPMuxAddress      *string          `json:"webrtcICEUDPMuxAddress,omitempty"`  // deprecated
	// WebRTCICETCPMuxAddress      *string          `json:"webrtcICETCPMuxAddress,omitempty"`  // deprecated
	// WebRTCICEHostNAT1To1IPs     *[]string        `json:"webrtcICEHostNAT1To1IPs,omitempty"` // deprecated
	// WebRTCICEServers            *[]string        `json:"webrtcICEServers,omitempty"`        // deprecated

	// // SRT server
	// SRT        bool   `json:"srt"`
	// SRTAddress string `json:"srtAddress"`

	// // Record (deprecated)
	// Record                *bool           `json:"record,omitempty"`                // deprecated
	// RecordPath            *string         `json:"recordPath,omitempty"`            // deprecated
	// RecordFormat          *RecordFormat   `json:"recordFormat,omitempty"`          // deprecated
	// RecordPartDuration    *StringDuration `json:"recordPartDuration,omitempty"`    // deprecated
	// RecordSegmentDuration *StringDuration `json:"recordSegmentDuration,omitempty"` // deprecated
	// RecordDeleteAfter     *StringDuration `json:"recordDeleteAfter,omitempty"`     // deprecated

	// // Path defaults
	// PathDefaults Path `json:"pathDefaults"`

	// // Paths
	// OptionalPaths map[string]*OptionalPath `json:"paths"`
	// Paths         map[string]*Path         `json:"-"` // filled by Check()
}
// === ^^^ copied from mediamtx's source code, file: internal/conf/conf.go ^^^ ===

// many many more fields not included, see doc
type SimplePathConfiguration struct {
	Name   string `json:"name"`
	Source string `json:"source"`

	SourceOnDemand bool `json:"sourceOnDemand"`
}

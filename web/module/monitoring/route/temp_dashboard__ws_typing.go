package route

import (
	"fmt"
	"time"

	wsUtil "noname001/web/base/comm/ws/util"

	cacheIntface "noname001/app/module/common/cache/intface"
)

const (
	td__wsReqCode__nodeInfoListing string = "/ni/l"
	td__wsReqCode__deviceInfoListing string = "/di/l"
	td__wsReqCode__streamInfoListing string = "/si/l"
	td__wsReqCode__streamInfoItem string = "/si/i"
)

// === REQ ===
type td__wsReqMessage__streamInfoItem struct {
	Header *wsUtil.BasicReqHeader `json:"_bh"`

	Payload struct{
		NodeID     string `json:"node_id"`
		StreamCode string `json:"stream_code"`
	 } `json:"_bp"`
}

// === REP ===
type td__wsRepPayload__nodeInfoListing struct {
	NodeInfoListing []*td__nodeInfo `json:"nil"`
}

type td__wsRepPayload__deviceInfoListing struct {
	DeviceInfoListing []*td__deviceInfo `json:"dil"`
}

type td__wsRepPayload__streamInfoListing struct {
	StreamInfoListing []*td__streamInfo `json:"sil"`
}

type td__wsRepPayload__streamInfoItem struct {
	StreamInfoItem *td__streamInfo `json:"sii"`
}

// === DF ===
type td__wsDatafeed__nodeStatusInfo struct {
	NodeStatusInfo *td__nodeStatusInfo `json:"nsi"`
}

type td__wsDatafeed__deviceStatusInfo struct {
	DeviceStatusInfo *td__deviceStatusInfo `json:"dsi"`
}

type td__wsDatafeed__streamStatusInfo struct {
	StreamStatusInfo *td__streamStatusInfo `json:"ssi"`
}

// === COMMON ===
type td__nodeInfo struct {
	ID             string    `json:"id"`
	LastActivityAt time.Time `json:"laa"`

	ResourceInfo *td__nodeResourceInfo `json:"nri"`
	StatusInfo   *td__nodeStatusInfo   `json:"nsi"`
}
type td__nodeStatusInfo struct {
	ID             string    `json:"id"`
	LastActivityAt time.Time `json:"laa"`

	TextualIndicator  string `json:"txid"`
	VisualIndicator   string `json:"vsid"`
	AuditoryIndicator string `json:"adid"`
}
type td__nodeResourceInfo struct {
	ID string `json:"id"`

	DisplayCPUPct     string `json:"dcp"`
	DisplayMemTotal   string `json:"dmt"`
	DisplayMemUsed    string `json:"dmu"`
	DisplayMemUsedPct string `json:"dmup"`
}


type td__deviceInfo struct {
	ID   string `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`

	StatusInfo *td__deviceStatusInfo `json:"dsi"`
}
type td__deviceStatusInfo struct {
	ID             string    `json:"id"`
	LastActivityAt time.Time `json:"laa"`

	TextualIndicator  string `json:"txid"`
	VisualIndicator   string `json:"vsid"`
	AuditoryIndicator string `json:"adid"`
}


type td__streamInfo struct {
	ID   string `json:"id"`
	Code string `json:"code"`
	// Name string `json:"name"`

	StatusInfo *td__streamStatusInfo `json:"ssi"`
}
type td__streamStatusInfo struct {
	ID             string    `json:"id"`
	LastActivityAt time.Time `json:"laa"`

	TextualIndicator  string `json:"txid"`
	VisualIndicator   string `json:"vsid"`
	AuditoryIndicator string `json:"adid"`
}

func (ni *td__nodeInfo) mapFromCachedNode(cachedNode *cacheIntface.CachedNode) {
	if cachedNode == nil { return }

	ni.ID = cachedNode.NodeSnapshot.ID
	// ni.State = string(cachedNode.NodeSnapshot.State)

	ni.LastActivityAt = cachedNode.LastActivityAt

	ni.ResourceInfo = &td__nodeResourceInfo{}
	ni.ResourceInfo.mapFromCachedNodeResource(cachedNode.NodeResource)

	ni.StatusInfo = &td__nodeStatusInfo{}
	ni.StatusInfo.mapFromCachedNodeStatus(cachedNode.NodeStatus)
}

func (nsi *td__nodeStatusInfo) mapFromCachedNodeStatus(cachedNodeStatus *cacheIntface.CachedNodeStatus) {
	if cachedNodeStatus == nil { return }

	nsi.ID = cachedNodeStatus.ID
	nsi.LastActivityAt = cachedNodeStatus.Timestamp

	nsi.TextualIndicator  = cachedNodeStatus.TextualIndicator
	nsi.VisualIndicator   = cachedNodeStatus.VisualIndicator
	nsi.AuditoryIndicator = cachedNodeStatus.AuditoryIndicator
}

func (nri *td__nodeResourceInfo) mapFromCachedNodeResource(cachedNodeResource *cacheIntface.CachedNodeResource) {
	if cachedNodeResource == nil { return }

	nri.ID = cachedNodeResource.ID

	nri.DisplayCPUPct = formatAsDisplayPercentage(cachedNodeResource.NodeResource.CPUPercent)
	nri.DisplayMemTotal = formatAsDisplaySize(cachedNodeResource.NodeResource.MemoryTotal)
	nri.DisplayMemUsed = formatAsDisplaySize(cachedNodeResource.NodeResource.MemoryUsed)
	nri.DisplayMemUsedPct = formatAsDisplayPercentage(cachedNodeResource.NodeResource.MemoryUsedPercent)
}

func (di *td__deviceInfo) mapFromCachedDevice(cachedDevice *cacheIntface.CachedDevice) {
	if cachedDevice == nil { return }

	di.ID = cachedDevice.DeviceSnapshot.Persistence.ID
	di.Code = cachedDevice.DeviceSnapshot.Persistence.Code
	di.Name = cachedDevice.DeviceSnapshot.Persistence.Name

	di.StatusInfo = &td__deviceStatusInfo{}
	di.StatusInfo.mapFromCachedDeviceStatus(cachedDevice.DeviceStatus)
}

func (dsi *td__deviceStatusInfo) mapFromCachedDeviceStatus(cachedDeviceStatus *cacheIntface.CachedDeviceStatus) {
	if cachedDeviceStatus == nil { return }

	dsi.ID = cachedDeviceStatus.DeviceID
	dsi.LastActivityAt = cachedDeviceStatus.Timestamp

	dsi.TextualIndicator  = cachedDeviceStatus.TextualIndicator
	dsi.VisualIndicator   = cachedDeviceStatus.VisualIndicator
	dsi.AuditoryIndicator = cachedDeviceStatus.AuditoryIndicator
}

func (si *td__streamInfo) mapFromCachedStream(cachedStream *cacheIntface.CachedStream) {
	if cachedStream == nil { return }

	si.ID = cachedStream.ID
	si.Code = cachedStream.Code

	si.StatusInfo = &td__streamStatusInfo{}
	si.StatusInfo.mapFromCachedStreamStatus(cachedStream.StreamStatus)
}
func (ssi *td__streamStatusInfo) mapFromCachedStreamStatus(cachedStreamStatus *cacheIntface.CachedStreamStatus) {
	if cachedStreamStatus == nil { return }

	ssi.ID = cachedStreamStatus.StreamID
	ssi.LastActivityAt = cachedStreamStatus.Timestamp

	ssi.TextualIndicator  = cachedStreamStatus.TextualIndicator
	ssi.VisualIndicator   = cachedStreamStatus.VisualIndicator
	ssi.AuditoryIndicator = cachedStreamStatus.AuditoryIndicator
}

func formatAsDisplayPercentage(v float64) (string) {
	return fmt.Sprintf("%.2f %%", v)
}

func formatAsDisplaySize(v uint64) (string) {
	inMB := float64(v) / 1024 / 1024

	return fmt.Sprintf("%.0f MB", inMB)
}

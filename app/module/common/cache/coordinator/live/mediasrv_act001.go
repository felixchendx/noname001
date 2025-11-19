package live

import (
	"time"
	"sync"

	mediamtxUtil "noname001/thirdparty/mediamtx/util"
)

// TODO: admin-area selectable ip to use
// TODO: use opFlag as main indicator

func (lc *LiveCache) refreshMediaServer(node *t_node) {
	var (
		mediasrvSnapshot = node.nodeSnapshot.AppSnapshot.TempMediasrvSnapshot

		needToResolveRTSPConnection bool = false
	)

	var (
		prevUsedIP    = node.mediaServer.ipToUse
		prevRTSPPort  = node.mediaServer.ports["rtsp"]

		stillSameUsedIP    = false
		stillSameRTSPPort  = false
	)

	node.mediaServer.ports     = mediasrvSnapshot.Ports
	node.mediaServer.authnPair = mediasrvSnapshot.AuthnPair

	for _, nodeIP := range node.nodeSnapshot.IPs {
		if prevUsedIP == nodeIP {
			stillSameUsedIP = true
			break
		}
	}
	stillSameRTSPPort = (prevRTSPPort == node.mediaServer.ports["rtsp"])

	needToResolveRTSPConnection = (!stillSameUsedIP || !stillSameRTSPPort)

	if needToResolveRTSPConnection {
		rtspPort := node.mediaServer.ports["rtsp"]

		pingTs, pingResults := time.Now(), make([]*t_pingResult, 0, 8)
		wg := new(sync.WaitGroup)
		for _, nodeIP := range node.nodeSnapshot.IPs {
			wg.Add(1)
			go func() {
				defer wg.Done()
				reachable, err := mediamtxUtil.TestRTSPConnWithRetry(nodeIP, rtspPort)
				pingResults = append(pingResults, &t_pingResult{
					ip       : nodeIP,
					reachable: reachable,
					err      : err,
				})
			}()
		}
		wg.Wait()

		node.mediaServer.pingTs      = pingTs
		node.mediaServer.pingResults = pingResults

		node.mediaServer.ipToUse = ""
		for _, pingResult := range node.mediaServer.pingResults {
			if pingResult.reachable {
				node.mediaServer.ipToUse = pingResult.ip
				break
			}
		}
	}

	if node.mediaServer.ipToUse == "" {
		node.mediaServer.opFlag = false
		node.mediaServer.opStatus = "unreachable_host"
		lc.logger.Warnf("%s: mediaserver on node '%s' unreachable", lc.logPrefix, node.id)

	} else {
		node.mediaServer.opFlag = true
		node.mediaServer.opStatus = "ok"
	}
}

func (lc *LiveCache) defunctMediaServer(node *t_node, reason string) {
	node.mediaServer.opFlag = false
	node.mediaServer.opStatus = reason

	node.mediaServer.ipToUse     = ""
	node.mediaServer.pingTs      = zero_time
	node.mediaServer.pingResults = make([]*t_pingResult, 0)
}

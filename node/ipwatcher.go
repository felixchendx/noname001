package node

import (
	"time"

	"noname001/internal/util/netutil"

	nodeTyping "noname001/node/base/typing"
	nodeConst  "noname001/node/constant"
)

const (
	_historyLimit = 100
)

var (
	// TEMP
	_simulateIPChange bool = false
	_ipChangeOn       int  = 0
)

type ipCollectionHistory_t struct {
	timestamp time.Time
	ips       []string
	err       error
}

func (node *Node_t) initIPWatcher() {
	node.ips = nil
	node.ipCollectionHistories = make([]*ipCollectionHistory_t, 0, _historyLimit)

	node.cronJobs["ipWatcher"], _ = node.cron.AddFunc(
		nodeConst.CROSSNODE_CRON_TIMING__NODE__IP_WATCHER,
		node._collectIPs,
	)
}

func (node *Node_t) startIPWatcher() {
	node._collectIPs()
}

func (node *Node_t) _collectIPs() {
	var (
		collectedIPs []string
		err          error

		tempChangeIP bool
	)

	if _simulateIPChange {
		_ipChangeOn += 1

		if _ipChangeOn >= 2 && _ipChangeOn <= 2 {
			tempChangeIP = true
		} else {
			tempChangeIP = false
		}

	} else {
		tempChangeIP = false
	}

	if tempChangeIP {
		collectedIPs = make([]string, 0)
		collectedIPs = append(collectedIPs, "1.2.3.4")
		collectedIPs = append(collectedIPs, "5.6.7.8")

	} else {
		// TODO: better way to detect ip change ?
		collectedIPs, err = netutil.CollectIPs()
		if err != nil {
			node.logger.Errorf("[%s] _collectIPs err, %s", node.logPrefix, err.Error())

			node._insertIPHistory(&ipCollectionHistory_t{
				timestamp: time.Now(),
				ips: nil,
				err: err,
			})

			return
		}
	}

	if node._isSameIPSet(node.ips, collectedIPs) {
		// do nothing

	} else {
		node.ips = collectedIPs

		node._insertIPHistory(&ipCollectionHistory_t{
			timestamp: time.Now(),
			ips: collectedIPs,
			err: nil,
		})

		node._announce(nodeTyping.NODE_EVENT_CODE__IP_CHANGE)
	}

	node.logger.Debugf("[%s] _collectIPs, curr ipset: %s", node.logPrefix, node.ips)
}

func (node *Node_t) _isSameIPSet(ipset1, ipset2 []string) (bool) {
	if len(ipset1) != len(ipset2) {
		node.logger.Debugf("[%s] _isSameIPSet, diff len: %s, %s", node.logPrefix, ipset1, ipset2)
		return false
	}

	ipCount := make(map[string]int)

	for _, ip1 := range ipset1 {
		ipCount[ip1] = 1
	}

	for _, ip2 := range ipset2 {
		ct, inMap := ipCount[ip2]
		if inMap {
			ipCount[ip2] = ct + 1
		} else {
			ipCount[ip2] = 1
		}
	}

	for ip, ct := range ipCount {
		_ = ip
		if ct != 2 {
			node.logger.Debugf("[%s] _isSameIPSet, diff ipCount: %s", node.logPrefix, ipCount)
			return false
		}
	}

	node.logger.Debugf("[%s] _isSameIPSet, same ipCount: %s", node.logPrefix, ipCount)
	return true
}

func (node *Node_t) _insertIPHistory(history *ipCollectionHistory_t) {
	node.ipCollectionHistories = append(node.ipCollectionHistories, history)
	
	if len(node.ipCollectionHistories) > _historyLimit {
		node.ipCollectionHistories = node.ipCollectionHistories[1:len(node.ipCollectionHistories)]
	}
}

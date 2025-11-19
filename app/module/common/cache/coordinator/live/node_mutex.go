package live

import (
	"slices"
)

func (lc *LiveCache) registerNode(node *t_node) {
	lc.nodesMutex.Lock()

	lc.nodes[node.id] = node
	lc._generateSortedNodes()

	lc.nodesMutex.Unlock()
}

func (lc *LiveCache) deregisterNode(node *t_node) {
	lc.nodesMutex.Lock()

	delete(lc.nodes, node.id)
	lc._generateSortedNodes()

	lc.nodesMutex.Unlock()
}

func (lc *LiveCache) _generateSortedNodes() {
	nodeIDs := make([]string, len(lc.nodes))
	i := 0
	for _nodeID, _ := range lc.nodes {
		nodeIDs[i] = _nodeID
		i++
	}

	slices.SortFunc(nodeIDs, _caseInsensitiveSort)

	sortedNodes := make([]*t_node, len(nodeIDs))
	for _idx, _nodeID := range nodeIDs {
		sortedNodes[_idx] = lc.nodes[_nodeID]
	}

	lc.sortedNodes = sortedNodes
}

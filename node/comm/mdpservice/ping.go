package mdpservice

func (rcvr *NodeReceiver) Ping(nodeID string) (bool, error) {
	nodeProviderID := nodeProviderID(nodeID)

	reply, err := rcvr.mdclient.Send(nodeProviderID, CMD__PING)
	if err != nil {
		// TODO
		rcvr.logger.Errorf("%s: Ping - send err %s", rcvr.logPrefix, err.Error())
		return false, err
	}

	if len(reply) > 0 && reply[0] == "pong" {
		return true, nil
	}

	return false, nil
}

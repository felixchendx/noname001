package mdpservice

func (rcvr *StreamReceiver) Ping(providerNodeID string) (bool, error) {
	streamProviderID := StreamProviderID(providerNodeID)

	reply, err := rcvr.mdclient.Send(streamProviderID, CMD_PING)
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

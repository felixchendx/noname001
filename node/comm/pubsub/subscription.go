package pubsub

// dilemma: multi topic subscriptions bad ?
// func (sub *NodeSubscriber) SubscribeOn(subCodes []SubscriptionCode) (*NodePubsubSubscription) {
// 	subscription := &NodePubsubSubscription{
// 		id: uuid.New().String(),
// 	}

// 	switch subCode {
// 	case SUBSCODE__NODE_LIVENESS:
// 		subscription.ChannelNodeLiveness = make(chan *ChannelNodeLiveness, 64)
// 	default:
// 	}
	
// 	sub.subscriptions[subscription.id] = subscription
// }

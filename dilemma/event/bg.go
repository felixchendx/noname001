package event

func (evHub *EventHub) bgEventListener() {
	evSub_bgEvent := evHub.NewSubscription(EVTREE__BACKGROUND)

	go func() {
		evListenerLoop:
		for {
			select {
			case <- evHub.Context.Done():
				break evListenerLoop

			case evAny, _ := <- evSub_bgEvent.Channel():
				ev, ok := evAny.(*BgEvent)
				if !ok { continue }

				if ev.Messages.HasError() {
					evHub.Logger.Errorf("%s.bgEv: %s", evHub.LogPrefix, ev.Messages.Dump())
				} else {
					evHub.Logger.Infof("%s.bgEv: %s", evHub.LogPrefix, ev)
				}
			}
		}
	}()
}
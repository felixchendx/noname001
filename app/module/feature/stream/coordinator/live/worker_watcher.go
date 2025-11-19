package live

import (
	"time"

	streamTyping "noname001/app/base/typing/stream"
)

func (liveStream *LiveStream) watcher() {
	var (
		// TODO: configurable
		watchTicker = time.NewTicker(3 * time.Second)
	)

	defer func() {
		watchTicker.Stop()
	}()

	watcherLoop:
	for {
		// selectCase:
		select {
		case <- liveStream.context.Done():
			break watcherLoop

		case <- watchTicker.C:
			liveStream.watcherRoutine()
		}
	}
}

func (liveStream *LiveStream) watcherRoutine() {
	switch liveStream.lDat.State {
	case streamTyping.LIVE_STATE__NEW         : // wait for user-invoked init
	case streamTyping.LIVE_STATE__INACTIVE    : // nothing to do
	case streamTyping.LIVE_STATE__INIT_BEGIN  : // wait for end of init sequence
	case streamTyping.LIVE_STATE__INIT_FAIL   : liveStream.execChan <- exec_init
	case streamTyping.LIVE_STATE__INIT_OK     : // all green
	case streamTyping.LIVE_STATE__RELOAD_BEGIN: // wat for end of reload sequence
	case streamTyping.LIVE_STATE__RELOAD_FAIL : liveStream.execChan <- exec_reload
	case streamTyping.LIVE_STATE__RELOAD_OK   : // all green
	case streamTyping.LIVE_STATE__DESTROY     : // bye
	case streamTyping.LIVE_STATE__BG_FAIL     : liveStream.execChan <- exec_reload
	default:
		return // explicit block
	}
}

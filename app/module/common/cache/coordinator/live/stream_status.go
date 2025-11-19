package live

import (
	"time"

	streamTyping "noname001/app/base/typing/stream"
)

func (lc *LiveCache) interpretStreamEventToStreamStatus(stream *t_stream, streamEv *streamTyping.LiveStreamEvent) {
	var streamStatus = &t_streamStatusInterpretation{
		timestamp: time.Now(),
	}

	switch streamEv.EventCode {
	case streamTyping.LIVE_STREAM_EVENT_CODE__DEACTIVATED:
		streamStatus.textualIndicator = "deactivated"
		streamStatus.visualIndicator  = visual_indicator__off

	case streamTyping.LIVE_STREAM_EVENT_CODE__INIT_BEGIN:
		streamStatus.textualIndicator = "initializing"
		streamStatus.visualIndicator  = visual_indicator__green_blink

	case streamTyping.LIVE_STREAM_EVENT_CODE__INIT_FAIL:
		streamStatus.textualIndicator = "init_fail"
		streamStatus.visualIndicator  = visual_indicator__red_blink

	case streamTyping.LIVE_STREAM_EVENT_CODE__INIT_OK:
		streamStatus.textualIndicator = "init_ok"
		streamStatus.visualIndicator  = visual_indicator__green_steady

	case streamTyping.LIVE_STREAM_EVENT_CODE__RELOAD_BEGIN:
		streamStatus.textualIndicator = "reloading"
		streamStatus.visualIndicator  = visual_indicator__green_blink

	case streamTyping.LIVE_STREAM_EVENT_CODE__RELOAD_FAIL:
		streamStatus.textualIndicator = "reload_fail"
		streamStatus.visualIndicator  = visual_indicator__red_blink

	case streamTyping.LIVE_STREAM_EVENT_CODE__RELOAD_OK:
		streamStatus.textualIndicator = "reload_ok"
		streamStatus.visualIndicator  = visual_indicator__green_steady

	case streamTyping.LIVE_STREAM_EVENT_CODE__DESTROYED:
		streamStatus.textualIndicator = "deleted"
		streamStatus.visualIndicator  = visual_indicator__off

	case streamTyping.LIVE_STREAM_EVENT_CODE__BG_FAIL:
		streamStatus.textualIndicator = "bg_fail"
		streamStatus.visualIndicator  = visual_indicator__red_blink

	default:
		// noop
	}

	stream.streamStatus = streamStatus
}

func (lc *LiveCache) interpretStreamStateToStreamStatus(stream *t_stream) {
	var streamStatus = &t_streamStatusInterpretation{
		timestamp: time.Now(),
	}

	switch stream.streamSnapshot.Live.State {
	case streamTyping.LIVE_STATE__NEW:
		streamStatus.textualIndicator = "new"
		streamStatus.visualIndicator  = visual_indicator__off

	case streamTyping.LIVE_STATE__INACTIVE:
		streamStatus.textualIndicator = "deactivated"
		streamStatus.visualIndicator  = visual_indicator__off

	case streamTyping.LIVE_STATE__INIT_BEGIN:
		streamStatus.textualIndicator = "initializing"
		streamStatus.visualIndicator  = visual_indicator__green_blink

	case streamTyping.LIVE_STATE__INIT_FAIL:
		streamStatus.textualIndicator = "init_fail"
		streamStatus.visualIndicator  = visual_indicator__red_blink

	case streamTyping.LIVE_STATE__INIT_OK:
		streamStatus.textualIndicator = "init_ok"
		streamStatus.visualIndicator  = visual_indicator__green_steady

	case streamTyping.LIVE_STATE__RELOAD_BEGIN:
		streamStatus.textualIndicator = "reloading"
		streamStatus.visualIndicator  = visual_indicator__green_blink

	case streamTyping.LIVE_STATE__RELOAD_FAIL:
		streamStatus.textualIndicator = "reload_fail"
		streamStatus.visualIndicator  = visual_indicator__red_blink

	case streamTyping.LIVE_STATE__RELOAD_OK:
		streamStatus.textualIndicator = "reload_ok"
		streamStatus.visualIndicator  = visual_indicator__green_steady

	case streamTyping.LIVE_STATE__DESTROY:
		streamStatus.textualIndicator = "deleted"
		streamStatus.visualIndicator  = visual_indicator__off

	case streamTyping.LIVE_STATE__BG_FAIL:
		streamStatus.textualIndicator = "bg_fail"
		streamStatus.visualIndicator  = visual_indicator__red_blink

	default:
		// noop
	}

	stream.streamStatus = streamStatus
}

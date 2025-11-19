package live

import (
	"fmt"

	// appConstant "noname001/app/constant"

	"noname001/node/commconf"

	streamTyping "noname001/app/base/typing/stream"
)

type t_execCode string

const (
	exec_init    t_execCode = "init"
	exec_reload  t_execCode = "reload"
	exec_destroy t_execCode = "destroy"

	exec_mark_streamer_fail t_execCode = "bg:streamer:fail"
)

// this executor acts as queue
// so that only ONE data / state altering stuffs can happen at any given time
// alternative to mutex i suppose
func (liveStream *LiveStream) executor() {
	executorLoop:
	for {
		select {
		// case <- liveStream.context.Done():
		// 	break executorLoop

		case execCode, _ := <- liveStream.execChan:
			switch execCode {
			case exec_init   :
				liveStream._init()

			case exec_reload :
				liveStream._reload()

			case exec_destroy:
				liveStream._destroy()
				break executorLoop


			case exec_mark_streamer_fail:
				liveStream._markStreamerFail()
			}
		}
	}
}

func (liveStream *LiveStream) _init() {
	switch liveStream.lDat.State {
	case streamTyping.LIVE_STATE__NEW         : // pass
	case streamTyping.LIVE_STATE__INACTIVE    : // pass
	case streamTyping.LIVE_STATE__INIT_BEGIN  : return // wait for end of init
	case streamTyping.LIVE_STATE__INIT_FAIL   : // pass
	case streamTyping.LIVE_STATE__INIT_OK     : // pass
	case streamTyping.LIVE_STATE__RELOAD_BEGIN: return // reloaded, cannot go back to init 
	case streamTyping.LIVE_STATE__RELOAD_FAIL : return // reloaded, cannot go back to init
	case streamTyping.LIVE_STATE__RELOAD_OK   : return // reloaded, cannot go back to init
	case streamTyping.LIVE_STATE__DESTROY     : return // destroyed, no going back
	case streamTyping.LIVE_STATE__BG_FAIL     : return // bg fail, means init is already ok, call reload instead
	default:
		return // explicit block
	}

	liveStream.lDat.FailState = streamTyping.LIVE_FAIL_STATE__NONE
	liveStream.lDat.State     = streamTyping.LIVE_STATE__INIT_BEGIN

	// TODO: announcing while id and code not initialized...
	liveStream._announce(streamTyping.LIVE_STREAM_EVENT_CODE__INIT_BEGIN)

	liveStream.resetAllData()

	failCode, err := liveStream.initializeStreamData()
	if err != nil {
		prevState := liveStream.lDat.State

		switch failCode {
		case "db"         : liveStream.lDat.FailState = streamTyping.LIVE_FAIL_STATE__DB
		case "dep:profile": liveStream.lDat.FailState = streamTyping.LIVE_FAIL_STATE__DEP_PROFILE
		default           : liveStream.lDat.FailState = streamTyping.LIVE_FAIL_STATE__OTHER
		}

		liveStream.lDat.State = streamTyping.LIVE_STATE__INIT_FAIL
		liveStream._announce(streamTyping.LIVE_STREAM_EVENT_CODE__INIT_FAIL)

		liveStream.insertErrorLog(prevState, liveStream.lDat.State, liveStream.lDat.FailState, err)

		return
	}

	// temp support for group state, group is to be removed
	var opMode = ""

	if liveStream.pDat.GroupState == "active" && liveStream.pDat.State == "active" {
		opMode = "temp:active"
	} else {
		opMode = "temp:inactive"
	}

	switch (opMode) {
	case "temp:active":
		failCode, err := liveStream.loadInputData()
		if err != nil {
			prevState := liveStream.lDat.State

			switch failCode {
			case "dep:device": liveStream.lDat.FailState = streamTyping.LIVE_FAIL_STATE__DEP_DEVICE
			case "dep:file"  : liveStream.lDat.FailState = streamTyping.LIVE_FAIL_STATE__DEP_FILE
			case "other"     : liveStream.lDat.FailState = streamTyping.LIVE_FAIL_STATE__OTHER
			default          : liveStream.lDat.FailState = streamTyping.LIVE_FAIL_STATE__OTHER
			}

			liveStream.lDat.State = streamTyping.LIVE_STATE__INIT_FAIL
			liveStream._announce(streamTyping.LIVE_STREAM_EVENT_CODE__INIT_FAIL)

			liveStream.insertErrorLog(prevState, liveStream.lDat.State, liveStream.lDat.FailState, err)

			return
		}

		liveStream.loadOutputData()

		liveStream.newFfmpegStreamer()
		go liveStream.startFfmpegStreamer()

		liveStream.lDat.FailState = streamTyping.LIVE_FAIL_STATE__NONE
		liveStream.lDat.State     = streamTyping.LIVE_STATE__INIT_OK
		liveStream._announce(streamTyping.LIVE_STREAM_EVENT_CODE__INIT_OK)

	case "temp:inactive":
		fallthrough
	default:
		liveStream.lDat.FailState = streamTyping.LIVE_FAIL_STATE__NONE
		liveStream.lDat.State     = streamTyping.LIVE_STATE__INACTIVE
		liveStream._announce(streamTyping.LIVE_STREAM_EVENT_CODE__DEACTIVATED)

		liveStream.stopFfmpegStreamer()
	}
}

func (liveStream *LiveStream) _reload() {
	switch liveStream.lDat.State {
	case streamTyping.LIVE_STATE__NEW         : return // call init instead
	case streamTyping.LIVE_STATE__INACTIVE    : // pass
	case streamTyping.LIVE_STATE__INIT_BEGIN  : return // wait for end of init sequence
	case streamTyping.LIVE_STATE__INIT_FAIL   : // pass // hmm...
	case streamTyping.LIVE_STATE__INIT_OK     : // pass
	case streamTyping.LIVE_STATE__RELOAD_BEGIN: return // wait for end of reload sequence
	case streamTyping.LIVE_STATE__RELOAD_FAIL : // pass
	case streamTyping.LIVE_STATE__RELOAD_OK   : // pass
	case streamTyping.LIVE_STATE__DESTROY     : return // destroyed, no going back
	case streamTyping.LIVE_STATE__BG_FAIL     : // pass
	default:
		return // explicit block
	}

	liveStream.lDat.FailState = streamTyping.LIVE_FAIL_STATE__NONE
	liveStream.lDat.State     = streamTyping.LIVE_STATE__RELOAD_BEGIN
	liveStream._announce(streamTyping.LIVE_STREAM_EVENT_CODE__RELOAD_BEGIN)

	liveStream.stopFfmpegStreamer()

	liveStream.resetAllData()

	failCode, err := liveStream.reloadStreamData()
	if err != nil {
		prevState := liveStream.lDat.State

		switch failCode {
		case "db"         : liveStream.lDat.FailState = streamTyping.LIVE_FAIL_STATE__DB
		case "dep:profile": liveStream.lDat.FailState = streamTyping.LIVE_FAIL_STATE__DEP_PROFILE
		default           : liveStream.lDat.FailState = streamTyping.LIVE_FAIL_STATE__OTHER
		}

		liveStream.lDat.State = streamTyping.LIVE_STATE__RELOAD_FAIL
		liveStream._announce(streamTyping.LIVE_STREAM_EVENT_CODE__RELOAD_FAIL)

		liveStream.insertErrorLog(prevState, liveStream.lDat.State, liveStream.lDat.FailState, err)

		return
	}


	// temp support for group state, group is to be removed
	var opMode = ""

	if liveStream.pDat.GroupState == "active" && liveStream.pDat.State == "active" {
		opMode = "temp:active"
	} else {
		opMode = "temp:inactive"
	}

	switch opMode {
	case "temp:active":
		failCode, err := liveStream.loadInputData()
		if err != nil {
			prevState := liveStream.lDat.State

			switch failCode {
			case "dep:device": liveStream.lDat.FailState = streamTyping.LIVE_FAIL_STATE__DEP_DEVICE
			case "dep:file"  : liveStream.lDat.FailState = streamTyping.LIVE_FAIL_STATE__DEP_FILE
			case "other"     : liveStream.lDat.FailState = streamTyping.LIVE_FAIL_STATE__OTHER
			default          : liveStream.lDat.FailState = streamTyping.LIVE_FAIL_STATE__OTHER
			}

			liveStream.lDat.State = streamTyping.LIVE_STATE__RELOAD_FAIL
			liveStream._announce(streamTyping.LIVE_STREAM_EVENT_CODE__RELOAD_FAIL)

			liveStream.insertErrorLog(prevState, liveStream.lDat.State, liveStream.lDat.FailState, err)

			return
		}

		liveStream.loadOutputData()

		liveStream.newFfmpegStreamer()
		go liveStream.startFfmpegStreamer()

		liveStream.lDat.FailState = streamTyping.LIVE_FAIL_STATE__NONE
		liveStream.lDat.State     = streamTyping.LIVE_STATE__RELOAD_OK
		liveStream._announce(streamTyping.LIVE_STREAM_EVENT_CODE__RELOAD_OK)

	case "temp:inactive":
		fallthrough
	default:
		liveStream.lDat.FailState = streamTyping.LIVE_FAIL_STATE__NONE
		liveStream.lDat.State     = streamTyping.LIVE_STATE__INACTIVE
		liveStream._announce(streamTyping.LIVE_STREAM_EVENT_CODE__DEACTIVATED)

		liveStream.stopFfmpegStreamer()
	}
}

func (liveStream *LiveStream) _destroy() {
	switch liveStream.lDat.State {
	case streamTyping.LIVE_STATE__DESTROY: return // already destroyed

	default:
		// pass
	}

	liveStream.lDat.FailState = streamTyping.LIVE_FAIL_STATE__NONE
	liveStream.lDat.State     = streamTyping.LIVE_STATE__DESTROY
	liveStream._announce(streamTyping.LIVE_STREAM_EVENT_CODE__DESTROYED)

	liveStream.cancel()
	// close(liveStream.execChan) // what happens to unclosed channel ?
}

func (liveStream *LiveStream) _markStreamerFail() {
	prevState := liveStream.lDat.State
	
	liveStream.lDat.FailState = streamTyping.LIVE_FAIL_STATE__BG_STREAMER
	liveStream.lDat.State     = streamTyping.LIVE_STATE__BG_FAIL
	liveStream._announce(streamTyping.LIVE_STREAM_EVENT_CODE__BG_FAIL)

	var err error
	if liveStream.streamerInstance != nil {
		err = fmt.Errorf(liveStream.streamerInstance.lastStderr)
	}

	liveStream.insertErrorLog(prevState, liveStream.lDat.State, liveStream.lDat.FailState, err)
}

func (liveStream *LiveStream) _announce(evCode streamTyping.LiveStreamEventCode) {
	liveStream.evHub.PublishLiveStreamEvent(
		evCode,
		commconf.ID(),
		liveStream.id, liveStream.code,
	)

	// TODO: move this to the tip of evhub
	_ = liveStream.commBundle.StreamPublisher.PublishStreamEvent(
		evCode,
		commconf.ID(),
		liveStream.id, liveStream.code,
		60,
	)
}

// also, implement queue buffer that merge same exec into one
// e.g. init, reload, reload, reload, destroy, destroy -> init, reload, destroy

// that buffer, also serves as observable pending queues

// and then, add a bit of delay to the end of each execution, enough to not seems glitchy

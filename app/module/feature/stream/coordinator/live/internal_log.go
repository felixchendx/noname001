package live

import (
	"time"

	streamTyping "noname001/app/base/typing/stream"
)

const (
	_historyLimit = 10
)

type t_internalLog struct {
	timestamp time.Time
	prevState string
	currState string
	failState string
	desc      string
}

func (liveStream *LiveStream) insertInternalLog(log *t_internalLog) {
	liveStream.logHistories = append(liveStream.logHistories, log)

	if len(liveStream.logHistories) > _historyLimit {
		liveStream.logHistories = liveStream.logHistories[1:len(liveStream.logHistories)]
	}
}

func (liveStream *LiveStream) insertErrorLog(
	prevState, currState streamTyping.LiveState,
	failState            streamTyping.LiveFailState,
	err                  error,
) {
	var errString = ""
	if err != nil { errString = err.Error() }
	
	liveStream.insertInternalLog(&t_internalLog{
		timestamp : time.Now(),
		prevState: string(prevState),
		currState: string(currState),
		failState: string(failState),
		desc     : errString,
	})
}

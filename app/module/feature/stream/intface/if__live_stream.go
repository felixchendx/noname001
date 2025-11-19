package intface

import (
	streamTyping "noname001/app/base/typing/stream"
)

type LiveStreamEventSubscription struct {
	ID string

	Channel chan streamTyping.LiveStreamEvent
}

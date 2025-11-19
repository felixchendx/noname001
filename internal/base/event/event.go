package event

import (
	"time"

	"noname001/app/base/messaging"
)

// Origin string
// OriginTree []string
// Tags []string

type BaseBgEvent struct {
	Action    string
	Timestamp time.Time
	Messages  *messaging.Messages
}

func NewBaseBgEvent(action string) (*BaseBgEvent) {
	ev := &BaseBgEvent{
		Action   : action,
		Timestamp: time.Now(),
		Messages : messaging.NewMessages(),
	}

	return ev
}

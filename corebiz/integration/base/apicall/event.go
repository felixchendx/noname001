package apicall

import (
	"time"

	"github.com/google/uuid"
)

// type ErrorCategory

type APICallEvent struct {
	id             string
	name           string

	goError        error
	apiError       APIErrorIntface

	serializedData []string

	beginTs        time.Time
	endTs          time.Time
}

func NewEvent(name string) (*APICallEvent) {
	return &APICallEvent{
		id: uuid.New().String(),
		name: name,
		serializedData: make([]string, 0),
		beginTs: time.Now(),
	}
}

func (ev *APICallEvent) EventID() (string) {
	return ev.id
}

func (ev *APICallEvent) EventName() (string) {
	return ev.name
}

func (ev *APICallEvent) MarkWithGoError(err error) {
	ev.endTs = time.Now()
	ev.goError = err
}

func (ev *APICallEvent) MarkWithAPIError(err APIErrorIntface) {
	ev.endTs = time.Now()
	ev.apiError = err
}

func (ev *APICallEvent) DumpThis(serialized ...string) {
	ev.serializedData = append(ev.serializedData, serialized...)
}

func (ev *APICallEvent) MarkAsEnded() {
	ev.endTs = time.Now()
}

func (ev *APICallEvent) SerializedData() ([]string) {
	return ev.serializedData
}

func (ev *APICallEvent) BeginTimestamp() (time.Time) {
	return ev.beginTs
}
func (ev *APICallEvent) EndTimestamp() (time.Time) {
	return ev.endTs
}
func (ev *APICallEvent) Elapsed() (time.Duration) {
	return ev.endTs.Sub(ev.beginTs)
}

// ================= VVV conform to APICallEventIntface VVV ================= //
func (ev *APICallEvent) IsConsideredError() (bool) {
	return ev.goError != nil || ev.apiError != nil
}

func (ev *APICallEvent) IsGoError() (bool) {
	return ev.goError != nil
}
func (ev *APICallEvent) GoError() (error) {
	return ev.goError
}

func (ev *APICallEvent) IsAPIError() (bool) {
	return ev.apiError != nil
}
func (ev *APICallEvent) APIError() (APIErrorIntface) {
	return ev.apiError
}

func (ev *APICallEvent) Error() (string) {
	if ev.goError != nil { return ev.goError.Error() }
	if ev.apiError != nil { return ev.apiError.SimpleError() }
	return ""
}

func (ev *APICallEvent) HasSerializedData() (bool) {
	return ev.serializedData != nil && len(ev.serializedData) > 0
}
// ================= ^^^ conform to APICallEventIntface ^^^ ================= //

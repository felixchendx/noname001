package store

import (
	"github.com/google/uuid"
)

const (
	DEFAULT_ERROR_MESSAGE = "Persistence Store has encountered internal error."
)

type PersistenceEvent struct {
	ID      string
	Name    string

	Message         string
	ConsideredError bool
	OriginalErr     error

	Logged          bool
}

func NewPersistenceEvent(name string) (*PersistenceEvent) {
	return &PersistenceEvent{uuid.New().String(), name, "", false, nil, false}
}

// func (pev *PersistenceEvent) Wrap(bubblinPev *PersistenceEvent) (*PersistenceEvent) {


// 	return pev
// }

func (pev *PersistenceEvent) MarkAsError(err error) {
	pev.Message = DEFAULT_ERROR_MESSAGE
	pev.ConsideredError = true
	pev.OriginalErr = err
}

// conform to interface StoreEventIdentifier
func (pev *PersistenceEvent) EventID() string {
	return pev.ID
}

// conform to interface StoreEventChecker
func (pev *PersistenceEvent) IsError() bool {
	return pev.ConsideredError
}
// conform to interface StoreEventChecker
func (pev *PersistenceEvent) OriErr() error {
	return pev.OriginalErr
}

// conform to interface StoreEventLogger
func (pev *PersistenceEvent) IsLogged() (bool) {
	return pev.Logged
}

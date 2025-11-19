package store

// this is the source interface, the other interfaces are for convenience
type StoreEventer interface {
	StoreEventIdentifier
	StoreEventChecker
	StoreEventLogger
}

type StoreEventIdentifier interface {
	EventID() string
}

type StoreEventChecker interface {
	IsError() bool
	OriErr() error
}

type StoreEventLogger interface {
	IsLogged() bool
}

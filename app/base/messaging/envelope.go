package messaging

type Envelope struct {
	From []string
	To   []string

	Messages *Messages
}

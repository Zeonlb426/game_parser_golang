package contracts

type EventName string

type Event interface {
	Validate() error
}

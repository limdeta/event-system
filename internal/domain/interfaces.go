package domain

type EventPublisher interface {
	Publish(event *Event) error
}

type EventHandler interface {
	Handle(event *Event) error
}

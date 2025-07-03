package application

import (
	"event-system/internal/domain"
)

type EventService struct {
	Validator domain.EventValidator
	Publisher domain.EventPublisher
}

func NewEventService(validator domain.EventValidator, publisher domain.EventPublisher) *EventService {
	return &EventService{
		Validator: validator,
		Publisher: publisher,
	}
}

// Validate and publish event
func (s *EventService) ProcessEvent(event *domain.Event) error {
	if err := s.Validator.Validate(event); err != nil {
		return err
	}
	return s.Publisher.Publish(event)
}

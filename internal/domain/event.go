package domain

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID        string
	Type      string
	Timestamp time.Time
	Payload   map[string]interface{}
}

func NewEvent(eventType string, payload map[string]interface{}) *Event {
	return &Event{
		ID:        uuid.New().String(),
		Type:      eventType,
		Timestamp: time.Now().UTC(),
		Payload:   payload,
	}
}

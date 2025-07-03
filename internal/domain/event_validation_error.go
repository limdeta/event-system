package domain

import (
	"fmt"
)

type EventValidationError struct {
	Reason string
}

func (e *EventValidationError) Error() string {
	return fmt.Sprintf("event validation error: %s", e.Reason)
}

func NewEventValidationError(reason string) *EventValidationError {
	return &EventValidationError{Reason: reason}
}

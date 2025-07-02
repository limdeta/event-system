package application

import (
	"event-system/internal/domain"
	"os"
	"testing"
	"time"
)

type fakePublisher struct {
	called bool
	event  *domain.Event
}

func (f *fakePublisher) Publish(e *domain.Event) error {
	f.called = true
	f.event = e
	return nil
}

func TestEventService_ProcessEvent_Success(t *testing.T) {
	// путь к схемам
	cwd, _ := os.Getwd()
	schemaDir := cwd + "/../schema"
	validator, err := domain.NewJSONSchemaValidator(schemaDir)
	if err != nil {
		t.Fatalf("failed to create validator: %v", err)
	}
	publisher := &fakePublisher{}
	service := NewEventService(validator, publisher)

	event := &domain.Event{
		ID:        "1",
		Type:      "order_status_notification",
		Timestamp: time.Now(),
		Payload: map[string]interface{}{
			"order_id": "12345",
			"status":   "packed",
			"user_id":  "u42",
			"message":  "Ваш заказ собран!",
		},
	}

	err = service.ProcessEvent(event)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !publisher.called {
		t.Error("expected publisher to be called")
	}
}

func TestEventService_ProcessEvent_ValidationError(t *testing.T) {
	cwd, _ := os.Getwd()
	schemaDir := cwd + "/../schema"
	validator, err := domain.NewJSONSchemaValidator(schemaDir)
	if err != nil {
		t.Fatalf("failed to create validator: %v", err)
	}
	publisher := &fakePublisher{}
	service := NewEventService(validator, publisher)

	event := &domain.Event{
		ID:        "1",
		Type:      "order_status_notification",
		Timestamp: time.Now(),
		Payload: map[string]interface{}{
			"order_id": "ABC",     // error: not all digits
			"status":   "unknown", // error: invalid enum
			"user_id":  "",
		},
	}

	err = service.ProcessEvent(event)
	if err == nil {
		t.Error("expected validation error, got nil")
	}
	if publisher.called {
		t.Error("publisher should NOT be called on validation error")
	}
}

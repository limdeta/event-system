package application

import (
	"encoding/json"
	"event-system/internal/domain"
	"event-system/internal/infrastructure"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestEventService_ProcessEvent_Success(t *testing.T) {
	mockPublisher, service := setupEventService(t)

	event := createValidOrderStatusEvent()
	err := service.ProcessEvent(event)

	assertNoError(t, err)
	assertPublisherCalled(t, mockPublisher, event)
}

func TestEventService_ProcessEvent_ValidationError(t *testing.T) {
	mockPublisher, service := setupEventService(t)

	event := createInvalidOrderStatusEvent()
	err := service.ProcessEvent(event)

	assertValidationError(t, err)
	assertPublisherNotCalled(t, mockPublisher)
}

func TestEventService_ProcessEvent_UnknownEventType(t *testing.T) {
	mockPublisher, service := setupEventService(t)

	event := createUnknownTypeEvent()
	err := service.ProcessEvent(event)

	assertSchemaNotFoundError(t, err)
	assertPublisherNotCalled(t, mockPublisher)
}

// === Test Helpers ===

func setupEventService(t *testing.T) (*FakePublisher, *EventService) {
	registry := createTestRegistry(t)
	validator := createTestValidatorWithRegistry(t, registry)
	mockPublisher := &FakePublisher{}
	service := NewEventService(validator, mockPublisher)

	return mockPublisher, service
}

func createValidOrderStatusEvent() *domain.Event {
	return &domain.Event{
		ID:        "test-event-123",
		Type:      "OrderStatusEvent",
		Timestamp: time.Now(),
		Payload: map[string]interface{}{
			"order_id": "12345",
			"status":   "packed",
			"user_id":  "u42",
			"message":  "Ваш заказ собран!",
		},
	}
}

func createInvalidOrderStatusEvent() *domain.Event {
	return &domain.Event{
		ID:        "test-event-456",
		Type:      "OrderStatusEvent",
		Timestamp: time.Now(),
		Payload: map[string]interface{}{
			"order_id": "ABC",     // error: not all digits
			"status":   "unknown", // error: invalid enum
			"user_id":  "",        // error: empty required field
		},
	}
}

func createUnknownTypeEvent() *domain.Event {
	return &domain.Event{
		ID:        "test-event-789",
		Type:      "UnknownEventType",
		Timestamp: time.Now(),
		Payload: map[string]interface{}{
			"some": "data",
		},
	}
}

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func assertValidationError(t *testing.T, err error) {
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
}

func assertSchemaNotFoundError(t *testing.T, err error) {
	if err == nil {
		t.Fatal("expected schema not found error, got nil")
	}
}

func assertPublisherCalled(t *testing.T, publisher *FakePublisher, expectedEvent *domain.Event) {
	if !publisher.called {
		t.Error("expected publisher to be called")
	}

	if publisher.event == nil {
		t.Error("expected event to be passed to publisher")
		return
	}

	if publisher.event.ID != expectedEvent.ID {
		t.Errorf("expected event ID '%s', got '%s'", expectedEvent.ID, publisher.event.ID)
	}

	if publisher.event.Type != expectedEvent.Type {
		t.Errorf("expected event type '%s', got '%s'", expectedEvent.Type, publisher.event.Type)
	}
}

func assertPublisherNotCalled(t *testing.T, publisher *FakePublisher) {
	if publisher.called {
		t.Error("publisher should NOT be called on validation error")
	}
}

// === Setup Helpers ===

func createTestRegistry(t *testing.T) *infrastructure.EventRegistry {
	tempDir := t.TempDir()

	channelsConfig := map[string]interface{}{
		"OrderStatusEvent": map[string]string{
			"endpoint": "order-topic",
			"schema":   "order_status_notification",
			"type":     "kafka",
		},
	}

	channelsData, _ := json.Marshal(channelsConfig)
	channelsPath := filepath.Join(tempDir, "channels.json")
	os.WriteFile(channelsPath, channelsData, 0644)

	registry, err := infrastructure.NewEventRegistryFromFile(channelsPath)
	if err != nil {
		t.Fatalf("failed to create test registry: %v", err)
	}

	return registry
}

func createTestValidatorWithRegistry(t *testing.T, registry *infrastructure.EventRegistry) *domain.JSONSchemaValidator {
	tempDir := t.TempDir()
	schemaDir := filepath.Join(tempDir, "schema")
	os.MkdirAll(schemaDir, 0755)

	orderStatusSchema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"order_id": map[string]interface{}{
				"type":    "string",
				"pattern": "^[0-9]+$",
			},
			"status": map[string]interface{}{
				"type": "string",
				"enum": []string{"created", "packed", "shipped", "delivered"},
			},
			"user_id": map[string]interface{}{
				"type":      "string",
				"minLength": 1,
			},
			"message": map[string]interface{}{
				"type": "string",
			},
		},
		"required": []string{"order_id", "status", "user_id"},
	}

	schemaData, _ := json.Marshal(orderStatusSchema)
	schemaPath := filepath.Join(schemaDir, "order_status_notification.schema.json")
	os.WriteFile(schemaPath, schemaData, 0644)

	validator, err := domain.NewJSONSchemaValidator(schemaDir, registry)
	if err != nil {
		t.Fatalf("failed to create test validator: %v", err)
	}

	return validator
}

// === Mock Publisher ===

type FakePublisher struct {
	called bool
	event  *domain.Event
}

func (f *FakePublisher) Publish(e *domain.Event) error {
	f.called = true
	f.event = e
	return nil
}

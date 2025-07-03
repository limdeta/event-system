package main

import (
	"bytes"
	"encoding/json"
	"event-system/internal/application"
	"event-system/internal/domain"
	"event-system/internal/infrastructure"
	iface "event-system/internal/interface"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestEventPublishingIntegration(t *testing.T) {
	// Given: configured system with validations, schema and mock publisher
	mockPublisher, eventHandler := setupEventSystem(t)

	// When: sending event by HTTP API
	response := sendOrderStatusEvent(t, eventHandler, "order-456", "confirmed")

	// Then: event is processed and published
	assertSuccessfulResponse(t, response)
	assertEventPublished(t, mockPublisher, "OrderStatusEvent", "order-456", "confirmed")

	t.Logf("✅ Integration test passed: OrderStatusEvent published successfully")
}

// === Test Helpers ===

func setupEventSystem(t *testing.T) (*MockPublisher, *iface.EventHandler) {
	registry := createTestEventRegistry(t)
	validator := createTestValidator(t, registry)
	mockPublisher := &MockPublisher{}
	service := application.NewEventService(validator, mockPublisher)
	eventHandler := iface.NewEventHandler(service)

	return mockPublisher, eventHandler
}

func sendOrderStatusEvent(t *testing.T, handler *iface.EventHandler, orderId, status string) *httptest.ResponseRecorder {
	eventPayload := map[string]interface{}{
		"id":        "test-event-123",
		"type":      "OrderStatusEvent",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"payload": map[string]interface{}{
			"orderId": orderId,
			"status":  status,
		},
	}

	payloadBytes, _ := json.Marshal(eventPayload)
	req := httptest.NewRequest("POST", "/event", bytes.NewReader(payloadBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.HandleEvent(w, req)

	return w
}

func assertSuccessfulResponse(t *testing.T, w *httptest.ResponseRecorder) {
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}
}

func assertEventPublished(t *testing.T, publisher *MockPublisher, expectedType, expectedOrderId, expectedStatus string) {
	if len(publisher.PublishedEvents) != 1 {
		t.Fatalf("expected 1 published event, got %d", len(publisher.PublishedEvents))
	}

	event := publisher.PublishedEvents[0].Event

	if event.Type != expectedType {
		t.Errorf("expected event type '%s', got '%s'", expectedType, event.Type)
	}

	payload := event.Payload
	if payload["orderId"] != expectedOrderId {
		t.Errorf("expected orderId '%s', got '%v'", expectedOrderId, payload["orderId"])
	}

	if payload["status"] != expectedStatus {
		t.Errorf("expected status '%s', got '%v'", expectedStatus, payload["status"])
	}
}

// === Setup Helpers ===

func createTestEventRegistry(t *testing.T) *infrastructure.EventRegistry {
	tempDir := t.TempDir()

	channelsConfig := map[string]any{
		"OrderStatusEvent": map[string]string{
			"endpoint": "order-topic",
			"schema":   "order_status_notification_test",
			"type":     "kafka",
		},
	}

	channelsData, _ := json.Marshal(channelsConfig)
	channelsPath := filepath.Join(tempDir, "channels.json")
	os.WriteFile(channelsPath, channelsData, 0644)

	registry, err := infrastructure.NewEventRegistryFromFile(channelsPath)
	if err != nil {
		t.Fatalf("failed to create event registry: %v", err)
	}

	// Verify registry setup
	topic, schema, err := registry.ResolveChannel("OrderStatusEvent")
	if err != nil {
		t.Fatalf("failed to resolve channel: %v", err)
	}

	if topic != "order-topic" || schema != "order_status_notification_test" {
		t.Fatalf("registry misconfigured: topic=%s, schema=%s", topic, schema)
	}

	return registry
}

func createTestValidator(t *testing.T, registry *infrastructure.EventRegistry) *domain.JSONSchemaValidator {
	tempDir := t.TempDir()
	schemaDir := filepath.Join(tempDir, "schema")
	os.MkdirAll(schemaDir, 0755)

	// Create test schema
	orderStatusSchema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"orderId": map[string]interface{}{
				"type": "string",
			},
			"status": map[string]interface{}{
				"type": "string",
				"enum": []string{"pending", "confirmed", "shipped", "delivered"},
			},
		},
		"required": []string{"orderId", "status"},
	}

	schemaData, _ := json.Marshal(orderStatusSchema)
	schemaPath := filepath.Join(schemaDir, "order_status_notification_test.schema.json")
	os.WriteFile(schemaPath, schemaData, 0644)

	validator, err := domain.NewJSONSchemaValidator(schemaDir, registry)
	if err != nil {
		t.Fatalf("failed to create validator: %v", err)
	}

	return validator
}

// === Mock Publisher ===

type MockPublisher struct {
	PublishedEvents []PublishedEvent
}

type PublishedEvent struct {
	Event *domain.Event
	Topic string
}

func (m *MockPublisher) Publish(event *domain.Event) error {
	m.PublishedEvents = append(m.PublishedEvents, PublishedEvent{
		Event: event,
		Topic: "", // TODO: get from EventRegistry
	})
	return nil
}

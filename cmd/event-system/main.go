package main

import (
	"event-system/internal/application"
	"event-system/internal/domain"
	iface "event-system/internal/interface"
	"log"
	"net/http"
)

type fakePublisher struct{}

func (f *fakePublisher) Publish(*domain.Event) error {
	// Implement the required logic or leave as a stub for testing
	return nil
}

func main() {
	validator, err := domain.NewJSONSchemaValidator("./internal/schema")
	if err != nil {
		log.Fatalf("failed to init validator: %v", err)
	}
	publisher := &fakePublisher{}
	service := application.NewEventService(validator, publisher)

	http.HandleFunc("/healthz", iface.HealthCheckHandler)
	http.HandleFunc("/readyz", iface.ReadyCheckHandler("localhost:9092"))

	eventHandler := iface.NewEventHandler(service)
	http.HandleFunc("/event", eventHandler.HandleEvent)

	log.Println("Listening on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

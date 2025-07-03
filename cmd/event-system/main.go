package main

import (
	"event-system/internal/application"
	"event-system/internal/domain"
	"event-system/internal/infrastructure"
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

	registry, err := infrastructure.NewEventRegistryFromFile("config/channels.json")
	if err != nil {
		log.Fatalf("failed to load event registry: %v", err)
	}

	validator, err := domain.NewJSONSchemaValidator("config/schema", registry)
	if err != nil {
		log.Fatalf("failed to init validator: %v", err)
	}

	// Проверяем что конфиг загрузился
	topic, schema, err := registry.ResolveChannel("OrderStatusEvent")
	if err != nil {
		log.Fatalf("failed to resolve channel: %v", err)
	}
	log.Printf("Loaded channel: OrderStatusEvent -> topic: %s, schema: %s", topic, schema)

	////////// Start Admin //////
	adminHandler := iface.NewAdminHandler(registry)

	mux := http.NewServeMux()
	mux.HandleFunc("/admin/reload-channels", adminHandler.ReloadChannels)
	mux.HandleFunc("/admin/channels", adminHandler.GetChannels)

	go func() {
		log.Println("Admin API started at :8081")
		log.Fatal(http.ListenAndServe(":8081", mux))
	}()

	//////////////////////////////

	// Kafka
	brokers := []string{"localhost:9092"} // Kafka brokers

	publisher := infrastructure.NewKafkaPublisher(brokers, registry)
	service := application.NewEventService(validator, publisher)

	// Topics for channels (development)
	allChannels := registry.GetAllChannels()
	var topics []string
	for _, channelInfo := range allChannels {
		if channelInfo.Type == "kafka" {
			topics = append(topics, channelInfo.Endpoint)
		}
	}

	if err := publisher.EnsureTopicsExist(topics); err != nil {
		log.Printf("Warning: failed to create topics: %v", err)
	}

	http.HandleFunc("/healthz", iface.HealthCheckHandler)
	// http.HandleFunc("/readyz", iface.ReadyCheckHandler("localhost:9092"))

	// Event handler
	eventHandler := iface.NewEventHandler(service)
	http.HandleFunc("/event", eventHandler.HandleEvent)

	log.Println("Event system started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

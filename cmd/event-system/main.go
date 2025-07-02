package main

import (
	iface "event-system/internal/interface"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/healthz", iface.HealthCheckHandler)
	http.HandleFunc("/readyz", iface.ReadyCheckHandler("localhost:9092"))
	log.Println("Listening on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

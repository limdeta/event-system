// internal/interface/http_readyz.go
package iface

import (
	"net/http"
	"time"

	"github.com/segmentio/kafka-go"
)

func ReadyCheckHandler(kafkaAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := kafka.Dial("tcp", kafkaAddr)
		if err != nil {
			http.Error(w, "Kafka unavailable", http.StatusServiceUnavailable)
			return
		}
		defer conn.Close()
		conn.SetDeadline(time.Now().Add(2 * time.Second))
		_, err = conn.Brokers()
		if err != nil {
			http.Error(w, "Kafka unavailable", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("READY"))
	}
}

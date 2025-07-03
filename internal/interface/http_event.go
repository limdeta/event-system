package iface

import (
	"encoding/json"
	"errors"
	"event-system/internal/application"
	"event-system/internal/domain"
	"io"
	"net/http"
)

type EventHandler struct {
	Service *application.EventService
}

func NewEventHandler(service *application.EventService) *EventHandler {
	return &EventHandler{Service: service}
}

func (h *EventHandler) HandleEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var event domain.Event
	if err := json.Unmarshal(body, &event); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.Service.ProcessEvent(&event); err != nil {
		var validationErr *domain.EventValidationError
		if errors.As(err, &validationErr) {
			http.Error(w, "Validation error: "+err.Error(), http.StatusBadRequest)
		} else {
			// Internal error
			http.Error(w, "Failed to process event: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

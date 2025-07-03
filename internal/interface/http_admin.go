package iface

import (
	"encoding/json"
	"event-system/internal/infrastructure"
	"net/http"
)

type AdminHandler struct {
	Registry *infrastructure.EventRegistry
}

func NewAdminHandler(reg *infrastructure.EventRegistry) *AdminHandler {
	return &AdminHandler{Registry: reg}
}

// POST /admin/reload-channels
func (h *AdminHandler) ReloadChannels(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := h.Registry.Reload(); err != nil {
		http.Error(w, "reload failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("reloaded"))
}

// GET /admin/channels
func (h *AdminHandler) GetChannels(w http.ResponseWriter, r *http.Request) {
	channels := h.Registry.GetAllChannels()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(channels)
}
